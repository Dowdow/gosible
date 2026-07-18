package config

type Task struct {
	Name     string   `yaml:"name"`
	Machines []string `yaml:"machines"`
	Actions  []Action `yaml:"actions"`
}
