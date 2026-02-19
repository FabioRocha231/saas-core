package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	valueobject "github.com/FabioRocha231/saas-core/internal/domain/value_object"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetUserByCpfInput struct {
	Cpf string
}

type UserDTO struct {
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

type GetUserOutputDTO struct {
	User UserDTO `json:"user"`
}

type GetUserByCpfUsecase struct {
	repo repository.UserRepository
	uuid ports.UUIDInterface
}

func NewGetUserByCpfUsecase(repo repository.UserRepository, uuid ports.UUIDInterface) *GetUserByCpfUsecase {
	return &GetUserByCpfUsecase{repo: repo, uuid: uuid}
}

func (uc *GetUserByCpfUsecase) Execute(ctx context.Context, input GetUserByCpfInput) (*GetUserOutputDTO, error) {
	if input.Cpf == "" {
		return nil, errx.New(errx.CodeInvalid, "invalid cpf")
	}

	cpf := valueobject.NewCpf(input.Cpf)
	if err := cpf.Validate(); err != nil {
		return nil, errx.New(errx.CodeInvalid, err.Error())
	}

	user, err := uc.repo.GetByCpf(ctx, cpf.Digits())
	if err != nil {
		return nil, err
	}

	return ToGetUserOutputDTO(user), nil
}

func ToGetUserOutputDTO(user *entity.User) *GetUserOutputDTO {
	return &GetUserOutputDTO{
		User: UserDTO{
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
		},
	}
}
