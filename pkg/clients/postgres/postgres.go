package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type DBConfig interface {
	DSN() string
}

func New(cfg DBConfig) *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
