package worker

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/vgrigalashvili/veemon/internal/repository/sqlc"
	"github.com/vgrigalashvili/veemon/pkg/mail"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	db     *db.Queries
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db *db.Queries, mailer mail.EmailSender) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},

			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	redisTaskProcessor := &RedisTaskProcessor{
		server: server,
		db:     db,
		mailer: mailer,
	}

	return redisTaskProcessor
}

func (rtp *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, rtp.ProcessTaskSendVerifyEmail)

	return rtp.server.Start(mux)
}

func (rtp *RedisTaskProcessor) Shutdown() {
	rtp.server.Shutdown()
}
