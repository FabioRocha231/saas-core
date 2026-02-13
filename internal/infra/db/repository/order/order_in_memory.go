package memoryorder

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type userStoreKey struct {
	UserID  string
	StoreID string
}

type Repo struct {
	mu sync.RWMutex

	byID map[string]*entity.Order

	// (userID, storeID) -> orderID (somente para Status=CREATED)
	activeDraftByUserIDAndStoreID map[userStoreKey]string
}

func New() repository.OrderRepository {
	return &Repo{
		byID:                          make(map[string]*entity.Order),
		activeDraftByUserIDAndStoreID: make(map[userStoreKey]string),
	}
}

func (r *Repo) Create(ctx context.Context, o *entity.Order) error {
	_ = ctx

	if o == nil {
		return errx.New(errx.CodeInvalid, "missing order")
	}
	if o.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if o.StoreID == "" {
		return errx.New(errx.CodeInvalid, "missing storeId")
	}
	if o.UserID == "" {
		return errx.New(errx.CodeInvalid, "missing userId")
	}
	if o.Status == "" {
		o.Status = entity.OrderCreated
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[o.ID]; ok {
		return errx.New(errx.CodeConflict, "order already exists")
	}

	// carrinho único: não pode existir outro draft ativo
	if o.Status == entity.OrderCreated {
		key := userStoreKey{UserID: o.UserID, StoreID: o.StoreID}
		if existingID := r.activeDraftByUserIDAndStoreID[key]; existingID != "" {
			return errx.New(errx.CodeConflict, "active draft already exists")
		}
		r.activeDraftByUserIDAndStoreID[key] = o.ID
	}

	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	o.UpdatedAt = now

	cp := cloneOrder(o)
	r.byID[cp.ID] = cp

	return nil
}

func (r *Repo) Update(ctx context.Context, o *entity.Order) error {
	_ = ctx

	if o == nil {
		return errx.New(errx.CodeInvalid, "missing order")
	}
	if o.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if o.StoreID == "" {
		return errx.New(errx.CodeInvalid, "missing storeId")
	}
	if o.UserID == "" {
		return errx.New(errx.CodeInvalid, "missing userId")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.byID[o.ID]
	if !ok || current == nil {
		return errx.New(errx.CodeNotFound, "order not found")
	}

	prevStatus := current.Status
	nextStatus := o.Status
	key := userStoreKey{UserID: o.UserID, StoreID: o.StoreID}

	// se estava CREATED e saiu de CREATED, libera o draft ativo
	if prevStatus == entity.OrderCreated && nextStatus != entity.OrderCreated {
		if id := r.activeDraftByUserIDAndStoreID[key]; id == o.ID {
			delete(r.activeDraftByUserIDAndStoreID, key)
		}
	}

	// se não era CREATED e virou CREATED, protege contra violação do carrinho único
	if prevStatus != entity.OrderCreated && nextStatus == entity.OrderCreated {
		if existingID := r.activeDraftByUserIDAndStoreID[key]; existingID != "" && existingID != o.ID {
			return errx.New(errx.CodeConflict, "active draft already exists")
		}
		r.activeDraftByUserIDAndStoreID[key] = o.ID
	}

	o.UpdatedAt = now
	if o.CreatedAt.IsZero() {
		o.CreatedAt = current.CreatedAt
	}

	cp := cloneOrder(o)
	r.byID[cp.ID] = cp

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	_ = ctx

	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	o, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || o == nil {
		return nil, errx.New(errx.CodeNotFound, "order not found")
	}

	return cloneOrder(o), nil
}

func (r *Repo) GetActiveDraftByUserIDAndStoreID(ctx context.Context, userID, storeID string) (*entity.Order, error) {
	_ = ctx

	if userID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing userId")
	}
	if storeID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing storeId")
	}

	key := userStoreKey{UserID: userID, StoreID: storeID}

	r.mu.RLock()
	id := r.activeDraftByUserIDAndStoreID[key]
	if id == "" {
		r.mu.RUnlock()
		return nil, errx.New(errx.CodeNotFound, "active draft not found")
	}
	o := r.byID[id]
	r.mu.RUnlock()

	if o == nil || o.Status != entity.OrderCreated {
		return nil, errx.New(errx.CodeNotFound, "active draft not found")
	}

	return cloneOrder(o), nil
}

// clone profundo do pedido (porque tem slices)
func cloneOrder(o *entity.Order) *entity.Order {
	if o == nil {
		return nil
	}
	cp := *o

	if o.Items != nil {
		cp.Items = make([]entity.OrderItem, len(o.Items))
		for i := range o.Items {
			cp.Items[i] = cloneOrderItem(o.Items[i])
		}
	}

	return &cp
}

func cloneOrderItem(it entity.OrderItem) entity.OrderItem {
	cp := it

	if it.Variants != nil {
		cp.Variants = make([]entity.OrderItemVariant, len(it.Variants))
		copy(cp.Variants, it.Variants)
	}
	if it.Addons != nil {
		cp.Addons = make([]entity.OrderItemAddon, len(it.Addons))
		copy(cp.Addons, it.Addons)
	}

	return cp
}
