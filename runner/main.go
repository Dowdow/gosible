package runner

import (
	"fmt"
	"os"

	"github.com/Dowdow/gosible/action"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Machine Machine
	Actions []Action
}

type Machine struct {
	Name     string
	Address  string
	User     string
	Password string
	Ssh      string
	Become   string
}

type Action struct {
	Name string
	Args action.Args
}

type Runner struct {
	config Config
}

func NewRunner(config *Config) *Runner {
	return &Runner{
		config: *config,
	}
}

func (r *Runner) Run(ch chan Event) {
	authMethods := []ssh.AuthMethod{}

	// Private key auth method
	if r.config.Machine.Ssh != "" {
		key, err := os.ReadFile(r.config.Machine.Ssh)
		if err != nil {
			ch <- Failed{Err: fmt.Errorf("cannot read private key : %v", err)}
			return
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			ch <- Failed{Err: fmt.Errorf("cannot parse private key : %v", err)}
			return
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Password auth method
	if r.config.Machine.Password != "" {
		authMethods = append(authMethods, ssh.Password(r.config.Machine.Password))
	}

	clientConfig := &ssh.ClientConfig{
		User:            r.config.Machine.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // dev
	}

	client, err := ssh.Dial("tcp", r.config.Machine.Address, clientConfig)
	if err != nil {
		ch <- Failed{Err: fmt.Errorf("cannot ssh dial : %v", err)}
		return
	}
	defer client.Close()

	executor := &sshExecutor{client: client}

	for _, action := range r.config.Actions {
		ch <- ActionStarted{Name: action.Name}

		stdout, stderr, err := action.Args.Run(executor)
		if err != nil {
			ch <- ActionCompleted{Stdout: stdout, Stderr: stderr, Success: false}
			ch <- Failed{Err: err}
			return
		}

		ch <- ActionCompleted{Stdout: stdout, Stderr: stderr, Success: true}
	}

	ch <- Done{}
}
