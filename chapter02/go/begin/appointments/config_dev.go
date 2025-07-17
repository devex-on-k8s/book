//go:build dev

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/contrib/testcontainers"
	"github.com/gofiber/fiber/v3"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// ConfigureApp configures the fiber app for development, including the database
// connection string. The connection string is retrieved from the PostgreSQL
// service.
// The service is contributed by the testcontainers contrib GoFiber package
// ().
func ConfigureApp(cfg fiber.Config) (*AppConfig, error) {
	// Define a context provider for the services startup.
	// The timeout is applied when the context is actually used during startup.
	startupCtx, startupCancel := context.WithCancel(context.Background())
	var startupTimeoutCancel context.CancelFunc
	cfg.ServicesStartupContextProvider = func() context.Context {
		// Cancel any previous timeout context
		if startupTimeoutCancel != nil {
			startupTimeoutCancel()
		}
		// Create a new timeout context
		ctx, cancel := context.WithTimeout(startupCtx, 10*time.Second)
		startupTimeoutCancel = cancel
		return ctx
	}

	// Define a context provider for the services shutdown.
	// The timeout is applied when the context is actually used during shutdown.
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	var shutdownTimeoutCancel context.CancelFunc
	cfg.ServicesShutdownContextProvider = func() context.Context {
		// Cancel any previous timeout context
		if shutdownTimeoutCancel != nil {
			shutdownTimeoutCancel()
		}
		// Create a new timeout context
		ctx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		shutdownTimeoutCancel = cancel
		return ctx
	}

	// Add the Postgres service to the app, including custom configuration.
	srv, err := setupPostgres(&cfg)
	if err != nil {
		if startupTimeoutCancel != nil {
			startupTimeoutCancel()
		}
		if shutdownTimeoutCancel != nil {
			shutdownTimeoutCancel()
		}
		startupCancel()
		shutdownCancel()
		return nil, fmt.Errorf("add postgres service: %w", err)
	}

	app := fiber.New(cfg)

	// Retrieve the Postgres service from GoFiber's state, using the service key.
	postgresSrv := fiber.MustGetService[*testcontainers.ContainerService[*postgres.PostgresContainer]](app.State(), srv.Key())

	connString, err := postgresSrv.Container().ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		if startupTimeoutCancel != nil {
			startupTimeoutCancel()
		}
		if shutdownTimeoutCancel != nil {
			shutdownTimeoutCancel()
		}
		startupCancel()
		shutdownCancel()
		return nil, fmt.Errorf("get postgres connection string: %w", err)
	}

	// Override the default database connection string with the one from the Testcontainers service.
	DB = connString

	return &AppConfig{
		App: app,
		StartupCancel: func() {
			if startupTimeoutCancel != nil {
				startupTimeoutCancel()
			}
			startupCancel()
		},
		ShutdownCancel: func() {
			if shutdownTimeoutCancel != nil {
				shutdownTimeoutCancel()
			}
			shutdownCancel()
		},
	}, nil
}

// setupPostgres adds a Postgres service to the app, including custom configuration to allow
// reusing the same container while developing locally.
func setupPostgres(cfg *fiber.Config) (*testcontainers.ContainerService[*postgres.PostgresContainer], error) {
	// Add the Postgres service to the app, including custom configuration.
	srv, err := testcontainers.AddService(cfg, testcontainers.NewModuleConfig(
		"postgres-db",
		"postgres:17.5-alpine",
		postgres.Run,
		postgres.BasicWaitStrategies(),
		postgres.WithDatabase(PostgresqlDatabase),
		postgres.WithUsername(PostgresqlUsername),
		postgres.WithPassword(PostgresqlPassword),
		tc.WithReuseByName("postgres-db-appointments"),
	))
	if err != nil {
		return nil, fmt.Errorf("add postgres service: %w", err)
	}

	return srv, nil
}
