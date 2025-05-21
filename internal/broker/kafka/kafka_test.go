package kafka

import (
	"context"
	"github.com/Adz-256/cheapVPN/internal/config/env"
	"github.com/Adz-256/cheapVPN/internal/models"
	"os"
	"testing"
	"time"
)

func TestReaderWriter(t *testing.T) {
	os.Setenv("PAYMENT_BROKERS", "localhost:9092")
	bs, err := env.NewPaymentBrokerConfig()
	if err != nil {
		t.Fatal(err)
	}
	broker, err := New(bs)
	if err != nil {
		t.Fatal(err)
	}

	r := broker.NewReader("payment_group", "payments")
	w := broker.NewWriter("payments")
	err = w.Write(context.Background(), models.BrokerMessage{
		Value: []byte("test"),
	})

	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Millisecond)
	msg, err := r.Read(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}
