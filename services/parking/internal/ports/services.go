package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

func String(key, value string) Field { return Field{Key: key, Value: value} }
func Err(err error) Field            { return Field{Key: "error", Value: err} }
func Any(key string, val interface{}) Field { return Field{Key: key, Value: val} }

// EventPublisher for domain events
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type Event struct {
	Type    string
	Payload map[string]interface{}
}

const (
	EventSessionStarted   = "parking.session.started"
	EventSessionEnded     = "parking.session.ended"
	EventSessionCancelled = "parking.session.cancelled"
	EventPaymentRequired  = "parking.payment.required"
)

// ProviderClient communicates with parking provider APIs
type ProviderClient interface {
	StartSession(ctx context.Context, req StartSessionRequest) (*StartSessionResponse, error)
	EndSession(ctx context.Context, req EndSessionRequest) (*EndSessionResponse, error)
	GetSessionStatus(ctx context.Context, providerID uuid.UUID, externalSessionID string) (*SessionStatusResponse, error)
}

type StartSessionRequest struct {
	ProviderID   uuid.UUID
	LocationID   uuid.UUID
	VehiclePlate string
	VehicleType  string
	UserRef      string
}

type StartSessionResponse struct {
	ExternalSessionID string
	EntryTime         string
	Status            string
}

type EndSessionRequest struct {
	ProviderID        uuid.UUID
	ExternalSessionID string
}

type EndSessionResponse struct {
	ExitTime string
	Duration int
	Amount   decimal.Decimal
	Currency string
}

type SessionStatusResponse struct {
	Status   string
	Duration int
	Amount   decimal.Decimal
}

// WalletClient for payment operations
type WalletClient interface {
	Pay(ctx context.Context, req PaymentRequest) (*PaymentResponse, error)
	GetWallet(ctx context.Context, userID uuid.UUID) (*WalletInfo, error)
}

type PaymentRequest struct {
	WalletID       uuid.UUID
	Amount         decimal.Decimal
	ProviderID     uuid.UUID
	ReferenceID    string
	Description    string
	IdempotencyKey string
}

type PaymentResponse struct {
	TransactionID uuid.UUID
	Status        string
}

type WalletInfo struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Balance  decimal.Decimal
	Currency string
	Status   string
}
