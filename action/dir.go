package action

import "fmt"

type DirArgs struct {
	Path string `json:"path"`
	Mod  string `json:"mod,omitempty"`
}

func (a *DirArgs) Validate() error {
	return nil
}

func (a *DirArgs) Prepare(resolvePath func(string) string, replaceEnv func(string) string) error {
	return nil
}

func (a *DirArgs) Run(executor Executor) (string, string, error) {
	command := fmt.Sprintf("mkdir -p %s", a.Path)
	if a.Mod != "" {
		command = fmt.Sprintf("%s -m %s", command, a.Mod)
	}

	stdout, stderr, err := executor.Run(command)
	if err != nil {
		return stdout, stderr, fmt.Errorf("[dir] %v\n", err)
	}

	return stdout, stderr, nil
}
