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

func (a *CopyArgs) Validate() error {
	a.Src = utils.ResolvePath(a.Src)
	return nil
}

func (a *CopyArgs) Run(session *ssh.Session, ch chan tea.Msg) error {
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}
	defer stdin.Close()

	srcInfo, err := os.Stat(a.Src)
	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}
	isFile := !srcInfo.IsDir()
	destDir := path.Dir(a.Dest)
	destBase := filepath.Base(a.Dest)

	cmd := fmt.Sprintf("scp -tr %s", destDir)
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	if isFile {
		f, err := os.Open(a.Src)
		if err != nil {
			return fmt.Errorf("[copy] %v\n", err)
		}
		defer f.Close()

		fmt.Fprintf(stdin, "C%04o %d %s\n", srcInfo.Mode().Perm(), srcInfo.Size(), destBase)
		io.Copy(stdin, f)
		fmt.Fprint(stdin, "\x00")
	} else {
		rootSrc := a.Src
		rootName := destBase
		err = filepath.WalkDir(rootSrc, func(curPath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(rootSrc, curPath)
			if err != nil {
				return err
			}

			// Root Directory
			if rel == "." {
				fmt.Fprintf(stdin, "D%04o 0 %s\n", info.Mode().Perm(), rootName)
				return nil
			}

			// Sub-Directory
			if d.IsDir() {
				fmt.Fprintf(stdin, "D%04o 0 %s\n", info.Mode().Perm(), filepath.Base(curPath))
				return nil
			}

			// File
			fmt.Fprintf(stdin, "C%04o %d %s\n", info.Mode().Perm(), info.Size(), filepath.Base(curPath))
			f, err := os.Open(curPath)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(stdin, f)
			if err != nil {
				return err
			}

			fmt.Fprint(stdin, "\x00")
			return nil
		})
		fmt.Fprint(stdin, "E\n")

		if err != nil {
			return fmt.Errorf("[copy] %v\n", err)
		}
	}

	stdin.Close()

	if err := session.Wait(); err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	return nil
}
