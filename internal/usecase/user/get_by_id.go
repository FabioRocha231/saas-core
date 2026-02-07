package usecase

import (
	"context"
	"errors"
	"time"

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

type GetUserByIdOutput struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Cpf             string     `json:"cpf"`
	Phone           string     `json:"phone"`
	Status          string     `json:"status"`
	Role            string     `json:"role"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
}

func NewGetUserByIdUsecase(repo repository.UserRepository, uuid ports.UUIDInterface, ctx context.Context) *GetUserByIdUsecase {
	return &GetUserByIdUsecase{repo: repo, uuid: uuid, context: ctx}
}

func (u *GetUserByIdUsecase) Execute(input GetUserByIdInput) (*GetUserByIdOutput, error) {
	isUuid := u.uuid.Validate(input.ID)

	if !isUuid {
		return nil, errors.New("invalid id")
	}

	user, err := u.repo.GetByID(u.context, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserByIdOutput{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Cpf:             user.Cpf,
		Phone:           user.Phone,
		Status:          user.Status.String(),
		Role:            user.Role.String(),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		DeletedAt:       user.DeletedAt,
		EmailVerifiedAt: user.EmailVerifiedAt,
		PhoneVerifiedAt: user.PhoneVerifiedAt,
		LastLoginAt:     user.LastLoginAt,
	}, nil
}
