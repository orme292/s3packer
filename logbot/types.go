package logbot

import (
	"strconv"

	"github.com/rs/zerolog"
)

type LogBot struct {
	Level       zerolog.Level
	FlagConsole bool
	FlagFile    bool
	Path        string
}

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

/*
ParseIntLevel takes an interface{} and returns a zerolog.Level.
If the interface{} is a string, it will attempt to convert it to an int.
When the interface{} is an int, it will attempt to convert it to a zerolog.Level.
*/
func ParseIntLevel(n any) zerolog.Level {
	switch v := n.(type) {
	case string:
		x, err := strconv.Atoi(n.(string))
		if err != nil {
			return BLAST
		}
		n = x
	case bool:
		return DEBUG
	case int:
		n = v
	default:
		return INFO
	}

	switch n {
	case 5:
		return PANIC
	case 4:
		return FATAL
	case 3:
		return ERROR
	case 2:
		return WARN
	case 1:
		return INFO
	case 0:
		return DEBUG
	case -1:
		return TRACE
	default:
		return INFO
	}
}
