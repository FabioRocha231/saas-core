package memorystoremenu

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type Repo struct {
	mu      sync.RWMutex
	byID    map[string]*entity.StoreMenu
	byStore map[string][]string // storeID -> []menuID
}

func New() repository.StoreMenuRepository {
	return &Repo{
		byID:    make(map[string]*entity.StoreMenu),
		byStore: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, m *entity.StoreMenu) error {
	_ = ctx

	if m == nil {
		return errx.New(errx.CodeInvalid, "missing menu")
	}
	if m.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if m.StoreID == "" {
		return errx.New(errx.CodeInvalid, "missing storeId")
	}
	if m.Name == "" {
		return errx.New(errx.CodeInvalid, "missing name")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[m.ID]; ok {
		return errx.New(errx.CodeConflict, "menu already exists")
	}

	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now

	cp := cloneStoreMenu(m)
	r.byID[cp.ID] = cp
	r.byStore[cp.StoreID] = append(r.byStore[cp.StoreID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.StoreMenu, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	m, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || m == nil {
		return nil, errx.New(errx.CodeNotFound, "menu not found")
	}
	return cloneStoreMenu(m), nil
}

func (r *Repo) ListByStoreID(ctx context.Context, storeID string) ([]*entity.StoreMenu, error) {
	_ = ctx

	if storeID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing storeId")
	}

	r.mu.RLock()
	ids := r.byStore[storeID]
	out := make([]*entity.StoreMenu, 0, len(ids))
	for _, id := range ids {
		if m := r.byID[id]; m != nil {
			out = append(out, cloneStoreMenu(m))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func cloneStoreMenu(m *entity.StoreMenu) *entity.StoreMenu {
	if m == nil {
		return nil
	}
	cp := *m
	return &cp
}
