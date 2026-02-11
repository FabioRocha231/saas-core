package entity

import "time"

type AddonOption struct {
	ID           string
	AddonGroupID string

	Name     string
	Price    int64 // centavos (preço do adicional)
	Order    int
	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
