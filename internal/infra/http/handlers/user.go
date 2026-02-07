package handlers

import (
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"net/http"
	"strings"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	portsRepository "github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Cpf      string  `json:"cpf"`
	Phone    string  `json:"phone"`
	StoreId  *string `json:"store_id"`
	Status   string  `json:"status"`
	Role     string  `json:"role"`
}

type UserHandler struct {
	userRepo  portsRepository.UserRepository
	storeRepo portsRepository.StoreRepository
	uuid      ports.UUIDInterface
}

func NewUserHandler(userRepo portsRepository.UserRepository, storeRepo portsRepository.StoreRepository, uuid ports.UUIDInterface) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		storeRepo: storeRepo,
		uuid:      uuid,
	}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondErr(ctx, errx.New(errx.CodeInternal, err.Error()))
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "name are required"))
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "email are required"))
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "password are required"))
		return
	}

	if strings.TrimSpace(req.Cpf) == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "cpf are required"))
		return
	}

	if strings.TrimSpace(req.Phone) == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "phone are required"))
		return
	}

	if strings.TrimSpace(req.Status) == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "status are required"))
		return
	}

	if strings.TrimSpace(req.Role) == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "role are required"))
		return
	}

	uc := usecase.NewCreateUserUsecase(h.userRepo, h.storeRepo, ctx, h.uuid)
	createUserInput := usecase.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Cpf:      req.Cpf,
		Password: req.Password,
		Phone:    req.Phone,
		Status:   req.Status,
		Role:     req.Role,
	}

	if req.StoreId != nil {
		createUserInput.StoreId = req.StoreId
	}

	usecaseOutput, err := uc.Execute(createUserInput)
	if err != nil {
		respondErr(ctx, err)
		return
	}

	respondOK(ctx, http.StatusCreated, usecaseOutput)
}

func (h *UserHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "missing id"))
		return
	}
	uc := usecase.NewGetUserByIdUsecase(h.userRepo, h.uuid, ctx)

	output, err := uc.Execute(usecase.GetUserByIdInput{ID: id})
	if err != nil {
		respondErr(ctx, err)
		return
	}

	respondOK(ctx, http.StatusOK, output)
}

func (h *UserHandler) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	if email == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "missing email"))
		return
	}

	uc := usecase.NewGetUserByEmailUsecase(h.userRepo, h.uuid, ctx)

	output, err := uc.Execute(usecase.GetUserByEmailInput{Email: email})
	if err != nil {
		respondErr(ctx, err)
		return
	}

	respondOK(ctx, http.StatusOK, output)
}

func (h *UserHandler) GetByCpf(ctx *gin.Context) {
	cpf := ctx.Param("cpf")

	if cpf == "" {
		respondErr(ctx, errx.New(errx.CodeInvalid, "missing cpf"))
		return
	}

	uc := usecase.NewGetUserByCpfUsecase(h.userRepo, h.uuid, ctx)

	output, err := uc.Execute(usecase.GetUserByCpfInput{Cpf: cpf})

	if err != nil {
		respondErr(ctx, err)
		return
	}

	respondOK(ctx, http.StatusOK, output)
}
