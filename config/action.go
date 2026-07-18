package config

import (
	"fmt"

	"github.com/Dowdow/gosible/action"
	yaml "github.com/goccy/go-yaml"
)

type Action struct {
	Id   string      `yaml:"id"`
	Name string      `yaml:"name"`
	Type string      `yaml:"type"`
	Args action.Args `yaml:"-"`
}

func (a *Action) argsFactory() (action.Args, error) {
	switch a.Type {
	case "copy":
		return &action.CopyArgs{}, nil
	case "dir":
		return &action.DirArgs{}, nil
	case "docker":
		return &action.DockerArgs{}, nil
	case "file":
		return &action.FileArgs{}, nil
	case "shell":
		var s action.ShellArgs = ""
		return &s, nil
	}

	return nil, fmt.Errorf("Unknown type: %s", a.Type)
}

func (a *Action) UnmarshalYAML(data []byte) error {
	var raw struct {
		Id   string          `yaml:"id"`
		Name string          `yaml:"name"`
		Type string          `yaml:"type"`
		Args yaml.RawMessage `yaml:"args"`
	}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}

	a.Id = raw.Id
	a.Name = raw.Name
	a.Type = raw.Type
	a.Args = nil

	if a.Type == "" {
		return nil
	}

	args, err := a.argsFactory()
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(raw.Args, args); err != nil {
		return err
	}

	a.Args = args

	return nil
}
