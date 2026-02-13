package usecase

import (
	"context"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type GetOrCreateDraftInput struct {
	UserID  string
	StoreID string
	MenuID  string
}

type Variant struct {
	VariantGroupID  string `json:"variant_group_id"`  // ItemVariantGroup.ID
	VariantGroup    string `json:"variant_group"`     // snapshot (ItemVariantGroup.Name)
	VariantOptionID string `json:"variant_option_id"` // VariantOption.ID
	OptionName      string `json:"option_name"`       // snapshot (VariantOption.Name)

	PriceDelta int64 `json:"price_delta"` // snapshot (VariantOption.PriceDelta)
}

type Addon struct {
	AddonGroupID  string `json:"addon_group_id"`  // ItemAddonGroup.ID
	AddonGroup    string `json:"addon_group"`     // snapshot (ItemAddonGroup.Name)
	AddonOptionID string `json:"addon_option_id"` // AddonOption.ID
	OptionName    string `json:"option_name"`     // snapshot (AddonOption.Name)

	Qty int64 `json:"qty"`

	UnitPrice int64 `json:"unit_price"` // snapshot (AddonOption.Price)
	LineTotal int64 `json:"line_total"` // UnitPrice * Qty
}

type Item struct {
	ID     string `json:"id"`
	ItemID string `json:"item_id"` // CategoryItem.ID (referência)
	Name   string `json:"name"`    // snapshot (CategoryItem.Name)

	Qty int64 `json:"qty"`

	BasePrice int64 `json:"base_price"` // snapshot (CategoryItem.BasePrice)
	LineTotal int64 `json:"line_total"` // (base + addons + delta variants) * qty

	Variants []Variant `json:"variants"`
	Addons   []Addon   `json:"addons"`

	Note string `json:"note"`
}

type Order struct {
	ID      string             `json:"id"`
	StoreID string             `json:"store_id"`
	MenuID  string             `json:"menu_id"`
	UserID  string             `json:"user_id"`
	Status  entity.OrderStatus `json:"status"`
	Items   []Item             `json:"items"`

	Subtotal int64 `json:"subtotal"`
	Fees     int64 `json:"fees"`
	Total    int64 `json:"total"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetOrCreateDraftOutput struct {
	Order   *Order `json:"order"`
	Created bool   `json:"-"`
}

type GetOrCreateDraftUsecase struct {
	ordersRepo repository.OrderRepository
	uuid       ports.UUIDInterface
	context    context.Context
}

func NewGetOrCreateDraftUsecase(
	orders repository.OrderRepository,
	ids ports.UUIDInterface,
	ctx context.Context,
) *GetOrCreateDraftUsecase {
	return &GetOrCreateDraftUsecase{
		ordersRepo: orders,
		uuid:       ids,
		context:    ctx,
	}
}

func (uc *GetOrCreateDraftUsecase) Execute(in GetOrCreateDraftInput) (*GetOrCreateDraftOutput, error) {
	if in.UserID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing userId")
	}
	if in.StoreID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing storeId")
	}

	// 1) tenta pegar draft ativo
	o, err := uc.ordersRepo.GetActiveDraftByUserStore(uc.context, in.UserID, in.StoreID)
	if err == nil && o != nil {
		return &GetOrCreateDraftOutput{Order: toOrderDTO(o), Created: false}, nil
	}
	if err != nil && !errx.Is(err, errx.CodeNotFound) {
		return nil, err
	}

	// 2) cria novo draft
	now := time.Now()
	newOrder := &entity.Order{
		ID:      uc.uuid.Generate(),
		StoreID: in.StoreID,
		MenuID:  in.MenuID,
		UserID:  in.UserID,
		Status:  entity.OrderCreated,
		Items:   []entity.OrderItem{},

		Subtotal: 0,
		Fees:     0,
		Total:    0,

		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.ordersRepo.Create(uc.context, newOrder); err != nil {
		// Race condition: dois requests tentando criar ao mesmo tempo
		if errx.Is(err, errx.CodeConflict) {
			existing, e2 := uc.ordersRepo.GetActiveDraftByUserStore(uc.context, in.UserID, in.StoreID)
			if e2 == nil && existing != nil {
				return &GetOrCreateDraftOutput{Order: toOrderDTO(existing), Created: true}, nil
			}
		}
		return nil, err
	}

	return &GetOrCreateDraftOutput{Order: toOrderDTO(newOrder), Created: true}, nil
}

func toOrderDTO(e *entity.Order) *Order {
	if e == nil {
		return nil
	}

	items := make([]Item, len(e.Items))
	for i, it := range e.Items {
		variants := make([]Variant, len(it.Variants))
		for j, v := range it.Variants {
			variants[j] = Variant{
				VariantGroupID:  v.VariantGroupID,
				VariantGroup:    v.VariantGroup,
				VariantOptionID: v.VariantOptionID,
				OptionName:      v.OptionName,
				PriceDelta:      int64(v.PriceDelta),
			}
		}

		addons := make([]Addon, len(it.Addons))
		for k, a := range it.Addons {
			addons[k] = Addon{
				AddonGroupID:  a.AddonGroupID,
				AddonGroup:    a.AddonGroup,
				AddonOptionID: a.AddonOptionID,
				OptionName:    a.OptionName,
				Qty:           a.Qty,
				UnitPrice:     int64(a.UnitPrice),
				LineTotal:     int64(a.LineTotal),
			}
		}

		items[i] = Item{
			ID:        it.ID,
			ItemID:    it.ItemID,
			Name:      it.Name,
			Qty:       it.Qty,
			BasePrice: int64(it.BasePrice),
			LineTotal: int64(it.LineTotal),
			Variants:  variants,
			Addons:    addons,
			Note:      it.Note,
		}
	}

	return &Order{
		ID:        e.ID,
		StoreID:   e.StoreID,
		MenuID:    e.MenuID,
		UserID:    e.UserID,
		Status:    e.Status,
		Items:     items,
		Subtotal:  int64(e.Subtotal),
		Fees:      int64(e.Fees),
		Total:     int64(e.Total),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
