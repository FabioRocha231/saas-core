package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, o *entity.Order) error
	Update(ctx context.Context, o *entity.Order) error
	GetByID(ctx context.Context, id string) (*entity.Order, error)

	// carrinho único
	GetActiveDraftByUserIDAndStoreID(ctx context.Context, userID, storeID string) (*entity.Order, error)
}
