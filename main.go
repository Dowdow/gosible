package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Dowdow/gosible/config"
	"github.com/Dowdow/gosible/converter"
	"github.com/Dowdow/gosible/runner"
)

var c = config.Config{}

func main() {
	args := os.Args[1:]
	argsSize := len(args)

	if argsSize == 0 {
		fmt.Println("Usage: gosible <config.json> <task-index> <machine.user>")
		os.Exit(0)
	}

	configFilePath := args[0]
	info, err := os.Stat(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid file: %s\n", configFilePath)
		os.Exit(1)
	}

	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is a directory\n", configFilePath)
		os.Exit(1)
	}

	extension := filepath.Ext(configFilePath)
	if extension != ".json" {
		fmt.Fprintf(os.Stderr, "Error: %s is not a JSON file (.json)\n", configFilePath)
		os.Exit(1)
	}

	err = c.ParseConfigFile(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: the configuration is not valid\n%v\n", err)
		os.Exit(1)
	}

	err = c.ValidateIds()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: the ids are not corrects\n%v\n", err)
		os.Exit(1)
	}

	// If no task index provided, show available tasks
	if argsSize < 2 {
		c.Print()
		os.Exit(0)
	}

	// Parse and check task index
	arg1 := args[1]
	taskIndex, err := strconv.Atoi(arg1)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: the task index must be a positive integer\n")
		os.Exit(1)
	}
	if !c.HasTask(taskIndex) {
		fmt.Fprint(os.Stderr, "Error: the task index must be between 0 and task size -1\n")
		os.Exit(1)
	}

	// If no machine.user provided, show available machines and users
	if argsSize < 3 {
		err = c.PrintTask(taskIndex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check the machine user combo
	machineUser := args[2]
	if !c.HasMachineUser(machineUser) {
		fmt.Fprintf(os.Stderr, "Error: the machine.user combo '%s' does not exists\n", machineUser)
		os.Exit(1)
	}

	// Convert the config to runner struct
	actions, machine, err := converter.ConvertConfig(&c, taskIndex, machineUser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Execute task
	runner.Run(*actions, *machine)
}
