package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/devex-on-k8s/book/appointments/config"
	"github.com/devex-on-k8s/book/appointments/db"
	"github.com/devex-on-k8s/book/appointments/web"
)

func main() {
	appConfig, err := config.New()
	if err != nil {
		log.Fatalf("Failed to create Fiber server: %v", err)
	}

	app := appConfig.App

	// cancel the startup and shutdown contexts in the main function, so that
	// the services are not left running after the main function returns.
	defer appConfig.StartupCancel()
	defer appConfig.ShutdownCancel()

	// add middlewares
	app.Use(logger.New())

	// create new server instance using a new database connection
	server := web.NewServer(db.New(config.DB))

	// add routes
	app.Get("/", server.Welcome)
	app.Get("/appointments", server.GetAllAppointments)
	app.Post("/appointments", server.CreateAppointment)
	app.Delete("/appointments", server.DeleteAllAppointments)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%v", config.APP_PORT)); err != nil {
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
