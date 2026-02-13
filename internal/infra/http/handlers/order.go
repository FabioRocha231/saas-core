package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/infra/http/helper"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/order"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo    repository.OrderRepository
	menuReadRepo repository.MenuReadRepository
	uuid         ports.UUIDInterface
}

type AddItemRequest struct {
	ItemID           string   `json:"item_id"`
	Qty              int64    `json:"qty"`
	VariantOptionIDs []string `json:"variant_option_ids"`
	Addons           []struct {
		OptionID string `json:"option_id"`
		Qty      int64  `json:"qty"`
	} `json:"addons"`
	Note string `json:"note"`
}

type UpdateItemQtyRequest struct {
	Qty int64 `json:"qty"`
}

func NewOrderHandler(
	orderRepo repository.OrderRepository,
	menuReadRepo repository.MenuReadRepository,
	uuid ports.UUIDInterface,
) *OrderHandler {
	return &OrderHandler{
		orderRepo:    orderRepo,
		menuReadRepo: menuReadRepo,
		uuid:         uuid,
	}
}

func (h *OrderHandler) Create(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	storeID := strings.TrimSpace(ctx.Param("storeId"))
	if storeID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing storeId"))
		return
	}

	uc := usecase.NewGetOrCreateDraftUsecase(h.orderRepo, h.uuid, ctx.Request.Context())

	out, err := uc.Execute(usecase.GetOrCreateDraftInput{
		UserID:  userID,
		StoreID: storeID,
		// MenuID: "" // se quiser rastrear depois, a gente injeta aqui
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	if out.Created {
		RespondOK(ctx, http.StatusCreated, out.Order)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}

func (h *OrderHandler) AddItem(ctx *gin.Context) {
	orderID := strings.TrimSpace(ctx.Param("orderId"))
	if orderID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing orderId"))
		return
	}

	var req AddItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, err)
		return
	}

	if strings.TrimSpace(req.ItemID) == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "item_id is required"))
		return
	}
	if req.Qty <= 0 {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "qty must be > 0"))
		return
	}

	uc := usecase.NewAddItem(h.orderRepo, h.menuReadRepo, h.uuid)

	addons := make([]usecase.AddonSelection, 0, len(req.Addons))
	for _, a := range req.Addons {
		addons = append(addons, usecase.AddonSelection{
			OptionID: a.OptionID,
			Qty:      a.Qty,
		})
	}

	out, err := uc.Execute(ctx, usecase.AddItemInput{
		OrderID:          orderID,
		ItemID:           req.ItemID,
		Qty:              req.Qty,
		VariantOptionIDs: req.VariantOptionIDs,
		Addons:           addons,
		Note:             req.Note,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}

func (h *OrderHandler) GetByID(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	orderID := strings.TrimSpace(ctx.Param("orderId"))
	if orderID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing orderId"))
		return
	}

	uc := usecase.NewGetOrderUsecase(h.orderRepo, h.uuid)

	out, err := uc.Execute(ctx, usecase.GetOrderInput{OrderID: orderID, UserID: userID})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}

func (h *OrderHandler) UpdateItemQty(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	orderID := strings.TrimSpace(ctx.Param("orderId"))
	itemID := strings.TrimSpace(ctx.Param("itemId"))
	if orderID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing orderId"))
		return
	}
	if itemID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing itemId"))
		return
	}

	var req UpdateItemQtyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, errx.New(errx.CodeInvalid, err.Error()))
		return
	}

	if req.Qty <= 0 {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "qty must be > 0"))
		return
	}

	uc := usecase.NewUpdateItemQtyUsecase(h.orderRepo, h.uuid)
	out, err := uc.Execute(ctx, usecase.UpdateItemQtyInput{
		UserID:  userID,
		OrderID: orderID,
		ItemID:  itemID,
		Qty:     req.Qty,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}

func (h *OrderHandler) RemoveItem(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	orderID := strings.TrimSpace(ctx.Param("orderId"))
	itemID := strings.TrimSpace(ctx.Param("itemId"))
	if orderID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing orderId"))
		return
	}
	if itemID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing itemId"))
		return
	}

	uc := usecase.NewRemoveItemUsecase(h.orderRepo, h.uuid)
	out, err := uc.Execute(ctx, usecase.RemoveItemInput{
		UserID:  userID,
		OrderID: orderID,
		ItemID:  itemID,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}

func (h *OrderHandler) PlaceOrder(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	orderID := strings.TrimSpace(ctx.Param("orderId"))
	if orderID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing orderId"))
		return
	}

	uc := usecase.NewPlaceOrderUsecase(h.orderRepo, h.uuid)
	out, err := uc.Execute(ctx, usecase.PlaceOrderInput{
		OrderID: orderID,
		UserID:  userID,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}
