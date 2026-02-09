package memorymenucategory

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type Repo struct {
	mu     sync.RWMutex
	byID   map[string]*entity.MenuCategory
	byMenu map[string][]string // menuID -> []categoryID
}

func New() repository.MenuCategoryRepository {
	return &Repo{
		byID:   make(map[string]*entity.MenuCategory),
		byMenu: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, c *entity.MenuCategory) error {
	_ = ctx

	if c == nil {
		return errx.New(errx.CodeInvalid, "missing category")
	}
	if c.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if c.MenuID == "" {
		return errx.New(errx.CodeInvalid, "missing menuId")
	}
	if c.Name == "" {
		return errx.New(errx.CodeInvalid, "missing name")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[c.ID]; ok {
		return errx.New(errx.CodeConflict, "category already exists")
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	cp := cloneCategory(c)
	r.byID[cp.ID] = cp
	r.byMenu[cp.MenuID] = append(r.byMenu[cp.MenuID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.MenuCategory, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	c, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || c == nil {
		return nil, errx.New(errx.CodeNotFound, "category not found")
	}
	return cloneCategory(c), nil
}

func (r *Repo) ListByMenuID(ctx context.Context, menuID string) ([]*entity.MenuCategory, error) {
	_ = ctx

	if menuID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing menuId")
	}

	r.mu.RLock()
	ids := r.byMenu[menuID]
	out := make([]*entity.MenuCategory, 0, len(ids))
	for _, id := range ids {
		if c := r.byID[id]; c != nil {
			out = append(out, cloneCategory(c))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func cloneCategory(c *entity.MenuCategory) *entity.MenuCategory {
	if c == nil {
		return nil
	}
	cp := *c
	return &cp
}
