package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type VariantOptionRepository interface {
	Create(ctx context.Context, o *entity.VariantOption) error
	GetByID(ctx context.Context, id string) (*entity.VariantOption, error)
	ListByVariantGroupID(ctx context.Context, groupID string) ([]*entity.VariantOption, error)
}
