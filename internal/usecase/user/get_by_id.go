package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetUserByIdUsecase struct {
	repo    repository.UserRepository
	uuid    ports.UUIDInterface
	context context.Context
}

type GetUserByIdInput struct {
	ID string
}

func NewGetUserByIdUsecase(repo repository.UserRepository, uuid ports.UUIDInterface, ctx context.Context) *GetUserByIdUsecase {
	return &GetUserByIdUsecase{repo: repo, uuid: uuid, context: ctx}
}

func (u *GetUserByIdUsecase) Execute(input GetUserByIdInput) (*GetUserOutputDTO, error) {
	isUuid := u.uuid.Validate(input.ID)

	if !isUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid id")
	}

	user, err := u.repo.GetByID(u.context, input.ID)
	if err != nil {
		return nil, err
	}

	return ToGetUserOutputDTO(user), nil
}
