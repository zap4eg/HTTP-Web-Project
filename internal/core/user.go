package core

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"firstName,omitempty" json:"firstName"`
	LastName  string             `bson:"lastName,omitempty" json:"lastName"`
}
