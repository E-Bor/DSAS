package dsas_logger

import (
	"log/slog"
	"os"
)

type LogLvl string

const (
	EnvLocal LogLvl = "LOCAL"
	EnvDev   LogLvl = "DEV"
	EnvProd  LogLvl = "PROD"
)

type CustomHandler struct {
	h     slog.Handler
	key   string
	value string
}

func SetDefaultSlog(logLvl LogLvl) {
	var logger *slog.Logger

	switch logLvl {
	case EnvLocal:
		logger = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case EnvDev:
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case EnvProd:
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelInfo,
					AddSource: true,
				},
			),
		)
	}
	slog.SetDefault(logger)
}
