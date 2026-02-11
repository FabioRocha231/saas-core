package handlers

import (
	"net/http"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/item_addon_group"
	"github.com/gin-gonic/gin"
)

type ItemAddonGroupHandler struct {
	itemAddonGroupRepo repository.ItemAddonGroupRepository
	categoryItemRepo   repository.CategoryItemRepository
	uuid               ports.UUIDInterface
}

type CreateItemAddonGroupRequest struct {
	Name      string `json:"name" binding:"required"`
	Required  bool   `json:"required"`
	MinSelect int    `json:"min_select" binding:"required"`
	MaxSelect int    `json:"max_select" binding:"required"`
	Order     int    `json:"order" binding:"required"`
	IsActive  bool   `json:"is_active" binding:"required"`
}

func NewItemAddonGroupHandler(
	itemAddonGroupRepo repository.ItemAddonGroupRepository,
	categoryItemRepo repository.CategoryItemRepository,
	uuid ports.UUIDInterface,
) *ItemAddonGroupHandler {
	return &ItemAddonGroupHandler{
		itemAddonGroupRepo: itemAddonGroupRepo,
		categoryItemRepo:   categoryItemRepo,
		uuid:               uuid,
	}
}

func (iah *ItemAddonGroupHandler) Create(ctx *gin.Context) {
	itemID := ctx.Param("categoryItemId")
	if itemID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "category_item_id is required"))
		return
	}

	var req CreateItemAddonGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInvalid, err.Error()))
		return
	}

	uc := usecase.NewCreateItemAddonGroupUseCase(ctx, iah.itemAddonGroupRepo, iah.categoryItemRepo, iah.uuid)
	input := usecase.CreateItemAddonGroupInput{
		CategoryItemID: itemID,
		Name:           req.Name,
		Required:       req.Required,
		MinSelect:      req.MinSelect,
		MaxSelect:      req.MaxSelect,
		Order:          req.Order,
		IsActive:       req.IsActive,
	}

	response, err := uc.Execute(input)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, response)
}

func (iah *ItemAddonGroupHandler) GetByID(ctx *gin.Context) {
	itemID := ctx.Param("id")
	if itemID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "id is required"))
		return
	}

	uc := usecase.NewGetItemAddonGroupByIDUseCase(ctx, iah.itemAddonGroupRepo, iah.uuid)
	input := usecase.GetItemAddonGroupByIDInput{
		ID: itemID,
	}

	response, err := uc.Execute(input)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, response)
}

func (iah *ItemAddonGroupHandler) ListByCategoryItemID(ctx *gin.Context) {
	itemID := ctx.Param("categoryItemId")
	if itemID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "category_item_id is required"))
		return
	}

	uc := usecase.NewListItemAddonGroupByCategoryItemIDUseCase(ctx, iah.itemAddonGroupRepo, iah.categoryItemRepo, iah.uuid)
	input := usecase.ListItemAddonGroupByCategoryItemIDInput{
		CategoryItemID: itemID,
	}

	response, err := uc.Execute(input)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, response)
}
