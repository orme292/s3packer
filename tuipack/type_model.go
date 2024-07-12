package tuipack

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	numListedResults = 23
	numErrorResults  = 7
)

type TuiModel struct {
	spinner    spinner.Model
	Header     string
	errors     []TuiResultMsg
	results    []TuiResultMsg
	isQuitting bool
}

var (
	appStyle = lipgloss.NewStyle().Margin(3, 3)
)

func NewTuiModel() TuiModel {

	s := spinner.New()
	s.Spinner = spinner.Meter
	s.Style = StyleSpinner

	return TuiModel{
		spinner: s,
		Header:  "Starting Up...",
		errors:  make([]TuiResultMsg, numErrorResults),
		results: make([]TuiResultMsg, numListedResults),
	}

}

func (m TuiModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.isQuitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case tea.QuitMsg:
		m.isQuitting = true
		return m, tea.Quit

	case TuiResultMsg:
		if msg.HeaderMsg != EMPTY {
			m.Header = msg.HeaderMsg
		}
		if msg.Icon == ScrnLfFailed || msg.Icon == ScrnLfUpload {
			m.errors = append(m.errors[1:], msg)
		} else {
			m.results = append(m.results[1:], msg)
		}

		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		return m, nil
	}

}

func (m TuiModel) View() string {

	var s string
	s = ScrnAppHeader
	s += StyleAppHeaderDesc.Render("  A profile-based S3 backup application  ")
	s += StyleAppHeaderNoLink.Render("  ")
	s += StyleAppHeaderLink.Render("https://github.com/orme292/s3packer")
	s += StyleAppHeaderNoLink.Render("  ")

	s += NEWLINE + NEWLINE

	for _, err := range m.errors {
		if err.Msg == EMPTY {
			s += StyleHelpMessage.Render("..........")
		}
		s += err.String() + NEWLINE
	}

	s += NEWLINE

	for _, res := range m.results {
		if res.Msg == EMPTY {
			s += EMPTY
		}
		s += res.String() + NEWLINE
	}

	s += "\n\n"

	if !m.isQuitting {
		s += m.spinner.View() + SPACE + SPACE + m.Header
	} else {
		s += EMPTY
	}

	s += NEWLINE + NEWLINE

	if m.isQuitting {
		s += NEWLINE
	} else {
		s += ScrnHelpMessage
	}

	return appStyle.Render(s)

}
