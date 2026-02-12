package handlers

import (
	"net/http"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/addon_option"
	"github.com/gin-gonic/gin"
)

type AddonOptionHandler struct {
	addonOptionRepo    repository.AddonOptionRepository
	itemAddonGroupRepo repository.ItemAddonGroupRepository
	uuid               ports.UUIDInterface
}

type CreateAddonOptionRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    int64  `json:"price" binding:"required"`
	Order    int    `json:"order" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

func NewAddonOptionHandler(
	addonOptionRepo repository.AddonOptionRepository,
	itemAddonGroupRepo repository.ItemAddonGroupRepository,
	uuid ports.UUIDInterface,
) *AddonOptionHandler {
	return &AddonOptionHandler{
		addonOptionRepo:    addonOptionRepo,
		itemAddonGroupRepo: itemAddonGroupRepo,
		uuid:               uuid,
	}
}

func (aoh *AddonOptionHandler) Create(ctx *gin.Context) {
	itemAddonGroupID := ctx.Param("itemAddonGroupId")
	if itemAddonGroupID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "item addon group id is required"))
		return
	}
	var req CreateAddonOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "invalid request body"))
		return
	}

	uc := usecase.NewCreateAddonOptionUseCase(aoh.addonOptionRepo, aoh.itemAddonGroupRepo, aoh.uuid, ctx)
	output, err := uc.Execute(usecase.CreateAddonOptionInput{
		ItemAddonGroupID: itemAddonGroupID,
		Name:             req.Name,
		Price:            req.Price,
		Order:            req.Order,
		IsActive:         req.IsActive,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, output)
}

func (aoh *AddonOptionHandler) GetByID(ctx *gin.Context) {
	addonOptionId := ctx.Param("id")
	if addonOptionId == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "addon option id is required"))
		return
	}
	uc := usecase.NewGetAddonOptionByIDUsecase(aoh.addonOptionRepo, aoh.uuid, ctx)
	output, err := uc.Execute(usecase.GetAddonOptionByIDInput{
		ID: addonOptionId,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (aoh *AddonOptionHandler) GetByItemAddonGroupID(ctx *gin.Context) {
	itemAddonGroupID := ctx.Param("itemAddonGroupId")
	if itemAddonGroupID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "item addon group id is required"))
		return
	}

	uc := usecase.NewListByItemAddonGroupIDUsecase(aoh.addonOptionRepo, aoh.itemAddonGroupRepo, aoh.uuid, ctx)
	output, err := uc.Execute(usecase.ListByItemAddonGroupIDInput{
		ItemAddonGroupID: itemAddonGroupID,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
