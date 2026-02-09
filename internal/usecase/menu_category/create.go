package usecase

import (
	"context"
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
	context          context.Context
}

func NewCreateMenuCategoryUsecase(menuCategoryRepo repository.MenuCategoryRepository, storeMenuRepo repository.StoreMenuRepository, uuid ports.UUIDInterface, ctx context.Context) *CreateMenuCategoryUseCase {
	return &CreateMenuCategoryUseCase{
		menuCategoryRepo: menuCategoryRepo,
		storeMenuRepo:    storeMenuRepo,
		uuid:             uuid,
		context:          ctx,
	}
}

func (uc *CreateMenuCategoryUseCase) Execute(input CreateMenuCategoryInput) (*CreateMenuCategoryOutput, error) {
	isValidUUID := uc.uuid.Validate(input.MenuID)
	if !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid menu id")
	}

	storeMenu, err := uc.storeMenuRepo.GetByID(uc.context, input.MenuID)
	if err != nil {
		return nil, err
	}

	menuCategory := &entity.MenuCategory{
		ID:        uc.uuid.Generate(),
		MenuID:    storeMenu.ID,
		Name:      input.Name,
		IsActive:  input.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createErr := uc.menuCategoryRepo.Create(uc.context, menuCategory)
	if createErr != nil {
		return nil, createErr
	}

	return &CreateMenuCategoryOutput{
		ID: menuCategory.ID,
	}, nil
}
