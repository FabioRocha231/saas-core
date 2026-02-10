package usecase

import (
	"context"
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
	context          context.Context
}

func NewCreateCategoryItemUsecase(
	categoryItemRepo repository.CategoryItemRepository,
	menuCategoryRepo repository.MenuCategoryRepository,
	uuid ports.UUIDInterface,
	context context.Context,
) *CreateCategoryItemUsecase {
	return &CreateCategoryItemUsecase{
		categoryItemRepo: categoryItemRepo,
		menuCategoryRepo: menuCategoryRepo,
		uuid:             uuid,
		context:          context,
	}
}

func (uc *CreateCategoryItemUsecase) Execute(input CreateCategoryItemInput) (*CreateCategoryItemOutput, error) {
	isValidUuid := uc.uuid.Validate(input.CategoryID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category id")
	}

	category, err := uc.menuCategoryRepo.GetByID(uc.context, input.CategoryID)
	if err != nil {
		return nil, errx.New(errx.CodeNotFound, "category not found")
	}

	itemCategory := &entity.CategoryItem{
		ID:          uc.uuid.Generate(),
		CategoryID:  category.ID,
		Name:        input.Name,
		Description: input.Description,
		BasePrice:   input.BasePrice,
		ImageURL:    input.ImageURL,
		IsActive:    input.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = uc.categoryItemRepo.Create(uc.context, itemCategory)
	if err != nil {
		return nil, err
	}

	return &CreateCategoryItemOutput{
		ID: itemCategory.ID,
	}, nil
}
