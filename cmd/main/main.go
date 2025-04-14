package main

import (
	"github.com/Adz-256/cheapVPN/internal/api"
	"github.com/Adz-256/cheapVPN/internal/config/env"
	"github.com/Adz-256/cheapVPN/internal/logger/slog"
	"github.com/Adz-256/cheapVPN/pkg/clients/postgres"
	"log"
)

func main() {
	cfg, err := env.New()
	if err != nil {
		log.Println("Error loading config:", err)
	}

	l := slog.NewLogger(cfg)

	pool := postgres.New(cfg)
	defer pool.Close()

	a := api.New(pool, l, cfg)

	log.Fatal(a.Run())
}
