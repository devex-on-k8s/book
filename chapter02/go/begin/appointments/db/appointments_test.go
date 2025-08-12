package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	// blank import to register the postgres driver
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/devex-on-k8s/book/appointments/types"
)

func setupDB(t *testing.T) string {
	t.Helper()

	ctx := context.Background()

	ctr, err := postgres.Run(
		ctx, "postgres:17.5-alpine",
		postgres.BasicWaitStrategies(),
		postgres.WithDatabase("appointments"),
		postgres.WithUsername("app"),
		postgres.WithPassword("pass"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	return connStr
}

func TestDB(t *testing.T) {
	db, err := New(setupDB(t))
	require.NoError(t, err)
	require.NotNil(t, db)

	t.Run("GetAll", func(t *testing.T) {
		appointments, err := GetAll(db)
		require.NoError(t, err)
		require.Empty(t, appointments)
	})

	t.Run("CreateAppointment", func(t *testing.T) {
		appointment := &types.Appointment{
			ID:              uuid.New().String(),
			PatientID:       "123",
			AppointmentDate: time.Now().UTC(),
		}

		err := CreateAppointment(db, appointment)
		require.NoError(t, err)

		appointments, err := GetAll(db)
		require.NoError(t, err)
		require.Len(t, appointments, 1)
		require.Equal(t, appointment.ID, appointments[0].ID)
		require.Equal(t, appointment.PatientID, appointments[0].PatientID)
		require.Equal(t, appointment.AppointmentDate, appointments[0].AppointmentDate)
	})

	t.Run("DeleteAll", func(t *testing.T) {
		err := DeleteAll(db)
		require.NoError(t, err)

		appointments, err := GetAll(db)
		require.NoError(t, err)
		require.Empty(t, appointments)
	})
}
