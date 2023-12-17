// Package s3pack provides functions for uploading files to s3.
// This file implements the ObjectList type and its methods. ObjectList is a slice of FileObject pointers. The methods
// are either convenience methods, like count(), or they iterate over the slice and call the corresponding FileObject
// method.
// https://github.com/orme292/s3packer is licensed under the MIT License.
package s3pack

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/orme292/s3packer/config"
)

/*
ObjectList is a slice of FileObject pointers (slices are inherently pointers). Most FileLists methods are just for
convenience -- they iterate over the elements of the slice and call the corresponding FileList method. Exceptions are noted below.

See FileList for more information
*/
type ObjectList []*FileObject

/*
UploadResult is a struct that holds the results of a concurrent upload. It is used with ConcurrentIterateAndExecute.

See ConcurrentIterateAndExecute for more information
*/
type UploadResult struct {
	Err         error
	UploadCount int
	IgnoreCount int
}

/*
NewObjectList is an ObjectList constructor. It takes a slice of paths and returns a slice FileObjects.
It calls NewFileObject on each path and appends the result to the slice of ObjectList. It then calls
SetIgnoreIfLocalNotExists, SetFileSizes, SetChecksum, and TagOrigins to fill in fields for each FileObject.

See NewFileObject for additional information
*/
func NewObjectList(c *config.Configuration, paths []string) (ol ObjectList, err error) {
	for _, path := range paths {
		fo, err := NewFileObject(c, path)
		if err != nil {
			return nil, err
		}
		ol = append(ol, fo)
	}

	ol.SetIgnoreIfLocalNotExists()
	ol.SetFileSizes()
	_ = ol.SetChecksum()
	ol.TagOrigins()
	for k, v := range c.Tags {
		ol.TagAll(k, v)
	}
	return
}

/*
FixRedundantKeys is an ObjectList method. It checks for duplicate occurrences of PrefixedName. If duplicates are found,
then it appends a counter to the end of the PrefixedName.

ALL occurrences are renamed. The first occurrence will get a -0 suffix, the second will get a -1 suffix, etc.
my-file.txt-0
...
my-file.txt-30

This is used when uploading individual files. There are multiple issues with this implementation.
*/
func (ol ObjectList) FixRedundantKeys() error {
	if len(ol) == 0 || len(ol) == 1 {
		return errors.New("FileList is empty or only contains one item")
	}

	ol[0].c.Logger.Debug("Fixing Redundant Keys...")
	occurrences := make(map[string]int)
	for index := range ol {
		if _, ok := occurrences[ol[index].PrefixedName]; ok {
			occurrences[ol[index].PrefixedName] += 1
		} else {
			occurrences[ol[index].PrefixedName] = 1
		}
	}

	for prefixedName, numOccurs := range occurrences {
		if numOccurs > 1 {
			counter := 0
			for index := range ol {
				if ol[index].PrefixedName == prefixedName {
					ol[index].PrefixedName = fmt.Sprintf("%s-%d", ol[index].PrefixedName, counter)
					counter++
				}
			}
		}
	}
	return nil
}

/*
IterateAndExecute is an ObjectList method. It takes a function that takes a FileObject pointer and returns an error.
It iterates over the ObjectList slice and calls the provided function on each FileObject pointer. If the function returns an
error, then it is returned and iteration stops
*/
func (ol ObjectList) IterateAndExecute(fn ObjectListIterationFunc) (err error) {
	for index := range ol {
		if err = fn(ol[index]); err != nil {
			return
		}
	}
	return
}

/*
ObjectListIterationFunc is a function type that takes a FileObject pointer and returns an error. It is used with
ObjectList.IterateAndExecute
*/
type ObjectListIterationFunc func(fo *FileObject) (err error)

type jobTracker struct {
	index int
	group int
}

