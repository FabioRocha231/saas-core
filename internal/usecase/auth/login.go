package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type LoginUsecase struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	jwtService   ports.JwtInterface
	passwordHash ports.PasswordHashInterface
	context      context.Context
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string `json:"token"`
}

func NewLoginUsecase(
	context context.Context,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwtService ports.JwtInterface,
	passwordHash ports.PasswordHashInterface,
) *LoginUsecase {
	return &LoginUsecase{
		context:      context,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		jwtService:   jwtService,
		passwordHash: passwordHash,
	}
}

func (l *LoginUsecase) Execute(input LoginInput) (*LoginOutput, error) {
	user, err := l.userRepo.GetByMail(l.context, input.Email)

	if err != nil {
		return nil, err
	}

	isEqual := l.passwordHash.Verify(user.Password, input.Password)
	if !isEqual {
		return nil, errx.New(errx.CodeNotFound, "user not found")
	}

	token, err := l.jwtService.Sign(user.ID, user.Role.String())
	if err != nil {
		return nil, err
	}

	err = l.sessionRepo.Create(l.context, &entity.Session{
		ID:        l.jwtService.GetJTI(token),
		UserID:    user.ID,
		ExpiresAt: l.jwtService.GetExpiresAt(token),
		Role:      user.Role.String(),
		CreatedAt: time.Now(),		
	})
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Token: token}, nil
}
