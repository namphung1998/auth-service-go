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

// CreateUserRequest contains data for creating a new user
type CreateUserRequest struct {
	Email    string
	Password string
}

// UserService defines the contract for interacting with a service
type UserService interface {
	Create(email, password string) error
}
