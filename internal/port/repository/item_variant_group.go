package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type ItemVariantGroupRepository interface {
	Create(ctx context.Context, g *entity.ItemVariantGroup) error
	GetByID(ctx context.Context, id string) (*entity.ItemVariantGroup, error)
	ListByCategoryItemID(ctx context.Context, itemID string) ([]*entity.ItemVariantGroup, error)
}
