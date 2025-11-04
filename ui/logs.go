package ui

import (
	"fmt"

	"github.com/Dowdow/gosible/runner"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	STATUS_SUCCESS = 0
	STATUS_ERROR   = 1
	STATUS_PENDING = 0
)

type actionView struct {
	name   string
	status int
}

type logsModel struct {
	runner  *runner.Runner
	ch      chan tea.Msg
	actions []actionView
}

func newLogsModel(config *runner.Config) logsModel {
	return logsModel{
		runner:  runner.NewRunner(config),
		ch:      make(chan tea.Msg),
		actions: make([]actionView, 0),
	}
}

func read(ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		return msg
	}
}

func (m logsModel) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			go m.runner.Run(m.ch)
			return nil
		},
		read(m.ch),
	)
}

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case runner.ActionStartMsg:
		m.actions = append(m.actions, actionView{
			name:   msg.Name,
			status: STATUS_PENDING,
		})
	case runner.ActionEndMsg:
	case runner.StdoutMsg:
		fmt.Println(msg.Msg)
	case runner.StderrMsg:
		fmt.Println(msg.Msg)
	case runner.ErrorMsg:
		close(m.ch)
		return m, func() tea.Msg {
			return errorMsg{err: msg.Error}
		}
	case runner.EndMsg:
		close(m.ch)
		return m, tea.Quit
	}

	return m, read(m.ch)
}

func (m logsModel) View() string {
	var str string
	for _, action := range m.actions {
		str = fmt.Sprintf("%s> %s\n", str, action.name)
	}

	return str
}
