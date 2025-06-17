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
	Role     string `json:"role"`
}

type UserAdmin struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Total int64  `json:"total"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
