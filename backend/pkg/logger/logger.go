package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(env string) {
	if env == "production" {
		log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	}
}

func Get() *zerolog.Logger {
	return &log
}

func FromCtx(ctx context.Context) *zerolog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zerolog.Logger); ok {
		return l
	}
	return &log
}

type ctxKey struct{}
