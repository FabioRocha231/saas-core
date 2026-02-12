package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListVariantOptionsByItemVariantGroupIDUsecase struct {
	variantOptionRepo    repository.VariantOptionRepository
	itemVariantGroupRepo repository.ItemVariantGroupRepository
	uuid                 ports.UUIDInterface
	context              context.Context
}

type ListVariantOptionsByItemVariantGroupIDInput struct {
	ItemVariantGroupID string
}

type VariantOption struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PriceDelta int64  `json:"price_delta"`
	IsDefault  bool   `json:"is_default"`
	Order      int    `json:"order"`
	IsActive   bool   `json:"is_active"`
}

type ListVariantOptionsByItemVariantGroupIDOutput struct {
	VariantOptions []VariantOption `json:"variant_options"`
}

func NewListByItemVariantGroupIDUsecase(
	variantOptionRepo repository.VariantOptionRepository,
	itemVariantGroupRepo repository.ItemVariantGroupRepository,
	uuid ports.UUIDInterface,
	ctx context.Context,
) *ListVariantOptionsByItemVariantGroupIDUsecase {
	return &ListVariantOptionsByItemVariantGroupIDUsecase{
		variantOptionRepo:    variantOptionRepo,
		itemVariantGroupRepo: itemVariantGroupRepo,
		uuid:                 uuid,
		context:              ctx,
	}
}

func (uc *ListVariantOptionsByItemVariantGroupIDUsecase) Execute(input ListVariantOptionsByItemVariantGroupIDInput) (*ListVariantOptionsByItemVariantGroupIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ItemVariantGroupID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid item variant group id")
	}

	variantOptions, err := uc.variantOptionRepo.ListByVariantGroupID(uc.context, input.ItemVariantGroupID)
	if err != nil {
		return nil, err
	}

	var output VariantOption
	var outputList []VariantOption
	for _, variantOption := range variantOptions {
		output = VariantOption{
			ID:         variantOption.ID,
			Name:       variantOption.Name,
			PriceDelta: variantOption.PriceDelta,
			IsDefault:  variantOption.IsDefault,
			Order:      variantOption.Order,
			IsActive:   variantOption.IsActive,
		}
		outputList = append(outputList, output)
	}

	return &ListVariantOptionsByItemVariantGroupIDOutput{
		VariantOptions: outputList,
	}, nil
}
