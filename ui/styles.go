package ui

import "github.com/charmbracelet/lipgloss"

const (
	colorBg       = "#11111B"
	colorOrange   = "#FAB387"
	colorGreen    = "#A6E3A1"
	colorRed      = "#F38BA8"
	colorLavender = "#B4BEFE"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorBg)).
			Background(lipgloss.Color(colorOrange)).
			Padding(0, 1)

	okStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorGreen))

	koStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorRed))

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorLavender))
)
