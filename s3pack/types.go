package s3pack

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

const (
	ChecksumSha256 = "SHA256"
	EmptyString    = ""
	Version        = "1.1.0"
)

const (
	ErrIgnoreObjectKeyEmpty      = "object key is empty"
	ErrIgnoreObjectErrorOnCheck  = "Error checking if object exists"
	ErrIgnoreObjectAlreadyExists = "Object of same name already exists."

	ErrLocalErrorOnCheck = "Error checking if local file exists"
	ErrLocalDoesNotExist = "Local file does not exist or is not accessible."

	ErrKeyNameMust = "A key naming method must be specified"
)

var (
	s3Error        smithy.APIError
	s3NotFound     *types.NotFound
	s3NoSuchKey    *types.NoSuchKey
	s3NoSuchBucket *types.NoSuchBucket
)
