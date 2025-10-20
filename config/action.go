package config

import (
	"encoding/json"
	"fmt"
)

type Action struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Machines []string `json:"machines"`
	Type     string   `json:"type"`
	Args     any      `json:"args"`
}

type ShellArgs string

type DirArgs struct {
	Path string `json:"path"`
	Mod  string `json:"mod,omitempty"`
}

func actionArgsFactory(t string) (any, error) {
	switch t {
	case "shell":
		return new(ShellArgs), nil
	case "dir":
		return &DirArgs{}, nil
	}

	return nil, fmt.Errorf("Unknown type: %s", t)
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
	a.Machines = alias.Machines
	a.Type = alias.Type
	a.Args = nil

	if alias.Type == "" {
		return nil
	}

	args, err := actionArgsFactory(alias.Type)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(alias.Args, &args); err != nil {
		return err
	}

	a.Args = args

	return nil
}
