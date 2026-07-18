package ui

import (
	"fmt"

	"github.com/Dowdow/gosible/config"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type tasksModel struct {
	list  list.Model
	tasks []config.Task
}

func newTasksModel(c *config.Config) tasksModel {
	items := make([]list.Item, 0, len(c.Tasks))
	for _, t := range c.Tasks {
		items = append(items, simpleItem{
			title: t.Name,
			desc:  fmt.Sprintf("%d actions", len(t.Actions)),
		})
	}

	l := newList(items, "Select a task to run")
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		}
	}

	return tasksModel{
		list:  l,
		tasks: c.Tasks,
	}
}

func (m tasksModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m tasksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		if !m.list.SettingFilter() {
			switch msg.String() {
			case "enter":
				index := m.list.GlobalIndex()
				if index < 0 || index >= len(m.tasks) {
					return m, func() tea.Msg {
						return errorMsg{err: fmt.Errorf("invalid task")}
					}
				}

				return m, func() tea.Msg {
					return selectedTask{index: index}
				}
			case "esc":
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tasksModel) View() string {
	return m.list.View()
}
