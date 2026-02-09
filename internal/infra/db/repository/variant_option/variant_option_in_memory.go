package memoryvariantoption

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
	byID    map[string]*entity.VariantOption
	byGroup map[string][]string // groupID -> []optionID
}

func New() repository.VariantOptionRepository {
	return &Repo{
		byID:    make(map[string]*entity.VariantOption),
		byGroup: make(map[string][]string),
	}
}

func (r *Repo) Create(ctx context.Context, o *entity.VariantOption) error {
	_ = ctx

	if o == nil {
		return errx.New(errx.CodeInvalid, "missing variant option")
	}
	if o.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if o.GroupID == "" {
		return errx.New(errx.CodeInvalid, "missing groupId")
	}
	if o.Name == "" {
		return errx.New(errx.CodeInvalid, "missing name")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[o.ID]; ok {
		return errx.New(errx.CodeConflict, "variant option already exists")
	}

	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	o.UpdatedAt = now

	cp := cloneVariantOption(o)
	r.byID[cp.ID] = cp
	r.byGroup[cp.GroupID] = append(r.byGroup[cp.GroupID], cp.ID)

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.VariantOption, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	o, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || o == nil {
		return nil, errx.New(errx.CodeNotFound, "variant option not found")
	}
	return cloneVariantOption(o), nil
}

func (r *Repo) ListByGroupID(ctx context.Context, groupID string) ([]*entity.VariantOption, error) {
	_ = ctx

	if groupID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing groupId")
	}

	r.mu.RLock()
	ids := r.byGroup[groupID]
	out := make([]*entity.VariantOption, 0, len(ids))
	for _, id := range ids {
		if o := r.byID[id]; o != nil {
			out = append(out, cloneVariantOption(o))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func cloneVariantOption(o *entity.VariantOption) *entity.VariantOption {
	if o == nil {
		return nil
	}
	cp := *o
	return &cp
}
