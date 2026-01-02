package ports

import (
	"context"

	"github.com/shopspring/decimal"
)

type PaymentGateway interface {
	ProcessTopUp(ctx context.Context, req TopUpRequest) (*TopUpResponse, error)
	ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error)
	ProcessRefund(ctx context.Context, req RefundRequest) (*RefundResponse, error)
}

type TopUpRequest struct {
	Amount        decimal.Decimal
	Currency      string
	PaymentMethod string
	Token         string
	UserID        string
	IdempotencyKey string
}

type TopUpResponse struct {
	TransactionID string
	Status        string
	Message       string
}

type PaymentRequest struct {
	Amount         decimal.Decimal
	Currency       string
	Description    string
	ReferenceID    string
	IdempotencyKey string
}

type PaymentResponse struct {
	TransactionID string
	Status        string
	Message       string
}

type RefundRequest struct {
	OriginalTransactionID string
	Amount                decimal.Decimal
	Reason                string
}

type RefundResponse struct {
	RefundID string
	Status   string
	Message  string
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type Event struct {
	Type      string
	Payload   map[string]interface{}
}

const (
	EventWalletCreated    = "wallet.created"
	EventTopUpCompleted   = "wallet.topup.completed"
	EventPaymentCompleted = "wallet.payment.completed"
	EventRefundCompleted  = "wallet.refund.completed"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func String(key, value string) Field { return Field{Key: key, Value: value} }
func Err(err error) Field            { return Field{Key: "error", Value: err} }
func Any(key string, value interface{}) Field { return Field{Key: key, Value: value} }
