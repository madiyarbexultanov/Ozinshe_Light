package models

type User struct {
	Id				int
	Name			string
	Email			string
	PasswordHash	string
}