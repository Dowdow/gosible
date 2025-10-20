package config

type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ssh      string `json:"ssh"`
	Become   string `json:"become"`
}
