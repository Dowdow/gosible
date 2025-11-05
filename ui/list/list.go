package list

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	ellipsis          = "…"
	titleStype        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#11111b")).Background(lipgloss.Color("#fab387")).Padding(0, 1)
	itemTitleStyle    = lipgloss.NewStyle().Bold(true)
	itemSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#fab387"))
)

type Item struct {
	Title string
	Desc  string
}

type Model struct {
	Title   string
	items   []Item
	current int
	width   int
	height  int
}

func New(items []Item) Model {
	return Model{
		items:   items,
		current: 0,
		width:   0,
		height:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.current > 0 {
				m.current--
			}
		case tea.KeyDown:
			if m.current < len(m.items)-1 {
				m.current++
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	str := fmt.Sprintf("%s\n\n", titleStype.Render(ansi.Truncate(m.Title, m.width-3, ellipsis)))

	for index, item := range m.items {
		title := ansi.Truncate(item.Title, m.width-3, ellipsis)
		desc := ansi.Truncate(item.Desc, m.width-3, ellipsis)
		if m.current == index {
			str += fmt.Sprintf("> %s\n  %s\n\n", itemSelectedStyle.Render(itemTitleStyle.Render(title)), itemSelectedStyle.Render(desc))
		} else {
			str += fmt.Sprintf("  %s\n  %s\n\n", itemTitleStyle.Render(title), desc)
		}
	}

	return str
}

func (m Model) SelectedIndex() int {
	return m.current
}
