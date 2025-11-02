package args

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

type CopyArgs struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

func (a CopyArgs) Pty() bool {
	return false
}

func (a CopyArgs) Run(session *ssh.Session) error {
	content, err := os.ReadFile(a.Src)
	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}
	defer stdin.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	cmd := fmt.Sprintf("scp -t %s", path.Dir(a.Dest))
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	fmt.Fprintf(stdin, "C0644 %d %s\n", len(content), path.Base(a.Dest))
	stdin.Write(content)
	fmt.Fprint(stdin, "\x00")
	stdin.Close()

	if err := session.Wait(); err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	return nil
}
