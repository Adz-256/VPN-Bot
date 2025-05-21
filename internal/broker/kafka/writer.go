package kafka

import (
	"context"
	"github.com/Adz-256/cheapVPN/internal/broker"

	"github.com/Adz-256/cheapVPN/internal/closer"
	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	k *kafka.Writer
}

func (b *Broker) NewWriter(topic string) broker.Publisher {
	k := &kafka.Writer{
		Addr:         kafka.TCP(b.brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
	}
	closer.Add(k.Close)
	
	return &Writer{k: k}
}

func (w *Writer) Write(ctx context.Context, msg models.BrokerMessage) error {
	return w.k.WriteMessages(ctx, messageToKafka(msg))
}
