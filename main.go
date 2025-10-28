package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Dowdow/gosible/config"
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
	if taskIndex < 0 || taskIndex >= len(c.Tasks) {
		fmt.Fprint(os.Stderr, "Error: the task index must be between 0 and task size -1\n")
		os.Exit(1)
	}

	// Of no machine.user provided, show available machines and users
	if argsSize < 3 {
		err = c.PrintTask(taskIndex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Merge idependant actions in the current task
	c.MergeActions(taskIndex)

	// Get the machine and User

	// Execute task

	/*
		for _, task := range c.Tasks {
			machineId := slices.IndexFunc(c.Inventory, func(machine config.Machine) bool {
				return slices.Contains(task.Machines, machine.Id)
			})

			fmt.Println(c.Inventory[machineId].Address)
			fmt.Println(c.Inventory[machineId].Name)
			fmt.Println(c.Inventory[machineId].Users[0].Ssh)

			key, err := os.ReadFile(c.Inventory[machineId].Users[0].Ssh)
			if err != nil {
				log.Fatalf("Impossible de lire la clé privée : %v", err)
			}

			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				log.Fatalf("Erreur lors du parsing de la clé : %v", err)
			}

			fmt.Println(c.Inventory[machineId].Users[0].User)

			clientConfig := &ssh.ClientConfig{
				User: c.Inventory[machineId].Users[0].User,
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ pour dev uniquement
			}

			client, err := ssh.Dial("tcp", c.Inventory[machineId].Address, clientConfig)
			if err != nil {
				log.Fatalf("Erreur de connexion SSH : %v", err)
			}
			defer client.Close()

			for _, action := range task.Actions {
				var command = action.Args
				if action.Id != "" {
					actionId := slices.IndexFunc(c.Actions, func(a config.Action) bool {
						return a.Id == action.Id
					})

					command = c.Actions[actionId].Args
				}

				session, err := client.NewSession()
				if err != nil {
					log.Fatalf("Erreur de création de session : %v", err)
				}
				defer session.Close()

				var stdout, stderr bytes.Buffer
				session.Stdout = &stdout
				session.Stderr = &stderr

				if err := session.Run(command); err != nil {
					log.Printf("Erreur d'exécution : %v\n", err)
				}

				fmt.Println("---- SORTIE STDOUT ----")
				fmt.Println(stdout.String())
				fmt.Println("---- SORTIE STDERR ----")
				fmt.Println(stderr.String())

				fmt.Println("✅ Commande exécutée avec succès !")
			}
			}
	*/
}
