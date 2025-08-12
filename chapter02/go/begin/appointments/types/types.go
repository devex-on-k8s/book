package types

import (
	"encoding/json"
	"time"
)

type Appointment struct {
	ID              string    `json:"id"`
	PatientID       string    `json:"patientId"`
	AppointmentDate time.Time `json:"appointmentDate"`
}

func (s Appointment) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
