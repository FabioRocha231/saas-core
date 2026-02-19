package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetUserByEmailInput struct {
	Email string
}

type GetUserByEmailUsecase struct {
	repo repository.UserRepository
	uuid ports.UUIDInterface
}

func NewGetUserByEmailUsecase(repo repository.UserRepository, uuid ports.UUIDInterface) *GetUserByEmailUsecase {
	return &GetUserByEmailUsecase{repo: repo, uuid: uuid}
}

func (u *GetUserByEmailUsecase) Execute(ctx context.Context, input GetUserByEmailInput) (*GetUserOutputDTO, error) {
	if input.Email == "" {
		return nil, errx.New(errx.CodeInvalid, "invalid email")
	}

	user, err := u.repo.GetByMail(ctx, input.Email)

	if err != nil {
		return nil, err
	}

	return ToGetUserOutputDTO(user), nil
}
