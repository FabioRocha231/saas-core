package usecase

import (
	"context"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	"time"
)

type GetUserByEmailInput struct {
	Email string
}

type GetUserByEmailOutput struct {
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

type GetUserByEmailUsecase struct {
	repo    repository.UserRepository
	uuid    ports.UUIDInterface
	context context.Context
}

func NewGetUserByEmailUsecase(repo repository.UserRepository, uuid ports.UUIDInterface, ctx context.Context) *GetUserByEmailUsecase {
	return &GetUserByEmailUsecase{repo: repo, uuid: uuid, context: ctx}
}

func (u *GetUserByEmailUsecase) Execute(input GetUserByEmailInput) (*GetUserByEmailOutput, error) {
	user, err := u.repo.GetByMail(u.context, input.Email)

	if err != nil {
		return nil, err
	}

	return &GetUserByEmailOutput{
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
