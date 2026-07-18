package config

type Machine struct {
	Id      string `yaml:"id"`
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Users   []User `yaml:"users"`
}
