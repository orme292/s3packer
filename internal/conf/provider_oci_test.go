package conf

import (
	"log"
	"testing"
)

func TestProviderOCIBuildFuncs(t *testing.T) {
	profile := newIncomingProfile()
	oci := &ProviderOCI{}

	err := oci.build(profile)
	if err != nil {
		log.Printf("Object ACL: %s, Storage: %s\n", profile.Google.ObjectACL, profile.Google.Storage)
		t.Fatal(failMsg("valid profile with sample data", "err = nil", "err != nil", err))
	}

	profile.OCI.Storage = testInvalidString
	err = oci.build(profile)
	if err == nil {
		t.Fatal(failMsg("invalid storage class", "err != nil", "err = nil"))
	}

	profile.OCI.Storage = "standard"
	profile.Provider.Profile = testEmptyString
	err = oci.build(profile)
	if err == nil {
		t.Fatal(failMsg("empty profile", "err != nil", "err = nil"))
	}
}

func TestProviderOCIStringMatches(t *testing.T) {
	profile := newIncomingProfile()
	oci := &ProviderOCI{}

	for storage := range ociStorageTiersMap {
		profile.OCI.Storage = camelCaseByChar(storage)
		err := oci.build(profile)
		if err != nil {
			t.Log(failMsg("storage tier check", "err = nil", "err != nil", err))
		}
	}

	for storage := range ociPutStorageTiersMap {
		profile.OCI.Storage = camelCaseByChar(storage)
		err := oci.build(profile)
		if err != nil {
			t.Log(failMsg("put storage tier check", "err = nil", "err != nil", err))
		}
	}
}
