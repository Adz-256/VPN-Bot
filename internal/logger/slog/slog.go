package slog

import (
	"log/slog"
	"os"
)

type LoggerConfig interface {
	Level() string
}

func NewDefaultLogger(cfg LoggerConfig) {
	var handler slog.Handler

	switch cfg.Level() {
	case "debug":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
