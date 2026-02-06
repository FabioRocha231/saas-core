package handlers

import (
	"net/http"
	"strings"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	portsRepository "github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/store"
	"github.com/gin-gonic/gin"
)

type CreateStoreRequest struct {
	Name string `json:"name"`
	Cnpj string `json:"cnpj"`
}

type HttpResponse struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
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
	var req CreateStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: err.Error(), Message: "invalid request"})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "name are required"})
		return
	}

	if strings.TrimSpace(req.Cnpj) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "cnpj are required"})
		return
	}

	uc := usecase.NewCreateStoreUsecase(sh.storeRepository, ctx, sh.uuid)
	output, err := uc.Execute(usecase.CreateStoreInput{Name: req.Name, Cnpj: req.Cnpj})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HttpResponse{Status: "error", Error: err.Error(), Message: "error creating store"})
		return
	}

	ctx.JSON(http.StatusCreated, HttpResponse{Status: "success", Message: "store created", Data: output})
}

func (h *StoreHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "id are required"})
		return
	}

	store, err := h.storeRepository.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HttpResponse{Status: "error", Error: err.Error(), Message: "error getting store"})
		return
	}

	ctx.JSON(http.StatusOK, HttpResponse{Status: "success", Message: "store found", Data: store})
}
