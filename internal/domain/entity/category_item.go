package entity

import "time"

type CategoryItem struct {
	ID         string
	CategoryID string

	Name        string
	Description string

	BasePrice int64 // centavos (pode ser 0 se o preço vier só por variação)
	ImageURL  string

	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
