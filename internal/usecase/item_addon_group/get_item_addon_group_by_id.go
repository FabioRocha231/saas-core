package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetItemAddonGroupByIDUseCase struct {
	itemAddonGroupRepo repository.ItemAddonGroupRepository
	uuid               ports.UUIDInterface
	context            context.Context
}

type GetItemAddonGroupByIDInput struct {
	ID string
}

type GetItemAddonGroupByIDOutput struct {
	ID             string
	CategoryItemID string

	Name      string
	Required  bool
	MinSelect int
	MaxSelect int
	Order     int
	IsActive  bool
}

func NewGetItemAddonGroupByIDUseCase(
	ctx context.Context,
	itemAddonGroupRepo repository.ItemAddonGroupRepository,
	uuid ports.UUIDInterface,
) *GetItemAddonGroupByIDUseCase {
	return &GetItemAddonGroupByIDUseCase{
		context:            ctx,
		itemAddonGroupRepo: itemAddonGroupRepo,
		uuid:               uuid,
	}
}

func (uc *GetItemAddonGroupByIDUseCase) Execute(input GetItemAddonGroupByIDInput) (*GetItemAddonGroupByIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid item addon group ID")
	}
	itemAddonGroup, err := uc.itemAddonGroupRepo.GetByID(uc.context, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetItemAddonGroupByIDOutput{
		ID:             itemAddonGroup.ID,
		CategoryItemID: itemAddonGroup.CategoryItemID,
		Name:           itemAddonGroup.Name,
		Required:       itemAddonGroup.Required,
		MinSelect:      itemAddonGroup.MinSelect,
		MaxSelect:      itemAddonGroup.MaxSelect,
		Order:          itemAddonGroup.Order,
		IsActive:       itemAddonGroup.IsActive,
	}, nil
}
