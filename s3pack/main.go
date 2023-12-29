package s3pack

import (
	"errors"
	"fmt"

	"github.com/orme292/s3packer/conf"
)

func IndividualFileUploader(a *conf.AppConfig, files []string) (err error, bytes int64, uploaded, ignored int) {
	exists, err := BucketExists(a)
	if err != nil {
		return
	} else if !exists {
		return errors.New(fmt.Sprintf("aws says %q does not exist", a.Bucket.Name)), 0, 0, 0
	}

	objList, err := NewObjectList(a, files)
	if err != nil {
		return
	}

	return objList.Upload(a)
}

func DirectoryUploader(a *conf.AppConfig, dirs []string) (err error, bytes int64, uploaded, ignored int) {
	exists, err := BucketExists(a)
	if err != nil {
		return
	} else if !exists {
		return errors.New(fmt.Sprintf("aws says %q does not exist", a.Bucket.Name)), 0, 0, 0
	}

	dirLists, err := NewRootList(a, dirs)
	if err != nil {
		return
	}

	return dirLists.Upload()
}
