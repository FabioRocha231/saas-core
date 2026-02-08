package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type StoreMenuRepository interface {
	Create(ctx context.Context, m *entity.StoreMenu) error
	GetByID(ctx context.Context, id string) (*entity.StoreMenu, error)
	ListByStoreID(ctx context.Context, storeID string) ([]*entity.StoreMenu, error)
}
