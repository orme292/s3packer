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

var googleObjectACLs = map[string]string{
	GCObjectACLAuthenticatedRead: GCObjectACLAuthenticatedRead,
	GCObjectACLPrivate:           GCObjectACLPrivate,
	GCObjectACLPublicRead:        GCObjectACLPublicRead,
	GCObjectACLProjectPrivate:    GCObjectACLProjectPrivate,
	GCObjectACLBucketOwnerFull:   GCObjectACLBucketOwnerFull,
	GCObjectACLBucketOwnerRead:   GCObjectACLBucketOwnerRead,
}

var googleBucketACLs = map[string]string{
	GCBucketACLAuthenticatedRead: GCBucketACLAuthenticatedRead,
	GCBucketACLPrivate:           GCBucketACLPrivate,
	GCBucketACLPublicRead:        GCBucketACLPublicRead,
	GCBucketACLPublicReadWrite:   GCBucketACLPublicReadWrite,
	GCBucketACLProjectPrivate:    GCBucketACLProjectPrivate,
}

var googleStorageClass = map[string]string{
	GCStorageStandard: GCStorageStandard,
	GCStorageNearline: GCStorageNearline,
	GCStorageColdline: GCStorageColdline,
	GCStorageArchive:  GCStorageArchive,
}

var googleLocationType = map[string]string{
	GCLocationTypeRegion: GCLocationTypeRegion,
	GCLocationTypeDual:   GCLocationTypeDual,
	GCLocationTypeMulti:  GCLocationTypeMulti,
}

func (gc *ProviderGoogle) build(inc *ProfileIncoming) error {

	gc.Project = inc.Google.Project

	err := gc.validate()
	if err != nil {
		return err
	}

	err = gc.matchObjectACL(inc.Google.ObjectACL)
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

	validAcl, ok := googleBucketACLs[tidyLowerString(acl)]
	if !ok {
		gc.BucketACL = GCBucketACLPrivate
		return fmt.Errorf("%s %q", InvalidGCBucketACL, acl)
	}
	gc.BucketACL = validAcl

	return nil

}

func (gc *ProviderGoogle) matchObjectACL(acl string) error {

	validAcl, ok := googleObjectACLs[tidyLowerString(acl)]
	if !ok {
		gc.ObjectACL = GCObjectACLPrivate
		return fmt.Errorf("%s %q", InvalidGCObjectACL, acl)
	}
	gc.ObjectACL = validAcl

	return nil

}

func (gc *ProviderGoogle) matchStorageClass(class string) error {

	validClass, ok := googleStorageClass[tidyUpperString(class)]
	if !ok {
		gc.Storage = GCStorageStandard
		return fmt.Errorf("%s %q", InvalidStorageClass, class)
	}
	gc.Storage = validClass

	return nil

}

func (gc *ProviderGoogle) matchLocationType(tp string) error {

	validType, ok := googleLocationType[tidyLowerString(tp)]
	if !ok {
		gc.LocationType = GCLocationTypeRegion
		return fmt.Errorf("%s %q", InvalidGCLocationType, tp)
	}
	gc.LocationType = validType

	return nil

}

func (gc *ProviderGoogle) validate() error {

	if gc.Project == Empty {
		return fmt.Errorf("project name cannot be empty")
	}

	if gc.ObjectACL == Empty {
		gc.ObjectACL = GCObjectACLPrivate
	}

	if gc.BucketACL == Empty {
		gc.BucketACL = GCBucketACLPrivate
	}

	return nil
}
