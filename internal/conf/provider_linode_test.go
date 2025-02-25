package conf

import (
	"testing"
)

func TestProviderLinodeBuildFuncs(t *testing.T) {
	profile := newIncomingProfile()
	linode := &ProviderLinode{}

	profile.Provider.Key = testSomeString
	profile.Provider.Secret = testSomeString
	err := linode.build(profile)
	if err != nil {
		t.Fatal(failMsg("linode profile with sample data", "err == nil", "err != nil", err))
	}

	profile.Provider.Key = testEmptyString
	err = linode.build(profile)
	if err == nil {
		t.Fatal(failMsg("no auth key", "err != nil", "err == nil"))
	}

	profile.Provider.Secret = testEmptyString
	err = linode.build(profile)
	if err == nil {
		t.Fatal(failMsg("no auth key or secret", "err != nil", "err == nil"))
	}

	profile.Provider.Key = testSomeString
	err = linode.build(profile)
	if err == nil {
		t.Fatal(failMsg("no auth secret", "err != nil", "err == nil"))
	}

}

func TestProviderLinodeStringMatches(t *testing.T) {
	profile := newIncomingProfile()
	linode := &ProviderLinode{}
	profile.Provider.Key = testSomeString
	profile.Provider.Secret = testSomeString

	for region := range linodeEndpointsMap {
		profile.Linode.Region = camelCaseByChar(region)
		err := linode.build(profile)
		if err != nil {
			t.Log(failMsg("linode region check", "err = nil", "err != nil", err))
		}
	}

	profile.Linode.Region = testInvalidString
	err := linode.build(profile)
	if err == nil {
		t.Log(failMsg("invalid linode region", "err != nil", "err = nil"))
	}
}
