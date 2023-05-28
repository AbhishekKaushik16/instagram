package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comments struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ParentType string             `bson:"parentType,omitempty"`
	ParentID   primitive.ObjectID `bson:"parentId"`
}
