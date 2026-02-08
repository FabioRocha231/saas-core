package repository

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, s *entity.Session) error
	GetByID(ctx context.Context, id string) (*entity.Session, error)
	DeleteExpired(ctx context.Context, now time.Time) error
}
