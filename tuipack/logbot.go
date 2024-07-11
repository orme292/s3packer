package tuipack

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	Screen  *tea.Program
}

type ScreenMsg struct {
	Msg  string
	Mark bool
}

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
			os.Exit(1)
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

func (lb *LogBot) route(lvl zerolog.Level, msg string) {

	if lb.Output.Screen && lb.Output.Console {
		lb.Output.Console = false
	}

	if lb.Output.Console || lb.Output.File {
		z := lb.BuildLogger(lvl)
		z.WithLevel(lvl).Msg(msg)
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

func (lb *LogBot) SendOutput(screen ScreenMsg, log string, lvl zerolog.Level, toLogger bool, toScreen bool) {

	if toScreen && lb.Screen != nil && lb.Output.Screen {
		lb.ToScreen(screen.Msg, screen.Mark)
	}
	if (lb.Output.File || lb.Output.Console) && toLogger {
		lb.route(lvl, log)
	}

}

func (lb *LogBot) ScreenQuit() {

	if lb.Output.Screen && lb.Screen != nil {
		lb.ToScreenHeader("s3packer exited.")
		lb.Screen.Quit()
	}

}

func (lb *LogBot) ToScreen(msg string, success bool) {

	if lb.Output.Screen && lb.Screen != nil {
		lb.Screen.Send(TuiResultMsg{
			IsSuccessful: success,
			Msg:          msg,
		})
	}

}

func (lb *LogBot) ResetHeader() {

	if lb.Output.Screen && lb.Screen != nil {
		lb.Screen.Send(TuiResultMsg{
			HeaderMsg: "Running...",
		})
	}

}

func (lb *LogBot) ToScreenHeader(header string) {

	if lb.Output.Screen && lb.Screen != nil {
		lb.Screen.Send(TuiResultMsg{
			HeaderMsg: header,
		})
	}

}
