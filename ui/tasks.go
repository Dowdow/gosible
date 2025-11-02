package ui

import (
	"fmt"

	"github.com/Dowdow/gosible/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type taskItem struct {
	index int
	title string
	desc  string
}

func (i taskItem) Title() string       { return i.title }
func (i taskItem) Description() string { return i.desc }
func (i taskItem) FilterValue() string { return i.title }

type tasksModel struct {
	list list.Model
}

func newTasksModel(tasks []config.Task) tasksModel {
	items := make([]list.Item, 0)
	for index, t := range tasks {
		items = append(items, taskItem{
			index: index,
			title: t.Name,
			desc:  fmt.Sprintf("%d actions", len(t.Actions)),
		})
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Select a task to run"

	return tasksModel{
		list: list,
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
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case tea.KeyEnter.String():
			item, ok := m.list.SelectedItem().(taskItem)
			if ok {
				return m, func() tea.Msg {
					return selectedTask{index: item.index}
				}
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tasksModel) View() string {
	return m.list.View()
}
