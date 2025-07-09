package config

import (
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(k *koanf.Koanf) *zerolog.Logger {
	var logger zerolog.Logger

	zerolog.TimeFieldFormat = time.RFC3339

	if k.String("app.env") == "production" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Logger()
	}

	log.Logger = logger
	return &logger

}
