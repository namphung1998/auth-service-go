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
func (u *User) Create(request internal.CreateUserRequest) error {
	taken, err := u.repo.IsEmailInUse(request.Email)
	if err != nil {
		return err
	}

	if taken {
		return internal.NewEmailInUseError(request.Email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.Create(request.Email, string(hash))
}
