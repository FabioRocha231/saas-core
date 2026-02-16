package handlers

import (
	"net/http"
	"strings"

	"github.com/FabioRocha231/saas-core/internal/domain/entity"
	"github.com/FabioRocha231/saas-core/internal/domain/errx"
	"github.com/FabioRocha231/saas-core/internal/infra/http/helper"
	ports "github.com/FabioRocha231/saas-core/internal/port"
	"github.com/FabioRocha231/saas-core/internal/port/repository"
	usecase "github.com/FabioRocha231/saas-core/internal/usecase/payment"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	orderRepo   repository.OrderRepository
	paymentRepo repository.PaymentRepository
	uuid        ports.UUIDInterface
}

func NewPaymentHandler(
	orderRepo repository.OrderRepository,
	paymentRepo repository.PaymentRepository,
	uuid ports.UUIDInterface,
) *PaymentHandler {
	return &PaymentHandler{
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
		uuid:        uuid,
	}
}

type CreatePaymentRequest struct {
	// por enquanto só MOCK, mas já deixamos a estrutura
	Method string `json:"method"`
	// para idempotência (opcional)
	IdempotencyKey string `json:"idempotency_key"`
}

func (h *PaymentHandler) CreateForOrder(ctx *gin.Context) {
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

	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, err)
		return
	}

	method := strings.TrimSpace(req.Method)
	if method == "" {
		method = string(entity.PaymentMethodMock)
	}
	if method != string(entity.PaymentMethodMock) {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "unsupported payment method (use MOCK)"))
		return
	}

	uc := usecase.NewCreatePaymentUsecase(h.orderRepo, h.paymentRepo, h.uuid)

	out, err := uc.Execute(ctx, usecase.CreatePaymentInput{
		OrderID:        orderID,
		UserID:         userID,
		Method:         entity.PaymentMethod(method),
		IdempotencyKey: strings.TrimSpace(req.IdempotencyKey),
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusCreated, out)
}

func (h *PaymentHandler) GetByID(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	paymentID := strings.TrimSpace(ctx.Param("paymentId"))
	if paymentID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing paymentId"))
		return
	}

	p, err := h.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		RespondErr(ctx, err)
		return
	}
	if p.UserID != userID {
		RespondErr(ctx, errx.New(errx.CodeForbidden, "payment does not belong to user"))
		return
	}

	RespondOK(ctx, http.StatusOK, p)
}

type FailPaymentRequest struct {
	Reason string `json:"reason"`
}

func (h *PaymentHandler) Confirm(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	paymentID := strings.TrimSpace(ctx.Param("paymentId"))
	if paymentID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing paymentId"))
		return
	}

	uc := usecase.NewConfirmPaymentUsecase(h.orderRepo, h.paymentRepo, h.uuid)

	out, err := uc.Execute(ctx, usecase.ConfirmPaymentInput{
		PaymentID: paymentID,
		UserID:    userID,
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}

func (h *PaymentHandler) Fail(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	paymentID := strings.TrimSpace(ctx.Param("paymentId"))
	if paymentID == "" {
		RespondErr(ctx, errx.New(errx.CodeInvalid, "missing paymentId"))
		return
	}

	var req FailPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		RespondErr(ctx, err)
		return
	}

	uc := usecase.NewFailPaymentUsecase(h.paymentRepo, h.uuid)

	out, err := uc.Execute(ctx, usecase.FailPaymentInput{
		PaymentID: paymentID,
		UserID:    userID,
		Reason:    strings.TrimSpace(req.Reason),
	})
	if err != nil {
		RespondErr(ctx, err)
		return
	}

	RespondOK(ctx, http.StatusOK, out)
}
