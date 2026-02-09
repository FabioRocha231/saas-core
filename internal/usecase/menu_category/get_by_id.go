package usecase

import (
	"context"

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
	context          context.Context
	uuid             ports.UUIDInterface
}

func NewGetMenuCategoryByIDUseCase(menuCategoryRepo repository.MenuCategoryRepository, uuid ports.UUIDInterface, ctx context.Context) *GetMenuCategoryByIDUseCase {
	return &GetMenuCategoryByIDUseCase{
		menuCategoryRepo: menuCategoryRepo,
		uuid:             uuid,
		context:          ctx,
	}
}

func (uc *GetMenuCategoryByIDUseCase) Execute(input GetMenuCategoryByIDInput) (*GetMenuCategoryByIDOutput, error) {
	isValidUUID := uc.uuid.Validate(input.ID)
	if !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid menu category id")
	}

	menuCategory, err := uc.menuCategoryRepo.GetByID(uc.context, input.ID)
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
