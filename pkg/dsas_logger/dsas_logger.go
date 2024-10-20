package dsas_logger

import (
	"log/slog"
	"os"
)

type CustomHandler struct {
	h     slog.Handler
	key   string
	value string
}

func SetDefaultSlog() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)

	slog.SetDefault(logger)
}
