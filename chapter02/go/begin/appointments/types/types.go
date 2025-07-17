package types

import (
	"encoding/json"
	"time"
)

type Appointment struct {
	Id              string    `json:"id"`
	PatientId       string    `json:"patientId"`
	AppointmentDate time.Time `json:"appointmentDate"`
}

func (s Appointment) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
