package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetVariantOptionByIDUsecase struct {
	variantOptionRepo repository.VariantOptionRepository
	uuid              ports.UUIDInterface
	context           context.Context
}

type GetVariantOptionByIDInput struct {
	ID string
}

type GetVariantOptionByIDOutput struct {
	ID             string
	VariantGroupID string
	Name           string
	PriceDelta     int64
	IsDefault      bool
	Order          int
	IsActive       bool
}

func NewGetVariantOptionByIDUsecase(
	variantOptionRepo repository.VariantOptionRepository,
	uuid ports.UUIDInterface,
	context context.Context,
) *GetVariantOptionByIDUsecase {
	return &GetVariantOptionByIDUsecase{
		variantOptionRepo: variantOptionRepo,
		uuid:              uuid,
		context:           context,
	}
}

func (uc *GetVariantOptionByIDUsecase) Execute(input GetVariantOptionByIDInput) (*GetVariantOptionByIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid variant option ID")
	}

	variantOption, err := uc.variantOptionRepo.GetByID(uc.context, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetVariantOptionByIDOutput{
		ID:             variantOption.ID,
		VariantGroupID: variantOption.VariantGroupID,
		Name:           variantOption.Name,
		PriceDelta:     variantOption.PriceDelta,
		IsDefault:      variantOption.IsDefault,
		Order:          variantOption.Order,
		IsActive:       variantOption.IsActive,
	}, nil
}
