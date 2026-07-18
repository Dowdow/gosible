package ui

import (
	"fmt"
	"strings"

	"github.com/Dowdow/gosible/runner"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	STATUS_SUCCESS = 0
	STATUS_ERROR   = 1
	STATUS_PENDING = 2
)

type actionView struct {
	name   string
	status int
	stdout string
	stderr string
}

type runnerEventMsg struct {
	event runner.Event
}

type logsModel struct {
	runner       *runner.Runner
	ch           chan runner.Event
	actions      []actionView
	totalActions int
	spinner      spinner.Model
	showStd      bool
	done         bool
}

func newLogsModel(config *runner.Config) logsModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return logsModel{
		runner:       runner.NewRunner(config),
		ch:           make(chan runner.Event),
		actions:      make([]actionView, 0),
		totalActions: len(config.Actions),
		spinner:      s,
		showStd:      false,
		done:         false,
	}
}

func read(ch chan runner.Event) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-ch
		if !ok {
			return nil
		}
		return runnerEventMsg{event: event}
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
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.showStd = !m.showStd
		case "enter":
			if m.done {
				return m, func() tea.Msg {
					return taskDoneMsg{}
				}
			}
		case "q":
			if m.done {
				return m, tea.Quit
			}
		}
	case runnerEventMsg:
		switch event := msg.event.(type) {
		case runner.ActionStarted:
			m.actions = append(m.actions, actionView{
				name:   event.Name,
				status: STATUS_PENDING,
			})
		case runner.ActionCompleted:
			last := &m.actions[len(m.actions)-1]
			last.stdout = event.Stdout
			last.stderr = event.Stderr
			if event.Success {
				last.status = STATUS_SUCCESS
			} else {
				last.status = STATUS_ERROR
			}
		case runner.Failed:
			close(m.ch)
			return m, func() tea.Msg {
				return errorMsg{err: event.Err}
			}
		case runner.Done:
			close(m.ch)
			m.done = true
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, read(m.ch)
}

func (m logsModel) View() string {
	var str strings.Builder
	fmt.Fprintf(&str, "%s ", titleStyle.Render("Executing"))

	if m.showStd {
		str.WriteString(titleStyle.Padding(0, 0).Render("TAB"))
	} else {
		str.WriteString("TAB")
	}
	str.WriteString(" to show outputs\n\n")

	for index, action := range m.actions {
		switch action.status {
		case STATUS_PENDING:
			str.WriteString(m.spinner.View())
		case STATUS_SUCCESS:
			str.WriteString(okStyle.Render("OK"))
		case STATUS_ERROR:
			str.WriteString(koStyle.Render("KO"))
		}

		fmt.Fprintf(&str, " [%d/%d] %s\n", index+1, m.totalActions, action.name)

		if m.showStd || action.status == STATUS_ERROR {
			fmt.Fprintf(&str, "%s\n", action.stdout)
			fmt.Fprintf(&str, "%s\n", action.stderr)
		}
	}

	if m.done {
		str.WriteString("\n")

		lastActionStatus := m.actions[len(m.actions)-1].status
		if lastActionStatus == STATUS_SUCCESS {
			str.WriteString(okStyle.Render("DONE"))
		} else {
			str.WriteString(koStyle.Render("DONE"))
		}

		str.WriteString(" Press ENTER to continue or Q to quit")
	}

	return str.String()
}
