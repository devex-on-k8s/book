package web

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"

	"github.com/devex-on-k8s/book/appointments/db"
	"github.com/devex-on-k8s/book/appointments/types"
)

// Server is the API Server struct
type Server struct {
	DB *sql.DB
}

// NewServer creates a newServer.
func NewServer(db *sql.DB) *Server {
	return &Server{
		DB: db,
	}
}

// GetAllAppointments returns all appointments.
func (s *Server) GetAllAppointments(ctx fiber.Ctx) error {
	appointments, err := db.GetAll(s.DB)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("There was an error executing the query %v", err))
	}

	log.Printf("Appointments retrieved from Database: %d", len(appointments))

	return ctx.JSON(appointments)
}

// DeleteAllAppointments delete all appointments.
func (s *Server) DeleteAllAppointments(ctx fiber.Ctx) error {
	err := db.DeleteAll(s.DB)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("There was an error executing the query %v", err))
	}

	return ctx.JSON(fiber.Map{
		"message": "All Appointments deleted from Database.",
	})
}

// CreateAppointment creates a new appointment.
func (s *Server) CreateAppointment(ctx fiber.Ctx) error {
	appointment := &types.Appointment{}
	if err := ctx.Bind().Body(appointment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("There was an error decoding the request body into the struct: %v", err))
	}

	err := db.CreateAppointment(s.DB, appointment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("There was an error creating the appointment: %v", err))
	}

	return ctx.JSON(appointment)
}

// Welcome returns a welcome message from the Appointments Service
func (s *Server) Welcome(ctx fiber.Ctx) error {
	welcome := Welcome{
		Message: "Welcome to the Appointments API!",
	}

	return ctx.JSON(welcome)
}

const (
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

type Welcome struct {
	Message string `json:"message"`
}
