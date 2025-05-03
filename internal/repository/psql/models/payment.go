package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID        int64           `db:"id,omitempty"`         // serial
	TransID   string          `db:"transaction_id"`       // TEXT
	UserID    int64           `db:"user_id"`              // INT REFERENCES Users(id)
	PlanID    int64           `db:"plan_id"`              // INT REFERENCES Plans(id)
	Amount    decimal.Decimal `db:"amount"`               // DECIMAL(10,2) CHECK > 0
	Status    string          `db:"status,omitempty"`     // TEXT CHECK IN ('pending', 'canceled', 'paid')
	Method    string          `db:"method"`               // TEXT
	CreatedAt time.Time       `db:"created_at,omitempty"` // timestamptz DEFAULT CURRENT_TIMESTAMP
	PaidAt    sql.NullTime    `db:"paid_at,omitempty"`    // может быть NULL
}
