package env

import (
	"errors"
	"github.com/Adz-256/cheapVPN/internal/config"
	"os"
)

const (
	paymentAccountIDEnv = "PAYMENT_ACCOUNT"
)

var (
	ErrNoPaymentAccount = errors.New("no payment account found")
)

type paymentConfig struct {
	accountID string
}

func NewPaymentConfig() (config.PaymentConfig, error) {
	acc := os.Getenv(paymentAccountIDEnv)
	if acc == "" {
		return nil, ErrNoPaymentAccount
	}

	return &paymentConfig{accountID: acc}, nil
}

func (cfg *paymentConfig) AccountID() string {
	return cfg.accountID
}
