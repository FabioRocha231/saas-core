package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	passwordHash ports.PasswordHashInterface
	jwtService   ports.JwtInterface
	userStore    repository.UserRepository
}

func NewAuthHandler(passwordHash ports.PasswordHashInterface, jwtService ports.JwtInterface, userStore repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		passwordHash: passwordHash,
		jwtService:   jwtService,
		userStore:    userStore,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInternal, err.Error()))
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

	uc := usecase.NewLoginUsecase(ctx, h.userStore, h.jwtService, h.passwordHash)
	output, err := uc.Execute(usecase.LoginInput{Email: req.Email, Password: req.Password})

	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, output)
}
