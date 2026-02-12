package handlers

import (
	"net/http"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/item_variant_group"
	"github.com/gin-gonic/gin"
)

type ItemVariantGroupHandler struct {
	itemVariantGroupRepo repository.ItemVariantGroupRepository
	categoryItemRepo     repository.CategoryItemRepository
	uuid                 ports.UUIDInterface
}

type CreateItemVariantGroupRequest struct {
	Name      string `json:"name" binding:"required"`
	Required  bool   `json:"required" binding:"required"`
	MinSelect int    `json:"min_select" binding:"required"`
	MaxSelect int    `json:"max_select" binding:"required"`
	Order     int    `json:"order" binding:"required"`
	IsActive  bool   `json:"is_active" binding:"required"`
}

func NewItemVariantGroupHandler(
	itemVariantGroupRepo repository.ItemVariantGroupRepository,
	categoryItemRepo repository.CategoryItemRepository,
	uuid ports.UUIDInterface,
) *ItemVariantGroupHandler {
	return &ItemVariantGroupHandler{
		itemVariantGroupRepo: itemVariantGroupRepo,
		categoryItemRepo:     categoryItemRepo,
		uuid:                 uuid,
	}
}

func (handler *ItemVariantGroupHandler) Create(ctx *gin.Context) {
	categoryItemID := ctx.Param("categoryItemId")
	if categoryItemID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "category item id is required"))
		return
	}

	var req CreateItemVariantGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInvalid, err.Error()))
		return
	}

	uc := usecase.NewCreateItemVariantGroupUseCase(ctx, handler.itemVariantGroupRepo, handler.categoryItemRepo, handler.uuid)
	input := usecase.CreateItemVariantGroupInput{
		CategoryItemID: categoryItemID,
		Name:           req.Name,
		Required:       req.Required,
		MinSelect:      req.MinSelect,
		MaxSelect:      req.MaxSelect,
		Order:          req.Order,
		IsActive:       req.IsActive,
	}

	output, err := uc.Execute(input)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, output)
}

func (handler *ItemVariantGroupHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "variant group id is required"))
		return
	}

	uc := usecase.NewGetItemVariantGroupByIDUseCase(handler.itemVariantGroupRepo, handler.uuid, ctx)
	output, err := uc.Execute(usecase.GetItemVariantGroupByIDInput{
		ID: id,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (handler *ItemVariantGroupHandler) ListByCategoryItemID(ctx *gin.Context) {
	categoryItemID := ctx.Param("categoryItemId")
	if categoryItemID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "category item id is required"))
		return
	}

	uc := usecase.NewListItemVariantGroupByCategoryItemIDUsecase(ctx, handler.itemVariantGroupRepo, handler.categoryItemRepo, handler.uuid)
	output, err := uc.Execute(usecase.ListItemVariantGroupByCategoryItemIDInput{
		CategoryItemID: categoryItemID,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
