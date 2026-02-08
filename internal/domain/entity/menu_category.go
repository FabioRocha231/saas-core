package entity

import "time"

type MenuCategory struct {
	ID       string
	MenuID   string
	Name     string
	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
