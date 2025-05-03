package env

import (
	"errors"
	"github.com/Adz-256/cheapVPN/internal/config"
	"os"
)

const (
	envToken = "BOT_TOKEN"
)

var (
	ErrNoToken = errors.New("no token provided")
)

type botConfig struct {
	token string
}

func NewBotConfig() (config.BotConfig, error) {
	token := os.Getenv(envToken)
	if token == "" {
		return nil, ErrNoToken
	}
	return &botConfig{token: token}, nil
}

func (cfg *botConfig) Token() string {
	return cfg.token
}
