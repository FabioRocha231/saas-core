package handlers

import (
	"net/http"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/variant_option"
	"github.com/gin-gonic/gin"
)

type VariantOptionHandler struct {
	variantOptionRepo    repository.VariantOptionRepository
	itemVariantGroupRepo repository.ItemVariantGroupRepository
	uuid                 ports.UUIDInterface
}

type CreateVariantOptionRequest struct {
	Name       string `json:"name" binding:"required"`
	PriceDelta *int64 `json:"price_delta" binding:"required"`
	IsDefault  *bool  `json:"is_default" binding:"required"`
	Order      *int   `json:"order" binding:"required"`
	IsActive   *bool  `json:"is_active" binding:"required"`
}

func NewVariantOptionHandler(
	variantOptionRepo repository.VariantOptionRepository,
	itemVariantGroupRepo repository.ItemVariantGroupRepository,
	uuid ports.UUIDInterface,
) *VariantOptionHandler {
	return &VariantOptionHandler{
		variantOptionRepo:    variantOptionRepo,
		itemVariantGroupRepo: itemVariantGroupRepo,
		uuid:                 uuid,
	}
}

func (handler *VariantOptionHandler) Create(ctx *gin.Context) {
	itemVariantGroupId := ctx.Param("itemVariantGroupId")
	if itemVariantGroupId == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "item variant group ID is required"))
		return
	}

	var req CreateVariantOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInvalid, err.Error()))
		return
	}

	uc := usecase.NewCreateVariantOptionUsecase(handler.variantOptionRepo, handler.itemVariantGroupRepo, handler.uuid, ctx.Request.Context())
	output, err := uc.Execute(usecase.CreateVariantOptionInput{
		ItemVariantGroupID: itemVariantGroupId,
		Name:               req.Name,
		PriceDelta:         *req.PriceDelta,
		IsDefault:          *req.IsDefault,
		Order:              *req.Order,
		IsActive:           *req.IsActive,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, output)
}

func (handler *VariantOptionHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "variant option ID is required"))
		return
	}

	uc := usecase.NewGetVariantOptionByIDUsecase(handler.variantOptionRepo, handler.uuid, ctx.Request.Context())
	output, err := uc.Execute(usecase.GetVariantOptionByIDInput{
		ID: id,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (handler *VariantOptionHandler) ListByItemVariantGroupID(ctx *gin.Context) {
	itemVariantGroupId := ctx.Param("itemVariantGroupId")
	if itemVariantGroupId == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "item variant group ID is required"))
		return
	}

	uc := usecase.NewListByItemVariantGroupIDUsecase(handler.variantOptionRepo, handler.itemVariantGroupRepo, handler.uuid, ctx.Request.Context())
	output, err := uc.Execute(usecase.ListVariantOptionsByItemVariantGroupIDInput{
		ItemVariantGroupID: itemVariantGroupId,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
