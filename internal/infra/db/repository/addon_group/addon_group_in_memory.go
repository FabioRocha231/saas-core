package memoryaddongroup

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
	byID   map[string]*entity.AddonGroup
	byItem map[string][]string // itemID -> []groupID
}

func New() repository.AddonGroupRepository {
	return &Repo{
		byID:   make(map[string]*entity.AddonGroup),
		byItem: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, g *entity.AddonGroup) error {
	_ = ctx

	if g == nil {
		return errx.New(errx.CodeInvalid, "missing addon group")
	}
	if g.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if g.ItemID == "" {
		return errx.New(errx.CodeInvalid, "missing itemId")
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
		return errx.New(errx.CodeConflict, "addon group already exists")
	}

	if g.CreatedAt.IsZero() {
		g.CreatedAt = now
	}
	g.UpdatedAt = now

	cp := cloneAddonGroup(g)
	r.byID[cp.ID] = cp
	r.byItem[cp.ItemID] = append(r.byItem[cp.ItemID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.AddonGroup, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	g, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || g == nil {
		return nil, errx.New(errx.CodeNotFound, "addon group not found")
	}
	return cloneAddonGroup(g), nil
}

func (r *Repo) ListByItemID(ctx context.Context, itemID string) ([]*entity.AddonGroup, error) {
	_ = ctx

	if itemID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing itemId")
	}

	r.mu.RLock()
	ids := r.byItem[itemID]
	out := make([]*entity.AddonGroup, 0, len(ids))
	for _, id := range ids {
		if g := r.byID[id]; g != nil {
			out = append(out, cloneAddonGroup(g))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func cloneAddonGroup(g *entity.AddonGroup) *entity.AddonGroup {
	if g == nil {
		return nil
	}
	cp := *g
	return &cp
}
