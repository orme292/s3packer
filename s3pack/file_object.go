// Package s3pack provides functions for uploading files to s3.
// This file implements the FileObject type and its methods. A FileObject contains information on a specific file.
// https://github.com/orme292/s3packer is licensed under the MIT License.
package s3pack

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/orme292/s3packer/config"
)

/*
FileObject is a struct that contains information about a single file. It is used to store the information needed to
build an S3 object.
*/
type FileObject struct {
	OriginPath      string
	OriginDirectory string
	RelativeRoot    string
	AbsolutePath    string
	BaseName        string

	FileSize int64
	Checksum string

	PrefixedName string
	Tags         string

	IsDirectoryPart bool
	Ignore          bool
	IgnoreString    string
	ShouldMultiPart bool
	IsUploaded      bool

	Group int

	c config.Configuration
}

/*
NewFileObject is a FileList constructor. It takes a path and returns a FileList.
*/
func NewFileObject(c *config.Configuration, path string) (fo *FileObject, err error) {
	abPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &FileObject{
		OriginPath:      path,
		OriginDirectory: filepath.Dir(abPath),
		AbsolutePath:    abPath,
		BaseName:        filepath.Base(path),
		PrefixedName:    EmptyString,
		c:               *c,
	}, nil
}

/*
IgnoreIfObjectExistsInBucket is a FileObject method. It calls ObjectExists on the FileObject's PrefixedName.
If it returns true, then it calls FileList.SetIgnore(ErrIgnoreObjectAlreadyExists). The entire function is bypassed
if the overwrite option is set to true.
*/
func (fo *FileObject) IgnoreIfObjectExistsInBucket() error {
	if fo.c.Options[config.ProfileOptionOverwrite].(bool) {
		return nil
	}
	exists, err := ObjectExists(&fo.c, fo.PrefixedName)
	if err != nil {
		return err
	} else if exists {
		fo.SetIgnore(ErrIgnoreObjectAlreadyExists)
	}
	return nil
}

/*
IgnoreIfLocalDoesNotExist is a FileObject method. It calls LocalFileExists on the FileObject's AbsolutePath.
If it returns false, then it calls FileList.SetIgnore(ErrIgnoreLocalNotFound).
*/
func (fo *FileObject) IgnoreIfLocalDoesNotExist() error {
	exists, err := LocalFileExists(fo.AbsolutePath)
	if err != nil {
		return err
	} else if !exists {
		fo.SetIgnore(ErrIgnoreLocalNotFound)
	}

	return nil
}

/*
SetNameMethodRelative is a FileObject method. The PrefixedName is constructed by removing the FileObject's RelativeRoot from
the FileObject's OriginDirectory and then appending the FileObject's BaseName with the prefix specified in the profile.

If the FileObject is not a part of a directory upload, then there is no relative path root, so we run SetNameMethodFlat().

For example, if the provided file name is "/home/users/forrest/mysql_backup.tar.gz", the prefix is "november-2023-", and the
relative root is "/home/users" then the PrefixedName will be set to "/forrest/november-2023-mysql_backup.tar.gz"
*/
func (fo *FileObject) SetNameMethodRelative() {
	if fo.IsDirectoryPart == true {
		var relativePath string
		if fo.c.Options[config.ProfileOptionOmitOriginDir].(bool) {
			relativePath = strings.Replace(fo.OriginDirectory, fo.RelativeRoot, EmptyString, 1)
		} else {
			relativePath = strings.Replace(fo.OriginDirectory, filepath.Dir(fo.RelativeRoot), EmptyString, 1)
		}
		fo.PrefixedName = AppendPathPrefix(&fo.c, fmt.Sprintf("/%s/%s", relativePath, AppendObjectPrefix(&fo.c, fo.BaseName)))
	} else {
		fo.PrefixedName = AppendPathPrefix(&fo.c, AppendObjectPrefix(&fo.c, fo.BaseName))
	}
}

/*
SetNameMethodAbsolute is a FileObject method. It sets the FileObject's PrefixedName to the FileObject's BaseName prefixed
with FileObject.OriginDirectory and the prefix specified in the profile.

For example, if the provided base file name is "mysql_backup.tar.gz", the origin directory is "/home/users/forrest",
and the prefix is "/2023/November/mysql/" then the PrefixedName will be set
to "/home/users/forrest/2023/November/mysql/mysql_backup.tar.gz".
*/
func (fo *FileObject) SetNameMethodAbsolute() {
	fo.PrefixedName = AppendPathPrefix(&fo.c, fmt.Sprintf("%s/%s", fo.OriginDirectory, AppendObjectPrefix(&fo.c, fo.BaseName)))
}

/*
SetChecksum is a FileObject method. It calls CalcChecksumSHA256 on the FileObject's AbsolutePath and sets the FileObject's
Checksum to the returned value. Checksums are attached as Tags to the S3 Object.
*/
func (fo *FileObject) SetChecksum() (err error) {
	if fo.Ignore {
		return
	}
	hash, err := CalcChecksumSHA256(fo.AbsolutePath)
	if err != nil {
		fo.c.Logger.Error(fmt.Sprintf("Error getting checksum for %q: %s", fo.BaseName, err.Error()))
		fo.DebugOutput()
		return
	}
	fo.Checksum = hash
	fo.Tag(ChecksumSha256, hash)
	return
}

/*
SetDirectoryPart is a FileList helper method. It sets the FileList's IsDirectoryPart bool to true.

This should not be called directly. It should only be called by the ObjectList.setAsDirectoryPart method. Individual
values will be ignored, only the first FileObject in an ObjectList will be checked.
*/
func (fo *FileObject) SetDirectoryPart() {
	fo.IsDirectoryPart = true
}

