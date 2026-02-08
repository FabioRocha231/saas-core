package memorysession

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
)

type Repo struct {
	mu   sync.RWMutex
	byID map[string]*entity.Session // jti -> session
}

func New() *Repo {
	return &Repo{
		byID: make(map[string]*entity.Session),
	}
}

func (r *Repo) Save(ctx context.Context, s *entity.Session) error {
	_ = ctx

	if s == nil {
		return errx.New(errx.CodeInvalid, "missing session")
	}
	if s.ID == "" {
		return errx.New(errx.CodeInvalid, "missing session id")
	}
	if s.UserID == "" {
		return errx.New(errx.CodeInvalid, "missing user id")
	}
	if s.ExpiresAt.IsZero() {
		return errx.New(errx.CodeInvalid, "missing expiresAt")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}

	cp := cloneSession(s)
	r.byID[cp.ID] = cp
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.Session, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing session id")
	}

	r.mu.RLock()
	s, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || s == nil {
		return nil, errx.New(errx.CodeNotFound, "session not found")
	}

	return cloneSession(s), nil
}

func (r *Repo) DeleteExpired(ctx context.Context, now time.Time) error {
	_ = ctx
	if now.IsZero() {
		now = time.Now()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for id, s := range r.byID {
		if s == nil {
			delete(r.byID, id)
			continue
		}
		if s.ExpiresAt.Before(now) || s.ExpiresAt.Equal(now) {
			delete(r.byID, id)
		}
	}
	return nil
}

func cloneSession(s *entity.Session) *entity.Session {
	if s == nil {
		return nil
	}
	cp := *s
	cp.RevokedAt = cloneTimePtr(s.RevokedAt)
	return &cp
}

func cloneTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := *t
	return &v
}
