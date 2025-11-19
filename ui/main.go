package ui

import (
	"fmt"

	"github.com/Dowdow/gosible/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type mainModel struct {
	config       *config.Config
	taskIndex    int
	machine      string
	user         string
	currentModel tea.Model
	err          error
}

func NewMainModel() mainModel {
	return mainModel{
		currentModel: newLoadingModel(),
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
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
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case parsedConfigMsg:
		if msg.err != nil {
			// Sequence quitting because of the spinner
			return m, tea.Sequence(
				tea.Println(PrintError(msg.err)),
				tea.Quit,
			)
		}
		m.config = msg.config
		m.currentModel = newTasksModel(m.config.Tasks)
		return m, m.currentModel.Init()
	case selectedTask:
		m.taskIndex = msg.index
		model, err := newMachinesModel(m.config, m.taskIndex)
		if err != nil {
			m.err = err
			return m, tea.Quit
		}
		m.currentModel = model
		return m, m.currentModel.Init()
	case selectedMachineUser:
		m.machine = msg.machine
		m.user = msg.user
		runnerConfig, err := m.config.Convert(m.taskIndex, m.machine, m.user)
		if err != nil {
			m.err = err
			return m, tea.Quit
		}
		m.currentModel = newLogsModel(runnerConfig)
		return m, m.currentModel.Init()
	case errorMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	str := m.currentModel.View()
	if m.err != nil {
		str += PrintError(m.err)
	}

	return str
}

func PrintError(err error) string {
	var style = lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#11111B")).
		Background(lipgloss.Color("#F38BA8"))

	return fmt.Sprintf("%s %s", style.Render("ERROR"), err.Error())
}
