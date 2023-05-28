package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Posts struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId,omitempty"`
	ImageUrl  string             `bson:"imageUrl,omitempty"`
	Caption   string             `bson:"caption,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
}
