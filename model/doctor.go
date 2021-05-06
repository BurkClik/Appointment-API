package model

type Doctor struct {
	Hospital   string `json:"hospital" bson:"hospital,omitempty"`
	Department string `json:"department" bson:"department,omitempty"`
	Title      string `json:"title" bson:"title,omitempty"`
	About      string `json:"about" bson:"about,omitempty"`
}
