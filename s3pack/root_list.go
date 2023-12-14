package s3pack

import (
	"github.com/orme292/s3packer/config"
)

/*
RootList is a slice of DirectoryList pointers (slices are inherently pointers).

Structure:
RootList => DirectoryList => DirectoryObject -> ObjectList => FileObject

RootLists are a slice of DirectryList pointers. DirectoryList pointers are a slice of DirectoryObject pointers.
DirectoryObject pointers reference an ObjectList, which is a slice of FileObject pointers.
*/
type RootList []DirectoryList

/*
NewRootList is a RootList constructor. It takes directories as a slice of strings and returns a RootList.
*/
func NewRootList(c *config.Configuration, dir []string) (dirLists RootList, err error) {
	for _, d := range dir {
		dirList, err := NewDirectoryList(c, d)
		if err != nil {
			return nil, err
		}
		dirLists = append(dirLists, dirList)
	}
	return
}

/*
IterateAndExecute is an RootList method. It takes a function that takes a DirectoryObject pointer and returns
an error. It iterates over the DirectoryList slice and calls the provided function on each slice. If the function returns an
error, then it is returned and iteration stops
*/
func (rList RootList) IterateAndExecute(fn IteratedDirListFunc) (err error) {
	for index := range rList {
		if err = fn(rList[index]); err != nil {
			return
		}
	}
	return
}

/*
IteratedDirListFunc is a function type that takes a RootList slice and returns an error. It is used with
RootList.IterateAndExecute
*/
type IteratedDirListFunc func(dirList DirectoryList) (err error)

/*
SetPrefixedNames is a RootList method. It iterates over the RootList slice and calls the
DirectoryList.SetPrefixedNames method on each DirectoryList pointer.

See DirectoryList.SetPrefixedNames for more information
*/
func (rList RootList) SetPrefixedNames() {
	_ = rList.IterateAndExecute(func(dirList DirectoryList) (err error) {
		dirList.SetPrefixedNames()
		return
	})
}

/*
Upload is a RootList method. It iterates over the RootList slice and calls the DirectoryList.Upload method
on each DirectoryList pointer.

See DirectoryList.Upload for more information
*/
func (rList RootList) Upload() (err error, uploaded, ignored int) {
	_ = rList.IterateAndExecute(func(dirList DirectoryList) (err error) {
		err, u, i := dirList.Upload()
		if err != nil {
			rList[0][0].c.Logger.Error(err.Error())
		}
		uploaded += u
		ignored += i
		return
	})
	return
}
