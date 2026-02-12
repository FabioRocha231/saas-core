package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateItemVariantGroupInput struct {
	CategoryItemID string
	Name           string
	Price          int
	Order          int
	Required       bool
	IsActive       bool
	MinSelect      int
	MaxSelect      int
}

type CreateItemVariantGroupOutput struct {
	ID string
}

type CreateItemVariantGroupUseCase struct {
	categoryItemRepo     repository.CategoryItemRepository
	itemVariantGroupRepo repository.ItemVariantGroupRepository
	uuid                 ports.UUIDInterface
	context              context.Context
}

func NewCreateItemVariantGroupUseCase(
	ctx context.Context,
	itemVariantGroupRepo repository.ItemVariantGroupRepository,
	categoryItemRepo repository.CategoryItemRepository,
	uuid ports.UUIDInterface,
) *CreateItemVariantGroupUseCase {
	return &CreateItemVariantGroupUseCase{
		context:              ctx,
		itemVariantGroupRepo: itemVariantGroupRepo,
		categoryItemRepo:     categoryItemRepo,
		uuid:                 uuid,
	}
}

func (uc *CreateItemVariantGroupUseCase) Execute(input CreateItemVariantGroupInput) (*CreateItemVariantGroupOutput, error) {
	isValidUuid := uc.uuid.Validate(input.CategoryItemID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category item id")
	}

	if _, err := uc.categoryItemRepo.GetByID(uc.context, input.CategoryItemID); err != nil {
		return nil, err
	}

	id := uc.uuid.Generate()
	now := time.Now()

	itemVariantGroup := &entity.ItemVariantGroup{
		ID:             id,
		CategoryItemID: input.CategoryItemID,
		Name:           input.Name,
		Order:          input.Order,
		Required:       input.Required,
		IsActive:       input.IsActive,
		MinSelect:      input.MinSelect,
		MaxSelect:      input.MaxSelect,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := uc.itemVariantGroupRepo.Create(uc.context, itemVariantGroup); err != nil {
		return nil, err
	}

	return &CreateItemVariantGroupOutput{ID: id}, nil
}
