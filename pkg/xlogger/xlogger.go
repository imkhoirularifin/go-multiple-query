package xlogger

import (
	"go-multiple-query/internal/config"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	Logger *zerolog.Logger
)

func Setup(cfg config.Config) {
	if cfg.IsDevelopment {
		l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
		Logger = &l
		return
	}
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	Logger = &l
}
