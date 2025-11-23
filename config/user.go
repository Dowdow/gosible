package config

import (
	"github.com/Dowdow/gosible/env"
)

type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ssh      string `json:"ssh"`
	Become   string `json:"become"`
}

func (u *User) Validate() error {
	u.Password = env.ReplaceEnv(u.Password)
	u.Become = env.ReplaceEnv(u.Become)

	return nil
}
