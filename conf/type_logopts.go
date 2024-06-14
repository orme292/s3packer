package conf

import (
	"fmt"

	"github.com/orme292/s3packer/logbot"
	"github.com/rs/zerolog"
)

// LogOpts contains the logging configuration, but not an instance of logbot.
type LogOpts struct {
	Level    zerolog.Level
	Console  bool
	File     bool
	Filepath string
}

func (lo *LogOpts) build(inc *ProfileIncoming) error {

	lo.Level = logbot.ParseIntLevel(inc.Logging.Level)
	lo.Console = inc.Logging.OutputToConsole
	lo.File = inc.Logging.OutputToFile
	lo.Filepath = inc.Logging.Path

	return lo.validate()

}

func (lo *LogOpts) validate() error {

	if lo.File && lo.Filepath == Empty {
		return fmt.Errorf("bad logging config: %s", ErrorLoggingFilepathNotSpecified)
	}
	return nil

}
