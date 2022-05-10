package models

type User struct {
	ID       int32  `json:"-"`
	UserName string `json:"login"`
	Password string `json:"password"`
}
