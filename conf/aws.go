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

// awsMatchACL will match the ACL string to the AWS ACL type. The constant values above are used to match the string.
func awsMatchACL(acl string) (cAcl types.ObjectCannedACL, err error) {
	acl = strings.ToLower(strings.Trim(acl, " "))
	awsCannedACLs := map[string]types.ObjectCannedACL{
		AwsACLPrivate:                types.ObjectCannedACLPrivate,
		AwsACLPublicRead:             types.ObjectCannedACLPublicRead,
		AwsACLPublicReadWrite:        types.ObjectCannedACLPublicReadWrite,
		AwsACLAuthenticatedRead:      types.ObjectCannedACLAuthenticatedRead,
		AwsACLAwsExecRead:            types.ObjectCannedACLAwsExecRead,
		AwsACLBucketOwnerRead:        types.ObjectCannedACLBucketOwnerRead,
		AwsACLBucketOwnerFullControl: types.ObjectCannedACLBucketOwnerFullControl,
	}

	cAcl, ok := awsCannedACLs[acl]
	if !ok {
		return types.ObjectCannedACLPrivate, errors.New(fmt.Sprintf("%s %q", InvalidAWSACL, acl))
	}
	return cAcl, nil
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

// awsMatchStorage will match the storage class string to the AWS storage class type. The constant values above are
// used to match the string.
func awsMatchStorage(class string) (sClass types.StorageClass, err error) {
	class = strings.ToUpper(strings.Trim(class, " "))
	awsStorageClasses := map[string]types.StorageClass{
		AwsClassStandard:           types.StorageClassStandard,
		AwsClassReducedRedundancy:  types.StorageClassReducedRedundancy,
		AwsClassGlacier:            types.StorageClassGlacier,
		AwsClassStandardIA:         types.StorageClassStandardIa,
		AwsClassOneZoneIA:          types.StorageClassOnezoneIa,
		AwsClassIntelligentTiering: types.StorageClassIntelligentTiering,
		AwsClassGlacierIR:          types.StorageClassGlacier,
		AwsClassDeepArchive:        types.StorageClassDeepArchive,
		AwsClassSnow:               types.StorageClassGlacier,
	}

	sClass, ok := awsStorageClasses[class]
	if !ok {
		return types.StorageClassStandard, errors.New(fmt.Sprintf("%s %q", InvalidStorageClass, class))
	}
	return sClass, nil
}

const (
	InvalidAWSACL                   = "invalid aws acl"
	ErrorAWSProfileAndKeys          = "both aws profile and keys are specified, use profile or keys"
	ErrorAWSKeyOrSecretNotSpecified = "profile should specified both key and secret"
)
