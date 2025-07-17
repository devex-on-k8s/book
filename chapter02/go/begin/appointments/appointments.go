package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// AppConfig holds the application configuration and cleanup functions.
type AppConfig struct {
	// App is the Fiber app instance.
	App *fiber.App

	// StartupCancel is the context cancel function for the services startup.
	StartupCancel context.CancelFunc

	// ShutdownCancel is the context cancel function for the services shutdown.
	ShutdownCancel context.CancelFunc
}

func main() {
	appConfig, err := NewFiberAppConfig()
	if err != nil {
		log.Fatalf("Failed to create Fiber server: %v", err)
	}

	app := appConfig.App

	// cancel the startup and shutdown contexts in the main function, so that
	// the services are not left running after the main function returns.
	defer appConfig.StartupCancel()
	defer appConfig.ShutdownCancel()

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%v", APP_PORT)); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	// Block the main thread until an interrupt is received
	<-quit

	log.Println("Gracefully shutting down...")
	err = app.Shutdown()
	if err != nil {
		log.Panic(err)
	}
}

// NewFiberAppConfig creates a new Fiber server AppConfig.
func NewFiberAppConfig() (*AppConfig, error) {
	cfg := fiber.Config{
		ErrorHandler: ErrorHandler,
	}

	// configure the app
	appConfig, err := ConfigureApp(cfg)
	if err != nil {
		return nil, fmt.Errorf("configure app: %w", err)
	}

	app := appConfig.App

	// add middlewares
	app.Use(logger.New())

	// create new server instance using a new database connection
	server := NewServer(NewDB())

	// add routes
	app.Get("/", server.Welcome)
	app.Get("/appointments", server.GetAllAppointments)
	app.Post("/appointments", server.CreateAppointment)
	app.Delete("/appointments", server.DeleteAllAppointments)

	return appConfig, nil
}

// server is the API server struct
type server struct {
	DB *sql.DB
}

// NewServer creates a newServer.
func NewServer(db *sql.DB) *server {
	return &server{
		DB: db,
	}
}

// GetAllAppointments returns all appointments.
func (s *server) GetAllAppointments(ctx fiber.Ctx) error {
	var query = "SELECT id, patientId, appointmentDate FROM Appointments a"
	var rows *sql.Rows
	var err error

	rows, err = s.DB.Query(query)

	if err != nil {
		log.Printf("There was an error executing the query %v", err)
	}

	defer rows.Close()
	appointments := []Appointment{}
	for rows.Next() {

		var appointment Appointment
		err = rows.Scan(&appointment.Id, &appointment.PatientId, &appointment.AppointmentDate)
		if err != nil {
			log.Printf("There was an error scanning the sql rows: %v", err)
		}
		appointments = append(appointments, appointment)

	}

	log.Printf("Appointments retrieved from Database: %d", len(appointments))

	return ctx.JSON(appointments)
}

// DeleteAllAppointments delete all appointments.
func (s *server) DeleteAllAppointments(ctx fiber.Ctx) error {
	var deleteStmt = "DELETE FROM Appointments"

	var err error

	_, err = s.DB.Exec(deleteStmt)

	if err != nil {
		log.Printf("There was an error executing the query %v", err)
	}

	log.Printf("All Appointments deleted from Database.")

	return ctx.JSON(fiber.Map{
		"message": "All Appointments deleted from Database.",
	})
}

// CreateAppointment creates a new appointment.
func (s *server) CreateAppointment(ctx fiber.Ctx) error {
	appointment := &Appointment{}
	if err := ctx.Bind().Body(appointment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("There was an error decoding the request body into the struct: %v", err))
	}

	appointment.Id = uuid.New().String()

	insertStmt := `insert into Appointments(id, patientId, appointmentDate) values($1, $2, $3)`

	_, err := s.DB.Exec(insertStmt, appointment.Id, appointment.PatientId, appointment.AppointmentDate)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("An error occurred while executing query: %v", err))
	}

	log.Printf("Appointment Stored in Database: %v", appointment)

	return ctx.JSON(appointment)
}

// Welcome returns a welcome message from the Appointments Service
func (s *server) Welcome(ctx fiber.Ctx) error {
	var welcome Welcome = Welcome{
		Message: "Welcome to the Appointments API!",
	}

	return ctx.JSON(welcome)
}

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func NewDB() *sql.DB {
	connStr := DB

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
	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("SQL migrations applied successfully")

	return db
}

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type Appointment struct {
	Id              string    `json:"id"`
	PatientId       string    `json:"patientId"`
	AppointmentDate time.Time `json:"appointmentDate"`
}

type Welcome struct {
	Message string `json:"message"`
}

func (s Appointment) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
