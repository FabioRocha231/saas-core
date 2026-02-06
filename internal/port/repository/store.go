package ports

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/store"
)

type StoreRepository interface {
	Create(ctx context.Context, s *store.Store) error
	GetByID(ctx context.Context, id string) (*store.Store, error)
	GetBySlug(ctx context.Context, slug string) (*store.Store, error)
}