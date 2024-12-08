// Package api is responsible for setting up and starting the API server.
// It configures middleware, routes, and database connections.
package api

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hibiken/asynq"
	"github.com/vgrigalashvili/veemon/config"
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/api/rest/handler"
	"github.com/vgrigalashvili/veemon/internal/api/rest/middleware"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/mail"
	"github.com/vgrigalashvili/veemon/internal/worker"
	"golang.org/x/sync/errgroup"

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

	// Configure database connection pooling.
	configureDatabaseConnection(db)

	// Run database migrations for required entities.
	runMigrations(db)

	// Create a RestHandler with necessary components.
	restHandler := &rest.RestHandler{
		API: api,
		DB:  db,
		SEC: appConfig.TokenSymmetricKey,
	}
	initializeHandlers(restHandler)

	// Create a root context with cancellation.
	ctx, cancel := context.WithCancel(context.Background())
	waitGroup, ctx := errgroup.WithContext(ctx)

	// instantiate smtp mailer.
	mailer := mail.NewSMTPMailer("localhost", "port", "username", "password", "from")
	// Start the task processor.
	redisAddr := appConfig.RedisAddress
	log.Printf("[DEBUG] Redis address: %s", redisAddr)
	runTaskProcessor(ctx, waitGroup, redisAddr, db, mailer)

	// Start the Fiber server in a separate goroutine.
	waitGroup.Go(func() error {
		log.Printf("[INFO] Starting Fiber server on port %s", appConfig.HttpPort)
		if err := api.Listen(appConfig.HttpPort); err != nil {
			log.Fatalf("[ERROR] Couldn't start server: %v", err)
			return err
		}
		return nil
	})

	// Handle graceful shutdown.
	handleGracefulShutdown(api, cancel, waitGroup)

	log.Println("[INFO] Server exited cleanly")
}

// configureDatabaseConnection configures database connection pooling.
func configureDatabaseConnection(db *gorm.DB) {
	pgDB, err := db.DB()
	if err != nil {
		log.Fatalf("[ERROR] Failed to get raw database connection: %v", err)
	}

	pgDB.SetMaxIdleConns(10)
	pgDB.SetMaxOpenConns(100)
	pgDB.SetConnMaxLifetime(time.Hour)
}

// runMigrations performs database migrations.
func runMigrations(db *gorm.DB) {
	log.Println("[INFO] Running database migrations...")
	if err := db.AutoMigrate(&domain.User{}, &domain.VerifyEmail{}); err != nil {
		log.Fatalf("[ERROR] AutoMigrate failed: %v", err)
	}
	log.Println("[INFO] Database migrations completed")
}

// initializeHandlers sets up API route handlers for the given RestHandler.
func initializeHandlers(rh *rest.RestHandler) {
	log.Println("[DEBUG] Initializing user handlers")
	handler.InitializeUserHandler(rh)
	log.Println("[INFO] User handlers initialized")
}

// runTaskProcessor starts the Redis task processor and manages its lifecycle.
func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, redisAddr string, db *gorm.DB, mailer mail.EmailSender) {
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	// mailer := mail.NewSMTPMailer()

	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, db, mailer)

	waitGroup.Go(func() error {
		log.Println("[INFO] Starting task processor...")
		if err := taskProcessor.Start(); err != nil {
			log.Fatalf("[ERROR] Failed to start task processor: %v", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Println("[INFO] Graceful shutdown of task processor...")
		taskProcessor.Shutdown()
		log.Println("[INFO] Task processor stopped.")
		return nil
	})
}

// handleGracefulShutdown manages the graceful shutdown of the server and its components.
func handleGracefulShutdown(api *fiber.App, cancel context.CancelFunc, waitGroup *errgroup.Group) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("[INFO] Shutting down server...")
	cancel()

	if err := api.Shutdown(); err != nil {
		log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
	}

	if err := waitGroup.Wait(); err != nil {
		log.Fatalf("[ERROR] Error during shutdown: %v", err)
	}
}
