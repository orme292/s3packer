package conf

import (
	"log"
	"testing"
)

func TestProviderGoogleBuildFuncs(t *testing.T) {
	profile := newIncomingProfile()
	google := &ProviderGoogle{}

	err := google.build(profile)
	if err != nil {
		log.Printf("Object ACL: %s, Storage: %s\n", profile.Google.ObjectACL, profile.Google.Storage)
		t.Fatal(failMsg("valid profile with sample data", "err = nil", "err != nil", err))
	}

	profile.Google.Project = testEmptyString
	err = google.build(profile)
	if err == nil {
		t.Fatal(failMsg("empty project", "err != nil", "err = nil"))
	}

	profile.Google.Project = testSomeString
	profile.Google.ObjectACL = testInvalidString
	err = google.build(profile)
	if err == nil {
		t.Fatal(failMsg("invalid object acl", "err != nil", "err = nil"))
	}
	profile.Google.ObjectACL = googleObjectACLs[GCObjectACLPrivate]

	profile.Google.BucketACL = testInvalidString
	err = google.build(profile)
	if err == nil {
		t.Fatal(failMsg("invalid bucket acl", "err != nil", "err = nil"))
	}

	profile.Google.BucketACL = googleBucketACLs[GCBucketACLPrivate]
	profile.Google.Storage = testInvalidString
	err = google.build(profile)
	if err == nil {
		t.Fatal(failMsg("invalid storage class", "err != nil", "err = nil"))
	}

	profile.Google.Storage = googleStorageClass[GCStorageStandard]
	profile.Google.LocationType = testInvalidString
	err = google.build(profile)
	if err == nil {
		t.Fatal(failMsg("invalid location type", "err != nil", "err = nil"))
	}

	profile.Google.LocationType = googleLocationType[GCLocationTypeRegion]
}

func TestProviderGoogleStringMatches(t *testing.T) {
	profile := newIncomingProfile()
	google := &ProviderGoogle{}

	for acl := range googleObjectACLs {
		profile.Google.ObjectACL = camelCaseByChar(acl)
		err := google.build(profile)
		if err != nil {
			t.Log(failMsg("object acl check", "err = nil", "err != nil", err))
		}
	}

	for acl := range googleBucketACLs {
		profile.Google.BucketACL = camelCaseByChar(acl)
		err := google.build(profile)
		if err != nil {
			t.Log(failMsg("bucket acl check", "err = nil", "err != nil", err))
		}
	}

	for storage := range googleStorageClass {
		profile.Google.Storage = camelCaseByChar(storage)
		err := google.build(profile)
		if err != nil {
			t.Log(failMsg("storage class check", "err = nil", "err != nil", err))
		}
	}

	for location := range googleLocationType {
		profile.Google.LocationType = camelCaseByChar(location)
		err := google.build(profile)
		if err != nil {
			t.Log(failMsg("location type check", "err = nil", "err != nil", err))
		}
	}
}
