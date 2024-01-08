package conf

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
)

// NewAppConfig will build a new AppConfig object out of the specified yaml file.
// It creates a readConfig object and called the load() method to open, read, and unmarshal the yaml file.
// There is also a versionOK() method that will eventually be fleshed out to make sure difference profile
// versions are unmarshalled correctly.
func NewAppConfig(file string) (a *AppConfig, err error) {
	r := &readConfig{}
	err = r.load(file)
	if err != nil {
		return nil, err
	}

	_, err = r.versionOK()
	if err != nil {
		return nil, err
	}

	a = &AppConfig{
		Log: &logbot.LogBot{
			Level:       logbot.ERROR,
			FlagConsole: true,
			FlagFile:    false,
		},
	}
	a.init()
	err = a.apply(r)
	return a, err
}

// init() sets the default values for the AppConfig object. It should be one of the first methods called after
// creating a new AppConfig object, but before the readConfig object values are transferred, or before apply() is
// called.
func (a *AppConfig) init() {
	a.Bucket = &Bucket{Create: false}
	a.Objects = &Objects{Naming: NamingAbsolute}
	a.Provider = &Provider{
		AwsACL:     types.ObjectCannedACLPrivate,
		AwsStorage: types.StorageClassStandard,
	}
	a.Opts = &Opts{
		MaxUploads: 5,
		MaxParts:   1,
		Overwrite:  OverwriteNever,
	}
	a.Tag = &TagOpts{
		ChecksumSHA256: true,
		Origins:        true,
	}
}

// apply will transfer the values from the readConfig object to the AppConfig object. It should be called after
// init() and after the readConfig object has been loaded and validated.
// Each struct of the AppConfig object is handled separately, and each should have its own readConfig method
// to handle validation of the values and transfer.
func (a *AppConfig) apply(r *readConfig) (err error) {
	a.LogOpts, err = r.getLogging()
	a.Log = &logbot.LogBot{
		Level:       a.LogOpts.Level,
		FlagConsole: a.LogOpts.Console,
		FlagFile:    a.LogOpts.File,
		Path:        a.LogOpts.Filepath,
	}
	if err != nil {
		return
	}
	a.Provider, err = r.getProvider()
	if err != nil {
		return
	}
	a.Bucket, err = r.getBucket()
	if err != nil {
		return
	}
	a.Objects, err = r.getObjects()
	if err != nil {
		return
	}
	a.Opts, err = r.getOpts()
	if err != nil {
		return
	}
	a.Tag, err = r.getTag()
	if err != nil {
		return
	}
	a.Tags, err = r.getValidTags()
	if err != nil {
		return
	}
	a.Files, a.Directories, err = r.getTargets()
	if err != nil {
		return
	}

	return
}
