package domain

import (
	"go-app-template/domain/value"
)

type User struct {
	UserID value.UserID `json:"user_id"`
	Name   value.Name   `json:"name"`
}

func NewUser(name string) *User {
	return &User{
		UserID: value.NewUserID(),
		Name:   value.NewName(name),
	}
}

func (u *User) ChangeName(name string) {
	u.Name = value.NewName(name)
}
