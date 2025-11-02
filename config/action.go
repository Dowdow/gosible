package config

import (
	"encoding/json"
	"fmt"

	"github.com/Dowdow/gosible/args"
)

type Action struct {
	Id   string    `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
	Args args.Args `json:"-"`
}

func (a *Action) argsFactory() (args.Args, error) {
	switch a.Type {
	case "shell":
		var s args.ShellArgs = ""
		return &s, nil
	case "copy":
		return &args.CopyArgs{}, nil
	case "dir":
		return &args.DirArgs{}, nil
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
