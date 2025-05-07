package kafka

import (
	"github.com/Adz-256/cheapVPN/internal/models"
	"github.com/segmentio/kafka-go"
)

func kafkaToMessage(msg kafka.Message) *models.BrokerMessage {
	return &models.BrokerMessage{
		Topic:  msg.Topic,
		Value:  msg.Value,
		Offset: msg.Offset,
		Time:   msg.Time,
	}
}

func messageToKafka(msg models.BrokerMessage) kafka.Message {
	return kafka.Message{
		Topic:  msg.Topic,
		Value:  msg.Value,
		Offset: msg.Offset,
		Time:   msg.Time,
	}
}
