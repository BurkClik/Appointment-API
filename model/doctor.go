package model

type Doctor struct {
	Hospital   string  `json:"hospital" bson:"hospital,omitempty"`
	Department string  `json:"department" bson:"department,omitempty"`
	Title      string  `json:"title" bson:"title,omitempty"`
	About      string  `json:"about" bson:"about,omitempty"`
	VoteCount  int     `json:"vote_count" bson:"vote_count,omitempty"`
	VoteRate   float32 `json:"vote_rate" bson:"vote_rate,omitempty"`
}
