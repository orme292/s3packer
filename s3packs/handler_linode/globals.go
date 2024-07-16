package handler_linode

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
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
