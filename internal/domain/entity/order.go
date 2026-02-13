package entity

import "time"

type OrderStatus string

const (
	OrderCreated  OrderStatus = "CREATED"
	OrderPlaced   OrderStatus = "PLACED"
	OrderPaid     OrderStatus = "PAID"
	OrderCanceled OrderStatus = "CANCELED"
)

// Dinheiro SEMPRE em centavos (int64)
type MoneyCents int64

type Order struct {
	ID      string
	StoreID string
	MenuID  string // rastreio do cardápio (opcional, mas útil)
	UserID  string

	Status OrderStatus

	Items []OrderItem

	Subtotal MoneyCents
	Fees     MoneyCents
	Total    MoneyCents

	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ID     string
	ItemID string // CategoryItem.ID (referência)
	Name   string // snapshot (CategoryItem.Name)

	Qty int64

	BasePrice MoneyCents // snapshot (CategoryItem.BasePrice)
	LineTotal MoneyCents // (base + addons + delta variants) * qty

	Variants []OrderItemVariant
	Addons   []OrderItemAddon

	Note string
}

type OrderItemVariant struct {
	VariantGroupID  string // ItemVariantGroup.ID
	VariantGroup    string // snapshot (ItemVariantGroup.Name)
	VariantOptionID string // VariantOption.ID
	OptionName      string // snapshot (VariantOption.Name)

	PriceDelta MoneyCents // snapshot (VariantOption.PriceDelta)
}

type OrderItemAddon struct {
	AddonGroupID  string // ItemAddonGroup.ID
	AddonGroup    string // snapshot (ItemAddonGroup.Name)
	AddonOptionID string // AddonOption.ID
	OptionName    string // snapshot (AddonOption.Name)

	Qty int64

	UnitPrice MoneyCents // snapshot (AddonOption.Price)
	LineTotal MoneyCents // UnitPrice * Qty
}

func (o *Order) RecalculateTotals() {
	var subtotal MoneyCents = 0

	for i := range o.Items {
		it := &o.Items[i]

		// addons por unidade
		var addonsPerUnit MoneyCents = 0
		for j := range it.Addons {
			ad := &it.Addons[j]
			ad.LineTotal = MoneyCents(int64(ad.UnitPrice) * ad.Qty)
			addonsPerUnit += ad.LineTotal
		}

		// variants delta por unidade
		var variantDelta MoneyCents = 0
		for _, v := range it.Variants {
			variantDelta += v.PriceDelta
		}

		unit := it.BasePrice + addonsPerUnit + variantDelta
		it.LineTotal = MoneyCents(int64(unit) * it.Qty)
		subtotal += it.LineTotal
	}

	o.Subtotal = subtotal
	o.Total = o.Subtotal + o.Fees
}
