package models

type User struct {
	UserName string `json:"login"`
	Password string `json:"password"`
}
