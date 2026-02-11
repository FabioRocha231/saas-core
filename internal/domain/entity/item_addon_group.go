package entity

import "time"

type ItemAddonGroup struct {
	ID             string
	CategoryItemID string

	Name      string
	Required  bool
	MinSelect int
	MaxSelect int
	Order     int
	IsActive  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
