package runner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/Dowdow/gosible/utils"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

type CopyArgs struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

func (a *CopyArgs) Validate() bool {
	a.Src = utils.ResolvePath(a.Src)

	return true
}

func (a *CopyArgs) Run(session *ssh.Session, ch chan tea.Msg) error {
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}
	defer stdin.Close()

	cmd := fmt.Sprintf("scp -tr %s", path.Dir(a.Dest))
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	err = filepath.WalkDir(a.Src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(filepath.Dir(a.Src), path)
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		// Directory
		if d.IsDir() {
			if rel != "." {
				fmt.Fprintf(stdin, "D%04o 0 %s\n", info.Mode().Perm(), filepath.Base(path))
			}
			return nil
		}

		// File
		size := info.Size()
		fmt.Fprintf(stdin, "C%04o %d %s\n", info.Mode().Perm(), size, filepath.Base(path))
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		io.Copy(stdin, f)
		f.Close()
		fmt.Fprint(stdin, "\x00")

		return nil
	})
	fmt.Fprint(stdin, "E\n")
	stdin.Close()

	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	return nil
}
