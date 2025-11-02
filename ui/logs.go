package ui

import (
	"github.com/Dowdow/gosible/runner"
	tea "github.com/charmbracelet/bubbletea"
)

type logsModel struct {
	config *runner.Config
}

func newLogsModel(config *runner.Config) logsModel {
	return logsModel{
		config: config,
	}
}

func (m logsModel) Init() tea.Cmd {
	// runner.Run(*m.config)
	return nil
}

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m logsModel) View() string {
	return ""
}
