package usecase

import (
	"context"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type LoginUsecase struct {
	userRepo     repository.UserRepository
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

func NewLoginUsecase(context context.Context, userRepo repository.UserRepository, jwtService ports.JwtInterface, passwordHash ports.PasswordHashInterface) *LoginUsecase {
	return &LoginUsecase{
		context:      context,
		userRepo:     userRepo,
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

	return &LoginOutput{Token: token}, nil
}
