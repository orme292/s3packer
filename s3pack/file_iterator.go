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
	objList ObjectList
	stage   struct {
		objListIndex int
		fo           *FileObject
		f            *os.File
	}
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
input and returns a s3manager.BatchUploadObject struct.
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
			fi.stage.fo.IsUploaded = true
			fi.stage.objListIndex += 1
			return f.Close()
		},
	}
}
