package kafka

import (
	"context"
	"github.com/Adz-256/cheapVPN/internal/broker"

	"github.com/Adz-256/cheapVPN/internal/closer"
	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/segmentio/kafka-go"
)

type Reader struct {
	k *kafka.Reader
}

func (b *Broker) NewReader(groupID string, topic string) broker.Consumer {
	k := kafka.NewReader(kafka.ReaderConfig{
		Brokers: b.brokers,
		GroupID: groupID,
		Topic:   topic,
	})

	k.Stats()

	closer.Add(k.Close)

	return &Reader{k: k}
}

func (r *Reader) Read(ctx context.Context) (*models.BrokerMessage, error) {
	msg, err := r.k.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	return kafkaToMessage(msg), nil
}
