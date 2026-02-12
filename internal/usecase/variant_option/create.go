package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateVariantOptionUsecase struct {
	variantOptionRepo    repository.VariantOptionRepository
	itemVariantGroupRepo repository.ItemVariantGroupRepository
	uuid                 ports.UUIDInterface
	context              context.Context
}

type CreateVariantOptionInput struct {
	ItemVariantGroupID string
	Name               string
	PriceDelta         int64
	IsDefault          bool
	Order              int
	IsActive           bool
}

type CreateVariantOptionOutput struct {
	ID string
}

func NewCreateVariantOptionUsecase(
	variantOptionRepo repository.VariantOptionRepository,
	itemVariantGroupRepo repository.ItemVariantGroupRepository,
	uuid ports.UUIDInterface,
	context context.Context,
) *CreateVariantOptionUsecase {
	return &CreateVariantOptionUsecase{
		variantOptionRepo:    variantOptionRepo,
		itemVariantGroupRepo: itemVariantGroupRepo,
		uuid:                 uuid,
		context:              context,
	}
}

func (uc *CreateVariantOptionUsecase) Execute(input CreateVariantOptionInput) (*CreateVariantOptionOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ItemVariantGroupID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid item variant group ID")
	}

	itemVariantGroup, err := uc.itemVariantGroupRepo.GetByID(uc.context, input.ItemVariantGroupID)
	if err != nil {
		return nil, err
	}

	variantOption := entity.VariantOption{
		ID:             uc.uuid.Generate(),
		VariantGroupID: itemVariantGroup.ID,
		Name:           input.Name,
		PriceDelta:     input.PriceDelta,
		IsDefault:      input.IsDefault,
		Order:          input.Order,
		IsActive:       input.IsActive,
	}

	err = uc.variantOptionRepo.Create(uc.context, &variantOption)
	if err != nil {
		return nil, err
	}

	return &CreateVariantOptionOutput{
		ID: variantOption.ID,
	}, nil
}
