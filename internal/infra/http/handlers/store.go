package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/infra/http/helper"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	portsRepository "github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/store"
	"github.com/gin-gonic/gin"
)

type CreateStoreRequest struct {
	Name string `json:"name"`
	Cnpj string `json:"cnpj"`
}

type StoreHandler struct {
	storeRepository portsRepository.StoreRepository
	uuid            ports.UUIDInterface
}

func NewStoreHandler(repo portsRepository.StoreRepository, uuid ports.UUIDInterface) *StoreHandler {
	return &StoreHandler{storeRepository: repo, uuid: uuid}
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

	uc := usecase.NewCreateStoreUsecase(sh.storeRepository, ctx, sh.uuid)
	output, err := uc.Execute(usecase.CreateStoreInput{
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
	id := ctx.Param("id")

	if id == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing id"))
		return
	}

	store, err := sh.storeRepository.GetByID(ctx, id)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, store)
}
