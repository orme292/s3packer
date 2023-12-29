package conf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	AwsACLPrivate                = "private"
	AwsACLPublicRead             = "public-read"
	AwsACLPublicReadWrite        = "public-read-write"
	AwsACLAuthenticatedRead      = "authenticated-read"
	AwsACLAwsExecRead            = "aws-exec-read"
	AwsACLBucketOwnerRead        = "bucket-owner-read"
	AwsACLBucketOwnerFullControl = "bucket-owner-full-control"
)

func awsMatchACL(acl string) (cAcl types.ObjectCannedACL, err error) {
	switch strings.ToLower(strings.Trim(acl, " ")) {
	case AwsACLPrivate:
		return types.ObjectCannedACLPrivate, nil
	case AwsACLPublicRead:
		return types.ObjectCannedACLPublicRead, nil
	case AwsACLPublicReadWrite:
		return types.ObjectCannedACLPublicReadWrite, nil
	case AwsACLAuthenticatedRead:
		return types.ObjectCannedACLAuthenticatedRead, nil
	case AwsACLAwsExecRead:
		return types.ObjectCannedACLAwsExecRead, nil
	case AwsACLBucketOwnerRead:
		return types.ObjectCannedACLBucketOwnerRead, nil
	case AwsACLBucketOwnerFullControl:
		return types.ObjectCannedACLBucketOwnerFullControl, nil
	default:
		return types.ObjectCannedACLPrivate, errors.New(fmt.Sprintf("%s %q", InvalidACL, acl))
	}
}

const (
	AwsClassStandard           = "STANDARD"
	AwsClassReducedRedundancy  = "REDUCED_REDUNDANCY"
	AwsClassGlacierIR          = "GLACIER_IR"
	AwsClassSnow               = "SNOW"
	AwsClassStandardIA         = "STANDARD_IA"
	AwsClassOneZoneIA          = "ONEZONE_IA"
	AwsClassIntelligentTiering = "INTELLIGENT_TIERING"
	AwsClassGlacier            = "GLACIER"
	AwsClassDeepArchive        = "DEEP_ARCHIVE"
)

func awsMatchStorage(class string) (sClass types.StorageClass, err error) {
	switch strings.ToUpper(strings.Trim(class, " ")) {
	case AwsClassStandard:
		return types.StorageClassStandard, nil
	case AwsClassReducedRedundancy:
		return types.StorageClassReducedRedundancy, nil
	case AwsClassGlacier:
		return types.StorageClassGlacier, nil
	case AwsClassStandardIA:
		return types.StorageClassStandardIa, nil
	case AwsClassOneZoneIA:
		return types.StorageClassOnezoneIa, nil
	case AwsClassIntelligentTiering:
		return types.StorageClassIntelligentTiering, nil
	case AwsClassGlacierIR:
		return types.StorageClassGlacierIr, nil
	case AwsClassDeepArchive:
		return types.StorageClassDeepArchive, nil
	case AwsClassSnow:
		return types.StorageClassSnow, nil
	default:
		return types.StorageClassStandard, errors.New(fmt.Sprintf("%s %q", InvalidStorageClass, class))
	}
}
