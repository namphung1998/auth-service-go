package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}
