package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetCategoryItemByIDInput struct {
	ID string
}

type GetCategoryItemByIDOutput struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   int64     `json:"base_price"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetCategoryItemByIDUsecase struct {
	categoryItemRepo repository.CategoryItemRepository
	context          context.Context
	uuid             ports.UUIDInterface
}

func NewGetCategoryItemByIDUsecase(categoryItemRepo repository.CategoryItemRepository, uuid ports.UUIDInterface, context context.Context) *GetCategoryItemByIDUsecase {
	return &GetCategoryItemByIDUsecase{
		categoryItemRepo: categoryItemRepo,
		context:          context,
		uuid:             uuid,
	}
}

func (uc *GetCategoryItemByIDUsecase) Execute(input GetCategoryItemByIDInput) (*GetCategoryItemByIDOutput, error) {
	isValidUuid := uc.uuid.Validate(input.ID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid id")
	}

	item, err := uc.categoryItemRepo.GetByID(uc.context, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetCategoryItemByIDOutput{
		ID:          item.ID,
		CategoryID:  item.CategoryID,
		Name:        item.Name,
		Description: item.Description,
		BasePrice:   item.BasePrice,
		ImageURL:    item.ImageURL,
		IsActive:    item.IsActive,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}
