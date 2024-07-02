package provider_v2

import (
	"github.com/orme292/s3packer/conf"
)

type Operator interface {
	BucketCreate() error
	BucketExists() (bool, error)
	BucketDelete() error

	ObjectExists(key string) (bool, error)
	ObjectDelete(key string) error
	ObjectUpload() error

	GetObjectTags(key string) (map[string]string, error)

	Support() *Supports
}

type OperGenFunc func(app *conf.AppConfig) (oper Operator, err error)

type Iterator interface {
	Pre() error
	Next() bool
	Prepare() error
	Post() error
}

type IterGenFunc func(app *conf.AppConfig) (iter Iterator, err error)
