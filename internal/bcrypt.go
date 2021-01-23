package internal

//go:generate mockgen -source=bcrypt.go -package=mock -destination=mock/bcrypt.go

// BcryptService provides a contract for working with bcrypt
type BcryptService interface {
	GenerateFromPassword(password []byte) ([]byte, error)
	CompareHashAndPassword(hash, password []byte) error
}
