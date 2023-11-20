package logbot

import (
	"github.com/rs/zerolog"
)

const (
	PANIC = zerolog.PanicLevel
	FATAL = zerolog.FatalLevel
	ERROR = zerolog.ErrorLevel
	WARN  = zerolog.WarnLevel
	INFO  = zerolog.InfoLevel
	DEBUG = zerolog.DebugLevel
	TRACE = zerolog.TraceLevel
	BLAST = zerolog.NoLevel
)

type LogBot struct {
	Level       zerolog.Level
	FlagConsole bool
	FlagFile    bool
	Path        string
}
