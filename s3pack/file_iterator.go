package s3pack

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/orme292/s3packer/config"
)

/*
FileIterator is used with the s3manager.BatchUploadIterator to process a list of files for upload to s3.
*/
type FileIterator struct {
	// objList is the ObjectList to be processed
	objList ObjectList

	/*
		stage is the struct that holds details on the current file being processed. objListIndex is the index of objList
		that is being used for stage.fo and stage.f. stage.fo is a pointer to the FileObject that is being processed.
		FileIterator.objList[FileIterator.stage.objListIndex] and FileIterator.stage.fo reference the same FileObject.
		stage.f is the contents of the file being processed, it is a pointer to an os.File.
	*/
	stage struct {
		objListIndex int
		fo           *FileObject
		f            *os.File
	}

	/*
		group is a number assigned and used by ObjectList.UploadHandler, a func that will create a goroutine to process
		the files in each group. The number of groups is determined by the configuration option "maxConcurrentUploads".
	*/
	group int
	err   error
	c     *config.Configuration
}

/*
NewFileIterator creates a new FileIterator struct, which implements the s3manager.BatchUploadIterator
interface.
*/
func NewFileIterator(c *config.Configuration, objList ObjectList, group int) s3manager.BatchUploadIterator {
	return &FileIterator{
		objList: objList,
		group:   group,
		c:       c,
	}
}

/*
Err returns the error that caused the iterator to stop.
*/
func (fi *FileIterator) Err() error {
	if fi.err != nil {
		fi.stage.fo.IsUploaded = false
		fmt.Println("Error:", fi.err)
		fmt.Println("Group:", fi.group)
		fmt.Println("Index:", fi.stage.objListIndex)
	}
	return fi.err
}

/*
Next returns true if there are more files to upload, or false if there are no more files to upload.

fi.stage.objListIndex is incremented by 1 until a FileObject is found where IsUploaded is false, Ignore is false, and
Group is equal to fi.group. Once a FileObject is found that meets those criteria, fi.stage.fo is set to the FileObject
the associated File is opened and the contents are stored in fi.stage.f and true is returned.
*/
func (fi *FileIterator) Next() bool {
	if len(fi.objList) == 0 {
		return false
	}

	for {
		if fi.stage.objListIndex >= len(fi.objList) {
			return false
		}
		if fi.objList[fi.stage.objListIndex].Group != fi.group {
			fi.stage.objListIndex += 1
			continue
		}
		if fi.objList[fi.stage.objListIndex].IsUploaded || fi.objList[fi.stage.objListIndex].Ignore {
			fi.c.Logger.Warn(fmt.Sprintf("Ignoring: GROUP(%d) INDEX(%d) %q, %s", fi.group, fi.stage.objListIndex, fi.objList[fi.stage.objListIndex].PrefixedName, fi.objList[fi.stage.objListIndex].IgnoreString))
			fi.stage.objListIndex += 1
			continue
		}
		break
	}

	f, err := os.Open(fi.objList[fi.stage.objListIndex].AbsolutePath)
	fi.err = err
	fi.stage.f = f
	fi.stage.fo = fi.objList[fi.stage.objListIndex]

	return fi.Err() == nil
}

/*
UploadObject returns a s3manager.BatchUploadObject struct for use with s3manager.UploadWithIterator. It takes no
input and returns a s3manager.BatchUploadObject struct. The details in the fi.stage struct are used to build the
UploadInput that will be passed to s3manager, which handles the actual upload.
*/
func (fi *FileIterator) UploadObject() s3manager.BatchUploadObject {
	f := fi.stage.f
	fi.c.Logger.Info(fmt.Sprintf("Uploading (%s) %q...", FileSizeString(fi.stage.fo.FileSize), fi.stage.fo.PrefixedName))
	return s3manager.BatchUploadObject{
		Object: &s3manager.UploadInput{
			Bucket:       aws.String(fi.c.Bucket[config.ProfileBucketName].(string)),
			Key:          aws.String(fi.stage.fo.PrefixedName),
			ACL:          aws.String(fi.c.Options[config.ProfileOptionACL].(string)),
			StorageClass: aws.String(fi.c.Options[config.ProfileOptionStorage].(string)),
			Tagging:      aws.String(fi.stage.fo.Tags),
			Body:         f,
		},
		After: func() error {
			// TODO: Add an explicit check to confirm the file is in the bucket.
			fi.stage.fo.IsUploaded = true
			fi.stage.objListIndex += 1
			return f.Close()
		},
	}
}
