package conf

import (
	"fmt"

	"github.com/orme292/s3packer/tuipack"
	"github.com/rs/zerolog"
)

// LogOpts contains the logging configuration, but not an instance of logbot.
type LogOpts struct {
	Level   zerolog.Level
	Screen  bool
	Console bool
	File    bool
	Logfile string
}

func (lo *LogOpts) build(inc *ProfileIncoming) error {

	lo.Level = tuipack.ParseLevel(inc.Logging.Level)
	lo.Screen = inc.Logging.Screen
	lo.Console = inc.Logging.Console
	lo.File = inc.Logging.File
	lo.Logfile = inc.Logging.Logfile

	return lo.validate()

}

func (lo *LogOpts) validate() error {

	if !lo.File && !lo.Console {
		lo.Screen = true
	}

	if lo.File && lo.Logfile == Empty {
		return fmt.Errorf("bad logging config: %s", ErrorLoggingFilepathNotSpecified)
	}
	return nil

}
