package usecase

import (
	"context"
	"strings"
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
) *CreateStoreMenuUsecase {
	return &CreateStoreMenuUsecase{
		storeRepository:     storeRepository,
		storeMenuRepository: storeMenuRepository,
		uuid:                uuid,
	}
}

func (uc *CreateStoreMenuUsecase) Execute(context context.Context, input CreateStoreMenuInput) (*CreateStoreMenuOutput, error) {
	menuName := strings.TrimSpace(input.Name)
	storeID := strings.TrimSpace(input.StoreID)

	if menuName == "" {
		return nil, errx.New(errx.CodeInvalid, "menu name are required")
	}
	if storeID == "" {
		return nil, errx.New(errx.CodeInvalid, "store id are required")
	}

	isValidUuid := uc.uuid.Validate(storeID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid store id")
	}

	store, err := uc.storeRepository.GetByID(context, storeID)
	if err != nil {
		return nil, err
	}

	storeMenu := &entity.StoreMenu{
		ID:        uc.uuid.Generate(),
		Name:      menuName,
		StoreID:   store.ID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.storeMenuRepository.Create(context, storeMenu); err != nil {
		return nil, err
	}

	return &CreateStoreMenuOutput{ID: storeMenu.ID}, nil
}
