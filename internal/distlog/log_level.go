package distlog

import (
	"strconv"

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

func ParseLevel(n any) zerolog.Level {

	var str string
	switch v := n.(type) {
	case int:
		str = strconv.Itoa(v)
	case bool:
		str = "2"
	case string:
		break
	default:
		str = "2"
	}

	lvl, err := zerolog.ParseLevel(str)
	if err != nil {
		return WARN
	}

	if lvl == zerolog.NoLevel {
		return WARN
	}

	return lvl

}
