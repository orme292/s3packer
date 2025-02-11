package conf

import (
	"fmt"

	"github.com/rs/zerolog"
	"s3p/internal/distlog"
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
	Skip     []string

	Log *distlog.LogBot
}

// NewAppConfig returns a new AppConfig object with preconfigured defaults.
func NewAppConfig() *AppConfig {

	return &AppConfig{
		Provider: &Provider{
			Is:     ProviderNameNone,
			AWS:    &ProviderAWS{},
			Google: &ProviderGoogle{},
			Linode: &ProviderLinode{},
			OCI:    &ProviderOCI{},
		},
		Opts: &Opts{
			MaxUploads:     1,
			FollowSymlinks: false,
			WalkDirs:       true,
			Overwrite:      OverwriteNever,
		},
		Bucket: &Bucket{
			Create: false,
		},
		Objects: &Objects{
			NamingType:  NamingNone,
			OmitRootDir: true,
		},
		Tags: make(Tags),
		TagOpts: &TagOpts{
			ChecksumSHA256: false,
			OriginPath:     false,
		},
		LogOpts: &LogOpts{
			Level:   zerolog.WarnLevel,
			Console: true,
			File:    false,
			Logfile: "/var/log/s3p.log",
		},
		Log: &distlog.LogBot{
			Level:   zerolog.WarnLevel,
			Output:  &distlog.LogOutput{},
			Logfile: "/var/log/s3p.log",
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
	ac.Log.Output.Console = ac.LogOpts.Console
	ac.Log.Output.File = ac.LogOpts.File
	ac.Log.Logfile = ac.LogOpts.Logfile

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

	if len(inc.Files) == 0 && len(inc.Dirs) == 0 {
		return fmt.Errorf("bad profile config: %s", ErrorNoFilesSpecified)
	}

	ac.Files = inc.Files
	ac.Dirs = inc.Dirs
	ac.Skip = inc.Skip

	ac.setGoogleExceptions()

	return nil

}

func (ac *AppConfig) setGoogleExceptions() {

	if ac.Provider.Is == ProviderNameGoogle {

		fmt.Println("Google Cloud Storage support is experimental")
		ac.Log.Warn("Google Cloud Storage support is experimental")

		if ac.Opts.MaxUploads > 1 {
			ac.Opts.MaxUploads = 1
			ac.Log.Warn("s3packer doesn't support parallel uploads with Google Cloud Storage")
		}

		if ac.Bucket.Create == true {
			if ac.Provider.Google.Project == Empty {
				ac.Log.Fatal("You have bucket creation enabled, but no project specified.")
			}
		}

	}

}
