package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/menu_category"
	"github.com/gin-gonic/gin"
)

type MenuCategoryHandler struct {
	menuCategoryRepository repository.MenuCategoryRepository
	storeMenuRepo          repository.StoreMenuRepository
	uuid                   ports.UUIDInterface
}

type CreateMenuCategoryRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func NewMenuCategoryHandler(menuCategoryRepository repository.MenuCategoryRepository, storeMenuRepo repository.StoreMenuRepository, uuid ports.UUIDInterface) *MenuCategoryHandler {
	return &MenuCategoryHandler{
		menuCategoryRepository: menuCategoryRepository,
		storeMenuRepo:          storeMenuRepo,
		uuid:                   uuid,
	}
}

func (mch *MenuCategoryHandler) Create(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	if strings.TrimSpace(menuId) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "menu id is required"))
		return
	}

	var req CreateMenuCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInternal, "internal server error"))
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "name are required"))
		return
	}

	uc := usecase.NewCreateMenuCategoryUsecase(mch.menuCategoryRepository, mch.storeMenuRepo, mch.uuid)
	output, err := uc.Execute(ctx, usecase.CreateMenuCategoryInput{Name: req.Name, IsActive: req.IsActive, MenuID: menuId})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, output)
}

func (mch *MenuCategoryHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if strings.TrimSpace(id) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "menu category id is required"))
		return
	}

	uc := usecase.NewGetMenuCategoryByIDUseCase(mch.menuCategoryRepository, mch.uuid, ctx)
	output, err := uc.Execute(usecase.GetMenuCategoryByIDInput{ID: id})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (mch *MenuCategoryHandler) ListByMenuID(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	if strings.TrimSpace(menuId) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "menu id is required"))
		return
	}

	uc := usecase.NewListMenuCategoriesByMenuIdUsecase(mch.storeMenuRepo, mch.menuCategoryRepository, ctx, mch.uuid)
	output, err := uc.Execute(usecase.ListMenuCategoriesByMenuIdInput{MenuID: menuId})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
