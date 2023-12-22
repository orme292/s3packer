package s3pack

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	app "github.com/orme292/s3packer/config"
)

type IterErrs []error

func (ie IterErrs) String() (s string) {
	if len(ie) == 0 {
		return EmptyString
	}

	for _, err := range ie {
		s += err.Error() + "\n"
	}
	return s
}

type Iterator interface {
	Err() error // usage tbd
	Next() bool
	Finish() error
	UploadObject() *BatchPutObject
}

type BatchPutObject struct {
	Object *s3.PutObjectInput
	Before func() error
	After  func() error
}

type ObjectIterator struct {
	ol    ObjectList
	stage struct {
		index int
		fo    *FileObject
		f     *os.File
	}
	group int
	err   error
	c     *app.Configuration
}

func NewObjectIterator(c *app.Configuration, ol ObjectList, g int) Iterator {
	return &ObjectIterator{
		ol:    ol,
		group: g,
		c:     c,
	}
}

func (oi *ObjectIterator) Err() error {
	if oi.err != nil {
		oi.stage.fo.IsUploaded = false
		oi.c.Logger.Debug(fmt.Sprintf("ObjectIterator.Err() called, err: %v", oi.err))
	}
	return oi.err
}

func (oi *ObjectIterator) Next() bool {
	if len(oi.ol) == 0 {
		return false
	}

	for {
		if oi.stage.index >= len(oi.ol) {
			return false
		}
		if oi.ol[oi.stage.index].Group != oi.group {
			oi.stage.index += 1
			continue
		}
		if oi.ol[oi.stage.index].IsUploaded || oi.ol[oi.stage.index].Ignore {
			oi.c.Logger.Warn(fmt.Sprintf("Ignoring %q, %s", oi.ol[oi.stage.index].PrefixedName,
				oi.ol[oi.stage.index].IgnoreString))
			oi.stage.index += 1
			continue
		}
		break
	}

	f, err := os.Open(oi.ol[oi.stage.index].AbsolutePath)
	oi.err = err
	oi.stage.f = f
	oi.stage.fo = oi.ol[oi.stage.index]

	return oi.Err() == nil
}

func (oi *ObjectIterator) Finish() error {
	return nil
}

func (oi *ObjectIterator) UploadObject() *BatchPutObject {
	f := oi.stage.f
	return &BatchPutObject{
		Before: func() error {
			oi.c.Logger.Info(fmt.Sprintf("Transferring (%s) %q...", FileSizeString(oi.stage.fo.FileSize),
				oi.stage.fo.PrefixedName))
			return nil
		},
		Object: &s3.PutObjectInput{
			ACL:               oi.c.GetACL(oi.c.Options[app.ProfileOptionACL].(string)),
			Body:              f,
			Bucket:            aws.String(oi.c.Bucket[app.ProfileBucketName].(string)),
			ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
			ChecksumSHA256:    aws.String(oi.stage.fo.Checksum),
			Key:               aws.String(oi.stage.fo.PrefixedName),
			StorageClass:      oi.c.GetStorageClass(oi.c.Options[app.ProfileOptionStorage].(string)),
			Tagging:           aws.String(oi.stage.fo.Tags),
		},
		After: func() error {
			oi.stage.fo.IsUploaded = true
			oi.stage.index += 1
			return f.Close()
		},
	}
}

func UploadWithIterator(c *app.Configuration, iter Iterator) (errs IterErrs) {
	svc, err := BuildUploader(c)
	if err != nil {
		errs = append(errs, err)
		return
	}

	for iter.Next() {
		object := iter.UploadObject()
		if object.Before != nil {
			if err = object.Before(); err != nil {
				errs = append(errs, err)
			}
		}
		_, err = svc.Upload(context.TODO(), object.Object)
		if err != nil {
			errs = append(errs, err)
		}
		if object.After == nil {
			continue
		}
		if err = object.After(); err != nil {
			errs = append(errs, err)
		}
	}
	if err = iter.Finish(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) == 0 {
		return nil
	}
	return
}
