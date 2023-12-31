package logbot

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
SetLogLevel takes a zerolog.Level and sets the global log level.
*/
func (lb *LogBot) SetLogLevel(l zerolog.Level) {
	zerolog.SetGlobalLevel(l)
}

/*
buildZ takes a zerolog.Level and returns a zerolog.Logger "object".
If LogBot.FlagFile is true, it will attempt to open or create the file specified by LogBot.Path.
A zerolog.MultiLevelWriter is created if both LogBot.FlagFile and LogBot.FlagConsole are true.

If LogBot.Level is DEBUG, the zerolog.Logger returned will have a caller field set. Otherwise, no caller will be logged.

Any zerolog.Logger returned will have a timestamp field set.
*/
func (lb *LogBot) buildZ(l zerolog.Level) (z zerolog.Logger) {
	if lb.FlagFile {
		logFile, err := os.OpenFile(lb.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o640)
		if err != nil {
			log.Fatal().Msg("Unable to open log file: " + err.Error())
			os.Exit(1)
		}
		if lb.FlagConsole {
			multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822}, logFile)
			z = zerolog.New(multi)
		} else {
			z = zerolog.New(logFile)
		}
	}
	if !lb.FlagFile && lb.FlagConsole {
		z = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822})
	}
	if lb.Level <= DEBUG {
		return z.Level(DEBUG).With().Timestamp().CallerWithSkipFrameCount(4).Logger()
	}
	return z.Level(l).With().Timestamp().Logger()
}

/*
route takes a zerolog.Level and a string.
It will build a zerolog.Logger "object" with the given zerolog.Level.
It will then log the string with the given zerolog.Level.
If the zerolog.Level is FATAL, the program will exit with status code 1.
*/
func (lb *LogBot) route(l zerolog.Level, s string) {
	if lb.FlagConsole || lb.FlagFile {
		z := lb.buildZ(lb.Level)
		z.WithLevel(l).Msg(s)
	}
	if l == FATAL || l == PANIC {
		os.Exit(1)
	}
}

/*
Blast takes a string and passes it to LogBot.route with zerolog.NoLevel.
This ensures the message is logged regardless of the global log level.
*/
func (lb *LogBot) Blast(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	lb.route(BLAST, msg)
}

func (lb *LogBot) Panic(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	lb.route(PANIC, msg)
}

/*
Fatal takes a string and passes it to LogBot.route with zerolog.FatalLevel.
*/
func (lb *LogBot) Fatal(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	lb.route(FATAL, msg)
}

/*
Error takes a string and passes it to LogBot.route with zerolog.ErrorLevel.
*/
func (lb *LogBot) Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	lb.route(ERROR, msg)
}

/*
Warn takes a string and passes it to LogBot.route with zerolog.WarnLevel.
*/
func (lb *LogBot) Warn(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	lb.route(WARN, msg)
}

/*
Info takes a string and passes it to LogBot.route with zerolog.InfoLevel.
*/
func (lb *LogBot) Info(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	lb.route(INFO, msg)
}

/*
Debug takes a string and passes it to LogBot.route with zerolog.DebugLevel.
*/
func (lb *LogBot) Debug(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	lb.route(DEBUG, msg)
}
