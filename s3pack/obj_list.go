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
ObjectList is a slice of FileObject pointers. Most FileLists methods are just for convenience -- they iterate over the
elements of the slice and call the corresponding FileList method. Exceptions are noted below.

See FileList for more information
*/
type ObjectList []*FileObject

type UploadResult struct {
	Err         error
	UploadCount int
	IgnoreCount int
}

/*
NewObjectList is an ObjectList constructor. It takes a slice of paths and returns a slice FileObjects.
It calls NewFileObject on each path and appends the result to the slice of ObjectList. It then calls
fixRedundantKeys, disregardIfLocalDoesNotExist, disregardIfExistsInBucket, and setFileSizes to sanitize the entire list.

See NewFileObject for additional information
*/
func NewObjectList(c *config.Configuration, paths []string) (objList ObjectList, err error) {
	for _, path := range paths {
		fo, err := NewFileObject(c, path)
		if err != nil {
			return nil, err
		}
		objList = append(objList, fo)
	}

	objList.SetIgnoreIfLocalNotExists()
	objList.SetFileSizes()
	_ = objList.SetChecksum()
	objList.TagOrigins()
	for k, v := range c.Tags {
		objList.TagAll(k, v)
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
func (objList ObjectList) FixRedundantKeys() error {
	if len(objList) == 0 || len(objList) == 1 {
		return errors.New("FileList is empty or only contains one item")
	}

	occurrences := make(map[string]int)
	for _, obj := range objList {
		if _, ok := occurrences[obj.PrefixedName]; ok {
			occurrences[obj.PrefixedName] += 1
		} else {
			occurrences[obj.PrefixedName] = 1
		}
	}

	for prefixedName, numOccurs := range occurrences {
		if numOccurs > 1 {
			counter := 0
			for index := range objList {
				if objList[index].PrefixedName == prefixedName {
					objList[index].PrefixedName = fmt.Sprintf("%s-%d", objList[index].PrefixedName, counter)
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
func (objList ObjectList) IterateAndExecute(fn IteratedObjectFunc) (err error) {
	for index := range objList {
		if err = fn(objList[index]); err != nil {
			return
		}
	}
	return
}

/*
IteratedObjectFunc is a function type that takes a FileObject pointer and returns an error. It is used with
ObjectList.IterateAndExecute
*/
type IteratedObjectFunc func(fo *FileObject) (err error)

/*
IgnoreIfObjectExistsInBucket is an ObjectList method. It iterates through each FileObject in the ObjectList and tries
to retrieve metadata from an S3 object of the same name (s3 key = FileObject.PrefixedName). If the object exists, then
the FileObject.Ignore field is set to true and the FileObject.IgnoreString field is set to ErrIgnoreObjectAlreadyExists.
*/
func (objList ObjectList) IgnoreIfObjectExistsInBucket() {
	if objList[0].c.Options[config.ProfileOptionOverwrite].(bool) || len(objList) == 0 {
		return
	}

	sess, _ := NewSession(&objList[0].c)

	svc := s3.New(sess, &aws.Config{})

	for index := range objList {
		_, err := svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(objList[index].c.Bucket[config.ProfileBucketName].(string)),
			Key:    aws.String(objList[index].PrefixedName),
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
					objList[index].SetIgnore(fmt.Sprintf("When checking for a duplicate object: an aws errored: %q", awsErr.Error()))
					continue
				}
			}
		}
		objList[index].SetIgnore(ErrIgnoreObjectAlreadyExists)
	}
}

/*
IgnoreIfLocalDoesNotExist is an ObjectList convenience method. It calls IgnoreIfLocalDoesNotExist on each FileObject
in the ObjectList slice.

See FileObject.IgnoreIfLocalDoesNotExist for more information.
*/
func (objList ObjectList) IgnoreIfLocalDoesNotExist() error {
	if len(objList) == 0 {
		return errors.New("FileList is empty")
	}

	if err := objList.IterateAndExecute(func(fo *FileObject) (err error) {
		return fo.IgnoreIfLocalDoesNotExist()
	}); err != nil {
		return err
	}
	return nil
}

/*
SetAsDirectoryPart is a ObjectList method. It calls the FileObject.SetDirectoryPart function on each FileObject in the
ObjectList slice.

See FileObject.SetAsDirectoryPart for more information.
*/
func (objList ObjectList) SetAsDirectoryPart() {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetDirectoryPart()
		return
	})
}

func (objList ObjectList) SetChecksum() (err error) {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
		_ = fo.SetChecksum()
		return
	})
	return
}

func (objList ObjectList) SetGroups() {
	for index, fo := range objList {
		fo.SetGroup(index % fo.c.Options[config.ProfileOptionsMaxConcurrent].(int))
	}
}

func (objList ObjectList) SetIgnoreIfLocalNotExists() {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetIgnoreIfLocalNotExists()
		return
	})
	return
}

