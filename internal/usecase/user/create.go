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
	Status   string
	Role     string
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if input.StoreId != nil {
		isValidUuid := uc.uuid.Validate(*input.StoreId)
		if !isValidUuid {
			return nil, errx.New(errx.CodeInvalid, "invalid store id")
		}
		store, err := uc.storeRepo.GetByID(uc.context, *input.StoreId)
		if err != nil {
			return nil, err
		}
		user.StoreId = &store.ID
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