/*
ConcurrentIterateAndExecute is an ObjectList method, and it is similar to ObjectList.IterateAndExecute. It takes an int
for the number of groups to split the ObjectList into and a function that takes a FileObject pointer and returns an
error. It iterates over the ObjectList slice and calls the provided function on each FileObject pointer. If the function
returns an error, then it is returned and iteration stops.

The difference between this and ObjectList.IterateAndExecute is that this method will split the ObjectList into groups
and execute the provided function concurrently. The number of groups is determined by the provided int.
*/
func (ol ObjectList) ConcurrentIterateAndExecute(groups int, fn ObjectListIterationFunc) (err []error) {
	if len(ol) < groups {
		groups = len(ol)
	}
	var wg sync.WaitGroup
	resultChan := make(chan UploadResult)

	var jobs []jobTracker
	for i := 0; i <= len(ol)-1; i++ {
		jobs = append(jobs, jobTracker{
			index: i,
			group: i % groups,
		})
	}

	for i := 0; i < groups; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, group int, j []jobTracker, fn ObjectListIterationFunc) {
			defer wg.Done()
			var err error
			for index := range j {
				if j[index].group == group {
					err = fn(ol[j[index].index])
					if err != nil {
						break
					}
				}
			}
			resultChan <- UploadResult{
				Err: err,
			}
		}(&wg, i, jobs, fn)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(resultChan)
	}(&wg)

	for result := range resultChan {
		if result.Err != nil {
			err = append(err, result.Err)
		}
	}
	return err
}

/*
GetTotalUploadedBytes is an ObjectList method. It adds the FileSize values of all FileObjects that have the
FileObject.IsUploaded field set to true and returns the total.
*/
func (ol ObjectList) GetTotalUploadedBytes() (total int64) {
	for index := range ol {
		if ol[index].IsUploaded {
			total += ol[index].FileSize
		}
	}
	return
}

/*
IgnoreIfObjectExistsInBucket is an ObjectList method. It iterates through each FileObject in the ObjectList and tries
to retrieve metadata from an S3 object of the same name (s3 key = FileObject.PrefixedName). If the object exists, then
the FileObject.Ignore field is set to true and the FileObject.IgnoreString field is set to ErrIgnoreObjectAlreadyExists.
*/
func (ol ObjectList) IgnoreIfObjectExistsInBucket() {
	if ol[0].c.Options[config.ProfileOptionOverwrite].(bool) || len(ol) == 0 {
		return
	}

	sess, _ := NewSession(ol[0].c)

	svc := s3.New(sess, &aws.Config{})

	for index := range ol {
		_, err := svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(ol[index].c.Bucket[config.ProfileBucketName].(string)),
			Key:    aws.String(ol[index].PrefixedName),
		})
		if err != nil {
			var awsErr awserr.Error
			if errors.As(err, &awsErr) {
				switch awsErr.Code() {
				case s3.ErrCodeNoSuchKey:
					continue
				default:
					if strings.Contains(awsErr.Error(), "status code: 404") {
						continue
					}
					ol[index].SetIgnore(fmt.Sprintf("When checking for a duplicate object: an aws errored: %q", awsErr.Error()))
					continue
				}
			}
		}
		ol[index].SetIgnore(ErrIgnoreObjectAlreadyExists)
	}
}

/*
IgnoreIfLocalDoesNotExist is an ObjectList convenience method. It calls IgnoreIfLocalDoesNotExist on each FileObject
in the ObjectList slice.

See FileObject.IgnoreIfLocalDoesNotExist for more information.
*/
func (ol ObjectList) IgnoreIfLocalDoesNotExist() error {
	if len(ol) == 0 {
		return errors.New("FileList is empty")
	}

	if err := ol.IterateAndExecute(func(fo *FileObject) (err error) {
		return fo.IgnoreIfLocalDoesNotExist()
	}); err != nil {
		return err
	}
	return nil
}

/*
SetAsDirectoryPart is a ObjectList method. It calls the FileObject.SetAsDirectoryPart function on each FileObject in the
ObjectList slice.

See FileObject.SetAsDirectoryPart for more information.
*/
func (ol ObjectList) SetAsDirectoryPart() {
	_ = ol.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetAsDirectoryPart()
		return
	})
}

/*
SetChecksum is a ObjectList method. It calls the FileObject.SetChecksum function on each FileObject in the
ObjectList slice.

See FileObject.SetChecksum for more information.
*/
func (ol ObjectList) SetChecksum() (err error) {
	_ = ol.ConcurrentIterateAndExecute(25, func(fo *FileObject) (err error) {
		fo.c.Logger.Debug(fmt.Sprintf("Setting checksum for %q", fo.BaseName))
		return fo.SetChecksum()
	})
	return
}

