package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

// strcuts to interact with service layer
type Payment struct {
	ID        int64
	TransID   string
	UserID    int64
	PlanID    int64
	Amount    decimal.Decimal
	Status    string
	Method    string
	CreatedAt time.Time
	PaidAt    sql.NullTime
}

type Plan struct {
	ID           int64
	Country      string
	DurationDays int
	Price        decimal.Decimal
	Description  string
}

type User struct {
	ID        int64
	UserID    int64
	Username  string
	IsAdmin   bool
	CreatedAt time.Time
}

type WgPeer struct {
	ID         int64
	Name       string
	UserID     int64
	PublicKey  string
	ConfigFile string
	ServerIP   string
	ProvidedIP string
	CreatedAt  time.Time
	EndAt      time.Time
	Blocked    bool
}

type BrokerMessage struct {
	Topic  string
	Value  []byte
	Offset int64
	Time   time.Time
}
