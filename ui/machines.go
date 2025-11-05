package ui

import (
	"fmt"
	"slices"

	"github.com/Dowdow/gosible/config"
	"github.com/Dowdow/gosible/ui/list"
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
				items = append(items, list.Item{

					Title: fmt.Sprintf("%s - %s", machine.Name, user.User),
					Desc:  machine.Address,
				})
				machineUsers = append(machineUsers, machineUser{
					machine: machine.Id,
					user:    user.User,
				})
			}
		}
	}

	list := list.New(items)
	list.Title = "Select a machine/user combo"

	return machinesModel{
		list:         list,
		machineUsers: machineUsers,
	}, nil
}

func (m machinesModel) Init() tea.Cmd {
	return m.list.Init()
}

func (m machinesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			index := m.list.SelectedIndex()
			if index < 0 || index >= len(m.machineUsers) {
				return m, func() tea.Msg {
					return errorMsg{err: fmt.Errorf("invalid machine/user combo")}
				}
			}

			machineUser := m.machineUsers[index]
			return m, func() tea.Msg {
				return selectedMachineUser{machine: machineUser.machine, user: machineUser.user}
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
