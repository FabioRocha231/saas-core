package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type MenuReadRepository interface {
	GetCategoryItemByID(ctx context.Context, id string) (*entity.CategoryItem, error)

	ListItemAddonGroupsByItemID(ctx context.Context, itemID string) ([]*entity.ItemAddonGroup, error)
	GetItemAddonGroupByID(ctx context.Context, id string) (*entity.ItemAddonGroup, error)
	GetAddonOptionByID(ctx context.Context, id string) (*entity.AddonOption, error)

	ListItemVariantGroupsByItemID(ctx context.Context, itemID string) ([]*entity.ItemVariantGroup, error)
	GetItemVariantGroupByID(ctx context.Context, id string) (*entity.ItemVariantGroup, error)
	GetVariantOptionByID(ctx context.Context, id string) (*entity.VariantOption, error)
}
