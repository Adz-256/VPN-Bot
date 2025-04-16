package models

import "time"

type WgPeer struct {
	ID         int64     `db:"id"`          // serial
	UserID     int64     `db:"user_id"`     // REFERENCES users(id)
	PublicKey  string    `db:"public_key"`  // TEXT UNIQUE
	ConfigFile string    `db:"config_file"` // TEXT
	ServerIP   string    `db:"server_ip"`   // inet
	ProvidedIP string    `db:"provided_ip"` // inet
	CreatedAt  time.Time `db:"created_at"`  // timestamp
}
