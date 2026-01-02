package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrWalletAlreadyExists  = errors.New("wallet already exists for this user")
	ErrWalletInactive       = errors.New("wallet is inactive")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrDuplicateTransaction = errors.New("duplicate transaction")
)

type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "active"
	WalletStatusInactive WalletStatus = "inactive"
	WalletStatusFrozen   WalletStatus = "frozen"
)

type Wallet struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	Status    WalletStatus    `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func NewWallet(userID uuid.UUID, currency string) *Wallet {
	now := time.Now().UTC()
	return &Wallet{
		ID:        uuid.New(),
		UserID:    userID,
		Balance:   decimal.Zero,
		Currency:  currency,
		Status:    WalletStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (w *Wallet) IsActive() bool {
	return w.Status == WalletStatusActive
}

func (w *Wallet) CanTransact() bool {
	return w.Status == WalletStatusActive
}

func (w *Wallet) Credit(amount decimal.Decimal) error {
	if !w.CanTransact() {
		return ErrWalletInactive
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	w.Balance = w.Balance.Add(amount)
	w.UpdatedAt = time.Now().UTC()
	return nil
}

func (w *Wallet) Debit(amount decimal.Decimal) error {
	if !w.CanTransact() {
		return ErrWalletInactive
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	if w.Balance.LessThan(amount) {
		return ErrInsufficientBalance
	}
	w.Balance = w.Balance.Sub(amount)
	w.UpdatedAt = time.Now().UTC()
	return nil
}

func (w *Wallet) HasSufficientBalance(amount decimal.Decimal) bool {
	return w.Balance.GreaterThanOrEqual(amount)
}

func (w *Wallet) Freeze() {
	w.Status = WalletStatusFrozen
	w.UpdatedAt = time.Now().UTC()
}

func (w *Wallet) Activate() {
	w.Status = WalletStatusActive
	w.UpdatedAt = time.Now().UTC()
}
