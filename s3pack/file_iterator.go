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
		fo *FileObject
		f  *os.File
	}
	err error
	c   *config.Configuration
}

/*
NewFileIterator creates a new FileIterator struct, which implements the s3manager.BatchUploadIterator
interface.
*/
func NewFileIterator(c *config.Configuration, objList ObjectList) s3manager.BatchUploadIterator {
	return &FileIterator{
		objList: objList,
		c:       c,
	}
}

/*
Err returns the error that caused the iterator to stop.
*/
func (fi *FileIterator) Err() error {
	fi.stage.fo.IsUploaded = false
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
		if fi.objList[0].IsUploaded || fi.objList[0].Ignore {
			fi.c.Logger.Warn(fmt.Sprintf("Ignoring %q, %s", fi.objList[0].PrefixedName, fi.objList[0].IgnoreString))
			if len(fi.objList) == 1 {
				return false
			}
			fi.objList = fi.objList[1:]
			continue
		} else {
			break
		}
	}

	f, err := os.Open(fi.objList[0].AbsolutePath)
	fi.err = err
	fi.stage.f = f
	fi.stage.fo = fi.objList[0]
	fi.objList = fi.objList[1:]

	return fi.Err() == nil
}

/*
UploadObject returns a s3manager.BatchUploadObject struct for use with s3manager.UploadWithIterator. It takes no
input and returns a s3manager.BatchUploadObject struct.
*/
func (fi *FileIterator) UploadObject() s3manager.BatchUploadObject {
	f := fi.stage.f
	fi.c.Logger.Info(fmt.Sprintf("Uploading %s...", fi.stage.fo.PrefixedName))
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
			return f.Close()
		},
	}
}
