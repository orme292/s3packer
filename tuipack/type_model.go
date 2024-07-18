package tuipack

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	numDefaultListedResults = 125
	numDefaultErrorResults  = 7
)

type TuiModel struct {
	spinner    spinner.Model
	Header     string
	errors     []TuiResultMsg
	results    []TuiResultMsg
	isQuitting bool
	Width      int
	Height     int
}

type TuiQuit struct{}

var (
	appStyle = lipgloss.NewStyle().Margin(3, 3)
)

func NewTuiModel() *TuiModel {

	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = StyleSpinner

	return &TuiModel{
		spinner: s,
		Header:  "Starting Up...",
		errors:  make([]TuiResultMsg, numDefaultErrorResults),
		results: make([]TuiResultMsg, numDefaultListedResults),
	}

}

func (m *TuiModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width - 6
		m.Height = msg.Height - 19

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.isQuitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case tea.QuitMsg, TuiQuit:
		m.isQuitting = true
		return m, tea.Quit

	case TuiResultMsg:
		if msg.HeaderMsg != EMPTY {
			m.Header = msg.HeaderMsg
		}

		switch msg.Icon {
		case ScrnLfFailed, ScrnLfUploadFailed, ScrnLfOperFailed:
			m.errors = append(m.errors[1:], msg)

		default:
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

	return m, nil

}

func (m *TuiModel) View() string {

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
		if len(err.String()) > m.Width {
			s += err.String()[:m.Width] + NEWLINE
		} else {
			s += err.String() + NEWLINE
		}

	}

	s += NEWLINE

	for i := 0; i < m.Height; i++ {
		if len(m.results[i].String()) > m.Width {
			s += m.results[i].String()[:m.Width] + NEWLINE
		} else {
			s += m.results[i].String() + NEWLINE
		}
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

	if m.isQuitting {
		return lipgloss.NewStyle().UnsetMargins().Render(s)
	}
	return appStyle.Render(s)

}
