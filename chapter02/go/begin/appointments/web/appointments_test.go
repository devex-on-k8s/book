package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/devex-on-k8s/book/appointments/config"
	"github.com/devex-on-k8s/book/appointments/db"
	"github.com/devex-on-k8s/book/appointments/types"
)

// testAppConfig returns a httptest.Server for testing.
func testAppConfig(t *testing.T) *config.AppConfig {
	t.Helper()

	appConfig, err := config.New()
	require.NoError(t, err)

	app := appConfig.App

	// create new server instance using a new database connection
	db, err := db.New(config.DB)
	require.NoError(t, err)

	server := NewServer(db)

	// add routes
	app.Get("/", server.Welcome)
	app.Get("/appointments", server.GetAllAppointments)
	app.Post("/appointments", server.CreateAppointment)
	app.Delete("/appointments", server.DeleteAllAppointments)

	return appConfig
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
		require.NoError(t, err)

		res, err = app.Test(req)
		require.NoError(t, err)

		defer res.Body.Close()

		var appointments []types.Appointment
		err = json.NewDecoder(res.Body).Decode(&appointments)
		require.NoError(t, err)

		// assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Empty(t, appointments)
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
		_, err = app.Test(req)
		require.NoError(t, err)

		// get
		req, err = http.NewRequest(http.MethodGet, "/appointments", nil)
		require.NoError(t, err)

		// assert
		res, err = app.Test(req)
		require.NoError(t, err)

		defer res.Body.Close()

		var appointments []types.Appointment
		err = json.NewDecoder(res.Body).Decode(&appointments)
		require.NoError(t, err)

		// assert
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Len(t, appointments, 1)
		require.NotEmpty(t, appointments[0].ID)
		require.Equal(t, appointments[0].PatientID, appointment.PatientID)
	})
}

func appointmentFake() types.Appointment {
	return types.Appointment{
		PatientID:       "test-patient",
		AppointmentDate: time.Now(),
	}
}
