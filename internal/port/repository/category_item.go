package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type CategoryItemRepository interface {
	Create(ctx context.Context, i *entity.CategoryItem) error
	GetByID(ctx context.Context, id string) (*entity.CategoryItem, error)
	ListByCategoryID(ctx context.Context, categoryID string) ([]*entity.CategoryItem, error)
}
