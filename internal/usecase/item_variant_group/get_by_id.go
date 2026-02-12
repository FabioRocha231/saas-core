package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetItemVariantGroupByIDInput struct {
	ID string
}

type GetItemVariantGroupByIDOutput struct {
	ID             string
	CategoryItemID string

	Name      string
	Required  bool
	MinSelect int
	MaxSelect int
	Order     int
	IsActive  bool
}

type GetItemVariantGroupByIDUseCase struct {
	itemVariantGroupRepo repository.ItemVariantGroupRepository
	uuid                 ports.UUIDInterface
	context              context.Context
}

func NewGetItemVariantGroupByIDUseCase(itemVariantGroupRepo repository.ItemVariantGroupRepository, uuid ports.UUIDInterface, ctx context.Context) *GetItemVariantGroupByIDUseCase {
	return &GetItemVariantGroupByIDUseCase{
		itemVariantGroupRepo: itemVariantGroupRepo,
		uuid:                 uuid,
		context:              ctx,
	}
}

func (uc *GetItemVariantGroupByIDUseCase) Execute(input GetItemVariantGroupByIDInput) (*GetItemVariantGroupByIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid id")
	}

	itemVariantGroup, err := uc.itemVariantGroupRepo.GetByID(uc.context, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetItemVariantGroupByIDOutput{
		ID:             itemVariantGroup.ID,
		CategoryItemID: itemVariantGroup.CategoryItemID,

		Name:      itemVariantGroup.Name,
		Required:  itemVariantGroup.Required,
		MinSelect: itemVariantGroup.MinSelect,
		MaxSelect: itemVariantGroup.MaxSelect,
		Order:     itemVariantGroup.Order,
		IsActive:  itemVariantGroup.IsActive,
	}, nil
}
