package usecase

import (
	"context"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	valueobject "github.com/FabioRocha231/saas-core/internal/domain/value_object"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	repositoryPorts "github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateStoreUsecase struct {
	storeRepository repositoryPorts.StoreRepository
	userRepo        repositoryPorts.UserRepository
	uuid            ports.UUIDInterface
}

type CreateStoreInput struct {
	Name    string
	Cnpj    string
	OwnerID string
}

type CreateStoreOutput struct {
	ID string `json:"id"`
}

func NewCreateStoreUsecase(storeRepository repositoryPorts.StoreRepository, userRepo repositoryPorts.UserRepository, uuid ports.UUIDInterface) *CreateStoreUsecase {
	return &CreateStoreUsecase{storeRepository: storeRepository, userRepo: userRepo, uuid: uuid}
}

func (uc *CreateStoreUsecase) Execute(ctx context.Context, input CreateStoreInput) (*CreateStoreOutput, error) {
	storeName := strings.TrimSpace(input.Name)
	storeCNPJ := strings.TrimSpace(input.Cnpj)
	storeOwnerID := strings.TrimSpace(input.OwnerID)

	if storeName == "" {
		return nil, errx.New(errx.CodeInvalid, "store name are required")
	}
	if storeCNPJ == "" {
		return nil, errx.New(errx.CodeInvalid, "store cnpj are required")
	}

	cnpj := valueobject.NewCnpj(storeCNPJ)
	if err := cnpj.Validate(); err != nil {
		return nil, errx.New(errx.CodeInvalid, err.Error())
	}

	if storeOwnerID == "" {
		return nil, errx.New(errx.CodeInvalid, "store owner id are required")
	}

	if isValidUUID := uc.uuid.Validate(storeOwnerID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid owner id")
	}

	if _, err := uc.userRepo.GetByID(ctx, storeOwnerID); err != nil {
		return nil, err
	}

	store := &entity.Store{
		Name:    storeName,
		Cnpj:    cnpj.Digits(),
		ID:      uc.uuid.Generate(),
		Slug:    storeName, // TODO: gerar slug de verdade
		IsOpen:  true,
		OwnerID: storeOwnerID,
	}

	err := uc.storeRepository.Create(ctx, store)
	if err != nil {
		return nil, err
	}

	return &CreateStoreOutput{ID: store.ID}, nil
}
