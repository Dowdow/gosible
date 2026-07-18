package action

import "fmt"

type ShellArgs string

func (a ShellArgs) Validate() error {
	return nil
}

func (a ShellArgs) Prepare(resolvePath func(string) string, replaceEnv func(string) string) error {
	return nil
}

func (a ShellArgs) Run(executor Executor) (string, string, error) {
	stdout, stderr, err := executor.Run(string(a))
	if err != nil {
		return stdout, stderr, fmt.Errorf("[shell] %v\n", err)
	}

	return stdout, stderr, nil
}
