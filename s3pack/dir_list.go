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
NewDirectoryList is a DirectoryList constructor. It takes a config and a directory as a string. It returns a
DirectoryList and an error.
*/
func NewDirectoryList(c *config.Configuration, dir string) (dirList DirectoryList, err error) {
	c.Logger.Info(fmt.Sprintf("Processing directory: %q", dir))

	subDirs, err := GetSubDirs(dir)
	if err != nil {
		return nil, err
	}

	for _, subDir := range subDirs {
		do, err := NewDirectoryObject(c, subDir)
		if err != nil {
			c.Logger.Error(err.Error())
		}
		dirList = append(dirList, do)
	}

	dirList.SetAsDirectoryPart()
	dirList.SetRelativeRoot(dir)

	return dirList, nil
}

/*
IterateAndExecute is an DirectoryList method. It takes a function that takes a DirectoryObject pointer and returns
an error. It iterates over the DirectoryList slice and calls the provided function on each DirectoryObject pointer. If the function returns an
error, then it is returned and iteration stops
*/
func (dirList DirectoryList) IterateAndExecute(fn IteratedDirObjectFunc) (err error) {
	for index := range dirList {
		if err = fn(dirList[index]); err != nil {
			return
		}
	}
	return
}

/*
IteratedDirObjectFunc is a function type that takes a DirectoryObject pointer and returns an error. It is used with
DirectoryList.IterateAndExecute
*/
type IteratedDirObjectFunc func(do *DirectoryObject) (err error)

/*
SetAsDirectoryPart is a DirectoryList method. It iterates over the DirectoryList slice and calls the
ObjectList.setAsDirectoryPart method on each DirectoryObject pointer.

See ObjectList.setAsDirectoryPart for more information
*/
func (dirList DirectoryList) SetAsDirectoryPart() {
	_ = dirList.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.objList.SetAsDirectoryPart()
		return
	})
}

/*
SetRelativeRoot is a DirectoryList method. It iterates over the DirectoryList slice and calls the
ObjectList.SetRelativeRoot method on each DirectoryObject pointer.

See ObjectList.SetRelativeRoot for more information
*/
func (dirList DirectoryList) SetRelativeRoot(dir string) {
	_ = dirList.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.objList.SetRelativeRoot(dir)
		return
	})
}

/*
SetPrefixedNames is a DirectoryList method. It iterates over the DirectoryList slice and calls the
ObjectList.setPrefixedNames method on each DirectoryObject pointer.

See ObjectList.setPrefixedNames for more information
*/
func (dirList DirectoryList) SetPrefixedNames() {
	_ = dirList.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.objList.SetPrefixedNames()
		return
	})
}

/*
TagAll is a DirectoryList method. It iterates over the DirectoryList slice and calls the DirectoryObject.TagAll method
on each DirectoryObject pointer.

See DirectoryObject.TagAll for more information
*/
func (dirList DirectoryList) TagAll(k, v string) {
	_ = dirList.IterateAndExecute(func(do *DirectoryObject) (err error) {
		do.TagAll(k, v)
		return
	})
}

/*
Upload is an DirectoryList method. It iterates over the DirectoryList slice and calls the DirectoryObject.Upload method
on each DirectoryObject pointer. It also executes the FixRedundantKeys method on the DirectoryList slice.

See DirectoryObject.Upload for more information
See DirectoryList.FixRedundantKeys for more information
*/
func (dirList DirectoryList) Upload() (err error, uploaded, ignored int) {
	dirList.SetPrefixedNames()
	for index := range dirList {
		dirList[index].config.Logger.Debug(fmt.Sprintf("Directory %s has %d objects", dirList[index].StartPath, dirList[index].CountObjects()))
		fErr, fUploaded, fIgnored := dirList[index].Upload()
		if fErr != nil {
			dirList[index].config.Logger.Error(fErr.Error())
			return
		}
		uploaded += fUploaded
		ignored += fIgnored
	}
	return
}

/*
Count is an DirectoryList method. It returns the number of DirectoryObjects in the DirectoryList slice.
*/
func (dirList DirectoryList) Count() (count int) {
	return len(dirList)
}

/*
CountObjects is an DirectoryList method. It returns the number of FileObjects in all ObjectLists in the DirectoryList slice.
*/
func (dirList DirectoryList) CountObjects() (count int) {
	for _, do := range dirList {
		count += do.CountObjects()
	}
	return
}

/*
DEBUG
*/

func (dirList DirectoryList) DebugOutput() {
	for _, do := range dirList {
		fmt.Println(do.StartPath)
	}
}
