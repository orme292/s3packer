package conf

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/orme292/s3packer/tuipack"
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
	Skip     []string

	Tui *tuipack.LogBot
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
			Screen:  false,
			Logfile: "/var/log/s3p.log",
		},
		Tui: &tuipack.LogBot{
			Level:   zerolog.WarnLevel,
			Output:  &tuipack.LogOutput{},
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

	ac.Tui.Level = ac.LogOpts.Level
	ac.Tui.Output.Screen = ac.LogOpts.Screen
	ac.Tui.Output.Console = ac.LogOpts.Console
	ac.Tui.Output.File = ac.LogOpts.File
	ac.Tui.Logfile = ac.LogOpts.Logfile

	if ac.Tui.Output.Screen {
		ac.Tui.Screen = tea.NewProgram(tuipack.NewTuiModel(), tea.WithAltScreen())
	}

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

	return nil

}
