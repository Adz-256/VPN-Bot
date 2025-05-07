package kafka

import (
	"github.com/Adz-256/cheapVPN/internal/config"
)

type Broker struct {
	brokers []string
}

func New(cfg config.KafkaConfig) *Broker {
	return &Broker{brokers: cfg.Brokers()}
}
