package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateItemAddonGroupUseCase struct {
	itemAddonGroupRepo repository.ItemAddonGroupRepository
	categoryItemRepo   repository.CategoryItemRepository
	uuid               ports.UUIDInterface
	context            context.Context
}

type CreateItemAddonGroupInput struct {
	CategoryItemID string
	Name           string
	Required       bool
	MinSelect      int
	MaxSelect      int
	Order          int
	IsActive       bool
}

type CreateItemAddonGroupResponse struct {
	ID string
}

func NewCreateItemAddonGroupUseCase(
	ctx context.Context,
	itemAddonGroupRepo repository.ItemAddonGroupRepository,
	categoryItemRepo repository.CategoryItemRepository,
	uuid ports.UUIDInterface,
) *CreateItemAddonGroupUseCase {
	return &CreateItemAddonGroupUseCase{
		context:            ctx,
		itemAddonGroupRepo: itemAddonGroupRepo,
		categoryItemRepo:   categoryItemRepo,
		uuid:               uuid,
	}
}

func (uc *CreateItemAddonGroupUseCase) Execute(input CreateItemAddonGroupInput) (*CreateItemAddonGroupResponse, error) {
	isValidUuid := uc.uuid.Validate(input.CategoryItemID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category item ID")
	}

	_, err := uc.categoryItemRepo.GetByID(uc.context, input.CategoryItemID)
	if err != nil {
		return nil, err
	}

	itemAddonGroup := &entity.ItemAddonGroup{
		ID:             uc.uuid.Generate(),
		CategoryItemID: input.CategoryItemID,
		Name:           input.Name,
		Required:       input.Required,
		MinSelect:      input.MinSelect,
		MaxSelect:      input.MaxSelect,
		Order:          input.Order,
		IsActive:       input.IsActive,
	}

	err = uc.itemAddonGroupRepo.Create(uc.context, itemAddonGroup)
	if err != nil {
		return nil, err
	}

	return &CreateItemAddonGroupResponse{
		ID: itemAddonGroup.ID,
	}, nil
}
