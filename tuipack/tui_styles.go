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
	StyleBgGrnFgBlk = gloss.NewStyle().
			Background(gloss.Color("10")).
			Foreground(gloss.Color("0"))
	StyleDefault = gloss.NewStyle().
			Foreground(gloss.Color("243"))
	StyleFgGreen = gloss.NewStyle().
			Foreground(gloss.Color("10"))
	StyleFgDrkRed = gloss.NewStyle().
			Foreground(gloss.Color("52"))
	StyleFgDrkYellow = gloss.NewStyle().
				Foreground(gloss.Color("184"))
	StyleFgRed = gloss.NewStyle().
			Foreground(gloss.Color("160"))
	StyleFgYellow = gloss.NewStyle().
			Foreground(gloss.Color("11"))

	StyleHelpMessage = gloss.NewStyle().
				Foreground(gloss.Color("241"))
	StyleSpinner = gloss.NewStyle().
			Foreground(gloss.Color("63"))
)

// Characters with Styles
var (
	ScrnNone           = EMPTY
	ScrnAppHeader      = StyleAppHeader.Render(" s3packer ")
	ScrnHelpMessage    = StyleHelpMessage.Render("Press CTRL+C or Q to exit early")
	ScrnLfCheck        = StyleFgGreen.Render("✓")
	ScrnLfDefault      = StyleDefault.Render(":")
	ScrnLfFailed       = StyleFgRed.Render("x")
	ScrnLfSkip         = StyleFgDrkYellow.Render("*")
	ScrnLfOperOK       = StyleFgYellow.Render("⇅")
	ScrnLfOperFailed   = StyleFgDrkRed.Render("*")
	ScrnLfUpload       = StyleBgGrnFgBlk.Render("✓")
	ScrnLfUploadFailed = StyleFgDrkRed.Render("⇋")
)
