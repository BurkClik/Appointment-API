package model

type Appointment struct {
	Hour        string `json:"hour" bson:"hour,omitempty"`
	Date        string `json:"date" bson:"date,omitempty"`
	PatientName string `json:"patientName" bson:"patient_name,omitempty"`
	DoctorName  string `json:"doctorName" bson:"doctor_name,omitempty"`
}