func (objList ObjectList) SetIgnoreIfObjExists() {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
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
func (objList ObjectList) SetFileSizes() {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
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
func (objList ObjectList) SetPrefixedNames() {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetPrefixedName()
		return
	})
}

/*
SetRelativeRoot is a ObjectList method. It calls the FileObject.SetRelativeRoot function on each FileObject in the
ObjectList slice.

See FileObject.SetRelativeRoot for more information.
*/
func (objList ObjectList) SetRelativeRoot(dir string) {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.SetRelativeRoot(dir)
		return
	})
}

func (objList ObjectList) ReturnTotalUploadedBytes() (total int64) {
	for index := range objList {
		if objList[index].IsUploaded {
			total += objList[index].FileSize
		}
	}
	return
}

/*
TagAll is a ObjectList method. It calls the FileObject.Tag function on each FileObject in the ObjectList slice.
It tags the FileObject with the key/value pair provided in the arguments.
*/
func (objList ObjectList) TagAll(k, v string) {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
		fo.Tag(k, v)
		return
	})
}

/*
TagOrigins is an ObjectList method. It calls the FileObject.Tag function on each FileObject in the ObjectList slice.
It tags the FileObject with the key "Origin" and the value of the FileObject's AbsolutePath.

See FileObject.Tag for more information.
*/
func (objList ObjectList) TagOrigins() {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
		if fo.c.Options["tagOrigins"].(bool) {
			fo.Tag("Origin", fo.AbsolutePath)
		}
		return
	})
}

/*
Upload is an ObjectList method. It creates a new s3manager.Uploader with BuildUploader, then creates a FileIterator
and passes it to the s3manager.Uploader.UploadWithIterator function. It returns an error, the number of files uploaded,
and the number of files ignored.
*/
func (objList ObjectList) Upload(c *config.Configuration) (err error, uploaded, ignored int) {
	if len(objList) == 0 {
		return nil, 0, 0
	}

	if err != nil {
		return
	}

	if !objList[0].IsDirectoryPart {
		objList.SetPrefixedNames()
		err = objList.FixRedundantKeys()
		if err != nil {
			return
		}
	}

	objList.SetIgnoreIfObjExists()
	objList.SetGroups()

	errs, _, _ := objList.UploadHandler(c)
	if len(errs) > 0 {
		for _, err := range errs {
			c.Logger.Error(fmt.Sprintf("Error during upload: %q", err.Error()))
		}
	}
	return nil, objList.CountUploaded(), objList.CountIgnored()
}

func (objList ObjectList) UploadHandler(c *config.Configuration) (err []error, uploaded, ignored int) {
	var wg sync.WaitGroup
	resultCh := make(chan UploadResult)

	svc, buErr := BuildUploader(c)
	if buErr != nil {
		err = append(err, buErr)
		return err, 0, 0
	}

	for i := 0; i < c.Options[config.ProfileOptionsMaxConcurrent].(int); i++ {
		wg.Add(1)
		go func(group int, svc *s3manager.Uploader, wg *sync.WaitGroup) {
			defer wg.Done()
			fi := NewFileIterator(&objList[0].c, objList, group)
			err := svc.UploadWithIterator(aws.BackgroundContext(), fi)

			resultCh <- UploadResult{
				Err:         err,
				UploadCount: objList.CountUploadedByGroup(group),
				IgnoreCount: objList.CountIgnoredByGroup(group),
			}
		}(i, svc, &wg)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(resultCh)
	}(&wg)

	for result := range resultCh {
		if result.Err != nil {
			err = append(err, result.Err)
		}
		uploaded += result.UploadCount
		ignored += result.IgnoreCount
	}
	return
}

/*
Count is an ObjectList method. It returns the number of FileObjects in the ObjectList slice.
*/
func (objList ObjectList) Count() (count int) {
	return len(objList)
}

/*
CountUploaded is an ObjectList method. It returns the number of FileObjects in the ObjectList slice that have the
FileObject.IsUploaded field set to true.
*/
func (objList ObjectList) CountUploaded() (count int) {
	for index := range objList {
		if objList[index].IsUploaded {
			count++
		}
	}
	return
}

func (objList ObjectList) CountUploadedByGroup(group int) (count int) {
	for index := range objList {
		if objList[index].IsUploaded && objList[index].Group == group {
			count++
		}
	}
	return
}

/*
CountIgnored is an ObjectList method. It returns the number of FileObjects in the ObjectList slice that have the
FileObject.Ignore field set to true.
*/
func (objList ObjectList) CountIgnored() (count int) {
	for index := range objList {
		if objList[index].Ignore {
			count++
		}
	}
	return
}

func (objList ObjectList) CountIgnoredByGroup(group int) (count int) {
	for index := range objList {
		if objList[index].Ignore && objList[index].Group == group {
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
func (objList ObjectList) DebugOutput() {
	_ = objList.IterateAndExecute(func(fo *FileObject) (err error) {
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
