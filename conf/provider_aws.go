package conf

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ProviderAWS represents the AWS provider configuration.
//
// Fields:
// - Profile: The profile name used for authentication.
// - ACL: The access control list for the storage objects.
// - Storage: The storage class for the objects.
// - Key: The AWS access key ID.
// - Secret: The AWS secret access key.
type ProviderAWS struct {
	Profile              string
	Key                  string
	Secret               string
	ACL                  types.ObjectCannedACL
	Storage              types.StorageClass
	AwsChecksumAlgorithm types.ChecksumAlgorithm
	AwsChecksumMode      types.ChecksumMode
}

func (aws *ProviderAWS) build(inc *ProfileIncoming) error {

	err := aws.matchACL(inc.AWS.ACL)
	if err != nil {
		return err
	}

	err = aws.matchStorage(inc.AWS.Storage)
	if err != nil {
		return err
	}

	aws.Key = inc.Provider.Key
	aws.Secret = inc.Provider.Secret
	aws.Profile = inc.Provider.Profile

	aws.AwsChecksumAlgorithm = types.ChecksumAlgorithmSha256
	aws.AwsChecksumMode = types.ChecksumModeEnabled

	return aws.validate()

}

// awsMatchACL will match the ACL string to the AWS ACL type. The constant values above are used to match the string.
func (aws *ProviderAWS) matchACL(acl string) error {

	awsCannedACLs := map[string]types.ObjectCannedACL{
		AwsACLPrivate:                types.ObjectCannedACLPrivate,
		AwsACLPublicRead:             types.ObjectCannedACLPublicRead,
		AwsACLPublicReadWrite:        types.ObjectCannedACLPublicReadWrite,
		AwsACLAuthenticatedRead:      types.ObjectCannedACLAuthenticatedRead,
		AwsACLAwsExecRead:            types.ObjectCannedACLAwsExecRead,
		AwsACLBucketOwnerRead:        types.ObjectCannedACLBucketOwnerRead,
		AwsACLBucketOwnerFullControl: types.ObjectCannedACLBucketOwnerFullControl,
	}

	validAcl, ok := awsCannedACLs[tidyString(acl)]
	if !ok {
		aws.ACL = types.ObjectCannedACLPrivate
		return fmt.Errorf("%s %q", InvalidAWSACL, acl)
	}
	aws.ACL = validAcl

	return nil

}

// awsMatchStorage will match the storage class string to the AWS storage class type. The constant values above are
// used to match the string.
func (aws *ProviderAWS) matchStorage(class string) error {

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

	validClass, ok := awsStorageClasses[tidyUpString(class)]
	if !ok {
		aws.Storage = types.StorageClassStandard
		return fmt.Errorf("%s %q", InvalidStorageClass, class)
	}
	aws.Storage = validClass

	return nil

}

func (aws *ProviderAWS) validate() error {

	if aws.Profile != Empty && (aws.Key != Empty || aws.Secret != Empty) {
		return fmt.Errorf("bad AWS config: %v", ErrorAWSProfileAndKeys)
	}
	if aws.Profile == Empty && (aws.Key == Empty || aws.Secret == Empty) {
		return fmt.Errorf("bad AWS config: %v", ErrorAWSMissingAuth)
	}
	if aws.Profile == Empty && aws.Key == Empty && aws.Secret == Empty {
		return fmt.Errorf("bad AWS config: %v", ErrorAWSAuthNeeded)
	}

	return nil

}
