package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Mail string `json:"mail" bson:"mail,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
	Token string `json:"token,omitempty" bson:"-"`
}
