package memoryuser

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type Repo struct {
	mu     sync.RWMutex
	byID   map[string]*entity.User
	byCpf  map[string]string // cpf -> id
	byMail map[string]string // email -> id
}

func New() *Repo {
	return &Repo{
		byID:   make(map[string]*entity.User),
		byCpf:  make(map[string]string),
		byMail: make(map[string]string),
	}
}

func (r *Repo) Create(ctx context.Context, u *entity.User) error {
	_ = ctx

	if u == nil {
		return errx.New(errx.CodeInvalid, "missing user")
	}
	if u.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if u.Cpf == "" {
		return errx.New(errx.CodeInvalid, "missing cpf")
	}
	if u.Email == "" {
		return errx.New(errx.CodeInvalid, "missing email")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[u.ID]; ok {
		return errx.New(errx.CodeConflict, "user already exists")
	}

	if existingID, ok := r.byCpf[u.Cpf]; ok && existingID != "" {
		return errx.New(errx.CodeConflict, "user already exists")
	}

	if existingID, ok := r.byMail[u.Email]; ok && existingID != "" {
		return errx.New(errx.CodeConflict, "user already exists")
	}

	if u.Status == "" {
		u.Status = entity.UserStatusActive
	}
	if u.Role == "" {
		u.Role = entity.UserRoleCostumer
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	cp := cloneUser(u)

	r.byID[cp.ID] = cp
	r.byCpf[cp.Cpf] = cp.ID
	r.byMail[cp.Email] = cp.ID

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	u, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || u == nil {
		return nil, errx.New(errx.CodeNotFound, "user not found")
	}
	return cloneUser(u), nil
}

func (r *Repo) GetByCpf(ctx context.Context, cpf string) (*entity.User, error) {
	_ = ctx

	if cpf == "" {
		return nil, errx.New(errx.CodeInvalid, "missing cpf")
	}

	r.mu.RLock()
	id, ok := r.byCpf[cpf]
	if !ok || id == "" {
		r.mu.RUnlock()
		return nil, errx.New(errx.CodeNotFound, "user not found")
	}
	u, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || u == nil {
		return nil, errx.New(errx.CodeNotFound, "user not found")
	}
	return cloneUser(u), nil
}

func (r *Repo) GetByMail(ctx context.Context, mail string) (*entity.User, error) {
	_ = ctx

	if mail == "" {
		return nil, errx.New(errx.CodeInvalid, "missing email")
	}

	r.mu.RLock()
	id, ok := r.byMail[mail]
	if !ok || id == "" {
		r.mu.RUnlock()
		return nil, errx.New(errx.CodeNotFound, "user not found")
	}
	u, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || u == nil {
		return nil, errx.New(errx.CodeNotFound, "user not found")
	}
	return cloneUser(u), nil
}

func cloneUser(u *entity.User) *entity.User {
	if u == nil {
		return nil
	}
	cp := *u
	cp.EmailVerifiedAt = cloneTimePtr(u.EmailVerifiedAt)
	cp.PhoneVerifiedAt = cloneTimePtr(u.PhoneVerifiedAt)
	cp.LastLoginAt = cloneTimePtr(u.LastLoginAt)
	cp.DeletedAt = cloneTimePtr(u.DeletedAt)
	return &cp
}

func cloneTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := *t
	return &v
}
