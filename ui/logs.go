package ui

import (
	"fmt"

	"github.com/Dowdow/gosible/runner"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	STATUS_SUCCESS = 0
	STATUS_ERROR   = 1
	STATUS_PENDING = 2
)

var (
	titleStype = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#11111b")).Background(lipgloss.Color("#fab387")).Padding(0, 1)
	okStyle    = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#11111B")).Background(lipgloss.Color("#A6E3A1"))
	koStyle    = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#11111B")).Background(lipgloss.Color("#F38BA8"))
)

type actionView struct {
	name   string
	status int
}

type logsModel struct {
	runner       *runner.Runner
	ch           chan tea.Msg
	actions      []actionView
	totalActions int
	spinner      spinner.Model
}

func newLogsModel(config *runner.Config) logsModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#B4BEFE"))
	return logsModel{
		runner:       runner.NewRunner(config),
		ch:           make(chan tea.Msg),
		actions:      make([]actionView, 0),
		totalActions: len(config.Actions),
		spinner:      s,
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
		m.spinner.Tick,
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
		if msg.Success {
			m.actions[len(m.actions)-1].status = STATUS_SUCCESS
		} else {
			m.actions[len(m.actions)-1].status = STATUS_ERROR
		}
	case runner.StdoutMsg:
		// fmt.Println(visualizeBytes(msg.Msg))
	case runner.StderrMsg:
		// fmt.Println(visualizeBytes(msg.Msg))
	case runner.ErrorMsg:
		close(m.ch)
		m.actions[len(m.actions)-1].status = STATUS_ERROR
		return m, func() tea.Msg {
			return errorMsg{err: msg.Error}
		}
	case runner.EndMsg:
		close(m.ch)
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, read(m.ch)
}

func (m logsModel) View() string {
	str := fmt.Sprintf("%s\n\n", titleStype.Render("Executing"))

	for index, action := range m.actions {
		switch action.status {
		case STATUS_PENDING:
			str += m.spinner.View()
		case STATUS_SUCCESS:
			str += okStyle.Render("OK")
		case STATUS_ERROR:
			str += koStyle.Render("KO")
		}
		str += fmt.Sprintf(" [%d/%d] %s\n", index+1, m.totalActions, action.name)
	}

	return str
}
