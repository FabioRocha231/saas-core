package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListItemAddonGroupByCategoryItemIDUseCase struct {
	itemAddonGroupRepo repository.ItemAddonGroupRepository
	categoryItemRepo   repository.CategoryItemRepository
	uuid               ports.UUIDInterface
	context            context.Context
}

type ListItemAddonGroupByCategoryItemIDInput struct {
	CategoryItemID string
}

type AddonGroup struct {
	ID             string
	CategoryItemID string

	Name      string
	Required  bool
	MinSelect int
	MaxSelect int
	Order     int
	IsActive  bool
}

type ListItemAddonGroupByCategoryItemIDOutput struct {
	AddonGroups []AddonGroup `json:"addon_groups"`
}

func NewListItemAddonGroupByCategoryItemIDUseCase(
	ctx context.Context,
	itemAddonGroupRepo repository.ItemAddonGroupRepository,
	categoryItemRepo repository.CategoryItemRepository,
	uuid ports.UUIDInterface,
) *ListItemAddonGroupByCategoryItemIDUseCase {
	return &ListItemAddonGroupByCategoryItemIDUseCase{
		context:            ctx,
		itemAddonGroupRepo: itemAddonGroupRepo,
		categoryItemRepo:   categoryItemRepo,
		uuid:               uuid,
	}
}

func (uc *ListItemAddonGroupByCategoryItemIDUseCase) Execute(input ListItemAddonGroupByCategoryItemIDInput) (*ListItemAddonGroupByCategoryItemIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.CategoryItemID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category item ID")
	}

	_, err := uc.categoryItemRepo.GetByID(uc.context, input.CategoryItemID)
	if err != nil {
		return nil, err
	}

	itemAddonGroups, err := uc.itemAddonGroupRepo.ListByCategoryItemID(uc.context, input.CategoryItemID)
	if err != nil {
		return nil, err
	}

	addonGroups := make([]AddonGroup, len(itemAddonGroups))
	for i, itemAddonGroup := range itemAddonGroups {
		addonGroups[i] = AddonGroup{
			ID:             itemAddonGroup.ID,
			CategoryItemID: itemAddonGroup.CategoryItemID,
			Name:           itemAddonGroup.Name,
			Required:       itemAddonGroup.Required,
			MinSelect:      itemAddonGroup.MinSelect,
			MaxSelect:      itemAddonGroup.MaxSelect,
			Order:          itemAddonGroup.Order,
			IsActive:       itemAddonGroup.IsActive,
		}
	}

	return &ListItemAddonGroupByCategoryItemIDOutput{
		AddonGroups: addonGroups,
	}, nil
}
