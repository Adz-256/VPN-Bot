package models

import "time"

type User struct {
	ID        int64     `db:"id"`         // serial
	ChatID    int64     `db:"chat_id"`    // BIGINT UNIQUE
	Username  string    `db:"username"`   // TEXT
	IsAdmin   bool      `db:"is_admin"`   // BOOL
	CreatedAt time.Time `db:"created_at"` // timestamptz
}
