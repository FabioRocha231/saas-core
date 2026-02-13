package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type RemoveItemInput struct {
	OrderID string
	UserID  string
	ItemID  string
}

type RemoveItemUsecase struct {
	OrderRepo repository.OrderRepository
	UUID      ports.UUIDInterface
}

func NewRemoveItemUsecase(orderRepo repository.OrderRepository, uuid ports.UUIDInterface) *RemoveItemUsecase {
	return &RemoveItemUsecase{OrderRepo: orderRepo, UUID: uuid}
}

func (uc *RemoveItemUsecase) Execute(ctx context.Context, in RemoveItemInput) (*Order, error) {
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

	idx := -1
	for i := range o.Items {
		if o.Items[i].ID == in.ItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, errx.New(errx.CodeNotFound, "order item not found")
	}

	o.Items = append(o.Items[:idx], o.Items[idx+1:]...)

	o.UpdatedAt = time.Now()
	o.RecalculateTotals()

	if err := uc.OrderRepo.Update(ctx, o); err != nil {
		return nil, err
	}

	return toOrderDTO(o), nil
}
