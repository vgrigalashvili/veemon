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
	"github.com/vgrigalashvili/veemon/internal/api/rest"
	"github.com/vgrigalashvili/veemon/internal/api/rest/handler"
	"github.com/vgrigalashvili/veemon/internal/config"
	"github.com/vgrigalashvili/veemon/internal/domain"
	"github.com/vgrigalashvili/veemon/internal/repository"
	"github.com/vgrigalashvili/veemon/internal/service"
	"github.com/vgrigalashvili/veemon/internal/token"

	"github.com/vgrigalashvili/veemon/internal/mail"
	"github.com/vgrigalashvili/veemon/internal/worker"
	"golang.org/x/sync/errgroup"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// performs initialization and starts the Fiber API server with the given configuration.
func StartServer(ac config.AppConfig) {

	// initialize a new Fiber app with configuration.
	api := fiber.New(fiber.Config{
		AppName:       "veemon api v1.0.0",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "veemon",
		BodyLimit:     1 * 1024, // Limit request body size to 1KB.
	})

	// initialize middlewares.
	api.Use(
		logger.New(), // logger middleware for request logging.

		// CORS middleware for handling cross-origin resource sharing (CORS).
		cors.New(cors.Config{
			AllowOrigins: "*", // allow all origins.
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		}),

		// limiter middleware for handling rate limiting.
		limiter.New(limiter.Config{
			Max:        100,             // max 100 requests per minute.
			Expiration: 1 * time.Minute, // Expiration time for rate limiting.
		}),
		// TODO: middleware.RequestIDGenerator, // Custom middleware to generate request IDs.
		// TODO: middleware.ResponseDurationLogger, // Custom middleware to log response duration.
	)

	// initialize services with the necessary components.

	// database connection using GORM and PostgreSQL driver.
	db, err := gorm.Open(postgres.Open(ac.DatabaseURI), &gorm.Config{})
	if err != nil {
		log.Fatalf("[ERROR] Database connection error: %v", err)
	}
	setupDatabaseConnection(db)
	runMigrations(db)

	// REST handler.
	tokenMaker, err := token.NewPasetoMaker(ac.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("[FATAL] error while creating Paseto maker: %v", err)
	}
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

	// root context with cancellation.
	ctx, cancel := context.WithCancel(context.Background())
	waitGroup, ctx := errgroup.WithContext(ctx)

	// smtp mailer.
	mailer := mail.NewSMTPMailer(ac.MailerHost, ac.MailerPort, ac.MailerUserName, ac.MailerPassword, "veemon")

	// initialize task processor.

	// redis client.
	redisAddr := ac.RedisAddress
	log.Printf("[DEBUG] redis address: %s", redisAddr)

	runTaskProcessor(ctx, waitGroup, redisAddr, db, mailer)

	// start the Fiber server in a separate goroutine.
	waitGroup.Go(func() error {
		if err := api.Listen(ac.HttpPort); err != nil {
			log.Fatalf("[ERROR] Couldn't start server: %v", err)
			return err
		}
		return nil
	})

	handleGracefulShutdown(api, cancel, waitGroup)
}

// performs database connection pooling setup.
func setupDatabaseConnection(db *gorm.DB) {
	pgDB, err := db.DB()
	if err != nil {
		log.Fatalf("[ERROR] failed to get raw database connection: %v", err)
	}
	pgDB.SetMaxIdleConns(10)
	pgDB.SetMaxOpenConns(100)
	pgDB.SetConnMaxLifetime(time.Hour)
}

// performs database migrations.
func runMigrations(db *gorm.DB) {
	if err := db.AutoMigrate(&domain.User{}, &domain.VerifyEmail{}); err != nil {
		log.Fatalf("[ERROR] auto migrate failed: %v", err)
	}
}

// performs initialization of the REST handlers.
func initializeHandler(rh *rest.RestHandler) {
	handler.InitializeUserHandler(rh)
	handler.InitializeAuthHandler(rh)
}

// performs initialization of "redis task processor" and manages its lifecycle.
func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, redisAddr string, db *gorm.DB, mailer mail.EmailSender) {
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, db, mailer)

	waitGroup.Go(func() error {
		if err := taskProcessor.Start(); err != nil {
			log.Fatalf("[ERROR] failed to start task processor: %v", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Println("[INFO] graceful shutdown of task processor...")
		taskProcessor.Shutdown()
		log.Println("[INFO] task processor stopped.")
		return nil
	})
}

// performs the graceful shutdown of the server and its components.
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
		log.Fatalf("[ERROR] error during shutdown: %v", err)
	}
}
