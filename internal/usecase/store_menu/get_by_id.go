package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetStoreMenuByIDUsecase struct {
	storeMenuRepository repository.StoreMenuRepository
	context             context.Context
	uuid                ports.UUIDInterface
}

type GetStoreMenuByIDInput struct {
	StoreMenuID string
}

type GetStoreMenuByIDOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StoreID   string    `json:"storeId"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewGetStoreMenuByIDUsecase(
	storeMenuRepository repository.StoreMenuRepository,
	ctx context.Context,
	uuid ports.UUIDInterface,
) *GetStoreMenuByIDUsecase {
	return &GetStoreMenuByIDUsecase{
		storeMenuRepository: storeMenuRepository,
		context:             ctx,
		uuid:                uuid,
	}
}

func (uc *GetStoreMenuByIDUsecase) Execute(input GetStoreMenuByIDInput) (*GetStoreMenuByIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.StoreMenuID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid storeMenuId")
	}

	storeMenu, err := uc.storeMenuRepository.GetByID(uc.context, input.StoreMenuID)
	if err != nil {
		return nil, err
	}

	return &GetStoreMenuByIDOutput{
		ID:        storeMenu.ID,
		Name:      storeMenu.Name,
		StoreID:   storeMenu.StoreID,
		IsActive:  storeMenu.IsActive,
		CreatedAt: storeMenu.CreatedAt,
		UpdatedAt: storeMenu.UpdatedAt,
	}, nil
}
