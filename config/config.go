package config

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Config struct {
	Inventory []Machine `json:"inventory"`
	Actions   []Action  `json:"actions"`
	Tasks     []Task    `json:"tasks"`
}

func (c *Config) Print() {
	fmt.Print("List of tasks:\n")
	for index, task := range c.Tasks {
		fmt.Printf(" > %d - %s (%d actions)\n", index, task.Name, len(task.Actions))
	}
}

func (c *Config) PrintTask(taskIndex int) error {
	if taskIndex < 0 || taskIndex >= len(c.Tasks) {
		return fmt.Errorf("the task index must be between 0 and task size -1")
	}

	task := c.Tasks[taskIndex]

	fmt.Printf("> %s\n", task.Name)
	for _, action := range task.Actions {
		if action.Id == "" {
			fmt.Printf("   - %s\n", action.Name)
			continue
		}
		for _, a := range c.Actions {
			if a.Id == action.Id {
				fmt.Printf("   - %s\n", a.Name)
				break
			}
		}
	}

	fmt.Print("\nList of available machines and user:\n")
	for _, machine := range c.Inventory {
		for _, user := range machine.Users {
			if len(task.Machines) == 0 || slices.Contains(task.Machines, machine.Id) || slices.Contains(task.Machines, fmt.Sprintf("%s.%s", machine.Id, user.User)) {
				fmt.Printf(" > %s.%s\n", machine.Id, user.User)
			}
		}
	}

	return nil
}

func (c *Config) ParseConfigFile(configFile string) error {
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configData, &c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) ValidateIds() error {
	machineUserIds := []string{}
	for _, machine := range c.Inventory {
		machineUserIds = append(machineUserIds, machine.Id)
		for _, user := range machine.Users {
			machineUserIds = append(machineUserIds, fmt.Sprintf("%s.%s", machine.Id, user.User))
		}
	}

	actionIds := []string{}
	for _, action := range c.Actions {
		actionIds = append(actionIds, action.Id)
	}

	for _, task := range c.Tasks {
		for _, action := range task.Actions {
			for _, machine := range action.Machines {
				if !slices.Contains(machineUserIds, string(machine)) {
					if strings.Contains(string(machine), ".") {
						return fmt.Errorf("Machine.User id not found: %s", machine)
					}
					return fmt.Errorf("Machine id not found: %s", machine)
				}
			}

			if action.Id != "" && !slices.Contains(actionIds, action.Id) {
				return fmt.Errorf("Action id not found: %s", action.Id)
			}
		}
	}

	return nil
}

func (c *Config) MergeActions(taskIndex int) error {
	if taskIndex < 0 || taskIndex >= len(c.Tasks) {
		return fmt.Errorf("the task index must be between 0 and task size -1")
	}

	task := c.Tasks[taskIndex]

	for _, action := range task.Actions {
		if action.Id != "" {
			for _, a := range c.Actions {
				if a.Id == action.Id {
					action.Id = ""
					action.Name = a.Name
					action.Machines = a.Machines
					action.Type = a.Type
					action.Args = a.Args
				}
			}
		}
	}

	return nil
}
