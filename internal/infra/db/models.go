package db

type User struct {
	Id       int
	Email    string
	Password string
}

type Account struct {
	Id     int
	Name   string
	Type   string
	UserId int
}
