package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// testAppConfig returns a httptest.Server for testing.
func testAppConfig(t *testing.T) *AppConfig {
	t.Helper()

	app, err := NewFiberAppConfig()
	require.NoError(t, err)

	return app
}

func Test_API(t *testing.T) {

	// test server
	appConfig := testAppConfig(t)
	defer appConfig.StartupCancel()
	defer appConfig.ShutdownCancel()

	app := appConfig.App

	t.Run("It should return empty when a GET request is made to '/appointments'", func(t *testing.T) {
		// prepare
		req, err := http.NewRequest(http.MethodDelete, "/appointments", nil)
		require.NoError(t, err)

		res, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		// arrange, act
		req, err = http.NewRequest(http.MethodGet, "/appointments", nil)
		res, err = app.Test(req)
		require.NoError(t, err)

		defer res.Body.Close()

		var appointments []Appointment
		json.NewDecoder(res.Body).Decode(&appointments)

		// assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, 0, len(appointments))

	})

	t.Run("It should return 200 when a POST request is made to '/appointments'", func(t *testing.T) {
		// prepare
		req, err := http.NewRequest(http.MethodDelete, "/appointments", nil)
		require.NoError(t, err)

		res, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		// arrange
		appointment := appointmentFake()

		appointmentAsBytes, _ := appointment.MarshalBinary()

		// act
		req, err = http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(appointmentAsBytes))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		// assert
		res, err = app.Test(req)
		require.NoError(t, err)

		// get
		req, err = http.NewRequest(http.MethodGet, "/appointments", nil)
		require.NoError(t, err)

		// assert
		res, err = app.Test(req)
		require.NoError(t, err)

		defer res.Body.Close()

		var appointments []Appointment
		json.NewDecoder(res.Body).Decode(&appointments)

		// assert
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, len(appointments), 1)
		require.NotEmpty(t, appointments[0].Id)
		require.Equal(t, appointments[0].PatientId, appointment.PatientId)
	})

}

func appointmentFake() Appointment {
	return Appointment{
		PatientId:       "test-patient",
		AppointmentDate: time.Now(),
	}
}
