package handlers

import (
	"net/http"
	"strings"

	valueobject "github.com/FabioRocha231/saas-core/internal/domain/value_object"
	"github.com/gin-gonic/gin"
)

type CreateStoreRequest struct {
	Name string `json:"name"`
	Cnpj string `json:"cnpj"`
}

type HttpResponse struct {
	Status string `json:"status"`
	Error string  `json:"error,omitempty"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func CreateStoreHandler(ctx *gin.Context) {
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

	cnpj := valueobject.NewCnpj(req.Cnpj)

	if err := cnpj.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: err.Error(), Message: "invalid cnpj"})
		return
	}

	ctx.JSON(http.StatusCreated, HttpResponse{Status: "success", Message: "store created", Data: req})
}

