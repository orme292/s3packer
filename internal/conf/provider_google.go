package conf

import (
	"fmt"
)

type ProviderGoogle struct {
	Project      string
	LocationType string
	Storage      string
	BucketACL    string
	ObjectACL    string
	ADC          string
}

func (gc *ProviderGoogle) build(inc *ProfileIncoming) error {

	gc.Project = inc.Google.Project

	_ = gc.validate()

	err := gc.matchObjectACL(inc.Google.ObjectACL)
	if err != nil {
		return err
	}

	err = gc.matchBucketACL(inc.Google.BucketACL)
	if err != nil {
		return err
	}

	err = gc.matchStorageClass(inc.Google.Storage)
	if err != nil {
		return err
	}

	err = gc.matchLocationType(inc.Google.LocationType)
	if err != nil {
		return err
	}

	return nil
}

func (gc *ProviderGoogle) matchBucketACL(acl string) error {

	bucketCannedACLs := map[string]string{
		GCBucketACLAuthenticatedRead: GCBucketACLAuthenticatedRead,
		GCBucketACLPrivate:           GCBucketACLPrivate,
		GCBucketACLPublicRead:        GCBucketACLPublicRead,
		GCBucketACLPublicReadWrite:   GCBucketACLPublicReadWrite,
		GCBucketACLProjectPrivate:    GCBucketACLProjectPrivate,
	}

	validAcl, ok := bucketCannedACLs[tidyLowerString(acl)]
	if !ok {
		gc.BucketACL = GCBucketACLPrivate
		return fmt.Errorf("%s %q", InvalidGCBucketACL, acl)
	}
	gc.BucketACL = validAcl

	return nil

}

func (gc *ProviderGoogle) matchObjectACL(acl string) error {

	objectCannedACL := map[string]string{
		GCObjectACLAuthenticatedRead: GCObjectACLAuthenticatedRead,
		GCObjectACLPrivate:           GCObjectACLPrivate,
		GCObjectACLPublicRead:        GCObjectACLPublicRead,
		GCObjectACLProjectPrivate:    GCObjectACLProjectPrivate,
		GCObjectACLBucketOwnerFull:   GCObjectACLBucketOwnerFull,
		GCObjectACLBucketOwnerRead:   GCObjectACLBucketOwnerRead,
	}

	validAcl, ok := objectCannedACL[tidyLowerString(acl)]
	if !ok {
		gc.ObjectACL = GCObjectACLPrivate
		return fmt.Errorf("%s %q", InvalidGCObjectACL, acl)
	}
	gc.ObjectACL = validAcl

	return nil

}

func (gc *ProviderGoogle) matchStorageClass(class string) error {

	storageClass := map[string]string{
		GCStorageStandard: GCStorageStandard,
		GCStorageNearline: GCStorageNearline,
		GCStorageColdline: GCStorageColdline,
		GCStorageArchive:  GCStorageArchive,
	}

	validClass, ok := storageClass[tidyUpperString(class)]
	if !ok {
		gc.Storage = GCStorageStandard
		return fmt.Errorf("%s %q", InvalidStorageClass, class)
	}
	gc.Storage = validClass

	return nil

}

func (gc *ProviderGoogle) matchLocationType(tp string) error {

	locationType := map[string]string{
		GCLocationTypeRegion: GCLocationTypeRegion,
		GCLocationTypeDual:   GCLocationTypeDual,
		GCLocationTypeMulti:  GCLocationTypeMulti,
	}

	validType, ok := locationType[tidyLowerString(tp)]
	if !ok {
		gc.LocationType = GCLocationTypeRegion
		return fmt.Errorf("%s %q", InvalidGCLocationType, tp)
	}
	gc.LocationType = validType

	return nil

}

func (gc *ProviderGoogle) validate() error {

	if gc.ObjectACL == Empty {
		gc.ObjectACL = GCObjectACLPrivate
	}

	if gc.BucketACL == Empty {
		gc.BucketACL = GCBucketACLPrivate
	}

	return nil
}
