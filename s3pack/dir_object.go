package s3pack

import (
	"path/filepath"

	"github.com/orme292/s3packer/config"
)

type DirectoryObject struct {
	StartPath string
	objList   ObjectList

	c *config.Configuration
}

func NewDirectoryObject(c *config.Configuration, path string) (do *DirectoryObject, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	files, err := GetFiles(absPath)
	if err != nil {
		return
	}

	list, err := NewObjectList(c, files)
	if err != nil {
		return nil, err
	}
	return &DirectoryObject{
		StartPath: absPath,
		objList:   list,
		c:         c,
	}, nil
}

/*
TagAll is an DirectoryObject method. It calls TagAll on the DirectoryObject's ObjectList.
*/
func (do *DirectoryObject) TagAll(k, v string) {
	do.objList.TagAll(k, v)
}

/*
Upload is an DirectoryObject method. It calls Upload on the DirectoryObject's ObjectList.
*/
func (do *DirectoryObject) Upload() (err error, uploaded, ignored int) {
	return do.objList.Upload(do.c)
}

/*
CountObjects is an DirectoryObject method. It returns the number of FileObjects in the DirectoryObject slice.
*/
func (do *DirectoryObject) CountObjects() (count int) {
	return len(do.objList)
}
