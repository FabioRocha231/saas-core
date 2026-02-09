package memorycategoryitem

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type Repo struct {
	mu         sync.RWMutex
	byID       map[string]*entity.CategoryItem
	byCategory map[string][]string // categoryID -> []itemID
}

func New() repository.CategoryItemRepository {
	return &Repo{
		byID:       make(map[string]*entity.CategoryItem),
		byCategory: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, i *entity.CategoryItem) error {
	_ = ctx

	if i == nil {
		return errx.New(errx.CodeInvalid, "missing item")
	}
	if i.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if i.CategoryID == "" {
		return errx.New(errx.CodeInvalid, "missing categoryId")
	}
	if i.Name == "" {
		return errx.New(errx.CodeInvalid, "missing name")
	}
	if i.BasePrice < 0 {
		return errx.New(errx.CodeInvalid, "basePrice must be >= 0")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[i.ID]; ok {
		return errx.New(errx.CodeConflict, "item already exists")
	}

	if i.CreatedAt.IsZero() {
		i.CreatedAt = now
	}
	i.UpdatedAt = now

	cp := cloneItem(i)
	r.byID[cp.ID] = cp
	r.byCategory[cp.CategoryID] = append(r.byCategory[cp.CategoryID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.CategoryItem, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	i, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || i == nil {
		return nil, errx.New(errx.CodeNotFound, "item not found")
	}
	return cloneItem(i), nil
}

func (r *Repo) ListByCategoryID(ctx context.Context, categoryID string) ([]*entity.CategoryItem, error) {
	_ = ctx

	if categoryID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing categoryId")
	}

	r.mu.RLock()
	ids := r.byCategory[categoryID]
	out := make([]*entity.CategoryItem, 0, len(ids))
	for _, id := range ids {
		if it := r.byID[id]; it != nil {
			out = append(out, cloneItem(it))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func cloneItem(i *entity.CategoryItem) *entity.CategoryItem {
	if i == nil {
		return nil
	}
	cp := *i
	return &cp
}
