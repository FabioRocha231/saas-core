package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	valueobject "github.com/FabioRocha231/saas-core/internal/domain/value_object"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	repositoryPorts "github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateStoreUsecase struct {
	storeRepository repositoryPorts.StoreRepository
	uuid            ports.UUIDInterface
	context         context.Context
}

type CreateStoreInput struct {
	Name    string
	Cnpj    string
	OwnerID string
}

type CreateStoreOutput struct {
	ID string `json:"id"`
}

func NewCreateStoreUsecase(storeRepository repositoryPorts.StoreRepository, ctx context.Context, uuid ports.UUIDInterface) *CreateStoreUsecase {
	return &CreateStoreUsecase{storeRepository: storeRepository, context: ctx, uuid: uuid}
}

func (uc *CreateStoreUsecase) Execute(input CreateStoreInput) (*CreateStoreOutput, error) {
	cnpj := valueobject.NewCnpj(input.Cnpj)
	if err := cnpj.Validate(); err != nil {
		return nil, errx.New(errx.CodeInvalid, "invalid cnpj")
	}

	store := &entity.Store{
		Name:    input.Name,
		Cnpj:    cnpj.Digits(),
		ID:      uc.uuid.Generate(),
		Slug:    input.Name,
		IsOpen:  true,
		OwnerID: input.OwnerID,
	}

	return &CreateStoreOutput{ID: store.ID}, uc.storeRepository.Create(uc.context, store)
}
