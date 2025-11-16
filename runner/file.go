package runner

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

type FileArgs struct {
	Dest    string   `json:"dest"`
	Content []string `json:"content"`
}

func (a *FileArgs) Validate() bool {
	return true
}

func (a *FileArgs) Run(session *ssh.Session, ch chan tea.Msg) error {
	file, err := os.CreateTemp("", "gosible-*")
	if err != nil {
		return fmt.Errorf("[file] %v\n", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	_, err = file.WriteString(strings.Join(a.Content, "\n"))
	if err != nil {
		return fmt.Errorf("[file] %v\n", err)
	}

	copy := CopyArgs{
		Src:  file.Name(),
		Dest: a.Dest,
	}

	copy.Validate()
	return copy.Run(session, ch)
}
