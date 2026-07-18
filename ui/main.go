package ui

import (
	"fmt"

	"github.com/Dowdow/gosible/config"
	tea "github.com/charmbracelet/bubbletea"
)

type errorMsg struct {
	err error
}

type parsedConfigMsg struct {
	config *config.Config
	err    error
}

type selectedTask struct {
	index int
}

type selectedMachineUser struct {
	machine string
	user    string
}

type taskDoneMsg struct{}

type backToTasksMsg struct{}

type mainModel struct {
	config       *config.Config
	taskIndex    int
	machine      string
	user         string
	currentModel tea.Model
}

func quitWithError(err error) tea.Cmd {
	return tea.Sequence(
		tea.ExitAltScreen,
		tea.Println(PrintError(err)),
		tea.Quit,
	)
}

func NewMainModel() mainModel {
	return mainModel{
		currentModel: newLoadingModel(),
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.currentModel.Init(),
		func() tea.Msg {
			c, err := config.ParseConfig()
			if err != nil {
				return parsedConfigMsg{err: err}
			}

			return parsedConfigMsg{config: c}
		},
	)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case parsedConfigMsg:
		if msg.err != nil {
			return m, quitWithError(msg.err)
		}
		m.config = msg.config
		m.currentModel = newTasksModel(m.config)
		return m, m.currentModel.Init()
	case selectedTask:
		m.taskIndex = msg.index
		model, err := newMachinesModel(m.config, m.taskIndex)
		if err != nil {
			return m, quitWithError(err)
		}
		m.currentModel = model
		return m, m.currentModel.Init()
	case selectedMachineUser:
		m.machine = msg.machine
		m.user = msg.user
		runnerConfig, err := m.config.Convert(m.taskIndex, m.machine, m.user)
		if err != nil {
			return m, quitWithError(err)
		}
		m.currentModel = newLogsModel(runnerConfig)
		return m, m.currentModel.Init()
	case taskDoneMsg:
		m.currentModel = newTasksModel(m.config)
		return m, m.currentModel.Init()
	case backToTasksMsg:
		m.currentModel = newTasksModel(m.config)
		return m, m.currentModel.Init()
	case errorMsg:
		return m, quitWithError(msg.err)
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	return m.currentModel.View()
}

func PrintError(err error) string {
	return fmt.Sprintf("%s %s", koStyle.Render("ERROR"), err.Error())
}
