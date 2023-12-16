package s3pack

import (
	"errors"
	"fmt"

	"github.com/orme292/s3packer/config"
)

func IndividualFileUploader(c *config.Configuration, files []string) (err error, bytes int64, uploaded, ignored int) {
	exists, err := BucketExists(c)
	if err != nil {
		return
	} else if !exists {
		return errors.New(fmt.Sprintf("AWS says %q does not exist", c.Bucket[config.ProfileBucketName].(string))), 0, 0, 0
	}

	objList, err := NewObjectList(c, files)
	if err != nil {
		return
	}

	return objList.Upload(c)
}

func DirectoryUploader(c *config.Configuration, dirs []string) (err error, bytes int64, uploaded, ignored int) {
	exists, err := BucketExists(c)
	if err != nil {
		return
	} else if !exists {
		return errors.New(fmt.Sprintf("AWS says %q does not exist", c.Bucket[config.ProfileBucketName].(string))), 0, 0, 0
	}

	dirLists, err := NewRootList(c, dirs)
	if err != nil {
		return
	}

	return dirLists.Upload()
}
