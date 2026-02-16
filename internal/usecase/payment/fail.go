package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type FailPaymentInput struct {
	PaymentID string
	UserID    string
	Reason    string
}

type FailPaymentOutput struct {
	Payment PaymentDTO `json:"payment"`
}

type FailPaymentUsecase struct {
	PaymentRepo repository.PaymentRepository
	UUID        ports.UUIDInterface
}

func NewFailPaymentUsecase(paymentRepo repository.PaymentRepository, uuid ports.UUIDInterface) *FailPaymentUsecase {
	return &FailPaymentUsecase{PaymentRepo: paymentRepo, UUID: uuid}
}

func (uc *FailPaymentUsecase) Execute(ctx context.Context, in FailPaymentInput) (*FailPaymentOutput, error) {
	if in.PaymentID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing payment id")
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
	if p.Status != entity.PaymentStatusPending {
		return nil, errx.New(errx.CodeConflict, "payment must be PENDING to fail")
	}

	p.Status = entity.PaymentStatusFailed
	p.UpdatedAt = time.Now()

	if err := uc.PaymentRepo.Update(ctx, p); err != nil {
		return nil, err
	}
	return &FailPaymentOutput{Payment: ToPaymentDTO(p)}, nil
}
