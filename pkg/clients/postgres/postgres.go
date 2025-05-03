package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type DBConfig interface {
	DSN() string
}

func New(ctx context.Context, cfg DBConfig) *pgxpool.Pool {

	pool, err := pgxpool.Connect(ctx, cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
