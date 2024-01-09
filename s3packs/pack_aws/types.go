package pack_aws

import (
	"context"
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

	ErrorCouldNotAssertObject = "could not assert object as *s3.PutObjectInput"
	ErrorNotImplemented       = "not implemented"
)

var (
	s3Error        smithy.APIError
	s3NotFound     *types.NotFound
	s3NoSuchKey    *types.NoSuchKey
	s3NoSuchBucket *types.NoSuchBucket
)

type AwsIterator struct {
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

type AwsOperator struct {
	ac     *conf.AppConfig
	client *s3.Client
	svc    *manager.Uploader
	ctl    *MultipartControl
}

type MultipartControl struct {
	uploadId string
	upload   map[int]*mpu
	cmo      *s3.CreateMultipartUploadOutput
	ctx      context.Context
	cancel   context.CancelFunc
	obj      *s3.PutObjectInput
	max      int
	retry    int
	temp     string
}

type mpu struct {
	index    int
	input    *s3.UploadPartInput
	output   *s3.UploadPartOutput
	cs       string
	data     []byte
	etag     string
	group    int
	err      error
	complete bool
}
