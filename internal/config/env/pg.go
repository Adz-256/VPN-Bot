package env

import (
	"errors"
	"os"
)

const (
	dsnENV = "DSN"
)

var (
	ErrNoDSN = errors.New("no DSN provided")
)

type PGConfig struct {
	dns string
}

func NewPGConfig() (*PGConfig, error) {
	dsn := os.Getenv(dsnENV)
	if dsn == "" {
		return nil, ErrNoDSN
	}

	return &PGConfig{dsn}, nil
}

func (c *PGConfig) DSN() string {
	return c.dns
}
