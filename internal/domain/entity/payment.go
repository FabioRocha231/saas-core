package entity

import "time"

type PaymentStatus string
type PaymentMethod string
type PaymentProvider string

const (
	PaymentStatusCreated  PaymentStatus = "CREATED"
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusCanceled PaymentStatus = "CANCELED"
)

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard  PaymentMethod = "DEBIT_CARD"
	PaymentMethodPix        PaymentMethod = "PIX"
	PaymentMethodMock       PaymentMethod = "MOCK"
)

const (
	PaymentProviderMock PaymentProvider = "MOCK"
)

func (p PaymentStatus) String() string {
	return string(p)
}

func (p PaymentMethod) String() string {
	return string(p)
}

func (p PaymentProvider) String() string {
	return string(p)
}

type Payment struct {
	ID      string
	OrderID string
	UserID  string
	StoreID string

	Method   PaymentMethod
	Provider PaymentProvider
	Status   PaymentStatus

	Amount   int64 // centavos (sempre = order.Total no momento da criação)
	Currency string

	IdempotencyKey string

	CreatedAt time.Time
	UpdatedAt time.Time
	PaidAt    *time.Time
}