/*
SetFileSize is a FileObject method. It calls GetFileSize on the FileObject's AbsolutePath and sets the FileObject's
FileSize to the returned value. If the FileSize is greater than 104857600 (100MB), then the FileObject's ShouldMultiPart is set
to true.
*/
func (fo *FileObject) SetFileSize() (size int64) {
	size, _ = GetFileSize(fo.AbsolutePath)
	fo.FileSize = size
	if size > 104857600 {
		fo.ShouldMultiPart = true
	}
	return
}

func (fo *FileObject) SetGroup(g int) {
	fo.Group = g
}

/*
SetIgnore is a FileList helper method. It sets the FileList's Ignore bool to true and IgnoreString string.
*/
func (fo *FileObject) SetIgnore(s string) {
	fo.Ignore = true
	fo.IgnoreString = s
}

/*
SetIgnoreIfLocalNotExists is a FileObject method. It calls LocalFileExists on the FileObject's AbsolutePath.
If LocalFileExists returns false, then the FileObject's Ignore is set to true and an IgnoreString is set.
*/
func (fo *FileObject) SetIgnoreIfLocalNotExists() {
	if fo.Ignore {
		return
	}
	exists, err := LocalFileExists(fo.AbsolutePath)
	if err != nil {
		fo.c.Logger.Error(fmt.Sprintf("Error checking if local file exists: %s, ignoring object", err.Error()))
		fo.SetIgnore("Error checking if local file exists")
		fo.DebugOutput()
		return
	}
	if !exists {
		fo.SetIgnore("Local file does not exist")
	}
}

/*
SetIgnoreIfObjExists is a FileObject method. If the object is not already ignored or the overwrite option is set, then
it calls ObjectExists on the FileObject's PrefixedName. If ObjectExists returns true, then the FileObject's Ignore is set
to true and an IgnoreString is set.
*/
func (fo *FileObject) SetIgnoreIfObjExists() {
	if fo.Ignore || fo.c.Options[config.ProfileOptionOverwrite].(bool) {
		return
	}
	exists, err := ObjectExists(&fo.c, fo.PrefixedName)
	if err != nil {
		fo.c.Logger.Error(fmt.Sprintf("Error checking if object exists: %s, ignoring object", err.Error()))
		fo.SetIgnore("Error checking if object exists")
		fo.DebugOutput()
		return
	}
	if exists {
		fo.SetIgnore("Object with same key already exists")
	}
}

/*
SetPrefixedName is a FileObject method. It calls the appropriate key naming method based on the keyNamingMethod option
in the profile. See the nameMethod* functions for more information.
*/
func (fo *FileObject) SetPrefixedName() {
	switch fo.c.Options["keyNamingMethod"] {
	case config.NameMethodRelative:
		fo.SetNameMethodRelative()
	case config.NameMethodAbsolute:
		fo.SetNameMethodAbsolute()
	default:
		fo.c.Logger.Fatal("A key naming method must be specified")
	}
}

/*
SetRelativeRoot is a FileObject method. It takes a string (dir) and sets the FileObject's RelativeRoot to the provided
string. To be used when building key names.
*/
func (fo *FileObject) SetRelativeRoot(dir string) {
	fo.RelativeRoot = dir
}

/*
Tag is a FileObject method. It takes a key (k) string and a value (v) string. It appends them to the FileList's Tags string
to be used as tags on the final S3 object.

Key=Value&Key2=Value2&Key3=Value3 ... (as so on)
*/
func (fo *FileObject) Tag(k, v string) {
	if fo.Tags == EmptyString {
		fo.Tags = fmt.Sprintf("%s=%s", k, v)
	} else {
		fo.Tags = fmt.Sprintf("%s&%s=%s", fo.Tags, k, v)
	}
}

/*
Upload is a FileObject method. The purpose will be to initiate a multipart upload for this FileObject.

NotImplemented
*/
func (fo *FileObject) Upload() {
	return
}

/* DEBUG */

/*
DebugOutput is an ObjectList method. It prints the AbsolutePath and IsDirectoryPart fields of each FileObject in the
ObjectList slice.

Change this as needed.
*/
func (fo *FileObject) DebugOutput() {
	fmt.Printf("\nOriginPath: %s\n", fo.OriginPath)
	fmt.Printf("OriginDirectory: %s\n", fo.OriginDirectory)
	fmt.Printf("RelativeRoot: %s\n", fo.RelativeRoot)
	fmt.Printf("AbsolutePath: %s\n", fo.AbsolutePath)
	fmt.Printf("BaseName: %s\n", fo.BaseName)
	fmt.Printf("FileSize: %d\n", fo.FileSize)
	fmt.Printf("Checksum: %s\n", fo.Checksum)
	fmt.Printf("PrefixedName: %s\n", fo.PrefixedName)
	fmt.Printf("Tags: %s\n", fo.Tags)
	fmt.Printf("IsDirectoryPart: %t\n", fo.IsDirectoryPart)
	fmt.Printf("Ignore: %t\n", fo.Ignore)
	fmt.Printf("IgnoreString: %s\n", fo.IgnoreString)
	fmt.Printf("ShouldMultiPart: %t\n", fo.ShouldMultiPart)
	fmt.Printf("IsUploaded: %t\n", fo.IsUploaded)
	fmt.Printf("Group: %d\n", fo.Group)
	fmt.Println()

}
