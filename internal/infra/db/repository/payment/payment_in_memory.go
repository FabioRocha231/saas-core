package memorypayment

import (
	"context"
	"sync"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type orderKey struct {
	OrderID string
	Key     string
}

type Repo struct {
	mu sync.RWMutex

	byID       map[string]*entity.Payment
	byOrder    map[string][]string
	byOrderKey map[orderKey]string
}

func New() repository.PaymentRepository {
	return &Repo{
		byID:       make(map[string]*entity.Payment),
		byOrder:    make(map[string][]string),
		byOrderKey: make(map[orderKey]string),
	}
}

func (r *Repo) Create(ctx context.Context, p *entity.Payment) error {
	_ = ctx

	if p == nil {
		return errx.New(errx.CodeInvalid, "missing payment")
	}
	if p.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}
	if p.OrderID == "" {
		return errx.New(errx.CodeInvalid, "missing orderId")
	}
	if p.UserID == "" {
		return errx.New(errx.CodeInvalid, "missing userId")
	}
	if p.StoreID == "" {
		return errx.New(errx.CodeInvalid, "missing storeId")
	}
	if p.Amount < 0 {
		return errx.New(errx.CodeInvalid, "amount must be >= 0")
	}
	if p.Currency == "" {
		p.Currency = "BRL"
	}
	if p.Status == "" {
		p.Status = entity.PaymentStatusCreated
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[p.ID]; ok {
		return errx.New(errx.CodeConflict, "payment already exists")
	}

	if p.IdempotencyKey != "" {
		k := orderKey{OrderID: p.OrderID, Key: p.IdempotencyKey}
		if existing := r.byOrderKey[k]; existing != "" {
			return errx.New(errx.CodeConflict, "payment already exists for idempotency key")
		}
		r.byOrderKey[k] = p.ID
	}

	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	cp := clonePayment(p)
	r.byID[cp.ID] = cp
	r.byOrder[cp.OrderID] = append(r.byOrder[cp.OrderID], cp.ID)

	return nil
}

func (r *Repo) Update(ctx context.Context, p *entity.Payment) error {
	_ = ctx

	if p == nil {
		return errx.New(errx.CodeInvalid, "missing payment")
	}
	if p.ID == "" {
		return errx.New(errx.CodeInvalid, "missing id")
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	cur, ok := r.byID[p.ID]
	if !ok || cur == nil {
		return errx.New(errx.CodeNotFound, "payment not found")
	}

	p.UpdatedAt = now
	if p.CreatedAt.IsZero() {
		p.CreatedAt = cur.CreatedAt
	}

	cp := clonePayment(p)
	r.byID[cp.ID] = cp
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*entity.Payment, error) {
	_ = ctx
	if id == "" {
		return nil, errx.New(errx.CodeInvalid, "missing id")
	}

	r.mu.RLock()
	p, ok := r.byID[id]
	r.mu.RUnlock()

	if !ok || p == nil {
		return nil, errx.New(errx.CodeNotFound, "payment not found")
	}
	return clonePayment(p), nil
}

func (r *Repo) GetByOrderAndKey(ctx context.Context, orderID, key string) (*entity.Payment, error) {
	_ = ctx
	if orderID == "" || key == "" {
		return nil, errx.New(errx.CodeInvalid, "missing orderId or key")
	}

	k := orderKey{OrderID: orderID, Key: key}

	r.mu.RLock()
	id := r.byOrderKey[k]
	if id == "" {
		r.mu.RUnlock()
		return nil, errx.New(errx.CodeNotFound, "payment not found")
	}
	p := r.byID[id]
	r.mu.RUnlock()

	if p == nil {
		return nil, errx.New(errx.CodeNotFound, "payment not found")
	}
	return clonePayment(p), nil
}

func (r *Repo) ListByOrderID(ctx context.Context, orderID string) ([]*entity.Payment, error) {
	_ = ctx
	if orderID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing orderId")
	}

	r.mu.RLock()
	ids := r.byOrder[orderID]
	out := make([]*entity.Payment, 0, len(ids))
	for _, id := range ids {
		if p := r.byID[id]; p != nil {
			out = append(out, clonePayment(p))
		}
	}
	r.mu.RUnlock()

	return out, nil
}

func clonePayment(p *entity.Payment) *entity.Payment {
	if p == nil {
		return nil
	}
	cp := *p
	if p.PaidAt != nil {
		t := *p.PaidAt
		cp.PaidAt = &t
	}
	return &cp
}
