package s3pack

import (
	"path/filepath"

	"github.com/orme292/s3packer/conf"
)

/*
DirectoryObject is a struct that contains information about each directory to be processed.

See RootList for more information and structure.
*/
type DirectoryObject struct {
	// StartPath is the absolute path to the directory
	StartPath string

	// ol is the ObjectList for the directory, naturally a pointer
	ol ObjectList

	// c is the application configuration
	a *conf.AppConfig
}

/*
NewDirectoryObject is an DirectoryObject constructor. It takes a path as a string and returns a DirectoryObject and an
error.

It scans the provided directory for all files and creates an ObjectList (ol) from the files.

See NewObjectList for additional information
*/
func NewDirectoryObject(a *conf.AppConfig, path string) (do *DirectoryObject, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	files, err := GetFiles(absPath)
	if err != nil {
		return
	}

	list, err := NewObjectList(a, files)
	if err != nil {
		return nil, err
	}
	return &DirectoryObject{
		StartPath: absPath,
		ol:        list,
		a:         a,
	}, nil
}

/*
TagAll is an DirectoryObject method. It calls TagAll on the DirectoryObject ObjectList.
*/
func (do *DirectoryObject) TagAll(k, v string) {
	do.ol.TagAll(k, v)
}

/*
Upload is an DirectoryObject method. It calls Upload on the DirectoryObject's ObjectList.
*/
func (do *DirectoryObject) Upload() (err error) {
	err, _, _, _ = do.ol.Upload(do.a)
	return
}

/*
CountFileObjects is an DirectoryObject method. It returns the number of FileObjects in the DirectoryObject's
ObjectList slice.
*/
func (do *DirectoryObject) CountFileObjects() (count int) {
	return len(do.ol)
}

/*
CountIgnored is an DirectoryObject method. It returns the number of FileObjects in the DirectoryObject's ObjectList
slice that are ignored.

See ObjectList.CountIgnored for more information
*/
func (do *DirectoryObject) CountIgnored() (count int) {
	return do.ol.CountIgnored()
}

/*
CountUploaded is an DirectoryObject method. It returns the number of FileObjects in the DirectoryObject's ObjectList slice
that are uploaded.

See ObjectList.CountUploaded for more information
*/
func (do *DirectoryObject) CountUploaded() (count int) {
	return do.ol.CountUploaded()
}

/*
GetTotalUploadedBytes is an DirectoryObject method. It returns the total number of bytes uploaded to S3 for all
FileObjects in the ObjectList slice.

See ObjectList.GetTotalUploadedBytes for more information
*/
func (do *DirectoryObject) GetTotalUploadedBytes() (total int64) {
	return do.ol.GetTotalUploadedBytes()
}
