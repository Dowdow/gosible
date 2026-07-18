package config

type User struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Ssh      string `yaml:"ssh"`
	Become   string `yaml:"become"`
}
