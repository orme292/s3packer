package pack_akamai

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
)

const (
	EmptyString = ""
)

const (
	ErrorCouldNotAssertObject = "could not assert object"
)

var (
	s3Error        smithy.APIError
	s3NotFound     *types.NotFound
	s3NoSuchKey    *types.NoSuchKey
	s3NoSuchBucket *types.NoSuchBucket
)

type AkamaiIterator struct {
	provider *conf.Provider
	svc      *manager.Uploader
	fol      objectify.FileObjList
	stage    struct {
		index int
		fo    *objectify.FileObj
		f     *os.File
	}
	group int
	err   error
	ac    *conf.AppConfig
}

type AkamaiOperator struct {
	ac     *conf.AppConfig
	client *s3.Client
	svc    *manager.Uploader
}
