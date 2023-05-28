package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	FullName string             `bson:"fullName"`
	Bio      string             `bson:"bio,omitempty"`
}
