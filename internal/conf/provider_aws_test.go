package conf

import (
	"log"
	"testing"
)

func TestProviderAWSBuildFuncs(t *testing.T) {
	profile := newIncomingProfile()
	aws := &ProviderAWS{}

	for acl := range awsCannedACLs {
		profile.AWS.ACL = camelCaseByChar(acl)
		for storage := range awsStorageClasses {
			profile.AWS.Storage = camelCaseByChar(storage)
			err := aws.build(profile)
			if err != nil {
				log.Printf("ACL: %s, Storage: %s\n", acl, storage)
				log.Println(err)
				t.Fail()
			}
		}
	}

	profile.AWS.ACL = testInvalidString
	err := aws.build(profile)
	if err == nil {
		log.Println("Provided invalid ACL. Expected: err != nil")
		t.Fail()
	}

	profile.AWS.ACL = AwsACLPrivate
	profile.AWS.Storage = testInvalidString
	err = aws.build(profile)
	if err == nil {
		log.Println("Provided invalid storage class. Expected: err != nil")
		t.Fail()
	}

	profile.Provider.Profile = testEmptyString
	profile.Provider.Key = testEmptyString
	profile.Provider.Secret = testEmptyString
	profile.AWS.Storage = AwsClassStandard
	err = aws.build(profile)
	if err == nil {
		log.Println("Provided no profile, key, or secret. Expected: err != nil")
		t.Fail()
	}

	profile.Provider.Key = testSomeString
	profile.Provider.Secret = testSomeString
	profile.Provider.Profile = testSomeString
	profile.AWS.ACL = AwsACLPrivate
	profile.AWS.Storage = AwsClassStandard
	err = aws.build(profile)
	if err == nil {
		log.Println("Provided profile, key, and secret. Expected: err != nil")
		t.Fail()
	}

	profile.Provider.Profile = testEmptyString
	err = aws.build(profile)
	if err != nil {
		log.Println("Provided key and secret without profile. Expected: err = nil")
		t.Fail()
	}

	profile.Provider.Profile = testSomeString
	profile.Provider.Key = testEmptyString
	profile.Provider.Secret = testEmptyString
	err = aws.build(profile)
	if err != nil {
		log.Println("Provided profile with no secret or key. Expected: err = nil")
		log.Println(err)
		t.Fail()
	}

	profile.Provider.Key = testSomeString
	err = aws.build(profile)
	if err == nil {
		log.Println("Provided profile and key with no secret. Expected: err != nil")
		t.Fail()
	}

	profile.Provider.Key = testEmptyString
	profile.Provider.Secret = testSomeString
	err = aws.build(profile)
	if err == nil {
		log.Println("Provided profile and secret with no key. Expected: err != nil")
		t.Fail()
	}

	profile.Provider.Profile = testEmptyString
	err = aws.build(profile)
	if err == nil {
		log.Println("Provider secret with no key or profile. Expected: err != nil")
		t.Fail()
	}

	profile.Provider.Key = testSomeString
	profile.Provider.Secret = testEmptyString
	err = aws.build(profile)
	if err == nil {
		log.Println("Provider key with no secret or profile. Expected: err != nil")
		t.Fail()
	}

}
