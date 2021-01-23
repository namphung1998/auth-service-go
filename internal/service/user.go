package service

import (
	"github.com/namphung1998/auth-service-go/internal"
)

// User implements internal.Service
type User struct {
	repo   internal.UserRepo
	jwt    internal.JWTService
	bcrypt internal.BcryptService
}

// NewUser returns a new User service object
func NewUser(repo internal.UserRepo, jwt internal.JWTService, bcrypt internal.BcryptService) *User {
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

	hash, err := u.bcrypt.GenerateFromPassword([]byte(password))
	if err != nil {
		return err
	}

	return u.repo.Create(email, string(hash))
}

// Login logs in a user
func (u *User) Login(email, password string) (internal.LoginResponse, error) {
	user, err := u.repo.Get(email)
	if err != nil {
		return internal.LoginResponse{}, err
	}

	if err := u.bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return internal.LoginResponse{}, internal.NewInvalidRequestError(email)
	}

	token, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		return internal.LoginResponse{}, err
	}

	return internal.LoginResponse{
		Token: token,
	}, nil
}
