package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListStoreMenuByStoreIDUsecase struct {
	storeMenuRepository repository.StoreMenuRepository
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
	uuid ports.UUIDInterface,
) *ListStoreMenuByStoreIDUsecase {
	return &ListStoreMenuByStoreIDUsecase{
		storeMenuRepository: storeMenuRepository,
		uuid:                uuid,
	}
}

func (uc *ListStoreMenuByStoreIDUsecase) Execute(context context.Context, input ListStoreMenuByStoreIDInput) (*ListStoreMenuByStoreIDOutput, error) {
	storeID := strings.TrimSpace(input.StoreID)

	if storeID == "" {
		return nil, errx.New(errx.CodeInvalid, "store id are required")
	}

	if isValidUuid := uc.uuid.Validate(storeID); !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid store id")
	}

	storeMenus, err := uc.storeMenuRepository.ListByStoreID(context, storeID)
	fmt.Println(storeMenus, err, "que diabo e isso")
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
