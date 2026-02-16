package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type CreatePaymentInput struct {
	OrderID string
	UserID  string

	Method         entity.PaymentMethod // sempre MOCK por agora
	IdempotencyKey string               // opcional
}

type PaymentDTO struct {
	ID             string     `json:"id"`
	OrderID        string     `json:"order_id"`
	UserID         string     `json:"user_id"`
	StoreID        string     `json:"store_id"`
	Method         string     `json:"method"`
	Provider       string     `json:"provider"`
	Status         string     `json:"status"`
	Amount         int64      `json:"amount"`
	Currency       string     `json:"currency"`
	IdempotencyKey string     `json:"idempotency_key"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	PaidAt         *time.Time `json:"paid_at"`
}

type CreatePaymentOutput struct {
	Payment PaymentDTO `json:"payment"`
}

type CreatePaymentUsecase struct {
	Orders   repository.OrderRepository
	Payments repository.PaymentRepository
	UUID     ports.UUIDInterface
}

func NewCreatePaymentUsecase(
	orders repository.OrderRepository,
	payments repository.PaymentRepository,
	uuid ports.UUIDInterface,
) *CreatePaymentUsecase {
	return &CreatePaymentUsecase{
		Orders:   orders,
		Payments: payments,
		UUID:     uuid,
	}
}

func (uc *CreatePaymentUsecase) Execute(ctx context.Context, in CreatePaymentInput) (*CreatePaymentOutput, error) {
	if in.OrderID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing order id")
	}
	if in.UserID == "" {
		return nil, errx.New(errx.CodeUnauthorized, "missing user")
	}

	if isValid := uc.UUID.Validate(in.OrderID); !isValid {
		return nil, errx.New(errx.CodeInvalid, "invalid order id")
	}
	if isValid := uc.UUID.Validate(in.UserID); !isValid {
		return nil, errx.New(errx.CodeInvalid, "invalid user id")
	}

	if in.Method == "" {
		in.Method = entity.PaymentMethodMock
	}

	o, err := uc.Orders.GetByID(ctx, in.OrderID)
	if err != nil {
		return nil, err
	}
	if o.UserID != in.UserID {
		return nil, errx.New(errx.CodeForbidden, "order does not belong to user")
	}
	if o.Status != entity.OrderPlaced {
		return nil, errx.New(errx.CodeConflict, "order must be PLACED to create payment")
	}
	if len(o.Items) == 0 {
		return nil, errx.New(errx.CodeInvalid, "order has no items")
	}

	// idempotência
	if in.IdempotencyKey != "" {
		existing, e := uc.Payments.GetByOrderAndKey(ctx, o.ID, in.IdempotencyKey)
		if e == nil && existing != nil {
			return &CreatePaymentOutput{Payment: ToPaymentDTO(existing)}, nil
		}
		if e != nil && !errx.Is(e, errx.CodeNotFound) {
			return nil, e
		}
	}

	o.RecalculateTotals() // garante amount correto
	now := time.Now()

	p := &entity.Payment{
		ID:             uc.UUID.Generate(),
		OrderID:        o.ID,
		UserID:         o.UserID,
		StoreID:        o.StoreID,
		Method:         in.Method,
		Provider:       entity.PaymentProviderMock,
		Status:         entity.PaymentStatusPending,
		Amount:         int64(o.Total),
		Currency:       "BRL",
		IdempotencyKey: in.IdempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := uc.Payments.Create(ctx, p); err != nil {
		if errx.Is(err, errx.CodeConflict) && in.IdempotencyKey != "" {
			existing, e := uc.Payments.GetByOrderAndKey(ctx, o.ID, in.IdempotencyKey)
			if e == nil && existing != nil {
				return &CreatePaymentOutput{Payment: ToPaymentDTO(existing)}, nil
			}
		}
		return nil, err
	}

	return &CreatePaymentOutput{Payment: ToPaymentDTO(p)}, nil
}

func ToPaymentDTO(p *entity.Payment) PaymentDTO {
	return PaymentDTO{
		ID:             p.ID,
		OrderID:        p.OrderID,
		UserID:         p.UserID,
		StoreID:        p.StoreID,
		Method:         p.Method.String(),
		Provider:       p.Provider.String(),
		Status:         p.Status.String(),
		Amount:         p.Amount,
		Currency:       p.Currency,
		IdempotencyKey: p.IdempotencyKey,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		PaidAt:         p.PaidAt,
	}
}
