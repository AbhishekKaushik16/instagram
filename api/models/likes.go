package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Likes struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	ParentType string             `bson:"parentType"`
	ParentID   primitive.ObjectID `bson:"parentId"`
	Likes      int                `bson:"likes"`
}
