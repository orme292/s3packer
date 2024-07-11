package tuipack

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const numListedResults = 30

type TuiResultMsg struct {
	IsSuccessful bool
	Msg          string
	HeaderMsg    string
}

type TuiModel struct {
	spinner    spinner.Model
	Header     string
	results    []TuiResultMsg
	isQuitting bool
}

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("15")).
			Foreground(lipgloss.Color("21")).Bold(true)
	greenCheckStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	pendingCheckStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	appStyle          = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

func (r TuiResultMsg) String() string {

	if len(r.Msg) == 0 {
		return ""
	}

	var s string
	if r.IsSuccessful {
		s = greenCheckStyle.Render("✓")
	} else {
		s = pendingCheckStyle.Render(":")
	}

	return fmt.Sprintf("%s %s", s, r.Msg)

}

func NewTuiModel() TuiModel {

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = spinnerStyle

	return TuiModel{
		spinner: s,
		Header:  "Starting Up...",
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

	case TuiResultMsg:
		if msg.HeaderMsg != "" {
			m.Header = msg.HeaderMsg
		}
		m.results = append(m.results[1:], msg)
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
	s = titleStyle.Render(" s3packer ")
	s += "  A profile-based S3 backup application\nhttps://github.com/orme292/s3packer"

	s += "\n\n"

	for _, res := range m.results {
		if res.Msg == "" {
			s += ""
		}
		s += res.String() + "\n"
	}

	s += "\n\n"

	if !m.isQuitting {
		s += m.spinner.View() + "  " + m.Header
	} else {
		s += ""
	}

	s += "\n\n"

	if m.isQuitting {
		s += "\n"
	} else {
		s += helpStyle.Render("Press CTRL+C or Q to exit early")
	}

	return appStyle.Render(s)

}
