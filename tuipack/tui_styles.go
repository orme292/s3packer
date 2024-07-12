package tuipack

import (
	gloss "github.com/charmbracelet/lipgloss"
)

// Styles
var (
	StyleAppHeader = gloss.NewStyle().
			Background(gloss.Color("15")).
			Foreground(gloss.Color("21")).
			Bold(true)
	StyleAppHeaderDesc = gloss.NewStyle().
				Background(gloss.Color("102")).
				Foreground(gloss.Color("50"))
	StyleAppHeaderLink = gloss.NewStyle().
				Background(gloss.Color("7")).
				Foreground(gloss.Color("27")).
				Underline(true)
	StyleAppHeaderNoLink = gloss.NewStyle().
				Background(gloss.Color("7")).
				Foreground(gloss.Color("27"))
	StyleFgGreen = gloss.NewStyle().
			Foreground(gloss.Color("10"))
	StyleFgRed = gloss.NewStyle().
			Foreground(gloss.Color("160"))
	StyleDefault = gloss.NewStyle().
			Foreground(gloss.Color("243"))
	StyleHelpMessage = gloss.NewStyle().
				Foreground(gloss.Color("241"))
	StyleSpinner = gloss.NewStyle().
			Foreground(gloss.Color("63"))
	StyleUpload = gloss.NewStyle().
			Foreground(gloss.Color("11")).
			Background(gloss.Color("0"))
)

// Characters with Styles
var (
	ScrnNone        = EMPTY
	ScrnAppHeader   = StyleAppHeader.Render(" s3packer ")
	ScrnHelpMessage = StyleHelpMessage.Render("Press CTRL+C or Q to exit early")
	ScrnLfCheck     = StyleFgGreen.Render("✓")
	ScrnLfDefault   = StyleDefault.Render(":")
	ScrnLfFailed    = StyleFgRed.Render("x")
	ScrnLfUpload    = StyleUpload.Render("▲")
)
