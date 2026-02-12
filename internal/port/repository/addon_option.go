package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type AddonOptionRepository interface {
	Create(ctx context.Context, o *entity.AddonOption) error
	GetByID(ctx context.Context, id string) (*entity.AddonOption, error)
	ListByAddonGroupID(ctx context.Context, groupID string) ([]*entity.AddonOption, error)
}
