package conf

import (
	"log"
	"strconv"
	"testing"
)

func TestNewAppConfig(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	err := app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("app config import with sample profile", "err == nil", "err != nil", err))
	}

	profile.Logging.File = true
	profile.Logging.Logfile = testEmptyString
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("file logging on with no filename", "err != nil", "err == nil"))
	}
	profile.Logging.File = false

	profile.Files = []string{}
	profile.Dirs = []string{}
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with no files or dirs", "err != nil", "err == nil"))
	}
}

func TestAppConfigAws(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	// build AWS profile
	profile.Provider.Use = ProviderNameAWS.String()
	profile.Provider.Profile = testSomeString
	profile.Provider.Key = testEmptyString
	profile.Provider.Secret = testEmptyString
	profile.AWS.ACL = AwsACLPrivate
	profile.AWS.Storage = AwsClassStandard

	err := app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with provider aws", "err == nil", "err != nil", err))
	}
}

func TestAppConfigLinode(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	// build Linode profile
	profile.Provider.Use = ProviderNameLinode.String()
	profile.Provider.Profile = testEmptyString
	profile.Provider.Key = testSomeString
	profile.Provider.Secret = testSomeString
	profile.Linode.Region = LinodeRegionNewark

	err := app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with provider linode", "err == nil", "err != nil", err))
	}
}

func TestAppConfigGoogle(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	// build Google profile
	profile.Provider.Use = ProviderNameGoogle.String()
	profile.Provider.Profile = testEmptyString
	profile.Provider.Key = testEmptyString
	profile.Provider.Secret = testEmptyString
	profile.Google.Project = testSomeString
	profile.Google.Storage = GCStorageStandard
	profile.Google.ObjectACL = GCObjectACLPrivate
	profile.Google.BucketACL = GCBucketACLPrivate

	err := app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with provider google", "err = nil", "err != nil", err))
	}
}

func TestAppConfigOCI(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	// build OCI profile
	profile.Provider.Use = ProviderNameOCI.String()
	profile.Provider.Profile = testSomeString
	profile.OCI.Compartment = testSomeString
	profile.OCI.Storage = OracleStorageTierStandard

	err := app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with provider oci", "err = nil", "err != nil", err))
	}
}

func TestAppConfigNone(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	// build empty profile
	profile.Provider.Use = ProviderNameNone.String()
	err := app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with provider none", "err != nil", "err = nil", err))
	}

	profile.Provider.Use = testSomeString
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with invalid provider", "err != nil", "err = nil", err))
	}
}

func TestAppConfigOptions(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	profile.Options.OverwriteObjects = testEmptyString
	err := app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with empty overwrite option", "err != nil", "err = nil", err))
	}

	profile.Options.OverwriteObjects = testInvalidString
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with invalid overwrite option", "err != nil", "err = nil", err))
	}

	profile.Options.OverwriteObjects = OverwriteNever.String()
	profile.Options.MaxUploads = -1
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with invalid max uploads", "err != nil", "err = nil", err))
	}

	profile.Options.OverwriteObjects = "true"
	profile.Options.MaxUploads = 1
	err = app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with alternate overwrite option", "err = nil", "err != nil", err))
	}

}

func TestAppConfigBucket(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	profile.Bucket.Name = testEmptyString
	err := app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with empty bucket name", "err != nil", "err = nil", err))
	}

	profile.Provider.Use = ProviderNameAWS.String()
	profile.Bucket.Name = testSomeString
	profile.Bucket.Region = testEmptyString
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with empty bucket region", "err != nil", "err = nil", err))
	}

	profile.Bucket.Region = testSomeString
	err = app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with valid config", "err = nil", "err != nil", err))
	}
}

func TestAppConfigObjects(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	profile.Objects.NamingType = NamingAbsolute.String()
	err := app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with absolute naming", "err = nil", "err != nil", err))
	}

	profile.Objects.NamingType = NamingRelative.String()
	err = app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with relative naming", "err = nil", "err != nil", err))
	}

	profile.Objects.NamingType = testInvalidString
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with invalid naming type", "err != nil", "err = nil", err))
	}
}

func TestAppConfigTags(t *testing.T) {
	app := NewAppConfig()
	profile := newIncomingProfile()

	profile.Tags = make(map[string]string)
	app.Tags = make(map[string]string)
	profile.Tags = map[string]string{
		"s3p-checksum-sha256": "true",
	}
	err := app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with reserved checksum tag", "err = nil", "err != nil", err))
	}

	profile.Tags = make(map[string]string)
	app.Tags = make(map[string]string)
	profile.Tags = map[string]string{
		"s3p-origin-path": "true",
	}
	err = app.ImportFromProfile(profile)
	if err == nil {
		t.Fatal(failMsg("import profile with reserved origin path tag", "err = nil", "err != nil", err))
	}

	profile.Tags = make(map[string]string)
	app.Tags = make(map[string]string)
	profile.Tags = map[string]string{
		"good-tag":    "good-value",
		"another-tag": "another-value",
	}
	err = app.ImportFromProfile(profile)
	if err != nil {
		t.Fatal(failMsg("import profile with valid tags", "err = nil", "err != nil", err))
	}

	tags := app.Tags.Get()
	if len(tags) != 2 {
		t.Fatal(failMsg("get appconfig tags length", "len(tags) = 2", strconv.Itoa(len(tags))))
	}
}

func TestAppConfigTagOpts(t *testing.T) {
	log.Println("TagOpts has no tests")
}
