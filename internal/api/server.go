// Package api is responsible for setting up and starting the API server.
// It configures middleware, routes, and database connections.
package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/vgrigalashvili/veemon/config"
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/api/rest/handler"
	"github.com/vgrigalashvili/veemon/internal/api/rest/middleware"
	"github.com/vgrigalashvili/veemon/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// StartServer initializes and starts the Fiber API server with the given configuration.
func StartServer(appConfig config.AppConfig) {
	log.Println("[INFO] Starting server initialization")

	// Create a new Fiber app with configuration.
	api := fiber.New(fiber.Config{
		AppName:       "veemon api v1.0.0",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "veemon",
		BodyLimit:     1 * 1024, // Limit request body size to 1KB.
	})

	// Set up middleware for the app.
	api.Use(
		logger.New(), // Logger middleware for request logging.
		cors.New(cors.Config{
			AllowOrigins: "*", // Allow all origins.
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		}),
		limiter.New(limiter.Config{
			Max:        100,             // Max 100 requests per minute.
			Expiration: 1 * time.Minute, // Expiration time for rate limiting.
		}),
		middleware.ResponseDurationLogger, // Custom middleware to log response duration.
	)

	// Initialize the database connection using GORM and PostgreSQL driver.
	db, err := gorm.Open(postgres.Open(appConfig.DatabaseURI), &gorm.Config{})
	if err != nil {
		log.Fatalf("[ERROR] Database connection error: %v", err)
	}
	log.Println("[INFO] Database connection established")

	// Get the underlying SQL database connection for further configuration.
	pgDB, err := db.DB()
	if err != nil {
		log.Fatalf("[ERROR] Failed to get raw database connection: %v", err)
	}

	// Configure database connection pooling.
	pgDB.SetMaxIdleConns(10)
	pgDB.SetMaxOpenConns(100)
	pgDB.SetConnMaxLifetime(time.Hour)

	// Run database migrations for the User entity.
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("[ERROR] AutoMigrate failed: %v", err)
	}
	// Create a RestHandler with necessary components.
	restHandler := &rest.RestHandler{
		API: api,
		DB:  db,
		SEC: appConfig.TokenSymmetricKey,
	}
	// Initialize API handlers for user-related routes.
	initializeHandlers(restHandler)

	// Start the Fiber server in a separate goroutine.
	go func() {
		log.Printf("[INFO] Starting Fiber server on port %s", appConfig.HttpPort)
		if err := api.Listen(appConfig.HttpPort); err != nil {
			log.Fatalf("[ERROR] Couldn't start server: %v", err)
		}
	}()

	// Graceful shutdown handling.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // Wait for a termination signal.
	log.Println("[INFO] Shutting down server...")

	// Shutdown the server gracefully.
	if err := api.Shutdown(); err != nil {
		log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exited cleanly")
}

// initializeHandlers sets up API route handlers for the given RestHandler.
func initializeHandlers(rh *rest.RestHandler) {
	log.Println("[DEBUG] Initializing user handlers")
	handler.InitializeUserHandler(rh)
	log.Println("[INFO] User handlers initialized")
}
