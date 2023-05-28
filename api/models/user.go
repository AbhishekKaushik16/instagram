package models

type User struct {
	ID       string `bson:"_id,omitempty"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	FullName string `bson:"fullName"`
	Bio      string `bson:"bio,omitempty"`
}
