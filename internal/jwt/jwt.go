package jwt

import (
	"time"

	lib "github.com/dgrijalva/jwt-go"
)

type service struct {
	signingMethod lib.SigningMethod
	secret        []byte
}

func NewService(signingMethod lib.SigningMethod, secret []byte) *service {
	return &service{
		signingMethod: signingMethod,
		secret:        secret,
	}
}

func (s *service) GenerateToken(userID string) (string, error) {
	token := lib.NewWithClaims(s.signingMethod, lib.StandardClaims{
		Subject:   userID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(s.secret)
}
