package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByCpf(ctx context.Context, cpf string) (*entity.User, error)
	GetByMail(ctx context.Context, mail string) (*entity.User, error)
}
