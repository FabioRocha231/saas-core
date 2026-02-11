package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetByItemAddonGroupIDUsecase struct {
	addonOptionRepo    repository.AddonOptionRepository
	itemAddonGroupRepo repository.ItemAddonGroupRepository
	uuid               ports.UUIDInterface
	context            context.Context
}

type GetByItemAddonGroupIDInput struct {
	ItemAddonGroupID string
}

type AddonOption struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Order     int       `json:"order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetByItemAddonGroupIDOutput struct {
	AddonOptions []AddonOption `json:"addon_options"`
}

func NewGetByItemAddonGroupIDUsecase(
	addonOptionRepo repository.AddonOptionRepository,
	itemAddonGroupRepo repository.ItemAddonGroupRepository,
	uuid ports.UUIDInterface,
	ctx context.Context,
) *GetByItemAddonGroupIDUsecase {
	return &GetByItemAddonGroupIDUsecase{
		addonOptionRepo:    addonOptionRepo,
		itemAddonGroupRepo: itemAddonGroupRepo,
		uuid:               uuid,
		context:            ctx,
	}
}

func (uc *GetByItemAddonGroupIDUsecase) Execute(input GetByItemAddonGroupIDInput) (*GetByItemAddonGroupIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ItemAddonGroupID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid item addon group id")
	}

	addonOptions, err := uc.addonOptionRepo.ListByGroupID(uc.context, input.ItemAddonGroupID)
	if err != nil {
		return nil, err
	}

	var output GetByItemAddonGroupIDOutput
	for _, addonOption := range addonOptions {
		output.AddonOptions = append(output.AddonOptions, AddonOption{
			ID:        addonOption.ID,
			GroupID:   addonOption.AddonGroupID,
			Name:      addonOption.Name,
			Price:     addonOption.Price,
			Order:     addonOption.Order,
			IsActive:  addonOption.IsActive,
			CreatedAt: addonOption.CreatedAt,
			UpdatedAt: addonOption.UpdatedAt,
		})
	}

	return &output, nil
}
