package bcrypt

import lib "golang.org/x/crypto/bcrypt"

type service struct{}

// NewService returns a new service object
func NewService() *service {
	return &service{}
}

func (*service) GenerateFromPassword(password []byte) ([]byte, error) {
	return lib.GenerateFromPassword(password, lib.DefaultCost)
}

func (*service) CompareHashAndPassword(hash, password []byte) error {
	return lib.CompareHashAndPassword(hash, password)
}
