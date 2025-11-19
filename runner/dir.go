package runner

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type DirArgs struct {
	Path string `json:"path"`
	Mod  string `json:"mod,omitempty"`
}

func (a *DirArgs) Validate() error {
	return nil
}

func (a *DirArgs) Run(session *ssh.Session) error {
	command := fmt.Sprintf("mkdir -p %s", a.Path)
	if a.Mod != "" {
		command = fmt.Sprintf("%s -m %s", command, a.Mod)
	}

	if err := session.Run(command); err != nil {
		return fmt.Errorf("[dir] %v\n", err)
	}

	return nil
}
