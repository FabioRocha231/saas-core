package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type UpdateItemQtyInput struct {
	OrderID string
	UserID  string
	ItemID  string
	Qty     int64
}

type UpdateItemQtyUsecase struct {
	OrderRepo repository.OrderRepository
	UUID      ports.UUIDInterface
}

func NewUpdateItemQtyUsecase(orderRepo repository.OrderRepository, uuid ports.UUIDInterface) *UpdateItemQtyUsecase {
	return &UpdateItemQtyUsecase{OrderRepo: orderRepo, UUID: uuid}
}

func (uc *UpdateItemQtyUsecase) Execute(ctx context.Context, in UpdateItemQtyInput) (*Order, error) {
	if in.OrderID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing orderId")
	}

	if in.UserID == "" {
		return nil, errx.New(errx.CodeUnauthorized, "missing user")
	}
	if in.ItemID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing itemId")
	}

	if isValidUUID := uc.UUID.Validate(in.ItemID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid item id")
	}

	if isValidUUID := uc.UUID.Validate(in.OrderID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid order id")
	}

	if isValidUUID := uc.UUID.Validate(in.UserID); !isValidUUID {
		return nil, errx.New(errx.CodeInvalid, "invalid user id")
	}

	if in.Qty <= 0 {
		return nil, errx.New(errx.CodeInvalid, "qty must be > 0")
	}

	o, err := uc.OrderRepo.GetByID(ctx, in.OrderID)
	if err != nil {
		return nil, err
	}
	if o.UserID != in.UserID {
		return nil, errx.New(errx.CodeForbidden, "order does not belong to user")
	}
	if o.Status != entity.OrderCreated {
		return nil, errx.New(errx.CodeConflict, "order is not editable")
	}

	found := false
	for i := range o.Items {
		if o.Items[i].ID == in.ItemID {
			o.Items[i].Qty = in.Qty
			found = true
			break
		}
	}
	if !found {
		return nil, errx.New(errx.CodeNotFound, "order item not found")
	}

	o.UpdatedAt = time.Now()
	o.RecalculateTotals()

	if err := uc.OrderRepo.Update(ctx, o); err != nil {
		return nil, err
	}

	return toOrderDTO(o), nil
}
