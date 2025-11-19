package runner

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type ShellArgs string

func (a ShellArgs) Validate() error {
	return nil
}

func (a ShellArgs) Run(session *ssh.Session) error {
	if err := session.Run(string(a)); err != nil {
		return fmt.Errorf("[shell] %v\n", err)
	}

	return nil
}
