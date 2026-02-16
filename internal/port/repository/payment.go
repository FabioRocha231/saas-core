package repository

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *entity.Payment) error
	Update(ctx context.Context, p *entity.Payment) error
	GetByID(ctx context.Context, id string) (*entity.Payment, error)

	// idempotência: 1 cobrança por (order + key)
	GetByOrderAndKey(ctx context.Context, orderID, key string) (*entity.Payment, error)

	ListByOrderID(ctx context.Context, orderID string) ([]*entity.Payment, error)
}
