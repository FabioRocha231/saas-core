package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListStoreMenuByStoreIDUsecase struct {
	storeMenuRepository repository.StoreMenuRepository
	context             context.Context
	uuid                ports.UUIDInterface
}

type ListStoreMenuByStoreIDInput struct {
	StoreID string
}

type StoreMenu struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	StoreID  string `json:"storeId"`
	IsActive bool   `json:"isActive"`
}

type ListStoreMenuByStoreIDOutput struct {
	Menus []StoreMenu `json:"menus"`
}

func NewListStoreMenuByStoreIDUsecase(
	storeMenuRepository repository.StoreMenuRepository,
	ctx context.Context,
	uuid ports.UUIDInterface,
) *ListStoreMenuByStoreIDUsecase {
	return &ListStoreMenuByStoreIDUsecase{
		storeMenuRepository: storeMenuRepository,
		context:             ctx,
		uuid:                uuid,
	}
}

func (uc *ListStoreMenuByStoreIDUsecase) Execute(input ListStoreMenuByStoreIDInput) (*ListStoreMenuByStoreIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.StoreID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid storeId")
	}

	storeMenus, err := uc.storeMenuRepository.ListByStoreID(uc.context, input.StoreID)
	if err != nil {
		return nil, err
	}

	var output ListStoreMenuByStoreIDOutput
	for _, menu := range storeMenus {
		output.Menus = append(output.Menus, StoreMenu{
			ID:       menu.ID,
			Name:     menu.Name,
			StoreID:  menu.StoreID,
			IsActive: menu.IsActive,
		})
	}

	return &output, nil
}
