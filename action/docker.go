package action

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
)

type DockerArgs struct {
	Src   string `json:"src"`
	Dest  string `json:"dest"`
	Image string `json:"image"`
	Tar   string `json:"tar,omitempty"`
	Pull  bool   `json:"pull,omitempty"`
	Clean bool   `json:"clean,omitempty"`
}

func (a *DockerArgs) Validate() error {
	// Validate tar name
	if a.Tar == "" {
		a.Tar = "image.tar"
	}
	if !bytes.HasSuffix([]byte(a.Tar), []byte(".tar")) {
		return fmt.Errorf("[docker] tar should be a *.tar file")
	}

	return nil
}

func (a *DockerArgs) Prepare(resolvePath func(string) string, replaceEnv func(string) string) error {
	a.Src = resolvePath(a.Src)
	return nil
}

func (a *DockerArgs) Run(executor Executor) (string, string, error) {
	srcTarPath := path.Join(a.Src, a.Tar)

	if a.Clean {
		defer func() error {
			cleanCmd := exec.Command("rm", srcTarPath)
			if err := cleanCmd.Run(); err != nil {
				return fmt.Errorf("[docker] %v\n", err)
			}
			return nil
		}()
	}

	args := []string{"build", "-t", a.Image, a.Src}
	if a.Pull {
		args = append(args, "--pull")
	}

	buildCmd := exec.Command("docker", args...)
	if err := buildCmd.Run(); err != nil {
		return "", "", fmt.Errorf("[docker] %v\n", err)
	}

	saveCmd := exec.Command("docker", "save", "-o", srcTarPath, a.Image)
	if err := saveCmd.Run(); err != nil {
		return "", "", fmt.Errorf("[docker] %v\n", err)
	}

	copy := CopyArgs{
		Src:  srcTarPath,
		Dest: path.Join(a.Dest, a.Tar),
	}
	return copy.Run(executor)
}
