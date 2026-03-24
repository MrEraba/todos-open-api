package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID     string
	Name   string
	Email  string
	Active bool
}

func NewUser(name string, email string) *User {

	return &User{
		ID:     uuid.New().String(),
		Name:   name,
		Email:  email,
		Active: false,
	}
}
