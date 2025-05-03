package env

import (
	"errors"
	"github.com/Adz-256/cheapVPN/internal/config"
	"os"
)

var (
	ErrNoLogLevel = errors.New("LOG_LEVEL environment variable not set")
)

type loggerConfig struct {
	level string
}

func NewLoggerConfig() (config.LoggerConfig, error) {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return nil, ErrNoLogLevel
	}

	return &loggerConfig{level: level}, nil
}

func (c loggerConfig) Level() string {
	return c.level
}
