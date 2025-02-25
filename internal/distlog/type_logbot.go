package distlog

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EMPTY   = ""
	SPACE   = " "
	NEWLINE = "\n"
)

type LogOutput struct {
	Console bool
	File    bool
	Screen  bool
}

type LogBot struct {
	Level   zerolog.Level
	Output  *LogOutput
	Logfile string
}

// exitFunc is a function pointer that normally points to os.Exit. We override
// it in tests so we can verify calls to RouteLogMsg that would exit.
var exitFunc = os.Exit

func (lb *LogBot) SetLogLevel(lvl zerolog.Level) {
	zerolog.SetGlobalLevel(lvl)
	lb.Level = lvl
}

func (lb *LogBot) BuildLogger(lvl zerolog.Level) zerolog.Logger {

	var z zerolog.Logger

	if lb.Output.File {

		logFile, err := os.OpenFile(lb.Logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o640)
		if err != nil {
			log.Fatal().Msg("Unable to open log file.")
			exitFunc(1)
		}

		if lb.Output.Console {
			multi := zerolog.MultiLevelWriter(
				zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822},
				logFile)
			z = zerolog.New(multi)
		} else {
			z = zerolog.New(logFile)
		}
	}

	if !lb.Output.File && lb.Output.Console {
		z = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822})
	}

	if lb.Level <= DEBUG {
		return z.Level(DEBUG).With().Timestamp().
			CallerWithSkipFrameCount(4).Logger()
	}

	return z.Level(lvl).With().Timestamp().Logger()

}

func (lb *LogBot) RouteLogMsg(lvl zerolog.Level, msg string) {

	if lb.Output.Console || lb.Output.File {
		z := lb.BuildLogger(lvl)
		z.WithLevel(lvl).Msg(msg)
	}

	if lvl == zerolog.FatalLevel {
		exitFunc(1)
	} else if lvl == zerolog.PanicLevel {
		panic(fmt.Sprintf("unrecoverable error: %s", msg))
	}
}

/*
Blast takes a string and passes it to LogBot.RouteLogMsg with zerolog.NoLevel.
This ensures the message is logged regardless of the global log level.
*/
func (lb *LogBot) Blast(format string, a ...any) {
	lb.RouteLogMsg(BLAST, getMsg(format, a...))
}

func (lb *LogBot) Panic(format string, a ...any) {
	lb.RouteLogMsg(PANIC, getMsg(format, a...))
}

/*
Fatal takes a string and passes it to LogBot.RouteLogMsg with zerolog.FatalLevel.
*/
func (lb *LogBot) Fatal(format string, a ...any) {
	lb.RouteLogMsg(FATAL, getMsg(format, a...))
}

/*
Error takes a string and passes it to LogBot.RouteLogMsg with zerolog.ErrorLevel.
*/
func (lb *LogBot) Error(format string, a ...any) {
	lb.RouteLogMsg(ERROR, getMsg(format, a...))
}

/*
Warn takes a string and passes it to LogBot.RouteLogMsg with zerolog.WarnLevel.
*/
func (lb *LogBot) Warn(format string, a ...any) {
	lb.RouteLogMsg(WARN, getMsg(format, a...))
}

/*
Info takes a string and passes it to LogBot.RouteLogMsg with zerolog.InfoLevel.
*/
func (lb *LogBot) Info(format string, a ...any) {
	lb.RouteLogMsg(INFO, getMsg(format, a...))
}

/*
Debug takes a string and passes it to LogBot.RouteLogMsg with zerolog.DebugLevel.
*/
func (lb *LogBot) Debug(format string, a ...any) {
	lb.RouteLogMsg(DEBUG, getMsg(format, a...))
}

func getMsg(format string, a ...any) string {
	var msg string
	if len(a) == 0 {
		msg = format
	} else {
		msg = fmt.Sprintf(format, a...)
	}
	return msg
}
