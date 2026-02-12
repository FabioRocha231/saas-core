package entity

import "time"

type VariantOption struct {
	ID             string
	VariantGroupID string

	Name       string
	PriceDelta int64
	IsDefault  bool
	Order      int
	IsActive   bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
