package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type MenuCategoryRepository interface {
	Create(ctx context.Context, c *entity.MenuCategory) error
	GetByID(ctx context.Context, id string) (*entity.MenuCategory, error)
	ListByMenuID(ctx context.Context, menuID string) ([]*entity.MenuCategory, error)
}
