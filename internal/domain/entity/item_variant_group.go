package entity

import "time"

type ItemVariantGroup struct {
	ID     string
	ItemID string

	Name      string
	Required  bool
	MinSelect int
	MaxSelect int
	Order     int
	IsActive  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
