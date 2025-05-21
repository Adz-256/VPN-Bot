package kafka

import (
	"fmt"
	"github.com/Adz-256/cheapVPN/internal/config"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	brokers []string
}

func New(cfg config.BrokerConfig) (*Broker, error) {
	b := cfg.Brokers()
	if len(b) < 0 {
		return nil, fmt.Errorf("brokers is empty")
	}

	if err := pingKafka(b); err != nil {
		return nil, fmt.Errorf("ping kafka failed %w", err)
	}

	return &Broker{brokers: cfg.Brokers()}, nil
}

func pingKafka(brokers []string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	defer conn.Close()
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	_, err = conn.Brokers()
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	return nil // OK
}
