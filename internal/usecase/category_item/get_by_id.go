package usecase

import (
	"context"
	"strings"
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
	uuid             ports.UUIDInterface
}

func NewGetCategoryItemByIDUsecase(categoryItemRepo repository.CategoryItemRepository, uuid ports.UUIDInterface) *GetCategoryItemByIDUsecase {
	return &GetCategoryItemByIDUsecase{
		categoryItemRepo: categoryItemRepo,
		uuid:             uuid,
	}
}

func (uc *GetCategoryItemByIDUsecase) Execute(context context.Context, input GetCategoryItemByIDInput) (*GetCategoryItemByIDOutput, error) {
	categoryItemID := strings.TrimSpace(input.ID)

	if categoryItemID == "" {
		return nil, errx.New(errx.CodeInvalid, "category item id are required")
	}

	isValidUuid := uc.uuid.Validate(categoryItemID)
	if !isValidUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid category item id")
	}

	item, err := uc.categoryItemRepo.GetByID(context, categoryItemID)
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
