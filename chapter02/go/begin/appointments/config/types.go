package config

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
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

// New creates a new Fiber server AppConfig.
func New() (*AppConfig, error) {
	cfg := fiber.Config{
		ErrorHandler: ErrorHandler,
	}

	// configure the app
	appConfig, err := ConfigureApp(cfg)
	if err != nil {
		return nil, fmt.Errorf("configure app: %w", err)
	}

	return appConfig, nil
}
