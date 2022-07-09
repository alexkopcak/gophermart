package models

type User struct {
	ID       int32  `json:"id"`
	UserName string `json:"login"`
	Password string `json:"password"`
}
