package model

type Review struct {
	PatientName string  `json:"patient_name" bson:"patient_name,omitempty"`
	Review      string  `json:"review" bson:"review,omitempty"`
	Vote        float32 `json:"vote" bson:"vote,omitempty"`
}
