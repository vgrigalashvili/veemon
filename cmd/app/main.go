package main

import (
	"log"

	"github.com/vgrigalashvili/veemon/api"
	"github.com/vgrigalashvili/veemon/internal/config"
)

func main() {
	log.Println("[INFO] veemon entry point!")

	appConfig, err := config.SetupEnvironment()
	if err != nil {
		log.Fatalf("[ERROR] Could not set up environment: %v", err)
	}
	log.Println("[INFO] development environment ready to run!")

	api.StartServer(appConfig)
	log.Println("[INFO] veemon application shutting down, falwell...")
}
