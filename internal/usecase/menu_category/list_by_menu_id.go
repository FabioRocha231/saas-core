package usecase

import (
	"context"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListMenuCategoriesByMenuIDInput struct {
	MenuID string
}

type MenuCategories struct {
	ID       string `json:"id"`
	MenuID   string `json:"menu_id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type ListMenuCategoriesByMenuIDOutput struct {
	Categories []*MenuCategories `json:"categories"`
}

type ListMenuCategoriesByMenuIdUseCase struct {
	storeMenuRepo    repository.StoreMenuRepository
	menuCategoryRepo repository.MenuCategoryRepository
	uuid             ports.UUIDInterface
}

func NewListMenuCategoriesByMenuIDUsecase(
	storeMenuRepo repository.StoreMenuRepository,
	menuCategoryRepo repository.MenuCategoryRepository,
	uuid ports.UUIDInterface,
) *ListMenuCategoriesByMenuIdUseCase {
	return &ListMenuCategoriesByMenuIdUseCase{
		storeMenuRepo:    storeMenuRepo,
		menuCategoryRepo: menuCategoryRepo,
		uuid:             uuid,
	}
}

func (uc *ListMenuCategoriesByMenuIdUseCase) Execute(context context.Context, input ListMenuCategoriesByMenuIDInput) (*ListMenuCategoriesByMenuIDOutput, error) {
	menuID := strings.TrimSpace(input.MenuID)
	if menuID == "" {
		return nil, errx.New(errx.CodeInvalid, "menu id are required")
	}

	if isValidUUID := uc.uuid.Validate(menuID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid menu id")
	}

	menu, err := uc.storeMenuRepo.GetByID(context, menuID)
	if err != nil {
		return nil, err
	}

	categories, err := uc.menuCategoryRepo.ListByMenuID(context, menu.ID)
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

	return &ListMenuCategoriesByMenuIDOutput{
		Categories: outputCategories,
	}, nil
}
