package models

import "github.com/shopspring/decimal"

type Plan struct {
	ID           int64           `db:"id"`            // serial
	Country      string          `db:"country"`       // text UNIQUE NOT NULL
	DurationDays int             `db:"duration_days"` // integer CHECK > 0
	Price        decimal.Decimal `db:"price"`         // DECIMAL(10,2)
	Description  string          `db:"description"`   // TEXT
}
