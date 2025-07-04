package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	VERSION            = getEnv("VERSION", "0.0.1-SNAPSHOT")
	SOURCE             = getEnv("SOURCE", "https://github.com/")
	APP_PORT           = getEnv("APP_PORT", "8081")
	PostgresqlHost     = getEnv("POSTGRES_HOST", "localhost")
	PostgresqlPort     = getEnv("POSTGRES_PORT", "5432")
	PostgresqlDatabase = getEnv("POSTGRES_DB", "postgres")
	PostgresqlUsername = getEnv("POSTGRES_USERNAME", "postgres")
	PostgresqlPassword = getEnv("POSTGRES_PASSWORD", "postgres")
)

// respondWithJSON is a helper function to write a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	chiServer := NewChiServer()

	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+APP_PORT, chiServer)
	log.Printf("Starting Appointments Service in Port: %s", APP_PORT)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

// getEnv returns the value of an environment variable, or a fallback value if not set.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// NewChiServer creates a new Chi server.
func NewChiServer() *chi.Mux {
	log.Printf("Starting Appointments Service in Port: %s", APP_PORT)

	// create new router
	r := chi.NewRouter()

	// add middlewares
	r.Use(middleware.Logger)

	// connect to database
	db := NewDB()

	// create new server
	server := NewServer(db)

	// add routes
	r.Get("/", server.Welcome)
	r.Get("/appointments", server.GetAllAppointments)
	r.Post("/appointments", server.CreateAppointment)
	r.Delete("/appointments", server.DeleteAllAppointments)

	return r
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
func (s *server) GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	var query = "SELECT id, patientId, category, appointmentDate FROM Appointments a"
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
		err = rows.Scan(&appointment.Id, &appointment.PatientId, &appointment.Category, &appointment.AppointmentDate)
		if err != nil {
			log.Printf("There was an error scanning the sql rows: %v", err)
		}
		appointments = append(appointments, appointment)

	}

	log.Printf("Appointments retrieved from Database: %d", len(appointments))
	respondWithJSON(w, http.StatusOK, appointments)
}

// DeleteAllAppointments delete all appointments.
func (s *server) DeleteAllAppointments(w http.ResponseWriter, r *http.Request) {
	var deleteStmt = "DELETE FROM Appointments"

	var err error

	_, err = s.DB.Exec(deleteStmt)

	if err != nil {
		log.Printf("There was an error executing the query %v", err)
	}

	log.Printf("All Appointments deleted from Database.")
	respondWithJSON(w, http.StatusOK, "")
}

// CreateAppointment creates a new appointment.
func (s *server) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment Appointment
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		log.Printf("There was an error decoding the request body into the struct: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	appointment.Id = uuid.New().String()

	insertStmt := `insert into Appointments(id, patientId, category, appointmentDate) values($1, $2, $3, $4)`

	_, err = s.DB.Exec(insertStmt, appointment.Id, appointment.PatientId, appointment.Category, appointment.AppointmentDate)

	if err != nil {
		log.Printf("An error occurred while executing query: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Appointment Stored in Database: %v", appointment)

	respondWithJSON(w, http.StatusOK, appointment)

}

// Welcome returns a welcome message from the Appointments Service
func (s *server) Welcome(w http.ResponseWriter, r *http.Request) {
	var welcome Welcome = Welcome{
		Message: "Welcome to the Appointments API!",
	}
	w.Header().Set(ContentType, ApplicationJson)
	json.NewEncoder(w).Encode(welcome)
}

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func NewDB() *sql.DB {
	connStr := "postgresql://" + PostgresqlUsername + ":" + PostgresqlPassword + "@" + PostgresqlHost + ":" + PostgresqlPort + "/" + PostgresqlDatabase + "?sslmode=disable"
	log.Printf("Connecting to Database: %s.", connStr)

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
	Category        string    `json:"category"`
	AppointmentDate time.Time `json:"appointmentDate"`
}

type Welcome struct {
	Message string `json:"message"`
}

func (s Appointment) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
