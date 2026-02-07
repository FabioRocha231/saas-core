package handlers

import (
	"net/http"
	"strings"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	portsRepository "github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Cpf      string `json:"cpf"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
	Role     string `json:"role"`
}

type UserHandler struct {
	userRepository portsRepository.UserRepository
	uuid           ports.UUIDInterface
}

func NewUserHandler(repo portsRepository.UserRepository, uuid ports.UUIDInterface) *UserHandler {
	return &UserHandler{userRepository: repo, uuid: uuid}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: err.Error(), Message: "invalid request"})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "name are required"})
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "email are required"})
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "password are required"})
		return
	}

	if strings.TrimSpace(req.Cpf) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "cpf are required"})
		return
	}

	if strings.TrimSpace(req.Phone) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "phone are required"})
		return
	}

	if strings.TrimSpace(req.Status) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "status are required"})
		return
	}

	if strings.TrimSpace(req.Role) == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "role are required"})
		return
	}

	uc := usecase.NewCreateUserUsecase(h.userRepository, ctx, h.uuid)
	usecaseOutput, err := uc.Execute(usecase.CreateUserInput{Name: req.Name, Email: req.Email, Cpf: req.Cpf, Password: req.Password, Phone: req.Phone, Status: req.Status, Role: req.Role})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HttpResponse{Status: "error", Error: err.Error(), Message: "error creating user"})
		return
	}

	ctx.JSON(http.StatusCreated, HttpResponse{Status: "succsess", Message: "user created", Data: usecaseOutput})
}

func (h *UserHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, HttpResponse{Status: "error", Error: "invalid request", Message: "id are required"})
		return
	}
	uc := usecase.NewGetUserByIdUsecase(h.userRepository, h.uuid, ctx)

	output, err := uc.Execute(usecase.GetUserByIdInput{ID: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HttpResponse{Status: "error", Error: err.Error(), Message: "error getting user"})
		return
	}

	ctx.JSON(http.StatusOK, HttpResponse{Status: "success", Message: "user found", Data: output})
}