/*
SetGroups is a ObjectList method. It splits FileObjects into groups based on the configuration value for
ProfileOptionsMaxConcurrent, then assigns each FileObject a group with FileObject.SetGroup.

See FileObject.SetGroup for more information.
*/
func (ol ObjectList) SetGroups() {
	for index, fo := range ol {
		fo.SetGroup(index % fo.c.Options[config.ProfileOptionsMaxConcurrent].(int))
	}
}

/*
SetIgnoreIfLocalNotExists is a ObjectList convenience method. It calls FileObject.SetIgnoreIfLocalNotExists on each
FileObject in the ObjectList.

See FileObject.SetIgnoreIfLocalNotExists for more information.
*/
func (ol ObjectList) SetIgnoreIfLocalNotExists() {
	_ = ol.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetIgnoreIfLocalNotExists()
		return
	})
	return
}

/*
SetIgnoreIfObjExists is a ObjectList convenience method. It calls FileObject.SetIgnoreIfObjExists on each FileObject
in the ObjectList.

See FileObject.SetIgnoreIfObjExists for more information.
*/
func (ol ObjectList) SetIgnoreIfObjExists() {
	_ = ol.ConcurrentIterateAndExecute(10, func(fo *FileObject) (err error) {
		fo.c.Logger.Debug(fmt.Sprintf("Checking if object exists %q", fo.PrefixedName))
		fo.SetIgnoreIfObjExists()
		return
	})
	return
}

/*
SetFileSizes is a ObjectList convenience method. It calls FileObject.SetFileSize on each FileObject in the
ObjectList slice.

See FileObject.SetFileSize for more information.
*/
func (ol ObjectList) SetFileSizes() {
	_ = ol.ConcurrentIterateAndExecute(25, func(fo *FileObject) (err error) {
		if !fo.Ignore {
			fo.SetFileSize()
		}
		return
	})
}

/*
SetPrefixedNames is an ObjectList method. It calls the FileObject.SetPrefixedName function on each FileObject in the
ObjectList slice.
*/
func (ol ObjectList) SetPrefixedNames() {
	_ = ol.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetPrefixedName()
		return
	})
}

/*
SetRelativeRoot is a ObjectList method. It calls the FileObject.SetRelativeRoot function on each FileObject in the
ObjectList slice.

See FileObject.SetRelativeRoot for more information.
*/
func (ol ObjectList) SetRelativeRoot(dir string) {
	_ = ol.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetRelativeRoot(dir)
		return
	})
}

/*
TagAll is a ObjectList method. It calls the FileObject.Tag function on each FileObject in the ObjectList slice.
It tags the FileObject with the key/value pair provided in the arguments.
*/
func (ol ObjectList) TagAll(k, v string) {
	_ = ol.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.Tag(k, v)
		return
	})
}

/*
TagOrigins is an ObjectList method. It calls the FileObject.Tag function on each FileObject in the ObjectList slice.
It tags the FileObject with the key "Origin" and the value of the FileObject's AbsolutePath.

See FileObject.Tag for more information.
*/
func (ol ObjectList) TagOrigins() {
	_ = ol.IterateAndExecute(func(fo *FileObject) (err error) {
		if fo.c.Options["tagOrigins"].(bool) {
			fo.Tag("Origin", fo.AbsolutePath)
		}
		return
	})
}

/*
Upload is an ObjectList method. It prepares each FileObject in the ObjectList to be uploaded by executing SetPrefixedNames,
FixRedundantKeys, SetIgnoreIfObjExists, and SetGroups. It then calls UploadHandler to handle concurrent uploads.

See ObjectList.UploadHandler for more information.
*/
func (ol ObjectList) Upload(c *config.Configuration) (err error, bytes int64, uploaded, ignored int) {
	if len(ol) == 0 {
		return nil, 0, 0, 0
	}

	if err != nil {
		return
	}

	if !ol[0].IsDirectoryPart {
		ol.SetPrefixedNames()
		err = ol.FixRedundantKeys()
		if err != nil {
			return
		}
	}

	ol.SetIgnoreIfObjExists()
	ol.SetGroups()

	errs := ol.UploadHandler(c)
	if len(errs) > 0 {
		for _, err := range errs {
			c.Logger.Error(fmt.Sprintf("Error during upload: %q", err.Error()))
		}
	}
	return nil, ol.GetTotalUploadedBytes(), ol.CountUploaded(), ol.CountIgnored()
}

