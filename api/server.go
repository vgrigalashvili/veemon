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
	"github.com/jackc/pgx/v5"
	swagger "github.com/swaggo/fiber-swagger"
	"github.com/vgrigalashvili/veemon/api/rest"
	"github.com/vgrigalashvili/veemon/api/rest/handler"
	"github.com/vgrigalashvili/veemon/internal/config"
	_ "github.com/vgrigalashvili/veemon/internal/docs"
	"github.com/vgrigalashvili/veemon/pkg/mail"
	"github.com/vgrigalashvili/veemon/pkg/token"
	"github.com/vgrigalashvili/veemon/pkg/worker"
	"golang.org/x/sync/errgroup"

	db "github.com/vgrigalashvili/veemon/internal/repository/sqlc"
)

func StartServer(ac config.AppConfig) {

	api := fiber.New(fiber.Config{
		AppName:       "veemon api v1.0.0",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "veemon",
		BodyLimit:     1 * 1024,
	})
	api.Get("/swagger/*", swagger.WrapHandler)

	log.Printf("[INFO] Starting Fiber with config: AppName=%s, CaseSensitive=%v, StrictRouting=%v, BodyLimit=%d",
		api.Config().AppName, api.Config().CaseSensitive, api.Config().StrictRouting, api.Config().BodyLimit)

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

	ctx, cancel := context.WithCancel(context.Background())
	waitGroup, ctx := errgroup.WithContext(ctx)

	conn, err := pgx.Connect(ctx, ac.DatabaseURI)
	if err != nil {
		log.Fatalf("[ERROR] failed to connect to the database: %v", err)
	}
	log.Println("[INFO] database connection established successfully")
	defer conn.Close(ctx)

	queries := db.New(conn)

	tokenMaker, err := token.NewPasetoMaker(ac.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("[FATAL] error while creating Paseto maker: %v", err)
	}

	restHandler := &rest.RestHandler{
		API:     api,
		Token:   tokenMaker,
		Querier: queries,
	}
	initializeHandler(restHandler)

	mailer := mail.NewSMTPMailer(ac.MailerHost, ac.MailerPort, ac.MailerUserName, ac.MailerPassword, "veemon")

	redisAddr := ac.RedisAddress
	log.Printf("[DEBUG] redis address: %s", redisAddr)
	runTaskProcessor(ctx, waitGroup, redisAddr, queries, mailer)

	waitGroup.Go(func() error {
		if err := api.Listen(ac.HttpPort); err != nil {
			log.Fatalf("[ERROR] Couldn't start server: %v", err)
			return err
		}
		return nil
	})

	handleGracefulShutdown(api, cancel, waitGroup)
}

func initializeHandler(rh *rest.RestHandler) {
	handler.InitializeAuthHandler(rh)
	handler.InitializeUserHandler(rh)
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, redisAddr string, db *db.Queries, mailer mail.EmailSender) {
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
