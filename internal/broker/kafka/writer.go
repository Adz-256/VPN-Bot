package kafka

import (
	"context"

	"github.com/Adz-256/cheapVPN/internal/closer"
	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	k *kafka.Writer
}

func (b *Broker) NewWriter(topic string) Writer {
	k := kafka.NewWriter(kafka.WriterConfig{
		Brokers: b.brokers,
		Topic:   topic,
	})

	closer.Add(k.Close)

	return Writer{k: k}
}

func (w *Writer) Write(ctx context.Context, msg models.BrokerMessage) error {
	return w.k.WriteMessages(ctx, messageToKafka(msg))
}
