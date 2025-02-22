package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RedisLogger struct{}

func NewLogger() *RedisLogger {
	return &RedisLogger{}
}

func (logger *RedisLogger) Print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *RedisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (logger *RedisLogger) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}

func (logger *RedisLogger) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}

func (logger *RedisLogger) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}

func (logger *RedisLogger) Error(args ...interface{}) {
	logger.Print(zerolog.ErrorLevel, args...)
}

func (logger *RedisLogger) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)
}
