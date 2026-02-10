package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ListCategoryItemsByCategoryIDInput struct {
	CategoryID string
}

type CategoryItem struct {
	ID          string
	CategoryID  string
	Name        string
	Description string
	BasePrice   int64
	ImageURL    string
	IsActive    bool
}

type ListCategoryItemsByCategoryIDOutput struct {
	Items []CategoryItem
}

type ListCategoryItemsByCategoryIDUsecase struct {
	categoryItemRepo repository.CategoryItemRepository
	menuCategoryRepo repository.MenuCategoryRepository
	uuid             ports.UUIDInterface
	context          context.Context
}

func NewListCategoryItemsByCategoryIDUsecase(
	categoryItemRepo repository.CategoryItemRepository,
	menuCategoryRepo repository.MenuCategoryRepository,
	uuid ports.UUIDInterface,
	context context.Context,
) *ListCategoryItemsByCategoryIDUsecase {
	return &ListCategoryItemsByCategoryIDUsecase{
		categoryItemRepo: categoryItemRepo,
		menuCategoryRepo: menuCategoryRepo,
		uuid:             uuid,
		context:          context,
	}
}

func (uc *ListCategoryItemsByCategoryIDUsecase) Execute(input ListCategoryItemsByCategoryIDInput) (*ListCategoryItemsByCategoryIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.CategoryID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category id")
	}

	_, err := uc.menuCategoryRepo.GetByID(uc.context, input.CategoryID)
	if err != nil {
		return nil, err
	}

	items, err := uc.categoryItemRepo.ListByCategoryID(uc.context, input.CategoryID)
	if err != nil {
		return nil, err
	}

	var outputItems []CategoryItem
	for _, item := range items {
		outputItems = append(outputItems, CategoryItem{
			ID:          item.ID,
			CategoryID:  item.CategoryID,
			Name:        item.Name,
			Description: item.Description,
			BasePrice:   item.BasePrice,
			ImageURL:    item.ImageURL,
			IsActive:    item.IsActive,
		})
	}

	return &ListCategoryItemsByCategoryIDOutput{
		Items: outputItems,
	}, nil
}
