package memorystore

import (
	"context"
	"sync"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/port/repository"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
)

type Repo struct {
	mu        sync.RWMutex
	byID      map[string]*entity.Store
	bySlug    map[string]string   // slug -> id
	byOwnerID map[string][]string // ownerID -> []id
}

func New() repository.StoreRepository {
	return &Repo{
		byID:      make(map[string]*entity.Store),
		bySlug:    make(map[string]string),
		byOwnerID: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, s *entity.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byID[s.ID]; exists {
		return errx.F(errx.CodeConflict, "id %s already exists", s.ID)
	}

	if id, exists := r.bySlug[s.Slug]; exists && id != "" {
		return errx.F(errx.CodeConflict, "Slug %s already exists", s.Slug)
	}

	cp := *s
	r.byID[cp.ID] = &cp
	r.bySlug[cp.Slug] = cp.ID
	r.byOwnerID[cp.OwnerID] = append(r.byOwnerID[cp.OwnerID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.Store, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.byID[id]
	if !ok {
		return nil, errx.New(errx.CodeNotFound, "Store not found")
	}

	cp := *s
	return &cp, nil
}

func (r *Repo) GetBySlug(ctx context.Context, slug string) (*entity.Store, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.bySlug[slug]
	if !ok || id == "" {
		return nil, errx.New(errx.CodeNotFound, "Store not found")
	}

	s, ok := r.byID[id]
	if !ok {
		return nil, errx.New(errx.CodeNotFound, "Store not found")
	}

	cp := *s
	return &cp, nil
}

func (r *Repo) CountByOwnerID(ctx context.Context, ownerID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.byOwnerID[ownerID]
	if !ok {
		return 0, nil
	}

	count := len(ids)

	return count, nil
}
