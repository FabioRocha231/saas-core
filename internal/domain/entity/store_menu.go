package entity

import "time"

type StoreMenu struct {
	ID       string
	StoreID  string
	Name     string
	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
