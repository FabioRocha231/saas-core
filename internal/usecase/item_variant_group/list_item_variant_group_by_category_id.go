package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListItemVariantGroupByCategoryItemIDUsecase struct {
	itemVariantGroupRepo repository.ItemVariantGroupRepository
	categoryItemRepo     repository.CategoryItemRepository
	uuid                 ports.UUIDInterface
	context              context.Context
}

type ListItemVariantGroupByCategoryItemIDInput struct {
	CategoryItemID string
}

type VariantGroup struct {
	ID             string
	CategoryItemID string

	Name      string
	Required  bool
	MinSelect int
	MaxSelect int
	Order     int
	IsActive  bool
}

type ListItemVariantGroupByCategoryItemIDOutput struct {
	VariantGroups []VariantGroup `json:"variant_groups"`
}

func NewListItemVariantGroupByCategoryItemIDUsecase(
	ctx context.Context,
	itemVariantGroupRepo repository.ItemVariantGroupRepository,
	categoryItemRepo repository.CategoryItemRepository,
	uuid ports.UUIDInterface,
) *ListItemVariantGroupByCategoryItemIDUsecase {
	return &ListItemVariantGroupByCategoryItemIDUsecase{
		context:              ctx,
		itemVariantGroupRepo: itemVariantGroupRepo,
		categoryItemRepo:     categoryItemRepo,
		uuid:                 uuid,
	}
}

func (uc *ListItemVariantGroupByCategoryItemIDUsecase) Execute(input ListItemVariantGroupByCategoryItemIDInput) (*ListItemVariantGroupByCategoryItemIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.CategoryItemID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category item ID")
	}

	_, err := uc.categoryItemRepo.GetByID(uc.context, input.CategoryItemID)
	if err != nil {
		return nil, err
	}

	variantGroupsEntities, err := uc.itemVariantGroupRepo.ListByCategoryItemID(uc.context, input.CategoryItemID)
	if err != nil {
		return nil, err
	}

	var variantGroups []VariantGroup
	for _, v := range variantGroupsEntities {
		variantGroups = append(variantGroups, VariantGroup{
			ID:             v.ID,
			CategoryItemID: v.CategoryItemID,

			Name:      v.Name,
			Required:  v.Required,
			MinSelect: v.MinSelect,
			MaxSelect: v.MaxSelect,
			Order:     v.Order,
			IsActive:  v.IsActive,
		})
	}

	return &ListItemVariantGroupByCategoryItemIDOutput{
		VariantGroups: variantGroups,
	}, nil
}
