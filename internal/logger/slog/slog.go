package slog

import (
	"log/slog"
	"os"
)

type LoggerConfig interface {
	ENV() string
}

func NewLogger(cfg LoggerConfig) *slog.Logger {
	var handler slog.Handler

	switch cfg.ENV() {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	return logger
}
