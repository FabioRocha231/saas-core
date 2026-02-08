package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"

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
	userRepo     portsRepository.UserRepository
	storeRepo    portsRepository.StoreRepository
	uuid         ports.UUIDInterface
	passwordHash ports.PasswordHashInterface
}

func NewUserHandler(
	userRepo portsRepository.UserRepository,
	storeRepo portsRepository.StoreRepository,
	uuid ports.UUIDInterface,
	passwordHash ports.PasswordHashInterface,
) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		storeRepo:    storeRepo,
		uuid:         uuid,
		passwordHash: passwordHash,
	}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInternal, err.Error()))
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "name are required"))
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "email are required"))
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "password are required"))
		return
	}

	if strings.TrimSpace(req.Cpf) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "cpf are required"))
		return
	}

	if strings.TrimSpace(req.Phone) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "phone are required"))
		return
	}

	if strings.TrimSpace(req.Status) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "status are required"))
		return
	}

	if strings.TrimSpace(req.Role) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "role are required"))
		return
	}

	uc := usecase.NewCreateUserUsecase(h.userRepo, h.storeRepo, ctx, h.uuid, h.passwordHash)
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
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, usecaseOutput)
}

func (h *UserHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing id"))
		return
	}
	uc := usecase.NewGetUserByIdUsecase(h.userRepo, h.uuid, ctx)

	output, err := uc.Execute(usecase.GetUserByIdInput{ID: id})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (h *UserHandler) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	if email == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing email"))
		return
	}

	uc := usecase.NewGetUserByEmailUsecase(h.userRepo, h.uuid, ctx)

	output, err := uc.Execute(usecase.GetUserByEmailInput{Email: email})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}

func (h *UserHandler) GetByCpf(ctx *gin.Context) {
	cpf := ctx.Param("cpf")

	if cpf == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing cpf"))
		return
	}

	uc := usecase.NewGetUserByCpfUsecase(h.userRepo, h.uuid, ctx)

	output, err := uc.Execute(usecase.GetUserByCpfInput{Cpf: cpf})

	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
