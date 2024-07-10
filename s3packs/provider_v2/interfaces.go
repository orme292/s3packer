package provider_v2

import (
	"github.com/orme292/s3packer/conf"
)

type Operator interface {
	BucketCreate() error
	BucketExists() (bool, error)
	BucketDelete() error

	ObjectDelete(key string) error
	ObjectExists(obj Object) (bool, error)
	ObjectUpload(obj Object) error
	GetObjectTags(key string) (map[string]string, error)

	Support() *Supports
}

type ObjectGenFunc func(job *Job) Object
type OperGenFunc func(app *conf.AppConfig) (oper Operator, err error)

type Object interface {
	Destroy() error
	Generate() error
	Post() error
	Pre() error
}
