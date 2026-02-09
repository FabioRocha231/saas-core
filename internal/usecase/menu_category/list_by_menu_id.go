package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListMenuCategoriesByMenuIdInput struct {
	MenuID string
}

type MenuCategories struct {
	ID       string `json:"id"`
	MenuID   string `json:"menu_id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type ListMenuCategoriesByMenuIdOutput struct {
	Categories []*MenuCategories `json:"categories"`
}

type ListMenuCategoriesByMenuIdUseCase struct {
	storeMenuRepo    repository.StoreMenuRepository
	menuCategoryRepo repository.MenuCategoryRepository
	context          context.Context
	uuid             ports.UUIDInterface
}

func NewListMenuCategoriesByMenuIdUsecase(
	storeMenuRepo repository.StoreMenuRepository,
	menuCategoryRepo repository.MenuCategoryRepository,
	context context.Context,
	uuid ports.UUIDInterface,
) *ListMenuCategoriesByMenuIdUseCase {
	return &ListMenuCategoriesByMenuIdUseCase{
		storeMenuRepo:    storeMenuRepo,
		menuCategoryRepo: menuCategoryRepo,
		context:          context,
		uuid:             uuid,
	}
}

func (uc *ListMenuCategoriesByMenuIdUseCase) Execute(input ListMenuCategoriesByMenuIdInput) (*ListMenuCategoriesByMenuIdOutput, error) {
	isValidUUID := uc.uuid.Validate(input.MenuID)
	if !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid menu id")
	}

	menu, err := uc.storeMenuRepo.GetByID(uc.context, input.MenuID)
	if err != nil {
		return nil, err
	}

	categories, err := uc.menuCategoryRepo.ListByMenuID(uc.context, menu.ID)
	if err != nil {
		return nil, err
	}

	var outputCategories []*MenuCategories
	for _, category := range categories {
		outputCategories = append(outputCategories, &MenuCategories{
			ID:       category.ID,
			MenuID:   category.MenuID,
			Name:     category.Name,
			IsActive: category.IsActive,
		})
	}

	return &ListMenuCategoriesByMenuIdOutput{
		Categories: outputCategories,
	}, nil
}
