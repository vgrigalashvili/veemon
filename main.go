// Package main is the entry point for the Veemon application.
// It initializes the application configuration and starts the API server.
package main

import (
	"log"

	"github.com/vgrigalashvili/veemon/config"
	"github.com/vgrigalashvili/veemon/internal/api"
)

func main() {
	log.Println("[INFO] veemon entry point!")

	// Load application configuration from environment variables or configuration files.
	appConfig, err := config.SetupEnvironment()
	if err != nil {
		log.Fatalf("[ERROR] Could not set up environment: %v", err)
	}
	log.Println("[INFO] development environment ready to run!")

	// Start the API server with the loaded configuration.
	api.StartServer(appConfig)
	log.Println("[INFO] veemon application shutting down, falwell...")
}
