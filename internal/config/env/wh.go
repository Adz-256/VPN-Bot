package env

import (
	"errors"
	"github.com/Adz-256/cheapVPN/internal/config"
	"os"
)

const (
	whAddressEnv = "WEBHOOK_ADDRESS"
	whPortEnv    = "WEBHOOK_PORT"
)

var (
	ErrNoWHAddress = errors.New("webhook address not set")
	ErrNoWHPort    = errors.New("webhook port not set")
)

type webhookConfig struct {
	addr string
	port string
}

func NewWebhookConfig() (config.WhConfig, error) {
	addr := os.Getenv(whAddressEnv)
	if addr == "" {
		return nil, ErrNoWHAddress
	}
	port := os.Getenv(whPortEnv)
	if port == "" {
		return nil, ErrNoWHPort
	}

	return &webhookConfig{addr: addr, port: port}, nil
}

func (c *webhookConfig) Address() string {
	return c.addr
}

func (c *webhookConfig) Port() string {
	return c.port
}
