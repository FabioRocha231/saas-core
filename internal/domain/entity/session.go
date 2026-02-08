package entity

import "time"

type Session struct {
	ID        string
	UserID    string
	Role      string
	ExpiresAt time.Time

	RevokedAt *time.Time
	CreatedAt time.Time
}
