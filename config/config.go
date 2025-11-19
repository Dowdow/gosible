package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Dowdow/gosible/env"
	"github.com/Dowdow/gosible/runner"
	"github.com/Dowdow/gosible/utils"
)

type Config struct {
	Inventory []Machine `json:"inventory"`
	Actions   []Action  `json:"actions"`
	Tasks     []Task    `json:"tasks"`
}

func (c *Config) HasTask(taskIndex int) bool {
	return taskIndex >= 0 && taskIndex < len(c.Tasks)
}

func (c *Config) HasMachineUser(machineId string, userId string) bool {
	for _, machine := range c.Inventory {
		for _, user := range machine.Users {
			if machineId == machine.Id && userId == user.User {
				return true
			}
		}
	}

	return false
}

func (c *Config) Validate() error {
	// Create all machine.user combo
	machineUserIds := []string{}
	for _, machine := range c.Inventory {
		machineUserIds = append(machineUserIds, machine.Id)
		for _, user := range machine.Users {
			machineUserIds = append(machineUserIds, fmt.Sprintf("%s.%s", machine.Id, user.User))
		}
	}

	actionIds := []string{}
	for _, action := range c.Actions {
		// Validate args
		err := action.Args.Validate()
		if err != nil {
			return fmt.Errorf("action args are not valid: %s", action.Id)
		}

		actionIds = append(actionIds, action.Id)
	}

	for _, task := range c.Tasks {
		// Check machine.user combo
		for _, machine := range task.Machines {
			if !slices.Contains(machineUserIds, string(machine)) {
				if strings.Contains(string(machine), ".") {
					return fmt.Errorf("machine.user id not found: %s", machine)
				}
				return fmt.Errorf("machine id not found: %s", machine)
			}
		}

		for _, action := range task.Actions {
			// Check action ids
			if action.Id != "" && !slices.Contains(actionIds, action.Id) {
				return fmt.Errorf("action id not found: %s", action.Id)
			}

			// Validate args
			if action.Id == "" {
				err := action.Args.Validate()
				if err != nil {
					return fmt.Errorf("action args are not valid: %s : %s", task.Name, action.Name)
				}
			}
		}
	}

	return nil
}

func (c *Config) Convert(taskIndex int, machineId string, userId string) (*runner.Config, error) {
	if !c.HasTask(taskIndex) {
		return nil, fmt.Errorf("the task index must be between 0 and task size -1")
	}
	if !c.HasMachineUser(machineId, userId) {
		return nil, fmt.Errorf("the machine/user combo '%s/%s' does not exists\n", machineId, userId)
	}

	// Convert config.Machine and config.User to runner.Machine
	runnerMachine := runner.Machine{}
	for _, machine := range c.Inventory {
		for _, user := range machine.Users {
			if machineId == machine.Id && userId == user.User {
				runnerMachine.Address = machine.Address
				runnerMachine.Name = machine.Name
				runnerMachine.User = user.User
				runnerMachine.Password = user.Password
				runnerMachine.Ssh = user.Ssh
				runnerMachine.Become = user.Become
			}
		}
	}

	// Convert config.Action to runner.Actions
	runnerActions := make([]runner.Action, 0)

	task := c.Tasks[taskIndex]
	for _, action := range task.Actions {
		// Specific action
		if action.Id == "" {
			runnerActions = append(runnerActions, runner.Action{
				Name: action.Name,
				Type: action.Type,
				Args: action.Args,
			})
		} else {
			// Generic action
			for _, a := range c.Actions {
				if a.Id == action.Id {
					runnerActions = append(runnerActions, runner.Action{
						Name: a.Name,
						Type: a.Type,
						Args: a.Args,
					})
				}
			}
		}
	}

	return &runner.Config{
		Machine: runnerMachine,
		Actions: runnerActions,
	}, nil
}

func ParseConfig() (*Config, error) {
	args := os.Args[1:]
	configFilePath := args[0]

	info, err := os.Stat(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("invalid file: %s", configFilePath)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", configFilePath)
	}

	extension := filepath.Ext(configFilePath)
	if extension != ".json" {
		return nil, fmt.Errorf("%s is not a JSON file (.json)", configFilePath)
	}

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read the config file: %v", err)
	}

	var c = Config{}
	err = json.Unmarshal(configData, &c)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshall the config file: %v", err)
	}

	utils.SetConfigDir(configFilePath)

	err = env.ParseEnv(utils.ResolvePath(".env"))
	if err != nil {
		return nil, fmt.Errorf("dotenv: %v", err)
	}

	err = c.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation: %v", err)
	}

	return &c, nil
}
