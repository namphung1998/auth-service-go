package internal

import "fmt"

//go:generate mockgen -source=service.go -package=mock -destination=mock/service.go

// InvalidRequestError is an error that is triggered when a user's input is invalid
type InvalidRequestError struct {
	detail string
}

func (e *InvalidRequestError) Error() string {
	return fmt.Sprintf("invalid request: %s", e.detail)
}

// NewInvalidRequestError returns a new InvalidRequestError
func NewInvalidRequestError(detail string) *InvalidRequestError {
	return &InvalidRequestError{detail: detail}
}

// EmailInUseError is an error that is triggered when the provided email is already used
type EmailInUseError struct {
	detail string
}

func (e *EmailInUseError) Error() string {
	return fmt.Sprintf("email in use: %s", e.detail)
}

// NewEmailInUseError returns a new EmailInUseError
func NewEmailInUseError(detail string) *EmailInUseError {
	return &EmailInUseError{detail: detail}
}

// UserNotFoundError is an error that is triggered when no user can be found with the provided email
type UserNotFoundError struct {
	email string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found with email: %s", e.email)
}

// NewUserNotFoundError returns a new UserNotFoundError
func NewUserNotFoundError(email string) *UserNotFoundError {
	return &UserNotFoundError{email: email}
}

// IncorrectPasswordError is an error that is triggered when the provided password is invalid
type IncorrectPasswordError struct {
	email string
}

func (e *IncorrectPasswordError) Error() string {
	return fmt.Sprintf("incorrect password for email: %s", e.email)
}

// NewIncorrectPasswordError returns a new IncorrectPasswordError
func NewIncorrectPasswordError(email string) *IncorrectPasswordError {
	return &IncorrectPasswordError{email: email}
}

// LoginResponse is returned from the UserService.Login method
type LoginResponse struct {
	Token string
}

// UserService defines the contract for interacting with a service
type UserService interface {
	Create(email, password string) error
	Login(email, password string) (LoginResponse, error)
}
