package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testServer returns a httptest.Server for testing.
func testServer() *httptest.Server {
	chiServer := NewChiServer()
	return httptest.NewServer(chiServer)
}

func Test_API(t *testing.T) {

	// test server
	ts := testServer()
	defer ts.Close()

	t.Run("It should return empty when a GET request is made to '/appointments'", func(t *testing.T) {
		// prepare
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/appointments", ts.URL), nil)

		client := &http.Client{}
		_, err := client.Do(req)

		assert.NoError(t, err)

		// arrange, act
		resp, err := http.Get(fmt.Sprintf("%s/appointments", ts.URL))

		assert.NoError(t, err)

		defer resp.Body.Close()

		var appointments []Appointment
		json.NewDecoder(resp.Body).Decode(&appointments)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 0, len(appointments))

	})

	t.Run("It should return 200 when a POST request is made to '/appointments'", func(t *testing.T) {
		// prepare
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/appointments", ts.URL), nil)

		client := &http.Client{}
		_, err := client.Do(req)

		assert.NoError(t, err)

		// arrange
		appointment := appointmentFake()

		appointmentAsBytes, _ := appointment.MarshalBinary()

		// act
		_, err = http.Post(fmt.Sprintf("%s/appointments", ts.URL), "application/json", bytes.NewBuffer(appointmentAsBytes))

		// assert
		assert.NoError(t, err)

		// get
		resp, err := http.Get(fmt.Sprintf("%s/appointments", ts.URL))
		// assert
		assert.NoError(t, err)

		defer resp.Body.Close()

		var appointments []Appointment
		json.NewDecoder(resp.Body).Decode(&appointments)

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, len(appointments), 1)
		assert.NotEmpty(t, appointments[0].Id)
		assert.Equal(t, appointments[0].PatientId, appointment.PatientId)
	})

	t.Run("It should return 200 when a POST request is made to '/appointments' with category", func(t *testing.T) {
		// prepare
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/appointments", ts.URL), nil)

		client := &http.Client{}
		_, err := client.Do(req)

		assert.NoError(t, err)

		// arrange
		appointment := appointmentWithCategoryFake()

		appointmentAsBytes, _ := appointment.MarshalBinary()

		// act
		_, err = http.Post(fmt.Sprintf("%s/appointments", ts.URL), "application/json", bytes.NewBuffer(appointmentAsBytes))

		// assert
		assert.NoError(t, err)

		// get
		resp, err := http.Get(fmt.Sprintf("%s/appointments", ts.URL))
		// assert
		assert.NoError(t, err)

		defer resp.Body.Close()

		var appointments []Appointment
		json.NewDecoder(resp.Body).Decode(&appointments)

		// assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, len(appointments), 1)
		assert.NotEmpty(t, appointments[0].Id)
		assert.Equal(t, appointments[0].PatientId, appointment.PatientId)
		assert.Equal(t, appointments[0].Category, appointment.Category)
	})

}

func appointmentFake() Appointment {
	return Appointment{
		PatientId:       "test-patient",
		AppointmentDate: time.Now(),
	}
}

func appointmentWithCategoryFake() Appointment {
	return Appointment{
		PatientId:       "test-patient",
		Category:        "checkup",
		AppointmentDate: time.Now(),
	}
}
