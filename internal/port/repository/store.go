package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type StoreRepository interface {
	Create(ctx context.Context, s *entity.Store) error
	GetByID(ctx context.Context, id string) (*entity.Store, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Store, error)
	CountByOwnerID(ctx context.Context, ownerID string) (int, error)
	ListByOwnerID(ctx context.Context, ownerID string) ([]*entity.Store, error)
}
