package conf

import (
	"github.com/orme292/s3packer/logbot"
)

// NewAppConfig will build a new AppConfig object out of the specified yaml file.
// It creates a readConfig object and called the loadProfile() method to open, read, and unmarshal the yaml file.
// There is also a validateVersion() method that will eventually be fleshed out to make sure difference profile
// versions are unmarshalled correctly.
func NewAppConfig(file string) (ac *AppConfig, err error) {
	rc := &readConfig{}
	err = rc.loadProfile(file)
	if err != nil {
		return nil, err
	}

	_, err = rc.validateVersion()
	if err != nil {
		return nil, err
	}

	ac = &AppConfig{}
	ac.setDefaults()
	err = ac.transpose(rc)
	return
}

// setDefaults() sets the default values for the AppConfig object.
//
// It should be one of the first methods called after
// creating a new AppConfig object, but before the readConfig object values are transferred, or before apply() is
// called.
func (ac *AppConfig) setDefaults() {
	ac.Bucket = &Bucket{Create: false}
	ac.Objects = &Objects{Naming: NamingAbsolute}
	ac.Opts = &Opts{
		MaxUploads: 5,
		MaxParts:   1,
		Overwrite:  OverwriteNever,
	}
	ac.Tag = &TagOpts{
		ChecksumSHA256: true,
		Origins:        true,
	}
	ac.Log = &logbot.LogBot{
		Level:       logbot.ERROR,
		FlagConsole: true,
		FlagFile:    false,
	}
}

// apply will transfer the values from the readConfig object to the AppConfig object. It should be called after
// init() and after the readConfig object has been loaded and validated.
// Each struct of the AppConfig object is handled separately, and each should have its own readConfig method
// to handle validation of the values and transfer.
func (ac *AppConfig) transpose(r *readConfig) (err error) {
	ac.LogOpts, err = r.transposeStructLogging()
	ac.Log = &logbot.LogBot{
		Level:       ac.LogOpts.Level,
		FlagConsole: ac.LogOpts.Console,
		FlagFile:    ac.LogOpts.File,
		Path:        ac.LogOpts.Filepath,
	}
	if err != nil {
		return
	}
	ac.Provider, err = r.transposeStructProvider()
	if err != nil {
		return
	}
	ac.Bucket, err = r.transposeStructBucket()
	if err != nil {
		return
	}
	ac.Objects, err = r.transposeStructObjects()
	if err != nil {
		return
	}
	ac.Opts, err = r.transposeStructOpts()
	if err != nil {
		return
	}
	ac.Tag, err = r.transposeStructTagOpts()
	if err != nil {
		return
	}
	ac.Tags, err = r.transposeStructTags()
	if err != nil {
		return
	}
	ac.Files, ac.Directories, err = r.transposeStructFileTargets()
	if err != nil {
		return
	}

	return
}
