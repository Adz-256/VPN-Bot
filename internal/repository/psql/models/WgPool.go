package models

import (
	"net"
	"time"
)

type WgPeer struct {
	ID         int64     `db:"id,omitempty"`         // serial
	UserID     int64     `db:"user_id"`              // REFERENCES users(id)
	PublicKey  string    `db:"public_key"`           // TEXT UNIQUE
	ConfigFile string    `db:"config_file"`          // TEXT
	ServerIP   net.IPNet `db:"server_ip"`            // inet
	ProvidedIP net.IPNet `db:"provided_ip"`          // inet
	CreatedAt  time.Time `db:"created_at,omitempty"` // timestamp
	EndAt      time.Time `db:"end_at,omitempty"`     // timestamp
}
