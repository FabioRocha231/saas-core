package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type PlaceOrderInput struct {
	OrderID string
	UserID  string
}

type PlaceOrderUsecase struct {
	OrderRepo repository.OrderRepository
	UUID      ports.UUIDInterface
}

func NewPlaceOrderUsecase(orderRepo repository.OrderRepository, uuid ports.UUIDInterface) *PlaceOrderUsecase {
	return &PlaceOrderUsecase{OrderRepo: orderRepo, UUID: uuid}
}

func (uc *PlaceOrderUsecase) Execute(ctx context.Context, in PlaceOrderInput) (*Order, error) {
	if in.OrderID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing orderId")
	}
	if in.UserID == "" {
		return nil, errx.New(errx.CodeUnauthorized, "missing user")
	}

	if isValidUUID := uc.UUID.Validate(in.OrderID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid order id")
	}

	if isValidUUID := uc.UUID.Validate(in.UserID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid user id")
	}

	o, err := uc.OrderRepo.GetByID(ctx, in.OrderID)
	if err != nil {
		return nil, err
	}

	if o.UserID != in.UserID {
		return nil, errx.New(errx.CodeForbidden, "order does not belong to user")
	}

	if o.Status != entity.OrderCreated {
		return nil, errx.New(errx.CodeConflict, "order cannot be placed")
	}

	if len(o.Items) == 0 {
		return nil, errx.New(errx.CodeInvalid, "order has no items")
	}

	// garante totals corretos no backend
	o.RecalculateTotals()

	o.Status = entity.OrderPlaced
	o.UpdatedAt = time.Now()

	if err := uc.OrderRepo.Update(ctx, o); err != nil {
		return nil, err
	}

	return toOrderDTO(o), nil
}
