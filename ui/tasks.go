package ui

import (
	"fmt"

	"github.com/Dowdow/gosible/config"
	"github.com/Dowdow/gosible/ui/list"
	tea "github.com/charmbracelet/bubbletea"
)

type tasksModel struct {
	list list.Model
}

func newTasksModel(tasks []config.Task) tasksModel {
	items := make([]list.Item, 0)
	for _, t := range tasks {
		items = append(items, list.Item{
			Title: t.Name,
			Desc:  fmt.Sprintf("%d actions", len(t.Actions)),
		})
	}

	list := list.New(items)
	list.Title = "Select a task to run"

	return tasksModel{
		list: list,
	}
}

func (m tasksModel) Init() tea.Cmd {
	return m.list.Init()
}

func (m tasksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, func() tea.Msg {
				return selectedTask{index: m.list.SelectedIndex()}
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
