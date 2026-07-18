package action

import (
	"fmt"
	"os"
	"strings"
)

type FileArgs struct {
	Dest    string   `yaml:"dest"`
	Content []string `yaml:"content"`
}

func (a *FileArgs) Validate() error {
	return nil
}

func (a *FileArgs) Prepare(resolvePath func(string) string, replaceEnv func(string) string) error {
	for index, line := range a.Content {
		a.Content[index] = replaceEnv(line)
	}
	return nil
}

func (a *FileArgs) Run(executor Executor) (string, string, error) {
	file, err := os.CreateTemp("", "gosible-*")
	if err != nil {
		return "", "", fmt.Errorf("[file] %v\n", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	_, err = file.WriteString(strings.Join(a.Content, "\n"))
	if err != nil {
		return "", "", fmt.Errorf("[file] %v\n", err)
	}

	copy := CopyArgs{
		Src:  file.Name(),
		Dest: a.Dest,
	}
	return copy.Run(executor)
}
