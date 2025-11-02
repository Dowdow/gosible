package args

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type ShellArgs string

func (a ShellArgs) Pty() bool {
	return true
}

func (a ShellArgs) Run(session *ssh.Session) error {
	session.Stdin = nil
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(string(a)); err != nil {
		return fmt.Errorf("[shell] %v\n", err)
	}

	return nil
}
