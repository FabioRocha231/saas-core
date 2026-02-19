package usecase

import (
	"context"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetStoreByIDUsecase struct {
	storeRepo repository.StoreRepository
	UUID      ports.UUIDInterface
}

type GetStoreByIDInput struct {
	StoreID string
}

type StoreDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	IsOpen  bool   `json:"is_open"`
	Cnpj    string `json:"cnpj"`
	OwnerID string `json:"owner_id"`
}

type GetStoreByIDOutput struct {
	Store StoreDTO `json:"store"`
}

func NewGetStoreByIDUsecase(
	storeRepo repository.StoreRepository,
	uuid ports.UUIDInterface,
) *GetStoreByIDUsecase {
	return &GetStoreByIDUsecase{
		storeRepo: storeRepo,
		UUID:      uuid,
	}
}

func (uc *GetStoreByIDUsecase) Execute(ctx context.Context, input GetStoreByIDInput) (*GetStoreByIDOutput, error) {
	storeID := strings.TrimSpace(input.StoreID)

	if storeID == "" {
		return nil, errx.New(errx.CodeInvalid, "store id are required")
	}

	isValidUuid := uc.UUID.Validate(storeID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid store id")
	}

	store, err := uc.storeRepo.GetByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return &GetStoreByIDOutput{
		Store: StoreDTO{
			ID:      store.ID,
			Name:    store.Name,
			Slug:    store.Slug,
			IsOpen:  store.IsOpen,
			Cnpj:    store.Cnpj,
			OwnerID: store.OwnerID,
		},
	}, nil
}
