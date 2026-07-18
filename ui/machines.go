package ui

import (
	"fmt"
	"slices"

	"github.com/Dowdow/gosible/config"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type machineUser struct {
	machine string
	user    string
}

type machinesModel struct {
	list         list.Model
	machineUsers []machineUser
}

func newMachinesModel(c *config.Config, taskIndex int) (machinesModel, error) {
	if !c.HasTask(taskIndex) {
		return machinesModel{}, fmt.Errorf("the task index must be between 0 and task size -1")
	}

	task := c.Tasks[taskIndex]

	items := make([]list.Item, 0)
	machineUsers := make([]machineUser, 0)
	for _, machine := range c.Inventory {
		for _, user := range machine.Users {
			if len(task.Machines) == 0 || slices.Contains(task.Machines, machine.Id) || slices.Contains(task.Machines, fmt.Sprintf("%s.%s", machine.Id, user.User)) {
				items = append(items, simpleItem{
					title: fmt.Sprintf("%s - %s", machine.Name, user.User),
					desc:  machine.Address,
				})
				machineUsers = append(machineUsers, machineUser{
					machine: machine.Id,
					user:    user.User,
				})
			}
		}
	}

	l := newList(items, "Select a machine/user combo")
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		}
	}

	return machinesModel{
		list:         l,
		machineUsers: machineUsers,
	}, nil
}

func (m machinesModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m machinesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		if !m.list.SettingFilter() {
			switch msg.String() {
			case "enter":
				index := m.list.GlobalIndex()
				if index < 0 || index >= len(m.machineUsers) {
					return m, func() tea.Msg {
						return errorMsg{err: fmt.Errorf("invalid machine/user combo")}
					}
				}

				mu := m.machineUsers[index]
				return m, func() tea.Msg {
					return selectedMachineUser{machine: mu.machine, user: mu.user}
				}
			case "esc":
				return m, func() tea.Msg {
					return backToTasksMsg{}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m machinesModel) View() string {
	return m.list.View()
}
