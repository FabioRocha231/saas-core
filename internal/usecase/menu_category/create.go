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

type CreateMenuCategoryInput struct {
	MenuID   string
	Name     string
	IsActive bool
}

type CreateMenuCategoryOutput struct {
	ID string `json:"id"`
}

type CreateMenuCategoryUseCase struct {
	menuCategoryRepo repository.MenuCategoryRepository
	storeMenuRepo    repository.StoreMenuRepository
	uuid             ports.UUIDInterface
}

func NewCreateMenuCategoryUsecase(
	menuCategoryRepo repository.MenuCategoryRepository,
	storeMenuRepo repository.StoreMenuRepository,
	uuid ports.UUIDInterface,
) *CreateMenuCategoryUseCase {
	return &CreateMenuCategoryUseCase{
		menuCategoryRepo: menuCategoryRepo,
		storeMenuRepo:    storeMenuRepo,
		uuid:             uuid,
	}
}

func (uc *CreateMenuCategoryUseCase) Execute(context context.Context, input CreateMenuCategoryInput) (*CreateMenuCategoryOutput, error) {
	menuID := strings.TrimSpace(input.MenuID)
	menuCategoryName := strings.TrimSpace(input.Name)

	if menuCategoryName == "" {
		return nil, errx.New(errx.CodeInvalid, "menu category name are required")
	}

	if menuID == "" {
		return nil, errx.New(errx.CodeInvalid, "menu id are required")
	}

	if isValidUUID := uc.uuid.Validate(menuID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid menu id")
	}

	storeMenu, err := uc.storeMenuRepo.GetByID(context, menuID)
	if err != nil {
		return nil, err
	}

	menuCategory := &entity.MenuCategory{
		ID:        uc.uuid.Generate(),
		MenuID:    storeMenu.ID,
		Name:      menuCategoryName,
		IsActive:  input.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createErr := uc.menuCategoryRepo.Create(context, menuCategory)
	if createErr != nil {
		return nil, createErr
	}

	return &CreateMenuCategoryOutput{
		ID: menuCategory.ID,
	}, nil
}
