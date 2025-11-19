package config

import (
	"encoding/json"
	"fmt"

	"github.com/Dowdow/gosible/runner"
)

type Action struct {
	Id   string      `json:"id"`
	Name string      `json:"name"`
	Type string      `json:"type"`
	Args runner.Args `json:"-"`
}

func (a *Action) argsFactory() (runner.Args, error) {
	switch a.Type {
	case "copy":
		return &runner.CopyArgs{}, nil
	case "dir":
		return &runner.DirArgs{}, nil
	case "docker":
		return &runner.DockerArgs{}, nil
	case "file":
		return &runner.FileArgs{}, nil
	case "shell":
		var s runner.ShellArgs = ""
		return &s, nil
	}

	return nil, fmt.Errorf("Unknown type: %s", a.Type)
}

func (a *Action) UnmarshalJSON(data []byte) error {
	type Alias Action
	alias := &struct {
		Args json.RawMessage `json:"args"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}

	a.Id = alias.Id
	a.Name = alias.Name
	a.Type = alias.Type
	a.Args = nil

	if a.Type == "" {
		return nil
	}

	args, err := a.argsFactory()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(alias.Args, args); err != nil {
		return err
	}

	a.Args = args

	return nil
}
