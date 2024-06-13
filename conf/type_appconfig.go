package conf

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
	"github.com/rs/zerolog"
)

type AppConfig struct {
	Provider *Provider
	Opts     *Opts
	Bucket   *Bucket
	Objects  *Objects
	TagOpts  *TagOpts
	Tags     Tags
	LogOpts  *LogOpts
	Paths    []string
	Files    []string
	Dirs     []string

	Log *logbot.LogBot
}

// NewAppConfig returns a new AppConfig object with preconfigured defaults.
func NewAppConfig() *AppConfig {

	return &AppConfig{
		Provider: &Provider{
			Is:     ProviderNameNone,
			AWS:    &ProviderAWS{},
			OCI:    &ProviderOCI{},
			Linode: &ProviderLinode{},
		},
		Opts: &Opts{
			MaxParts:   10,
			MaxUploads: 5,
			Overwrite:  OverwriteNever,
		},
		Bucket: &Bucket{
			Create: false,
		},
		Objects: &Objects{
			NamingType: NamingNone,
		},
		Tags: make(Tags),
		TagOpts: &TagOpts{
			ChecksumSHA256:       true,
			AwsChecksumAlgorithm: types.ChecksumAlgorithmSha256,
			AwsChecksumMode:      types.ChecksumModeEnabled,
			OriginPath:           true,
		},
		LogOpts: &LogOpts{
			Level:    zerolog.ErrorLevel,
			Console:  true,
			File:     false,
			Filepath: "/var/log/s3packer.log",
		},
		Log: &logbot.LogBot{
			Level:       zerolog.ErrorLevel,
			FlagConsole: true,
			FlagFile:    false,
			Path:        "/var/log/s3packer.log",
		},
	}

}

func (ac *AppConfig) ImportFromProfile(inc *ProfileIncoming) error {

	var err error

	err = ac.LogOpts.build(inc)
	if err != nil {
		return err
	}

	ac.Log.Level = ac.LogOpts.Level
	ac.Log.FlagConsole = ac.LogOpts.Console
	ac.Log.FlagFile = ac.LogOpts.File
	ac.Log.Path = ac.LogOpts.Filepath

	err = ac.Provider.build(inc)
	if err != nil {
		return err
	}

	err = ac.Opts.build(inc)
	if err != nil {
		return err
	}

	err = ac.Bucket.build(inc, ac.Provider.Is)
	if err != nil {
		return err
	}

	err = ac.Objects.build(inc)
	if err != nil {
		return err
	}

	err = ac.Tags.build(inc.Tags)
	if err != nil {
		return err
	}

	err = ac.TagOpts.build(inc)
	if err != nil {
		return err
	}

	if len(inc.Files) <= 0 && len(inc.Dirs) <= 0 {
		return fmt.Errorf("bad profile config: %s", ErrorNoFilesSpecified)
	}

	ac.Files = inc.Files
	ac.Dirs = inc.Dirs

	return nil

}
