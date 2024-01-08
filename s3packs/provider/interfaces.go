package provider

import (
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
)

/*
Operator is an interface used to handle operations that will be performed
in the provider's implementation of S3. The following fields and methods
are required to implement this interface:

These methods are required:

			CreateBucket() (err error)
				- CreateBucket creates a new bucket at the provider.
			Get(key string) (obj *GetObject, err error)
				- Get returns a GetObject struct for a given key.
	            - This method is required, but does not need to be implemented.
			ObjectExists(key string) (exists bool, err error)
				- ObjectExists returns true if the object exists at the provider.
			Upload(po PutObject) (err error)
				- Upload uploads an object to the provider.
			UploadMultipart(po PutObject) (err error)
				- UploadMultipart uploads an object to the provider using a multipart
		          upload.
				- This method is required, but does not need to be implemented. Make
				  sure SupportsMultipartUploads returns false if this is not implemented.
			BucketExists() (exists bool, errs Errs)
				- BucketExists returns true if the bucket exists at the provider.

			SupportsMultipartUploads() bool
				- SupportsMultipartUploads returns true if the provider supports
				  multipart uploads.

These fields are required:

	ac *conf.AppConfig
*/
type Operator interface {
	CreateBucket() (err error)

	Get(key string) (obj *GetObject, err error)
	Upload(po PutObject) (err error)
	UploadMultipart(po PutObject) (err error)

	SupportsMultipartUploads() bool

	BucketExists() (exists bool, errs Errs)
	ObjectExists(key string) (exists bool, err error)
}

type GetObject interface {
	Object() Object
	Before()
	After()
}

/*
Iterator is an interface used to handle iterating over a list of files to
perform operations. The purpose is to iterator through a FileObjList and
upload each file to the provider's S3 service. The following fields and
methods are required to implement this interface:

These methods are required:

	First() (err error)
		- A function that runs before the iterator loop begins.
	Next() bool
		- Next should return if there is another iteration.
	Prepare() *PutObject
		- Prepare returns a PutObject struct for the file that will be
		  uploaded next.
	Final() (err error)
		- Final is called after the last iteration.
	Err() (err error)
		- Err returns the error from the last iteration.
	MarkIgnore(s string)
		- MarkIgnore marks the current FileObj as ignored and sets the string to s
*/
type Iterator interface {
	First() (err error)
	Next() bool
	Prepare() *PutObject
	Final() (err error)
	Err() (err error)
	MarkIgnore(s string)
}

/*
IteratorFunc is a function that should build and return an Iterator.
It's not intended to be a constructor function; instead, the IteratorFunc
ideally calls a constructor function which returns a pointer to an
object that implements the Iterator interface, which is then returned
by the IteratorFunc.

For Example:

	Func AwsIteratorFunc(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter provider.Iterator, err error) {
		return NewIterator(ac, fol, grp)
	}

	Func NewIterator(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter *AwsIterator, err error) {
		...
*/
type IteratorFunc func(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter Iterator, err error)
