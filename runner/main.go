package runner

import (
	"fmt"
	"os"

	"github.com/Dowdow/gosible/args"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Config struct {
	Machine Machine
	Actions []Action
}

type Action struct {
	Name string
	Type string
	Args args.Args
}

type Machine struct {
	Name     string
	Address  string
	User     string
	Password string
	Ssh      string
	Become   string
}

type RunnerMessage struct {
	text   string
	stdout string
	stderr string
}

func Run(config Config) error {
	authMethods := []ssh.AuthMethod{}

	// Private key auth method
	if config.Machine.Ssh != "" {
		key, err := os.ReadFile(config.Machine.Ssh)
		if err != nil {
			return fmt.Errorf("cannot read private key : %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("cannot parse private key : %v", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Password auth method
	if config.Machine.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Machine.Password))
	}

	clientConfig := &ssh.ClientConfig{
		User:            config.Machine.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // dev
	}

	client, err := ssh.Dial("tcp", config.Machine.Address, clientConfig)
	if err != nil {
		return fmt.Errorf("cannot ssh dial : %v", err)
	}
	defer client.Close()

	for _, action := range config.Actions {
		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("cannot create ssh session : %v", err)
		}
		defer session.Close()

		if action.Args.Pty() {
			modes := ssh.TerminalModes{
				ssh.ECHO:          0,
				ssh.ICANON:        0,
				ssh.ISIG:          0,
				ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_OSPEED: 14400,
			}

			fd := int(os.Stdin.Fd())
			width, height, _ := term.GetSize(fd)

			if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
				return fmt.Errorf("cannot request pty : %v", err)
			}
		}

		err = action.Args.Run(session)
		if err != nil {
			return err
		}
	}

	return nil
}
