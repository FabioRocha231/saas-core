package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetStoreMenuByIDUsecase struct {
	storeMenuRepository repository.StoreMenuRepository
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
	uuid ports.UUIDInterface,
) *GetStoreMenuByIDUsecase {
	return &GetStoreMenuByIDUsecase{
		storeMenuRepository: storeMenuRepository,
		uuid:                uuid,
	}
}

func (uc *GetStoreMenuByIDUsecase) Execute(context context.Context, input GetStoreMenuByIDInput) (*GetStoreMenuByIDOutput, error) {
	storeMenuID := strings.TrimSpace(input.StoreMenuID)
	if storeMenuID == "" {
		return nil, errx.New(errx.CodeInvalid, "store menu id are required")
	}

	isValidUuid := uc.uuid.Validate(storeMenuID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid store menu id")
	}

	storeMenu, err := uc.storeMenuRepository.GetByID(context, storeMenuID)
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
