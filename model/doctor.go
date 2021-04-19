package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Doctor struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Mail string `json:"mail" bson:"mail,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
	City string `json:"city" bson:"city,omitempty"`
	Hospital string `json:"hospital" bson:"hospital,omitempty"`
	Department string `json:"department" bson:"department,omitempty"`
	Title string `json:"title" bson:"title,omitempty"`
	About string `json:"about" bson:"about,omitempty"`
	Token string `json:"token,omitempty" bson:"-"`
}