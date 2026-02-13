package usecase

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
)

type AddItemInput struct {
	OrderID string
	ItemID  string
	Qty     int64

	VariantOptionIDs []string
	Addons           []AddonSelection

	Note string
}

type AddonSelection struct {
	OptionID string
	Qty      int64
}

type AddItem struct {
	OrdersRepo repository.OrderRepository
	MenuRepo   repository.MenuReadRepository
	UUID       ports.UUIDInterface
}

func NewAddItem(
	ordersRepo repository.OrderRepository,
	menuRepo repository.MenuReadRepository,
	uuid ports.UUIDInterface,
) *AddItem {
	return &AddItem{OrdersRepo: ordersRepo, MenuRepo: menuRepo, UUID: uuid}
}

func (uc *AddItem) Execute(ctx context.Context, in AddItemInput) (*Order, error) {
	if in.OrderID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing orderId")
	}
	if in.ItemID == "" {
		return nil, errx.New(errx.CodeInvalid, "missing itemId")
	}
	if in.Qty <= 0 {
		return nil, errx.New(errx.CodeInvalid, "qty must be > 0")
	}

	o, err := uc.OrdersRepo.GetByID(ctx, in.OrderID)
	if err != nil {
		return nil, err
	}
	if o.Status != entity.OrderCreated {
		return nil, errx.New(errx.CodeConflict, "order is not editable")
	}

	item, err := uc.MenuRepo.GetCategoryItemByID(ctx, in.ItemID)
	if err != nil {
		return nil, err
	}
	if !item.IsActive {
		return nil, errx.New(errx.CodeConflict, "item is inactive")
	}

	// --- Carrega grupos do item (pra validar pertencimento + regras)
	addonGroups, err := uc.MenuRepo.ListItemAddonGroupsByItemID(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	addonGroupSet := make(map[string]*entity.ItemAddonGroup, len(addonGroups))
	for _, g := range addonGroups {
		if g != nil && g.IsActive {
			addonGroupSet[g.ID] = g
		}
	}

	variantGroups, err := uc.MenuRepo.ListItemVariantGroupsByItemID(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	variantGroupSet := make(map[string]*entity.ItemVariantGroup, len(variantGroups))
	for _, g := range variantGroups {
		if g != nil && g.IsActive {
			variantGroupSet[g.ID] = g
		}
	}

	// --- Monta snapshot de Variants + valida pertencimento
	var (
		variants    []entity.OrderItemVariant
		varCountByG = map[string]int{}
	)

	for _, optID := range uniqueStrings(in.VariantOptionIDs) {
		if strings.TrimSpace(optID) == "" {
			continue
		}

		opt, err := uc.MenuRepo.GetVariantOptionByID(ctx, optID)
		if err != nil {
			return nil, err
		}
		if !opt.IsActive {
			return nil, errx.New(errx.CodeConflict, "variant option is inactive")
		}

		g, ok := variantGroupSet[opt.VariantGroupID]
		if !ok || g == nil {
			return nil, errx.New(errx.CodeInvalid, "variant option not allowed for item")
		}

		// snapshot
		variants = append(variants, entity.OrderItemVariant{
			VariantGroupID:  g.ID,
			VariantGroup:    g.Name,
			VariantOptionID: opt.ID,
			OptionName:      opt.Name,
			PriceDelta:      entity.MoneyCents(opt.PriceDelta),
		})

		varCountByG[g.ID]++
	}

	// --- Monta snapshot de Addons + valida pertencimento
	var (
		addons      []entity.OrderItemAddon
		addCountByG = map[string]int{}
	)

	for _, a := range in.Addons {
		if strings.TrimSpace(a.OptionID) == "" {
			return nil, errx.New(errx.CodeInvalid, "missing addon optionId")
		}
		if a.Qty <= 0 {
			return nil, errx.New(errx.CodeInvalid, "addon qty must be > 0")
		}

		opt, err := uc.MenuRepo.GetAddonOptionByID(ctx, a.OptionID)
		if err != nil {
			return nil, err
		}
		if !opt.IsActive {
			return nil, errx.New(errx.CodeConflict, "addon option is inactive")
		}

		g, ok := addonGroupSet[opt.AddonGroupID]
		if !ok || g == nil {
			return nil, errx.New(errx.CodeInvalid, "addon option not allowed for item")
		}

		addons = append(addons, entity.OrderItemAddon{
			AddonGroupID:  g.ID,
			AddonGroup:    g.Name,
			AddonOptionID: opt.ID,
			OptionName:    opt.Name,
			Qty:           a.Qty,
			UnitPrice:     entity.MoneyCents(opt.Price),
		})

		// regra min/max normalmente conta "opções selecionadas", não qty
		addCountByG[g.ID]++
	}

	// --- Valida regras por grupo (Required/Min/Max)
	if err := validateVariantGroups(variantGroupSet, varCountByG); err != nil {
		return nil, err
	}
	if err := validateAddonGroups(addonGroupSet, addCountByG); err != nil {
		return nil, err
	}

	// --- Merge automático: se mesma combinação (item + variants + addons + note), soma qty
	newSig := signature(item.ID, variants, addons, in.Note)

	for i := range o.Items {
		if signatureFromExisting(o.Items[i]) == newSig {
			o.Items[i].Qty += in.Qty
			o.UpdatedAt = time.Now()
			o.RecalculateTotals()
			if err := uc.OrdersRepo.Update(ctx, o); err != nil {
				return nil, err
			}
			return toOrderDTO(o), nil
		}
	}

	// --- Se não existe igual, cria linha nova
	newItem := entity.OrderItem{
		ID:        uc.UUID.Generate(),
		ItemID:    item.ID,
		Name:      item.Name,
		Qty:       in.Qty,
		BasePrice: entity.MoneyCents(item.BasePrice),
		Variants:  variants,
		Addons:    addons,
		Note:      in.Note,
	}

	o.Items = append(o.Items, newItem)
	o.UpdatedAt = time.Now()
	o.RecalculateTotals()

	if err := uc.OrdersRepo.Update(ctx, o); err != nil {
		return nil, err
	}
	return toOrderDTO(o), nil
}

func validateVariantGroups(groups map[string]*entity.ItemVariantGroup, count map[string]int) error {
	for id, g := range groups {
		if g == nil || !g.IsActive {
			continue
		}
		c := count[id]

		// Required + MinSelect
		if g.Required && c == 0 {
			return errx.New(errx.CodeInvalid, "missing required variant selection")
		}
		if g.MinSelect > 0 && c < g.MinSelect {
			return errx.New(errx.CodeInvalid, "variant selection below minimum")
		}
		// MaxSelect == 0 => sem limite
		if g.MaxSelect > 0 && c > g.MaxSelect {
			return errx.New(errx.CodeInvalid, "variant selection above maximum")
		}
	}
	return nil
}

func validateAddonGroups(groups map[string]*entity.ItemAddonGroup, count map[string]int) error {
	for id, g := range groups {
		if g == nil || !g.IsActive {
			continue
		}
		c := count[id]

		if g.Required && c == 0 {
			return errx.New(errx.CodeInvalid, "missing required addon selection")
		}
		if g.MinSelect > 0 && c < g.MinSelect {
			return errx.New(errx.CodeInvalid, "addon selection below minimum")
		}
		if g.MaxSelect > 0 && c > g.MaxSelect {
			return errx.New(errx.CodeInvalid, "addon selection above maximum")
		}
	}
	return nil
}

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// Assinatura determinística (ordena selections)
func signature(itemID string, vars []entity.OrderItemVariant, adds []entity.OrderItemAddon, note string) string {
	varIDs := make([]string, 0, len(vars))
	for _, v := range vars {
		varIDs = append(varIDs, v.VariantOptionID)
	}
	sort.Strings(varIDs)

	addParts := make([]string, 0, len(adds))
	for _, a := range adds {
		addParts = append(addParts, a.AddonOptionID+":"+strconv.FormatInt(a.Qty, 10))
	}
	sort.Strings(addParts)

	return itemID + "||v:" + strings.Join(varIDs, ",") + "||a:" + strings.Join(addParts, ",") + "||n:" + strings.TrimSpace(note)
}

func signatureFromExisting(it entity.OrderItem) string {
	return signature(it.ItemID, it.Variants, it.Addons, it.Note)
}
