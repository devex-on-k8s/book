//go:build !dev

package main

import (
	"github.com/gofiber/fiber/v3"
)

// ConfigureApp configures the fiber app, including the database connection string.
// The connection string is built using the environment variables for Postgres.
// If the environment variables are not set, the connection string falls back to
// a default connection string targeting localhost and the default database values.
// See the env.go file for the default values.
func ConfigureApp(cfg fiber.Config) (*AppConfig, error) {
	app := fiber.New(cfg)

	DB = "postgres://" + PostgresqlUsername + ":" + PostgresqlPassword + "@" + PostgresqlHost + ":" + PostgresqlPort + "/" + PostgresqlDatabase + "?sslmode=disable"

	return &AppConfig{
		App:            app,
		StartupCancel:  func() {}, // No-op for production
		ShutdownCancel: func() {}, // No-op for production
	}, nil
}
