package models

type User struct {
	Login    string
	Password string
	Name     string
}

type UserDB struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
