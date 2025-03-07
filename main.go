package main

import (
	"log"

	"github.com/vgrigalashvili/veemon/api"
	"github.com/vgrigalashvili/veemon/internal/config"
	"github.com/vgrigalashvili/veemon/pkg/mqtt"
)

// @title			veemon API
// @version		1.0
// @description	This is the API for the Veemon application.
// @host			localhost:3000
func main() {
	log.Println("[INFO] veemon entry point!")

	appConfig, err := config.SetupEnvironment()
	if err != nil {
		log.Fatalf("[ERROR] Could not set up environment: %v", err)
	}
	log.Println("[INFO] development environment ready to run!")

	// Start MQTT client and subscribe to the heartbeat topic concurrently.
	go func() {
		// Use "tcp://localhost:1883" if you have mapped the container's port 1883 to localhost.
		// If you run this inside Docker (or via Docker network), you might use "tcp://rabbitmq:1883".
		brokerURL := "tcp://localhost:1883"
		clientID := "veemon-client"
		mqtt.Connect(brokerURL, clientID)

		// Subscribe to heartbeat messages.
		heartbeatTopic := "Lift/+/events/heartbeat"
		mqtt.SubscribeHeartbeat(heartbeatTopic)
	}()

	api.StartServer(appConfig)
	log.Println("[INFO] veemon application shutting down, falwell...")
}
