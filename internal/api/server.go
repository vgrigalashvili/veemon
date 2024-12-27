// Package api is responsible for setting up and starting the API server.
// It provides middleware configuration, route handling, database integration,
// and task processing for the Veemon application.

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
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/api/rest/handler"
	"github.com/vgrigalashvili/veemon/internal/config"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/mail"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
	"github.com/vgrigalashvili/veemon/internal/token"
	"github.com/vgrigalashvili/veemon/internal/worker"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// StartServer initializes and starts the API server.
// It sets up middleware, routes, database connections, and task processing.
func StartServer(ac config.AppConfig) {

	api := fiber.New(fiber.Config{
		AppName:       "veemon api v1.0.0",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "veemon",
		BodyLimit:     1 * 1024,
	})

	// Log configuration details for debugging.
	log.Printf("[INFO] Starting Fiber with config: AppName=%s, CaseSensitive=%v, StrictRouting=%v, BodyLimit=%d",
		api.Config().AppName, api.Config().CaseSensitive, api.Config().StrictRouting, api.Config().BodyLimit)

	// Configure middleware for the API server.
	api.Use(
		logger.New(),
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		}),
		limiter.New(limiter.Config{
			Max:        100,
			Expiration: 1 * time.Minute,
		}),
	)

	// Set up the database connection.
	db, err := gorm.Open(postgres.Open(ac.DatabaseURI), &gorm.Config{})
	if err != nil {
		log.Fatalf("[ERROR] Database connection error: %v", err)
	}
	setupDatabaseConnection(db)
	runMigrations(db)

	// Initialize the token maker.
	tokenMaker, err := token.NewPasetoMaker(ac.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("[FATAL] error while creating Paseto maker: %v", err)
	}

	// Initialize services and handlers.
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(tokenMaker, userService)

	restHandler := &rest.RestHandler{
		API:         api,
		Token:       tokenMaker,
		AuthService: authService,
		UserService: userService,
		DB:          db,
	}
	initializeHandler(restHandler)

	// Set up context and error group for managing goroutines.
	ctx, cancel := context.WithCancel(context.Background())
	waitGroup, ctx := errgroup.WithContext(ctx)

	// Initialize the mailer.
	mailer := mail.NewSMTPMailer(ac.MailerHost, ac.MailerPort, ac.MailerUserName, ac.MailerPassword, "veemon")

	// Set up the task processor.
	redisAddr := ac.RedisAddress
	log.Printf("[DEBUG] redis address: %s", redisAddr)
	runTaskProcessor(ctx, waitGroup, redisAddr, db, mailer)

	// Start the API server.
	waitGroup.Go(func() error {
		if err := api.Listen(ac.HttpPort); err != nil {
			log.Fatalf("[ERROR] Couldn't start server: %v", err)
			return err
		}
		return nil
	})

	// Handle graceful shutdown.
	handleGracefulShutdown(api, cancel, waitGroup)
}

// setupDatabaseConnection configures the database connection pool settings.
func setupDatabaseConnection(db *gorm.DB) {
	pgDB, err := db.DB()
	if err != nil {
		log.Fatalf("[ERROR] failed to get raw database connection: %v", err)
	}
	pgDB.SetMaxIdleConns(10)
	pgDB.SetMaxOpenConns(100)
	pgDB.SetConnMaxLifetime(time.Hour)
}

// runMigrations runs the database migrations to set up necessary schemas.
func runMigrations(db *gorm.DB) {
	if err := db.AutoMigrate(&domain.User{}, &domain.VerifyEmail{}); err != nil {
		log.Fatalf("[ERROR] auto migrate failed: %v", err)
	}
}

// initializeHandler sets up the REST API handlers for the application.
func initializeHandler(rh *rest.RestHandler) {
	handler.InitializeUserHandler(rh)
	handler.InitializeAuthHandler(rh)
}

// runTaskProcessor starts the task processor for handling background tasks.
func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, redisAddr string, db *gorm.DB, mailer mail.EmailSender) {
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, db, mailer)

	// Start the task processor.
	waitGroup.Go(func() error {
		if err := taskProcessor.Start(); err != nil {
			log.Fatalf("[ERROR] failed to start task processor: %v", err)
			return err
		}
		return nil
	})

	// Handle task processor shutdown.
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Println("[INFO] graceful shutdown of task processor...")
		taskProcessor.Shutdown()
		log.Println("[INFO] task processor stopped.")
		return nil
	})
}

// handleGracefulShutdown manages graceful shutdown of the API server and other resources.
func handleGracefulShutdown(api *fiber.App, cancel context.CancelFunc, waitGroup *errgroup.Group) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("[INFO] shutting down server...")
	cancel()

	if err := api.Shutdown(); err != nil {
		log.Fatalf("[ERROR] server forced to shutdown: %v", err)
	}

	if err := waitGroup.Wait(); err != nil {
		log.Printf("[WARN] shutdown completed with errors: %v", err)
	} else {
		log.Println("[INFO] graceful shutdown completed successfully.")
	}
}
