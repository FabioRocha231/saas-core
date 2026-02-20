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

type CreateCategoryItemInput struct {
	CategoryID  string
	Name        string
	Description string
	BasePrice   int64
	ImageURL    string
	IsActive    bool
}

type CreateCategoryItemOutput struct {
	ID string
}

type CreateCategoryItemUsecase struct {
	categoryItemRepo repository.CategoryItemRepository
	menuCategoryRepo repository.MenuCategoryRepository
	uuid             ports.UUIDInterface
}

func NewCreateCategoryItemUsecase(
	categoryItemRepo repository.CategoryItemRepository,
	menuCategoryRepo repository.MenuCategoryRepository,
	uuid ports.UUIDInterface,
) *CreateCategoryItemUsecase {
	return &CreateCategoryItemUsecase{
		categoryItemRepo: categoryItemRepo,
		menuCategoryRepo: menuCategoryRepo,
		uuid:             uuid,
	}
}

func (uc *CreateCategoryItemUsecase) Execute(context context.Context, input CreateCategoryItemInput) (*CreateCategoryItemOutput, error) {
	categoryID := strings.TrimSpace(input.CategoryID)
	categoryItemName := strings.TrimSpace(input.Name)
	categoryItemDescription := strings.TrimSpace(input.Description)

	if categoryItemName == "" {
		return nil, errx.New(errx.CodeInvalid, "category item name are required")
	}
	if categoryItemDescription == "" {
		return nil, errx.New(errx.CodeInvalid, "category item description are required")
	}

	if categoryID == "" {
		return nil, errx.New(errx.CodeInvalid, "category id are required")
	}

	isValidUuid := uc.uuid.Validate(categoryID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category id")
	}

	category, err := uc.menuCategoryRepo.GetByID(context, categoryID)
	if err != nil {
		return nil, err
	}

	itemCategory := &entity.CategoryItem{
		ID:          uc.uuid.Generate(),
		CategoryID:  category.ID,
		Name:        categoryItemName,
		Description: categoryItemDescription,
		BasePrice:   input.BasePrice,
		ImageURL:    input.ImageURL,
		IsActive:    input.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = uc.categoryItemRepo.Create(context, itemCategory)
	if err != nil {
		return nil, err
	}

	return &CreateCategoryItemOutput{
		ID: itemCategory.ID,
	}, nil
}
