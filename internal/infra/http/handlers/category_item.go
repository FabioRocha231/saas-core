package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/category_item"
	"github.com/gin-gonic/gin"
)

type CategoryItemHandler struct {
	categoryItemRepo repository.CategoryItemRepository
	menuCategoryRepo repository.MenuCategoryRepository
	uuid             ports.UUIDInterface
}

type CreateCategoryItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BasePrice   int64  `json:"base_price"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

func NewCategoryItemHandler(categoryItemRepo repository.CategoryItemRepository, menuCategoryRepo repository.MenuCategoryRepository, uuid ports.UUIDInterface) *CategoryItemHandler {
	return &CategoryItemHandler{
		categoryItemRepo: categoryItemRepo,
		menuCategoryRepo: menuCategoryRepo,
		uuid:             uuid,
	}
}

func (cih *CategoryItemHandler) Create(ctx *gin.Context) {
	categoryID := ctx.Param("categoryId")
	if strings.TrimSpace(categoryID) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "category id is required"))
		return
	}

	var req CreateCategoryItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "invalid request body"))
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "name are required"))
		return
	}

	if strings.TrimSpace(req.Description) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "description are required"))
		return
	}

	uc := usecase.NewCreateCategoryItemUsecase(cih.categoryItemRepo, cih.menuCategoryRepo, cih.uuid)
	output, err := uc.Execute(
		ctx,
		usecase.CreateCategoryItemInput{
			Name:        req.Name,
			Description: req.Description,
			BasePrice:   req.BasePrice,
			ImageURL:    req.ImageURL,
			IsActive:    req.IsActive,
			CategoryID:  categoryID,
		},
	)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, output)
}

func (cih *CategoryItemHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if strings.TrimSpace(id) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "category item id is required"))
		return
	}

	uc := usecase.NewGetCategoryItemByIDUsecase(cih.categoryItemRepo, cih.uuid)
	output, err := uc.Execute(ctx, usecase.GetCategoryItemByIDInput{ID: id})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (cih *CategoryItemHandler) ListByCategoryID(ctx *gin.Context) {
	categoryID := ctx.Param("categoryId")
	if strings.TrimSpace(categoryID) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "category id is required"))
		return
	}

	uc := usecase.NewListCategoryItemsByCategoryIDUsecase(cih.categoryItemRepo, cih.menuCategoryRepo, cih.uuid)
	output, err := uc.Execute(ctx, usecase.ListCategoryItemsByCategoryIDInput{CategoryID: categoryID})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
