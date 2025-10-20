package config

type Task struct {
	Name     string   `json:"name"`
	Machines []string `json:"machines"`
	Actions  []Action `json:"actions"`
}
