package usecase

import (
	"context"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetUserByEmailInput struct {
	Email string
}

type GetUserByEmailUsecase struct {
	repo    repository.UserRepository
	uuid    ports.UUIDInterface
	context context.Context
}

func NewGetUserByEmailUsecase(repo repository.UserRepository, uuid ports.UUIDInterface, ctx context.Context) *GetUserByEmailUsecase {
	return &GetUserByEmailUsecase{repo: repo, uuid: uuid, context: ctx}
}

func (u *GetUserByEmailUsecase) Execute(input GetUserByEmailInput) (*GetUserOutputDTO, error) {
	user, err := u.repo.GetByMail(u.context, input.Email)

	if err != nil {
		return nil, err
	}

	return ToGetUserOutputDTO(user), nil
}
