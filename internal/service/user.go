package service

import (
	"fmt"

	"github.com/namphung1998/auth-service-go/internal"
	"golang.org/x/crypto/bcrypt"
)

// User implements internal.Service
type User struct {
	repo internal.UserRepo
	jwt  internal.JWTService
}

func NewUser(repo internal.UserRepo, jwt internal.JWTService) *User {
	return &User{
		repo: repo,
		jwt:  jwt,
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

// Login logs in a user
func (u *User) Login(email, password string) (string, error) {
	user, err := u.repo.Get(email)
	if err != nil {
		return "", err
	}

	fmt.Println(user.ID)
	fmt.Println(user.Password)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", internal.NewInvalidRequestError(email)
	}

	return u.jwt.GenerateToken(user.ID)
}
