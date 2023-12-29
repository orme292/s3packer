package s3pack

import (
	"github.com/orme292/s3packer/conf"
)

/*
RootList is a slice of DirectoryList pointers (slices are inherently pointers).

Structure:
RootList => DirectoryList => DirectoryObject -> ObjectList => FileObject
RootList is a slice of DirectoryList Pointers.
DirectoryList is a slice containing DirectoryObject pointers.
DirectoryObject is a struct with an ObjectList pointer.
ObjectList is a slice of FileObject pointers.
FileObject is a struct that holds information about a file.
*/
type RootList []DirectoryList

/*
NewRootList is a RootList constructor. It takes directory paths as a slice of strings and returns a RootList.

See NewDirectoryList for more information
*/
func NewRootList(a *conf.AppConfig, dir []string) (dl RootList, err error) {
	for _, d := range dir {
		dirList, err := NewDirectoryList(a, d)
		if err != nil {
			return nil, err
		}
		dl = append(dl, dirList)
	}
	return
}

/*
IterateAndExecute is an RootList method. It takes a function that takes a DirectoryList pointer and returns
an error. It iterates over the DirectoryList slice and calls the provided function on each slice. If the function returns an
error, then it is returned and iteration stops.
*/
func (rl RootList) IterateAndExecute(fn DirListIterationFunc) (err error) {
	for index := range rl {
		if err = fn(rl[index]); err != nil {
			return
		}
	}
	return
}

/*
DirListIterationFunc is a function type that takes a DirectoryList slice and returns an error. It is used with
RootList.IterateAndExecute
*/
type DirListIterationFunc func(dl DirectoryList) (err error)

/*
SetPrefixedNames is a RootList method. It iterates over the RootList slice and calls the
DirectoryList.SetPrefixedNames method on each DirectoryList pointer.

See DirectoryList.SetPrefixedNames for more information
*/
func (rl RootList) SetPrefixedNames() {
	_ = rl.IterateAndExecute(func(dl DirectoryList) (err error) {
		dl.SetPrefixedNames()
		return
	})
}

/*
Upload is a RootList method. It iterates over the RootList slice and calls the DirectoryList.Upload method
on each DirectoryList pointer.

See DirectoryList.Upload for more information
*/
func (rl RootList) Upload() (err error, bytes int64, uploaded, ignored int) {
	_ = rl.IterateAndExecute(func(dl DirectoryList) (err error) {
		err = dl.Upload()
		if err != nil {
			rl[0][0].a.Log.Error(err.Error())
		}
		return
	})
	return err, rl.GetTotalUploadedBytes(), rl.CountUploaded(), rl.CountIgnored()
}

/*
CountIgnored is a RootList method. It returns the number of FileObjects in the slice with Ignore set to true. It
iterates over the RootList slice and calls the DirectoryList.CountFileObjects method on each DirectoryList pointer.

See DirectoryList.CountIgnored for more information
*/
func (rl RootList) CountIgnored() (count int) {
	for _, dl := range rl {
		count += dl.CountIgnored()
	}
	return
}

/*
CountUploaded is a RootList method. It returns the number of FileObjects in the slice with IsUploaded set to true. It
iterates over the RootList slice and calls the DirectoryList.CountFileObjects method on each DirectoryList pointer.

See DirectoryList.CountUploaded for more information
*/
func (rl RootList) CountUploaded() (count int) {
	for _, dl := range rl {
		count += dl.CountUploaded()
	}
	return
}

/*
GetTotalUploadedBytes is a RootList method. It returns the total number of bytes uploaded. It iterates over the RootList
slice and calls the DirectoryList.GetTotalUploadedBytes method on each DirectoryList pointer.

See DirectoryList.GetTotalUploadedBytes for more information
*/
func (rl RootList) GetTotalUploadedBytes() (total int64) {
	for _, dl := range rl {
		total += dl.GetTotalUploadedBytes()
	}
	return
}
