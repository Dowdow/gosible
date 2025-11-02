package ui

import (
	"fmt"
	"slices"

	"github.com/Dowdow/gosible/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type machineItem struct {
	machine string
	user    string
	title   string
	desc    string
}

func (i machineItem) Title() string       { return i.title }
func (i machineItem) Description() string { return i.desc }
func (i machineItem) FilterValue() string { return i.title }

type machinesModel struct {
	list list.Model
}

func newMachinesModel(c *config.Config, taskIndex int) (machinesModel, error) {
	if !c.HasTask(taskIndex) {
		return machinesModel{}, fmt.Errorf("the task index must be between 0 and task size -1")
	}

	task := c.Tasks[taskIndex]

	items := make([]list.Item, 0)
	for _, machine := range c.Inventory {
		for _, user := range machine.Users {
			if len(task.Machines) == 0 || slices.Contains(task.Machines, machine.Id) || slices.Contains(task.Machines, fmt.Sprintf("%s.%s", machine.Id, user.User)) {
				items = append(items, machineItem{
					machine: machine.Id,
					user:    user.User,
					title:   fmt.Sprintf("%s - %s", machine.Name, user.User),
					desc:    machine.Address,
				})
			}
		}
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Select a machine/user combo"

	return machinesModel{
		list: list,
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
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case tea.KeyEnter.String():
			item, ok := m.list.SelectedItem().(machineItem)
			if ok {
				return m, func() tea.Msg {
					return selectedMachineUser{machine: item.machine, user: item.user}
				}
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m machinesModel) View() string {
	return m.list.View()
}
