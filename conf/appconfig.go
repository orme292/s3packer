package conf

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
)

func New(file string) (a *AppConfig, err error) {
	r := &readConfig{}
	err = r.load(file)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n\n", r)

	_, err = r.versionOK()
	if err != nil {
		return nil, err
	}

	a = &AppConfig{
		Log: &logbot.LogBot{
			Level:       logbot.WARN,
			FlagConsole: true,
			FlagFile:    false,
		},
	}
	a.init()
	err = a.apply(r)
	return a, err
}

/*
init() sets defaults
*/
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
	a.Tag = &Tag{
		ChecksumSHA256: true,
		Origins:        true,
	}
}

/*
apply() reads the values from the readConfig object into the AppConfig object.
*/
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
	// Done: Objects
	a.Objects, err = r.getObjects()
	if err != nil {
		return
	}
	// Done: Options
	a.Opts, err = r.getOpts()
	if err != nil {
		return
	}
	// Done: Tagging
	a.Tag, err = r.getTag()
	if err != nil {
		return
	}
	// Done: Tags
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
