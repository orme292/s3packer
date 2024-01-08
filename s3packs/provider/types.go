package provider

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/s3packs/objectify"
)

const (
	EmptyPath   = ""
	EmptyString = ""

	DisregardGroups = -1

	MultipartThreshold = 10000000
	ObjectExists       = "object already exists"
)

type Object struct {
	F        os.File
	Key      string
	Checksum string
	Err      types.Error
}

type PutObject struct {
	Before func() error
	Object func() any
	After  func() error
	Fo     func() *objectify.FileObj
	Output func() Object
}

type Errs struct {
	Each []error
}

func (e *Errs) Add(err error) {
	e.Each = append(e.Each, err)
}

func (e *Errs) Append(errs Errs) {
	e.Each = append(e.Each, errs.Each...)
}

func (e *Errs) Release() {
	e.Each = e.Each[:0]
	e.Each = nil
}
