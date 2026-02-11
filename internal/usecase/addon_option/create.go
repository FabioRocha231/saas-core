package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateAddonOptionUseCase struct {
	addonOptionRepo    repository.AddonOptionRepository
	itemAddonGroupRepo repository.ItemAddonGroupRepository
	uuid               ports.UUIDInterface
	context            context.Context
}

type CreateAddonOptionInput struct {
	ItemAddonGroupID string
	Name             string
	Price            int64
	Order            int
	IsActive         bool
}

type CreateAddonOptionOutput struct {
	ID string
}

func NewCreateAddonOptionUseCase(
	addonOptionRepo repository.AddonOptionRepository,
	itemAddonGroupRepo repository.ItemAddonGroupRepository,
	uuid ports.UUIDInterface,
	ctx context.Context,
) *CreateAddonOptionUseCase {
	return &CreateAddonOptionUseCase{
		addonOptionRepo:    addonOptionRepo,
		itemAddonGroupRepo: itemAddonGroupRepo,
		uuid:               uuid,
		context:            ctx,
	}
}

func (uc *CreateAddonOptionUseCase) Execute(input CreateAddonOptionInput) (*CreateAddonOptionOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ItemAddonGroupID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid item addon group id")
	}

	addonGroup, err := uc.itemAddonGroupRepo.GetByID(uc.context, input.ItemAddonGroupID)
	if err != nil {
		return nil, errx.New(errx.CodeNotFound, "item addon group not found")
	}

	addonOption := &entity.AddonOption{
		ID:       uc.uuid.Generate(),
		AddonGroupID:  addonGroup.ID,
		Name:     input.Name,
		Price:    input.Price,
		Order:    input.Order,
		IsActive: input.IsActive,
	}

	err = uc.addonOptionRepo.Create(uc.context, addonOption)
	if err != nil {
		return nil, err
	}

	return &CreateAddonOptionOutput{ID: addonOption.ID}, nil
}
