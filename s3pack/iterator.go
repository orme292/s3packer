package s3pack

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/conf"
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
	a     *conf.AppConfig
}

func NewObjectIterator(a *conf.AppConfig, ol ObjectList, g int) Iterator {
	return &ObjectIterator{
		ol:    ol,
		group: g,
		a:     a,
	}
}

func (oi *ObjectIterator) Err() error {
	if oi.err != nil {
		oi.stage.fo.IsUploaded = false
		oi.a.Log.Debug("ObjectIterator.Err() called, err: %v", oi.err)
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
			oi.a.Log.Warn("Ignoring %q, %s", oi.ol[oi.stage.index].PrefixedName,
				oi.ol[oi.stage.index].IgnoreString)
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
			oi.a.Log.Info("Transferring (%s) %q...", FileSizeString(oi.stage.fo.FileSize),
				oi.stage.fo.PrefixedName)
			return nil
		},
		Object: &s3.PutObjectInput{
			ACL:               oi.a.Provider.AwsACL,
			Body:              f,
			Bucket:            aws.String(oi.a.Bucket.Name),
			ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
			ChecksumSHA256:    aws.String(oi.stage.fo.Checksum),
			Key:               aws.String(oi.stage.fo.PrefixedName),
			StorageClass:      oi.a.Provider.AwsStorage,
			Tagging:           aws.String(oi.stage.fo.Tags),
		},
		After: func() error {
			oi.stage.fo.IsUploaded = true
			oi.stage.index += 1
			return f.Close()
		},
	}
}

func UploadWithIterator(a *conf.AppConfig, iter Iterator) (errs IterErrs) {
	svc, err := BuildUploader(a)
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
		_, err = svc.Upload(context.Background(), object.Object)
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
