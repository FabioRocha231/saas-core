package memoryitemvariantgroup

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type Repo struct {
	mu             sync.RWMutex
	byID           map[string]*entity.ItemVariantGroup
	ByCategoryItem map[string][]string // itemID -> []groupID
}

func New() repository.ItemVariantGroupRepository {
	return &Repo{
		byID:           make(map[string]*entity.ItemVariantGroup),
		ByCategoryItem: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, g *entity.ItemVariantGroup) error {
	_ = ctx

	if g == nil {
		return errx.New(errx.CodeInvalid, "missing variant group")
	}
	if g.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if g.CategoryItemID == "" {
		return errx.New(errx.CodeInvalid, "missing category item ID")
	}
	if g.Name == "" {
		return errx.New(errx.CodeInvalid, "missing name")
	}
	if g.MinSelect < 0 || g.MaxSelect < 0 || (g.MaxSelect > 0 && g.MinSelect > g.MaxSelect) {
		return errx.New(errx.CodeInvalid, "invalid min/max select")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[g.ID]; ok {
		return errx.New(errx.CodeConflict, "variant group already exists")
	}

	if g.CreatedAt.IsZero() {
		g.CreatedAt = now
	}
	g.UpdatedAt = now

	cp := cloneVariantGroup(g)
	r.byID[cp.ID] = cp
	r.ByCategoryItem[cp.CategoryItemID] = append(r.ByCategoryItem[cp.CategoryItemID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.ItemVariantGroup, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	g, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || g == nil {
		return nil, errx.New(errx.CodeNotFound, "variant group not found")
	}
	return cloneVariantGroup(g), nil
}

func (r *Repo) ListByCategoryItemID(ctx context.Context, itemID string) ([]*entity.ItemVariantGroup, error) {
	_ = ctx

	if itemID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing category item ID")
	}

	r.mu.RLock()
	ids := r.ByCategoryItem[itemID]
	out := make([]*entity.ItemVariantGroup, 0, len(ids))
	for _, id := range ids {
		if g := r.byID[id]; g != nil {
			out = append(out, cloneVariantGroup(g))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func cloneVariantGroup(g *entity.ItemVariantGroup) *entity.ItemVariantGroup {
	if g == nil {
		return nil
	}
	cp := *g
	return &cp
}
