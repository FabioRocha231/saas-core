package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	valueobject "github.com/FabioRocha231/saas-core/internal/domain/value_object"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreateUserUsecase struct {
	userRepo     repository.UserRepository
	storeRepo    repository.StoreRepository
	uuid         ports.UUIDInterface
	passwordHash ports.PasswordHashInterface
	context      context.Context
}

type CreateUserInput struct {
	Name     string
	Email    string
	Cpf      string
	Password string
	StoreId  *string
	Phone    string
	UserType string
}

type CreateUserOutput struct {
	ID string `json:"id"`
}

func NewCreateUserUsecase(
	userRepo repository.UserRepository,
	storeRepo repository.StoreRepository,
	ctx context.Context,
	uuid ports.UUIDInterface,
	passwordHash ports.PasswordHashInterface,
) *CreateUserUsecase {
	return &CreateUserUsecase{
		userRepo:     userRepo,
		storeRepo:    storeRepo,
		uuid:         uuid,
		passwordHash: passwordHash,
		context:      ctx,
	}
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
		Phone:     input.Phone,
		Status:    entity.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	switch input.UserType {
	case string(entity.UserKindCustomer):
		user.Role = entity.UserRoleCostumer
	case string(entity.UserKindStore):
		user.Role = entity.UserRoleStoreOwner
	default:
		return nil, errx.New(errx.CodeInvalid, "invalid user type")
	}

	hash, err := uc.passwordHash.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hash

	if err := uc.userRepo.Create(uc.context, user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{ID: user.ID}, nil
}
