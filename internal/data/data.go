package data

import (
	"context"

	"github.com/namphung1998/auth-service-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	db *mongo.Database
}

func NewUserRepo(db *mongo.Database) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) IsEmailInUse(email string) (bool, error) {
	err := r.db.Collection("users").FindOne(context.Background(), bson.M{
		"email": email,
	}).Err()

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *userRepo) Create(email, password string) error {
	_, err := r.db.Collection("users").InsertOne(context.Background(), models.User{
		Email:    email,
		Password: password,
	})

	return err
}
