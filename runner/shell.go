package runner

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

type ShellArgs string

func (a ShellArgs) Pty() bool {
	return true
}

func (a ShellArgs) Run(session *ssh.Session, ch chan tea.Msg) error {
	if err := session.Run(string(a)); err != nil {
		return fmt.Errorf("[shell] %v\n", err)
	}

	fmt.Println("OUI")

	return nil
}
