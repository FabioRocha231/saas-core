package usecase

import (
	"context"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	valueobject "github.com/FabioRocha231/saas-core/internal/domain/value_object"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateUserUsecase struct {
	repo    repository.UserRepository
	uuid    ports.UUIDInterface
	context context.Context
}

type CreateUserInput struct {
	Name     string
	Email    string
	Cpf      string
	Password string
	Phone    string
	Status   string
	Role     string
}

type CreateUserOutput struct {
	ID string `json:"id"`
}

func NewCreateUserUsecase(repo repository.UserRepository, ctx context.Context, uuid ports.UUIDInterface) *CreateUserUsecase {
	return &CreateUserUsecase{repo: repo, uuid: uuid, context: ctx}
}

func (uc *CreateUserUsecase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	cpf := valueobject.NewCpf(input.Cpf)
	if err := cpf.Validate(); err != nil {
		return nil, errx.New(errx.CodeInvalid, "invalid cpf")
	}
	user := &entity.User{
		ID:        uc.uuid.Generate(),
		Name:      input.Name,
		Email:     input.Email,
		Cpf:       cpf.Digits(),
		Password:  input.Password,
		Phone:     input.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	roleValue, okRole := entity.UserRoleMap[input.Role]
	if !okRole || roleValue == "" {
		return nil, errx.New(errx.CodeInvalid, "invalid user role")
	}
	user.Role = roleValue

	statusValue, okValue := entity.UserStatusMap[input.Status]
	if !okValue || statusValue == "" {
		return nil, errx.New(errx.CodeInvalid, "invalid user status")
	}
	user.Status = statusValue

	if err := uc.repo.Create(uc.context, user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{ID: user.ID}, nil
}
