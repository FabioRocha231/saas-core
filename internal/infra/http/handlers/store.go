package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/infra/http/helper"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/store"
	"github.com/gin-gonic/gin"
)

type CreateStoreRequest struct {
	Name string `json:"name"`
	Cnpj string `json:"cnpj"`
}

type StoreHandler struct {
	storeRepo repository.StoreRepository
	userRepo  repository.UserRepository
	uuid      ports.UUIDInterface
}

func NewStoreHandler(
	storeRepo repository.StoreRepository,
	userRepo repository.UserRepository,
	uuid ports.UUIDInterface,
) *StoreHandler {
	return &StoreHandler{
		storeRepo: storeRepo,
		userRepo:  userRepo,
		uuid:      uuid,
	}
}

// TODO pedir endereço para cadastro
func (sh *StoreHandler) Create(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	var req CreateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, err)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "name are required"))
		return
	}

	if strings.TrimSpace(req.Cnpj) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "cnpj are required"))
		return
	}

	uc := usecase.NewCreateStoreUsecase(sh.storeRepo, sh.userRepo, sh.uuid)
	output, err := uc.Execute(ctx, usecase.CreateStoreInput{
		Name:    req.Name,
		Cnpj:    req.Cnpj,
		OwnerID: userID,
	})

	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, output)
}

func (sh *StoreHandler) GetByID(ctx *gin.Context) {
	storeID := strings.TrimSpace(ctx.Param("id"))

	if storeID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing store id"))
		return
	}

	uc := usecase.NewGetStoreByIDUsecase(sh.storeRepo, sh.uuid)
	output, err := uc.Execute(ctx, usecase.GetStoreByIDInput{StoreID: storeID})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
