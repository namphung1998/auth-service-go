package internal

//go:generate mockgen -source=data.go -package=mock -destination=mock/data.go

// UserRepo defines a contract for accessing data related to users
type UserRepo interface {
	IsEmailInUse(email string) (bool, error)
	Create(email, password string) error
}
