package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty"`
	Mail     string             `json:"mail" bson:"mail,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
	City     string             `json:"city" bson:"city,omitempty"`
	Gender   string             `json:"gender" bson:"gender,omitempty"`
	Birthday time.Time          `json:"birthday" bson:"birthday,omitempty"`
	IsDoctor bool               `json:"is_doctor" bson:"is_doctor"`
	Doctor   Doctor             `json:"doctor" bson:"doctor,omitempty"`
	Token    string             `json:"token,omitempty" bson:"-"`
}
