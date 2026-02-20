package usecase

import (
	"context"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetMenuCategoryByIDInput struct {
	ID string
}

type GetMenuCategoryByIDOutput struct {
	ID       string `json:"id"`
	MenuID   string `json:"menu_id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type GetMenuCategoryByIDUseCase struct {
	menuCategoryRepo repository.MenuCategoryRepository
	uuid             ports.UUIDInterface
}

func NewGetMenuCategoryByIDUsecase(
	menuCategoryRepo repository.MenuCategoryRepository,
	uuid ports.UUIDInterface,
) *GetMenuCategoryByIDUseCase {
	return &GetMenuCategoryByIDUseCase{
		menuCategoryRepo: menuCategoryRepo,
		uuid:             uuid,
	}
}

func (uc *GetMenuCategoryByIDUseCase) Execute(context context.Context, input GetMenuCategoryByIDInput) (*GetMenuCategoryByIDOutput, error) {
	menuCategoryID := strings.TrimSpace(input.ID)
	if menuCategoryID == "" {
		return nil, errx.New(errx.CodeInvalid, "menu category id are required")
	}

	isValidUUID := uc.uuid.Validate(menuCategoryID)
	if !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid menu category id")
	}

	menuCategory, err := uc.menuCategoryRepo.GetByID(context, menuCategoryID)
	if err != nil {
		return nil, err
	}

	return &GetMenuCategoryByIDOutput{
		ID:       menuCategory.ID,
		MenuID:   menuCategory.MenuID,
		Name:     menuCategory.Name,
		IsActive: menuCategory.IsActive,
	}, nil
}
