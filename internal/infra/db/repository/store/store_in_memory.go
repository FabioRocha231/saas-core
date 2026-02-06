package memorystore

import (
	"context"
	"fmt"
	"sync"

	"github.com/FabioRocha231/saas-core/internal/domain/store"
)

type Repo struct {
	mu     sync.RWMutex
	byID   map[string]*store.Store
	bySlug map[string]string // slug -> id
}

func New() *Repo {
	return &Repo{
		byID:   make(map[string]*store.Store),
		bySlug: make(map[string]string),
	}
}

func (r *Repo) Create(ctx context.Context, s *store.Store) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// garante unicidade por ID (se você gerar ID fora)
	if _, exists := r.byID[s.ID]; exists {
		return fmt.Errorf("id %s already exists", s.ID)
	}

	// garante unicidade por slug
	if id, exists := r.bySlug[s.Slug]; exists && id != "" {
		return fmt.Errorf("Slug %s already exists", s.Slug)
	}

	cp := *s
	r.byID[cp.ID] = &cp
	r.bySlug[cp.Slug] = cp.ID

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*store.Store, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.byID[id]
	if !ok {
		return nil, fmt.Errorf("Store not found!")
	}

	cp := *s
	return &cp, nil
}

func (r *Repo) GetBySlug(ctx context.Context, slug string) (*store.Store, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.bySlug[slug]
	if !ok || id == "" {
		return nil, fmt.Errorf("Store not found!")
	}

	s, ok := r.byID[id]
	if !ok {
		return nil, fmt.Errorf("Store not found!")
	}

	cp := *s
	return &cp, nil
}
