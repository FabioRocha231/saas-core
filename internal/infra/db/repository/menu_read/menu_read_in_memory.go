package memorymenuread

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type Repo struct {
	items        repository.CategoryItemRepository
	addonGroups  repository.ItemAddonGroupRepository
	addonOptions repository.AddonOptionRepository
	varGroups    repository.ItemVariantGroupRepository
	varOptions   repository.VariantOptionRepository
}

func New(
	items repository.CategoryItemRepository,
	addonGroups repository.ItemAddonGroupRepository,
	addonOptions repository.AddonOptionRepository,
	varGroups repository.ItemVariantGroupRepository,
	varOptions repository.VariantOptionRepository,
) repository.MenuReadRepository {
	return &Repo{
		items:        items,
		addonGroups:  addonGroups,
		addonOptions: addonOptions,
		varGroups:    varGroups,
		varOptions:   varOptions,
	}
}

func (r *Repo) GetCategoryItemByID(ctx context.Context, id string) (*entity.CategoryItem, error) {
	return r.items.GetByID(ctx, id)
}

func (r *Repo) ListItemAddonGroupsByItemID(ctx context.Context, itemID string) ([]*entity.ItemAddonGroup, error) {
	return r.addonGroups.ListByCategoryItemID(ctx, itemID)
}

func (r *Repo) GetItemAddonGroupByID(ctx context.Context, id string) (*entity.ItemAddonGroup, error) {
	return r.addonGroups.GetByID(ctx, id)
}

func (r *Repo) GetAddonOptionByID(ctx context.Context, id string) (*entity.AddonOption, error) {
	return r.addonOptions.GetByID(ctx, id)
}

func (r *Repo) ListItemVariantGroupsByItemID(ctx context.Context, itemID string) ([]*entity.ItemVariantGroup, error) {
	return r.varGroups.ListByCategoryItemID(ctx, itemID)
}

func (r *Repo) GetItemVariantGroupByID(ctx context.Context, id string) (*entity.ItemVariantGroup, error) {
	return r.varGroups.GetByID(ctx, id)
}

func (r *Repo) GetVariantOptionByID(ctx context.Context, id string) (*entity.VariantOption, error) {
	return r.varOptions.GetByID(ctx, id)
}
