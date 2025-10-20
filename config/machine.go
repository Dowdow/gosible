package config

type Machine struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Users   []User `json:"users"`
}
