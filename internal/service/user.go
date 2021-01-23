package service

import "github.com/namphung1998/auth-service-go/internal"

// User implements internal.Service
type User struct {
	repo internal.UserRepo
}

// CreateUser registers a new user
func (u *User) CreateUser(email, password string) error {
	panic("not implemented")
}
