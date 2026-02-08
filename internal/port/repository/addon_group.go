package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type AddonGroupRepository interface {
	Create(ctx context.Context, g *entity.AddonGroup) error
	GetByID(ctx context.Context, id string) (*entity.AddonGroup, error)
	ListByItemID(ctx context.Context, itemID string) ([]*entity.AddonGroup, error)
}
