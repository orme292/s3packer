package s3pack

import (
	"fmt"

	"github.com/orme292/s3packer/config"
)

/*
DirectoryList is a slice of DirectoryObject pointers (slices are inherently pointers).

See RootList for more information and structure.
*/
type DirectoryList []*DirectoryObject

/*
NewDirectoryList is a DirectoryList constructor. It takes a configuration and a directory as a string. It returns a
DirectoryList and an error.
*/
func NewDirectoryList(c *config.Configuration, dir string) (dl DirectoryList, err error) {
	c.Logger.Info(fmt.Sprintf("Processing directory: %q", dir))

	subDirs, err := GetSubDirs(dir)
	if err != nil {
		return nil, err
	}

	for index := range subDirs {
		do, err := NewDirectoryObject(c, subDirs[index])
		if err != nil {
			c.Logger.Error(err.Error())
		}
		dl = append(dl, do)
	}

	dl.SetAsDirectoryPart()
	dl.SetRelativeRoot(dir)

	return dl, nil
}

/*
IterateAndExecute is an DirectoryList method. It takes a function that takes a DirectoryObject pointer and returns
an error. It iterates over the DirectoryList slice and calls the provided function on each DirectoryObject pointer. If the function returns an
error, then it is returned and iteration stops
*/
func (dl DirectoryList) IterateAndExecute(fn DirObjectIterationFunc) (err error) {
	for index := range dl {
		if err = fn(dl[index]); err != nil {
			return
		}
	}
	return
}

/*
DirObjectIterationFunc is a function type that takes a DirectoryObject pointer and returns an error. It is used with
DirectoryList.IterateAndExecute
*/
type DirObjectIterationFunc func(do *DirectoryObject) (err error)

/*
SetAsDirectoryPart is a DirectoryList method. It iterates over the DirectoryList slice and calls the
SetAsDirectoryPart method on each ObjectList in a DirectoryObject pointer.

See ObjectList.SetAsDirectoryPart for more information
*/
func (dl DirectoryList) SetAsDirectoryPart() {
	_ = dl.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.ol.SetAsDirectoryPart()
		return
	})
}

/*
SetRelativeRoot is a DirectoryList method. It iterates over the DirectoryList slice and calls the
SetRelativeRoot method on each ObjectList in a DirectoryObject pointer.

See ObjectList.SetRelativeRoot for more information
*/
func (dl DirectoryList) SetRelativeRoot(dir string) {
	_ = dl.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.ol.SetRelativeRoot(dir)
		return
	})
}

/*
SetPrefixedNames is a DirectoryList method. It iterates over the DirectoryList slice and calls the
SetPrefixedNames method on each ObjectList in a DirectoryObject pointer.

See ObjectList.setPrefixedNames for more information
*/
func (dl DirectoryList) SetPrefixedNames() {
	_ = dl.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.ol.SetPrefixedNames()
		return
	})
}

/*
TagAll is a DirectoryList method. It iterates over the DirectoryList slice and calls the DirectoryObject.TagAll method
on each DirectoryObject pointer.

See DirectoryObject.TagAll for more information
*/
func (dl DirectoryList) TagAll(k, v string) {
	_ = dl.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.TagAll(k, v)
		return
	})
}

/*
Upload is an DirectoryList method. It iterates over the DirectoryList slice and calls the DirectoryObject.Upload method
on each DirectoryObject pointer. It also executes the SetPrefixedNames method on the DirectoryList slice.

See DirectoryObject.Upload for more information
See DirectoryList.SetPrefixedNames for more information
*/
func (dl DirectoryList) Upload() (err error) {
	dl.SetPrefixedNames()
	for index := range dl {
		dl[index].c.Logger.Debug(fmt.Sprintf("Directory %s has %d objects", dl[index].StartPath, dl[index].CountFileObjects()))
		fErr := dl[index].Upload()
		if fErr != nil {
			dl[index].c.Logger.Error(fErr.Error())
			return
		}
	}
	return
}

/*
Count is an DirectoryList method. It returns the number of DirectoryObjects in the DirectoryList slice.
*/
func (dl DirectoryList) Count() (count int) {
	return len(dl)
}

/*
CountFileObjects is an DirectoryList method. It returns the number of FileObjects in all ObjectLists in
the DirectoryList slice.
*/
func (dl DirectoryList) CountFileObjects() (count int) {
	for _, do := range dl {
		count += do.CountFileObjects()
	}
	return
}

/*
CountIgnored is an DirectoryList method. It returns the number of FileObjects in the slice with Ignore set to true. It
iterates over the DirectoryList slice and calls the DirectoryObject.CountIgnored method on each DirectoryObject pointer.

See DirectoryObject.CountIgnored for more information
*/
func (dl DirectoryList) CountIgnored() (count int) {
	for _, do := range dl {
		count += do.CountIgnored()
	}
	return
}

/*
CountUploaded is an DirectoryList method. It returns the number of FileObjects in the slice with IsUploaded set to true.
It iterates over the DirectoryList slice and calls the DirectoryObject.CountUploaded method on each DirectoryObject
pointer.

See DirectoryObject.CountUploaded for more information
*/
func (dl DirectoryList) CountUploaded() (count int) {
	for _, do := range dl {
		count += do.CountUploaded()
	}
	return
}

/*
GetTotalUploadedBytes is an DirectoryList method. It returns the total number of bytes uploaded. It iterates over the
DirectoryList slice and calls the DirectoryObject.GetTotalUploadedBytes method on each DirectoryObject pointer.

See DirectoryObject.GetTotalUploadedBytes for more information
*/
func (dl DirectoryList) GetTotalUploadedBytes() (total int64) {
	for _, do := range dl {
		total += do.GetTotalUploadedBytes()
	}
	return
}

/*
DEBUG
*/

func (dl DirectoryList) DebugOutput() {
	for _, do := range dl {
		fmt.Println(do.StartPath)
	}
}
