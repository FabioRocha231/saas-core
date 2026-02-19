package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/store_menu"
	"github.com/gin-gonic/gin"
)

type StoreMenuHandler struct {
	storeRepository     repository.StoreRepository
	storeMenuRepository repository.StoreMenuRepository
	uuid                ports.UUIDInterface
}

type CreateStoreMenuRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"isActive,omitempty"`
}

func NewStoreMenuHandler(
	storeRepo repository.StoreRepository,
	storeMenuRepo repository.StoreMenuRepository,
	uuid ports.UUIDInterface,
) *StoreMenuHandler {
	return &StoreMenuHandler{
		storeRepository:     storeRepo,
		storeMenuRepository: storeMenuRepo,
		uuid:                uuid,
	}
}

func (smh *StoreMenuHandler) Create(ctx *gin.Context) {
	var req CreateStoreMenuRequest
	var storeId string = ctx.Param("storeId")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, err)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "name are required"))
		return
	}

	if strings.TrimSpace(storeId) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "storeId are required"))
		return
	}

	uc := usecase.NewCreateStoreMenuUsecase(smh.storeRepository, smh.storeMenuRepository, smh.uuid)
	output, err := uc.Execute(ctx, usecase.CreateStoreMenuInput{Name: req.Name, StoreID: storeId})

	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, output)
}

func (smh *StoreMenuHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "menu id are required"))
		return
	}

	uc := usecase.NewGetStoreMenuByIDUsecase(smh.storeMenuRepository, smh.uuid)
	output, err := uc.Execute(ctx, usecase.GetStoreMenuByIDInput{StoreMenuID: id})

	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (smh *StoreMenuHandler) ListByStoreID(ctx *gin.Context) {
	storeId := ctx.Param("storeId")

	if storeId == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "storeId are required"))
		return
	}

	uc := usecase.NewListStoreMenuByStoreIDUsecase(smh.storeMenuRepository, smh.uuid)
	output, err := uc.Execute(ctx, usecase.ListStoreMenuByStoreIDInput{StoreID: storeId})

	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
