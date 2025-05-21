package env

import (
	"errors"
	"os"
	"strings"

	"github.com/Adz-256/cheapVPN/internal/config"
)

var (
	ErrNoBrokerInstance = errors.New("no Broker Instance provided")
)

type brokerConfig struct {
	brokers []string
}

func NewPaymentBrokerConfig() (config.BrokerConfig, error) {
	b := os.Getenv("PAYMENT_BROKERS")
	if b == "" {
		return nil, ErrNoBrokerInstance
	}

	brokers := strings.Split(b, ",")

	return &brokerConfig{brokers: brokers}, nil
}

func (cfg *brokerConfig) Brokers() []string {
	return cfg.brokers
}
