package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateStoreMenuUsecase struct {
	storeRepository     repository.StoreRepository
	storeMenuRepository repository.StoreMenuRepository
	uuid                ports.UUIDInterface
	context             context.Context
}

type CreateStoreMenuInput struct {
	Name    string
	StoreID string
}

type CreateStoreMenuOutput struct {
	ID string `json:"id"`
}

func NewCreateStoreMenuUsecase(
	storeRepository repository.StoreRepository,
	storeMenuRepository repository.StoreMenuRepository,
	uuid ports.UUIDInterface,
	ctx context.Context,
) *CreateStoreMenuUsecase {
	return &CreateStoreMenuUsecase{
		storeRepository:     storeRepository,
		storeMenuRepository: storeMenuRepository,
		uuid:                uuid,
		context:             ctx,
	}
}

func (uc *CreateStoreMenuUsecase) Execute(input CreateStoreMenuInput) (*CreateStoreMenuOutput, error) {
	isValidUuid := uc.uuid.Validate(input.StoreID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid storeId")
	}

	store, err := uc.storeRepository.GetByID(uc.context, input.StoreID)
	if err != nil {
		return nil, err
	}

	storeMenu := &entity.StoreMenu{
		ID:        uc.uuid.Generate(),
		Name:      input.Name,
		StoreID:   store.ID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.storeMenuRepository.Create(uc.context, storeMenu); err != nil {
		return nil, err
	}

	return &CreateStoreMenuOutput{ID: storeMenu.ID}, nil
}
