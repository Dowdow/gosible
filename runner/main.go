package runner

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

type Action struct {
	Name string
	Type string
	Args any
}

type Machine struct {
	Name     string
	Address  string
	User     string
	Password string
	Ssh      string
	Become   string
}

func Run(actions []Action, machine Machine) error {
	// Prio clé ssh si pas vide
	// Lecture mot de passe si pas vide
	// Si aucun des deux erreur

	key, err := os.ReadFile(machine.Ssh)
	if err != nil {
		log.Fatalf("Impossible de lire la clé privée : %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Erreur lors du parsing de la clé : %v", err)
	}

	clientConfig := &ssh.ClientConfig{
		User: machine.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // dev
	}

	client, err := ssh.Dial("tcp", machine.Address, clientConfig)
	if err != nil {
		log.Fatalf("Erreur de connexion SSH : %v", err)
	}
	defer client.Close()

	for _, action := range actions {
		session, err := client.NewSession()
		if err != nil {
			log.Fatalf("Erreur de création de session : %v", err)
		}
		defer session.Close()

		var stdout, stderr bytes.Buffer
		session.Stdout = &stdout
		session.Stderr = &stderr

		command, ok := action.Args.(string)
		if !ok {
			continue
		}

		if err := session.Run(command); err != nil {
			log.Printf("Erreur d'exécution : %v\n", err)
		}

		fmt.Println("---- SORTIE STDOUT ----")
		fmt.Println(stdout.String())
		fmt.Println("---- SORTIE STDERR ----")
		fmt.Println(stderr.String())
	}

	return nil
}
