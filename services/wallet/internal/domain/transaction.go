package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeTopUp    TransactionType = "topup"
	TransactionTypePayment  TransactionType = "payment"
	TransactionTypeRefund   TransactionType = "refund"
	TransactionTypeTransfer TransactionType = "transfer"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusRefunded  TransactionStatus = "refunded"
)

type Transaction struct {
	ID             uuid.UUID         `json:"id"`
	WalletID       uuid.UUID         `json:"wallet_id"`
	Type           TransactionType   `json:"type"`
	Amount         decimal.Decimal   `json:"amount"`
	BalanceBefore  decimal.Decimal   `json:"balance_before"`
	BalanceAfter   decimal.Decimal   `json:"balance_after"`
	ReferenceID    string            `json:"reference_id"`
	ProviderID     *uuid.UUID        `json:"provider_id,omitempty"`
	Status         TransactionStatus `json:"status"`
	Description    string            `json:"description"`
	IdempotencyKey string            `json:"idempotency_key"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func NewTransaction(
	walletID uuid.UUID,
	txType TransactionType,
	amount decimal.Decimal,
	balanceBefore decimal.Decimal,
	referenceID string,
	idempotencyKey string,
	description string,
) *Transaction {
	now := time.Now().UTC()
	return &Transaction{
		ID:             uuid.New(),
		WalletID:       walletID,
		Type:           txType,
		Amount:         amount,
		BalanceBefore:  balanceBefore,
		BalanceAfter:   balanceBefore,
		ReferenceID:    referenceID,
		Status:         TransactionStatusPending,
		Description:    description,
		IdempotencyKey: idempotencyKey,
		Metadata:       make(map[string]string),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (t *Transaction) Complete(balanceAfter decimal.Decimal) {
	t.Status = TransactionStatusCompleted
	t.BalanceAfter = balanceAfter
	t.UpdatedAt = time.Now().UTC()
}

func (t *Transaction) Fail() {
	t.Status = TransactionStatusFailed
	t.UpdatedAt = time.Now().UTC()
}

func (t *Transaction) SetProvider(providerID uuid.UUID) {
	t.ProviderID = &providerID
}

func (t *Transaction) AddMetadata(key, value string) {
	if t.Metadata == nil {
		t.Metadata = make(map[string]string)
	}
	t.Metadata[key] = value
}

func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

type PaymentMethod struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Provider  string    `json:"provider"`
	Token     string    `json:"-"`
	LastFour  string    `json:"last_four,omitempty"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPaymentMethod(userID uuid.UUID, methodType, provider, token string) *PaymentMethod {
	return &PaymentMethod{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      methodType,
		Provider:  provider,
		Token:     token,
		IsDefault: false,
		CreatedAt: time.Now().UTC(),
	}
}
