package core

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type ZerologGORMLogger struct {
	log zerolog.Logger
	logger.Config
}

func (z *ZerologGORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	z.LogLevel = level
	return z
}

func (z *ZerologGORMLogger) Debug(ctx context.Context, s string, i ...interface{}) {
	z.log.Debug().Msgf(s, i...)
}

func (z *ZerologGORMLogger) Info(ctx context.Context, s string, i ...interface{}) {
	z.log.Info().Msgf(s, i...)
}

func (z *ZerologGORMLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	z.log.Warn().Msgf(s, i...)
}

func (z *ZerologGORMLogger) Error(ctx context.Context, s string, i ...interface{}) {
	z.log.Error().Msgf(s, i...)
}

func (z *ZerologGORMLogger) Fatal(ctx context.Context, s string, i ...interface{}) {
	z.log.Fatal().Msgf(s, i...)
}

func (z *ZerologGORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && z.LogLevel >= logger.Error:
		sql, rows := fc()
		z.log.Error().Err(err).Str("elapsed", elapsed.String()).Int64("rows", rows).Msg(sql)
	case elapsed > z.SlowThreshold && z.LogLevel >= logger.Warn:
		sql, rows := fc()
		z.log.Warn().Str("elapsed", elapsed.String()).Int64("rows", rows).Msg(sql)
	}
}

// Helper function to create a new logger
func NewZerologGORMLogger(debug bool, config logger.Config) logger.Interface {
	var zlog zerolog.Logger
	if debug {
		zlog = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	} else {
		zlog = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	return &ZerologGORMLogger{
		log:    zlog,
		Config: config,
	}
}
