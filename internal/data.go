package internal

//go:generate mockgen -source=data.go -package=mock -destination=mock/data.go

// User mirrors the users collection
type User struct {
	ID       string `bson:"_id"`
	Password string `bson:"password"`
}

// UserRepo defines a contract for accessing data related to users
type UserRepo interface {
	IsEmailInUse(email string) (bool, error)
	Create(email, password string) error
	Get(email string) (User, error)
}
