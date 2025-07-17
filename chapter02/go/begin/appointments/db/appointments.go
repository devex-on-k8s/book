package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	"github.com/devex-on-k8s/book/appointments/types"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// New creates a new database connection and applies migrations.
func New(connStr string) *sql.DB {
	// Open a new database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Configure goose to use embedded migrations
	goose.SetBaseFS(embedMigrations)

	// Set up goose
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set goose dialect: %v", err)
	}

	// Run migrations
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("SQL migrations applied successfully")

	return db
}

// GetAll returns all appointments.
func GetAll(db *sql.DB) ([]types.Appointment, error) {
	query := "SELECT id, patientId, appointmentDate FROM Appointments a"
	var rows *sql.Rows
	var err error

	rows, err = db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	defer rows.Close()
	appointments := []types.Appointment{}
	for rows.Next() {
		var appointment types.Appointment
		err = rows.Scan(&appointment.ID, &appointment.PatientID, &appointment.AppointmentDate)
		if err != nil {
			return nil, fmt.Errorf("scan rows: %w", err)
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

// DeleteAll delete all appointments.
func DeleteAll(db *sql.DB) error {
	deleteStmt := "DELETE FROM Appointments"

	var err error

	_, err = db.Exec(deleteStmt)
	if err != nil {
		return fmt.Errorf("execute delete: %w", err)
	}

	return nil
}

// CreateAppointment creates a new appointment.
func CreateAppointment(db *sql.DB, appointment *types.Appointment) error {
	appointment.ID = uuid.New().String()

	insertStmt := `insert into Appointments(id, patientId, appointmentDate) values($1, $2, $3)`

	_, err := db.Exec(insertStmt, appointment.ID, appointment.PatientID, appointment.AppointmentDate)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	return nil
}