/*
UploadHandler is an ObjectList method. It takes a Configuration pointer and returns a slice of errors, the number of
uploaded files, and the number of ignored files.

It uses sync.WaitGroup and a channel to handle concurrent uploads. A single s3manager.Uploader is created and passed to
each goroutine. Each goroutine is assigned a group number, which is used to determine which FileObject it will upload.
The number of goroutines is determined by the ProfileOptionsMaxConcurrent configuration value.
*/
func (ol ObjectList) UploadHandler(c *config.Configuration) (err []error) {
	var wg sync.WaitGroup
	resultChan := make(chan UploadResult)

	svc, buErr := BuildUploader(c)
	if buErr != nil {
		err = append(err, buErr)
		return err
	}

	for i := 0; i < c.Options[config.ProfileOptionsMaxConcurrent].(int); i++ {
		wg.Add(1)
		go func(group int, svc *s3manager.Uploader, wg *sync.WaitGroup) {
			defer wg.Done()
			fi := NewFileIterator(ol[0].c, ol, group)
			err := svc.UploadWithIterator(aws.BackgroundContext(), fi)

			resultChan <- UploadResult{
				Err:         err,
				UploadCount: ol.CountUploadedByGroup(group),
				IgnoreCount: ol.CountIgnoredByGroup(group),
			}
		}(i, svc, &wg)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(resultChan)
	}(&wg)

	for result := range resultChan {
		if result.Err != nil {
			err = append(err, result.Err)
		}
	}
	return
}

/*
Count is an ObjectList method. It returns the number of FileObjects in the ObjectList slice.
*/
func (ol ObjectList) Count() (count int) {
	return len(ol)
}

/*
CountUploaded is an ObjectList method. It returns the number of FileObjects in the ObjectList slice that have the
FileObject.IsUploaded field set to true.
*/
func (ol ObjectList) CountUploaded() (count int) {
	for index := range ol {
		if ol[index].IsUploaded {
			count++
		}
	}
	return
}

/*
CountUploadedByGroup is an ObjectList method. It returns the number of FileObjects in the ObjectList slice that have the
FileObject.IsUploaded field set to true and the FileObject.Group field set to the provided group number.
*/
func (ol ObjectList) CountUploadedByGroup(group int) (count int) {
	for index := range ol {
		if ol[index].IsUploaded && ol[index].Group == group {
			count++
		}
	}
	return
}

/*
CountIgnored is an ObjectList method. It returns the number of FileObjects in the ObjectList slice that have the
FileObject.Ignore field set to true.
*/
func (ol ObjectList) CountIgnored() (count int) {
	for index := range ol {
		if ol[index].Ignore {
			count++
		}
	}
	return
}

/*
CountIgnoredByGroup is an ObjectList method. It returns the number of FileObjects in the ObjectList slice that have the
FileObject.Ignore field set to true and the FileObject.Group field set to the provided group number.
*/
func (ol ObjectList) CountIgnoredByGroup(group int) (count int) {
	for index := range ol {
		if ol[index].Ignore && ol[index].Group == group {
			count++
		}
	}
	return
}

/* DEBUG */

/*
DebugOutput is an ObjectList method. It prints the AbsolutePath and IsDirectoryPart fields of each FileObject in the
ObjectList slice.

Change this as needed.
*/
func (ol ObjectList) DebugOutput() {
	_ = ol.IterateAndExecute(func(fo *FileObject) (err error) {
		fmt.Println()
		fmt.Printf("AbsolutePath: %q\n", fo.AbsolutePath)
		fmt.Printf("PrefixedName: %q\n", fo.PrefixedName)
		fmt.Printf("Group: %d\n", fo.Group)
		fmt.Printf("IsUploaded: %v\n", fo.IsUploaded)
		fmt.Printf("Ignored: %v\n", fo.Ignore)
		fmt.Printf("IgnoreString: %q\n", fo.IgnoreString)
		fmt.Println()
		return
	})
}
