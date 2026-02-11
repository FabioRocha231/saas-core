package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type ItemAddonGroupRepository interface {
	Create(ctx context.Context, g *entity.ItemAddonGroup) error
	GetByID(ctx context.Context, id string) (*entity.ItemAddonGroup, error)
	ListByCategoryItemID(ctx context.Context, itemID string) ([]*entity.ItemAddonGroup, error)
}
