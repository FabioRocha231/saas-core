package entity

import "time"

type VariantOption struct {
	ID      string
	GroupID string

	Name       string
	PriceDelta int64 // soma no preço final (centavos), pode ser negativo
	IsDefault  bool
	Order      int
	IsActive   bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
