package service

import (
	"github.com/namphung1998/auth-service-go/internal"
	"golang.org/x/crypto/bcrypt"
)

// User implements internal.Service
type User struct {
	repo internal.UserRepo
}

func NewUser(repo internal.UserRepo) *User {
	return &User{
		repo: repo,
	}
}

// Create registers a new user
func (u *User) Create(email, password string) error {
	taken, err := u.repo.IsEmailInUse(email)
	if err != nil {
		return err
	}

	if taken {
		return internal.NewEmailInUseError(email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.Create(email, string(hash))
}
