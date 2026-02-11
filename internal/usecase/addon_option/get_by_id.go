package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetAddonOptionByIDUseCase struct {
	addonOptionRepo repository.AddonOptionRepository
	uuid            ports.UUIDInterface
	context         context.Context
}

type GetAddonOptionByIDInput struct {
	ID string
}

type GetAddonOptionByIDOutput struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Order     int       `json:"order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetAddonOptionByIDUsecase(
	addonOptionRepo repository.AddonOptionRepository,
	uuid ports.UUIDInterface,
	ctx context.Context,
) *GetAddonOptionByIDUseCase {
	return &GetAddonOptionByIDUseCase{
		addonOptionRepo: addonOptionRepo,
		uuid:            uuid,
		context:         ctx,
	}
}

func (uc *GetAddonOptionByIDUseCase) Execute(input GetAddonOptionByIDInput) (*GetAddonOptionByIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid addon option id")
	}

	addonOption, err := uc.addonOptionRepo.GetByID(uc.context, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetAddonOptionByIDOutput{
		ID:        addonOption.ID,
		GroupID:   addonOption.AddonGroupID,
		Name:      addonOption.Name,
		Price:     addonOption.Price,
		Order:     addonOption.Order,
		IsActive:  addonOption.IsActive,
		CreatedAt: addonOption.CreatedAt,
		UpdatedAt: addonOption.UpdatedAt,
	}, nil
}
