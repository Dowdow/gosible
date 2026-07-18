package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type simpleItem struct {
	title string
	desc  string
}

func (i simpleItem) Title() string       { return i.title }
func (i simpleItem) Description() string { return i.desc }
func (i simpleItem) FilterValue() string { return i.title }

func newList(items []list.Item, title string) list.Model {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color(colorOrange)).
		BorderForeground(lipgloss.Color(colorOrange))
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color(colorOrange)).
		BorderForeground(lipgloss.Color(colorOrange))

	l := list.New(items, d, 0, 0)
	l.Title = title
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)

	filterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorOrange))
	l.FilterInput.PromptStyle = filterStyle
	l.FilterInput.TextStyle = filterStyle
	l.FilterInput.Cursor.Style = filterStyle

	return l
}
