package main

import (
	"log"

	"github.com/Adz-256/cheapVPN/internal/api"
	"github.com/Adz-256/cheapVPN/internal/config/env"
	"github.com/Adz-256/cheapVPN/internal/logger/slog"
	"github.com/Adz-256/cheapVPN/internal/webhook/smee"
	"github.com/Adz-256/cheapVPN/pkg/clients/postgres"
)

func main() {
	cfg, err := env.New()
	if err != nil {
		log.Println("Error loading config:", err)
	}

	l := slog.NewLogger(cfg)

	pool := postgres.New(cfg)
	defer pool.Close()

	paymentsCh := make(chan map[string]any, 1024)

	a := api.New(pool, l, cfg, paymentsCh)

	go smee.New(cfg.WebhookAddress(), cfg.WebhookPort()).Run(paymentsCh)

	log.Fatal(a.Run())
}
