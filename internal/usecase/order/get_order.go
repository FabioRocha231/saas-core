package usecase

import (
	"context"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetOrderInput struct {
	OrderID string
	UserID  string
}

type GetOrder struct {
	OrderRepo repository.OrderRepository
	UUID      ports.UUIDInterface
}

func NewGetOrderUsecase(orderRepo repository.OrderRepository, uuid ports.UUIDInterface) *GetOrder {
	return &GetOrder{OrderRepo: orderRepo, UUID: uuid}
}

func (uc *GetOrder) Execute(ctx context.Context, in GetOrderInput) (*Order, error) {
	if in.OrderID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing orderId")
	}
	if in.UserID == "" {
		return nil, errx.New(errx.CodeUnauthorized, "missing userId")
	}

	if isValid := uc.UUID.Validate(in.OrderID); !isValid {
		return nil, errx.New(errx.CodeInvalid, "invalid orderId")
	}
	if isValid := uc.UUID.Validate(in.UserID); !isValid {
		return nil, errx.New(errx.CodeInvalid, "invalid userId")
	}

	o, err := uc.OrderRepo.GetByID(ctx, in.OrderID)
	if err != nil {
		return nil, err
	}

	if o.UserID != in.UserID {
		return nil, errx.New(errx.CodeForbidden, "order does not belong to user")
	}

	return toOrderDTO(o), nil
}
