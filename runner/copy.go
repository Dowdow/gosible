package runner

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Dowdow/gosible/utils"
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

func (a *CopyArgs) Run(session *ssh.Session) error {
	srcInfo, err := os.Stat(a.Src)
	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}
	defer stdin.Close()

	cmd := fmt.Sprintf("scp -tr %s", path.Dir(a.Dest))
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("[copy] %v\n", err)
	}

	if !srcInfo.IsDir() {
		f, err := os.Open(a.Src)
		if err != nil {
			return fmt.Errorf("[copy] %v\n", err)
		}
		defer f.Close()

		fmt.Fprintf(stdin, "C%04o %d %s\n", srcInfo.Mode().Perm(), srcInfo.Size(), filepath.Base(a.Dest))
		io.Copy(stdin, f)
		fmt.Fprint(stdin, "\x00")
	} else {
		var remainingDepth int = 0
		var latestRelative string = ""
		err = filepath.WalkDir(a.Src, func(curPath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(a.Src, curPath)
			if err != nil {
				return err
			}

			// root directory
			if rel == "." {
				fmt.Fprintf(stdin, "D%04o 0 %s\n", info.Mode().Perm(), filepath.Base(a.Dest))
				return nil
			}

			// sub directory
			if d.IsDir() {
				depth := pathDepthDiff(latestRelative, rel)
				remainingDepth += depth

				if depth < 1 {
					toClosed := -depth + 1
					for range toClosed {
						fmt.Fprint(stdin, "E\n") // Close previous folders
					}
				}

				latestRelative = rel
				fmt.Fprintf(stdin, "D%04o 0 %s\n", info.Mode().Perm(), filepath.Base(curPath))
				return nil
			}

			// file
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
		for range remainingDepth {
			fmt.Fprint(stdin, "E\n")
		}

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

func pathDepthDiff(p1, p2 string) int {
	d1 := pathDepth(p1)
	d2 := pathDepth(p2)
	return d2 - d1
}

func pathDepth(p string) int {
	clean := filepath.Clean(p)

	if clean == "." {
		return 0
	}

	parts := strings.Split(clean, string(filepath.Separator))
	return len(parts)
}
