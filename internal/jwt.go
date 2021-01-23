package internal

//go:generate mockgen -source=jwt.go -package=mock -destination=mock/jwt.go

// JWTService defines a contract for working with JWTs
type JWTService interface {
	GenerateToken(userID string) (string, error)
}
