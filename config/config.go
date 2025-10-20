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
