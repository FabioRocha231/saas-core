package memoryaddonoption

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
	byID    map[string]*entity.AddonOption
	byGroup map[string][]string // groupID -> []optionID
}

func New() repository.AddonOptionRepository {
	return &Repo{
		byID:    make(map[string]*entity.AddonOption),
		byGroup: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, o *entity.AddonOption) error {
	_ = ctx

	if o == nil {
		return errx.New(errx.CodeInvalid, "missing addon option")
	}
	if o.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if o.AddonGroupID == "" {
		return errx.New(errx.CodeInvalid, "missing groupId")
	}
	if o.Name == "" {
		return errx.New(errx.CodeInvalid, "missing name")
	}
	if o.Price < 0 {
		return errx.New(errx.CodeInvalid, "price must be >= 0")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[o.ID]; ok {
		return errx.New(errx.CodeConflict, "addon option already exists")
	}

	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	o.UpdatedAt = now

	cp := cloneAddonOption(o)
	r.byID[cp.ID] = cp
	r.byGroup[cp.AddonGroupID] = append(r.byGroup[cp.AddonGroupID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.AddonOption, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	o, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || o == nil {
		return nil, errx.New(errx.CodeNotFound, "addon option not found")
	}
	return cloneAddonOption(o), nil
}

func (r *Repo) ListByAddonGroupID(ctx context.Context, groupID string) ([]*entity.AddonOption, error) {
	_ = ctx

	if groupID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing groupId")
	}

	r.mu.RLock()
	ids := r.byGroup[groupID]
	out := make([]*entity.AddonOption, 0, len(ids))
	for _, id := range ids {
		if o := r.byID[id]; o != nil {
			out = append(out, cloneAddonOption(o))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func cloneAddonOption(o *entity.AddonOption) *entity.AddonOption {
	if o == nil {
		return nil
	}
	cp := *o
	return &cp
}
