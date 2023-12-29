// Package s3pack provides functions for uploading files to s3.
// This file implements the FileObject type and its methods. A FileObject contains information on a specific file.
// https://github.com/orme292/s3packer is licensed under the MIT License.
package s3pack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/conf"
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
	IsUploaded      bool

	Group int

	a *conf.AppConfig
}

/*
NewFileObject is a FileList constructor. It takes a path and returns a FileList. Basic FileObject fields are prefilled
based on the provided path.
*/
func NewFileObject(a *conf.AppConfig, path string) (fo *FileObject, err error) {
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
		a:               a,
	}, nil
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
		if fo.a.Objects.OmitOriginDirectory {
			relativePath = strings.TrimPrefix(
				strings.Replace(fo.OriginDirectory, fo.RelativeRoot, EmptyString, 1), "/")
		} else {
			relativePath = strings.TrimPrefix(
				strings.Replace(fo.OriginDirectory, filepath.Dir(fo.RelativeRoot), EmptyString, 1), "/")
		}
		fo.PrefixedName = strings.TrimPrefix(
			AppendPathPrefix(fo.a, fmt.Sprintf("%s/%s", relativePath, AppendObjectPrefix(fo.a, fo.BaseName))), "/")
	} else {
		fo.PrefixedName = strings.TrimPrefix(
			AppendPathPrefix(fo.a, AppendObjectPrefix(fo.a, fo.BaseName)), "/")
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
	fo.PrefixedName = AppendPathPrefix(fo.a, fmt.Sprintf("%s/%s", fo.OriginDirectory, AppendObjectPrefix(fo.a, fo.BaseName)))
}

/*
SetChecksum is a FileObject method. It calls CalcChecksumSHA256 on the FileObject's AbsolutePath and sets the FileObject's
Checksum to the returned value. Checksums are attached as Tags to the S3 Object.
*/
func (fo *FileObject) SetChecksum() (err error) {
	if fo.Ignore {
		return
	}
	fo.a.Log.Debug("Calculating checksum for %q", fo.BaseName)
	hash, err := CalcChecksumSHA256(fo.AbsolutePath)
	if err != nil {
		fo.a.Log.Error("Error getting checksum for %q: %s", fo.BaseName, err.Error())
		fo.DebugOutput()
		return
	}
	fo.Checksum = hash
	fo.Tag(ChecksumSha256, hash)
	return
}

/*
SetIsDirectoryPart is a FileList helper method. It sets the FileList's IsDirectoryPart bool to true.

This should not be called directly. It should only be called by the ObjectList.setAsDirectoryPart method. Individual
values will be ignored, only the first FileObject in an ObjectList will be checked.
*/
func (fo *FileObject) SetIsDirectoryPart() {
	fo.IsDirectoryPart = true
}

/*
SetFileSize is a FileObject method. It calls GetFileSize on the FileObject's AbsolutePath and sets the FileObject's
FileSize to the returned value.
*/
func (fo *FileObject) SetFileSize() (size int64) {
	size, err := GetFileSize(fo.AbsolutePath)
	if err != nil {
		fo.a.Log.Error("Error getting file size for %q: %s", fo.BaseName, err.Error())
	}
	fo.FileSize = size
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
		fo.a.Log.Error("%s: %s, ignoring.", ErrLocalErrorOnCheck, err.Error())
		fo.SetIgnore(ErrLocalErrorOnCheck)
		fo.DebugOutput()
		return
	}
	if !exists {
		fo.SetIgnore(ErrLocalDoesNotExist)
	}
}

/*
SetIgnoreIfObjExistsInBucket is a FileObject method. If the object is not already ignored or the overwrite option is set, then
it calls ObjectExists on the FileObject's PrefixedName. If ObjectExists returns true, then the FileObject's Ignore is set
to true and an IgnoreString is set.
*/
func (fo *FileObject) SetIgnoreIfObjExistsInBucket() {
	client, _ := BuildClient(fo.a)
	fo.SetIgnoreIfObjExistsInBucketWithClient(client)
}

func (fo *FileObject) SetIgnoreIfObjExistsInBucketWithClient(client *s3.Client) {
	if fo.Ignore || (fo.a.Opts.Overwrite == conf.OverwriteAlways) {
		return
	}
	exists, err := ObjectExistsWithClient(fo.a, fo.PrefixedName, client)
	if err != nil {
		fo.a.Log.Error("%s: %s, ignoring.", ErrIgnoreObjectErrorOnCheck, err.Error())
		fo.SetIgnore(ErrIgnoreObjectErrorOnCheck)
		return
	}
	if exists {
		fo.SetIgnore(ErrIgnoreObjectAlreadyExists)
	}
}

/*
SetPrefixedName is a FileObject method. It calls the appropriate key naming method based on the keyNamingMethod option
in the profile. See the nameMethod* functions for more information.
*/
func (fo *FileObject) SetPrefixedName() {
	switch fo.a.Objects.Naming {
	case conf.NamingRelative:
		fo.SetNameMethodRelative()
	case conf.NamingAbsolute:
		fo.SetNameMethodAbsolute()
	default:
		fo.a.Log.Fatal(ErrKeyNameMust)
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
*/
func (fo *FileObject) Upload() (uploaded bool, err error) {
	svc, err := BuildUploader(fo.a)
	if err != nil {
		return false, err
	}
	return fo.UploadWithProvided(svc)
}

func (fo *FileObject) UploadWithProvided(svc *manager.Uploader) (uploaded bool, err error) {
	f, err := os.Open(fo.AbsolutePath)
	if err != nil {
		return false, err
	}

	if fo.Ignore || fo.IsUploaded {
		return false, nil
	}

	fo.a.Log.Info("Uploading %s...", fo.PrefixedName)
	_, err = svc.Upload(context.Background(), &s3.PutObjectInput{
		ACL:               fo.a.Provider.AwsACL,
		Body:              f,
		Bucket:            aws.String(fo.a.Bucket.Name),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
		ChecksumSHA256:    aws.String(fo.Checksum),
		Key:               aws.String(fo.PrefixedName),
		StorageClass:      fo.a.Provider.AwsStorage,
		Tagging:           aws.String(fo.Tags),
	})
	if err != nil {
		return false, err
	}

	fo.IsUploaded = true
	return true, nil
}

/* DEBUG */

/*
DebugOutput is an ObjectList method. It prints the AbsolutePath and IsDirectoryPart fields of each FileObject in the
ObjectList slice.

Change this as needed.
*/
func (fo *FileObject) DebugOutput() {
	fmt.Println()
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
	fmt.Printf("IsUploaded: %t\n", fo.IsUploaded)
	fmt.Printf("Group: %d\n", fo.Group)
	fmt.Println()

}
