package broker

import (
	"context"

	"github.com/Adz-256/cheapVPN/internal/models"
)

type Consumer interface {
	Read(ctx context.Context) (*models.BrokerMessage, error)
}

type Publisher interface {
	Write(ctx context.Context, msg models.BrokerMessage) error
}

type Broker interface {
	NewReader(groupID string, topic string) Consumer
	NewWriter(topic string) Publisher
}
