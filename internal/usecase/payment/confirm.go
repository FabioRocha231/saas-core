package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type ConfirmPaymentInput struct {
	PaymentID string
	UserID    string
}

type ConfirmPaymentUsecase struct {
	OrderRepo   repository.OrderRepository
	PaymentRepo repository.PaymentRepository
	UUID        ports.UUIDInterface
}

type ConfirmPaymentOutput struct {
	Payment PaymentDTO `json:"payment"`
}

func NewConfirmPaymentUsecase(
	orders repository.OrderRepository,
	payments repository.PaymentRepository,
	uuid ports.UUIDInterface,
) *ConfirmPaymentUsecase {
	return &ConfirmPaymentUsecase{
		OrderRepo:   orders,
		PaymentRepo: payments,
		UUID:        uuid,
	}
}

func (uc *ConfirmPaymentUsecase) Execute(ctx context.Context, in ConfirmPaymentInput) (*ConfirmPaymentOutput, error) {
	if in.PaymentID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing paymentId")
	}
	if in.UserID == "" {
		return nil, errx.New(errx.CodeUnauthorized, "missing user")
	}

	if isValid := uc.UUID.Validate(in.PaymentID); !isValid {
		return nil, errx.New(errx.CodeInvalid, "invalid payment id")
	}
	if isValid := uc.UUID.Validate(in.UserID); !isValid {
		return nil, errx.New(errx.CodeInvalid, "invalid user id")
	}

	p, err := uc.PaymentRepo.GetByID(ctx, in.PaymentID)
	if err != nil {
		return nil, err
	}
	if p.UserID != in.UserID {
		return nil, errx.New(errx.CodeForbidden, "payment does not belong to user")
	}
	if p.Status == entity.PaymentStatusPaid {
		return &ConfirmPaymentOutput{Payment: ToPaymentDTO(p)}, nil
	}
	if p.Status != entity.PaymentStatusPending {
		return nil, errx.New(errx.CodeConflict, "payment must be PENDING to confirm")
	}

	now := time.Now()
	p.Status = entity.PaymentStatusPaid
	p.PaidAt = &now
	p.UpdatedAt = now

	if err := uc.PaymentRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	// marca pedido como PAID
	o, err := uc.OrderRepo.GetByID(ctx, p.OrderID)
	if err != nil {
		return nil, err
	}
	if o.Status == entity.OrderPlaced {
		o.Status = entity.OrderPaid
		o.UpdatedAt = now
		if err := uc.OrderRepo.Update(ctx, o); err != nil {
			return nil, err
		}
	}

	return &ConfirmPaymentOutput{Payment: ToPaymentDTO(p)}, nil
}
