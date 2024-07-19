package tuipack

import (
	"fmt"
)

type TuiResultMsg struct {
	Icon      string
	Msg       string
	HeaderMsg string
}

func (r TuiResultMsg) String() string {

	if r.Msg == "" {
		return r.Msg
	}

	var s string

	if r.Icon == EMPTY {
		r.Icon = ScrnLfDefault
	}
	s = r.Icon

	return fmt.Sprintf("%s %s", s, r.Msg)

}
