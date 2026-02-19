package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"

	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetUserByIdUsecase struct {
	repo repository.UserRepository
	uuid ports.UUIDInterface
}

type GetUserByIdInput struct {
	ID string
}

func NewGetUserByIdUsecase(repo repository.UserRepository, uuid ports.UUIDInterface) *GetUserByIdUsecase {
	return &GetUserByIdUsecase{repo: repo, uuid: uuid}
}

func (u *GetUserByIdUsecase) Execute(ctx context.Context, input GetUserByIdInput) (*GetUserOutputDTO, error) {
	if input.ID == "" {
		return nil, errx.New(errx.CodeInvalid, "invalid user id")
	}

	isUuid := u.uuid.Validate(input.ID)

	if !isUuid {
		return nil, errx.New(errx.CodeInvalid, "invalid user id")
	}

	user, err := u.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return ToGetUserOutputDTO(user), nil
}
