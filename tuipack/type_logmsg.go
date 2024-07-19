package tuipack

import (
	"fmt"

	"github.com/rs/zerolog"
)

const (
	EMPTY   = ""
	SPACE   = " "
	NEWLINE = "\n"
)

type LogMsg struct {
	Level   zerolog.Level
	LogMsg  string
	ScrMsg  string
	ScrIcon string
}

func NewLogMsg(scrMsg, scrIcon string, level zerolog.Level, logMsg string) *LogMsg {
	return &LogMsg{
		Level:   level,
		LogMsg:  logMsg,
		ScrMsg:  scrMsg,
		ScrIcon: scrIcon,
	}
}

func NewLogMsgB() *LogMsg {
	return &LogMsg{}
}

func (lm *LogMsg) setMessages(msg string) {
	lm.ScrMsg = msg
	lm.LogMsg = msg
}

func (lm *LogMsg) SetMsgUpload(name string) *LogMsg {
	lm.Level = INFO
	lm.LogMsg = fmt.Sprintf("Uploading %s", name)
	lm.ScrMsg = EMPTY
	return lm
}
