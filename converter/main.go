package converter

import (
	"fmt"

	"github.com/Dowdow/gosible/config"
	"github.com/Dowdow/gosible/runner"
)

func ConvertConfig(c *config.Config, taskIndex int, machineUser string) (*[]runner.Action, *runner.Machine, error) {
	if !c.HasTask(taskIndex) {
		return nil, nil, fmt.Errorf("the task index must be between 0 and task size -1")
	}
	if !c.HasMachineUser(machineUser) {
		return nil, nil, fmt.Errorf("the machine.user combo '%s' does not exists\n", machineUser)
	}

	// Convert config.Machine and config.User to runner.Machine
	runnerMachine := runner.Machine{}
	for _, machine := range c.Inventory {
		for _, user := range machine.Users {
			if machineUser == fmt.Sprintf("%s.%s", machine.Id, user.User) {
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

	return &runnerActions, &runnerMachine, nil
}
