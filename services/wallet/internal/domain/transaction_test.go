package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestNewTransaction(t *testing.T) {
	walletID := uuid.New()
	amount := decimal.NewFromFloat(100.00)
	balanceBefore := decimal.NewFromFloat(500.00)

	tx := NewTransaction(
		walletID,
		TransactionTypeTopUp,
		amount,
		balanceBefore,
		"ref-123",
		"idem-key-456",
		"Test top-up",
	)

	if tx.ID == uuid.Nil {
		t.Error("expected transaction ID to be set")
	}
	if tx.WalletID != walletID {
		t.Errorf("expected walletID %v, got %v", walletID, tx.WalletID)
	}
	if tx.Type != TransactionTypeTopUp {
		t.Errorf("expected type topup, got %s", tx.Type)
	}
	if !tx.Amount.Equal(amount) {
		t.Errorf("expected amount %s, got %s", amount.String(), tx.Amount.String())
	}
	if !tx.BalanceBefore.Equal(balanceBefore) {
		t.Errorf("expected balance before %s, got %s", balanceBefore.String(), tx.BalanceBefore.String())
	}
	if !tx.BalanceAfter.Equal(balanceBefore) {
		t.Error("expected balance after to equal balance before initially")
	}
	if tx.Status != TransactionStatusPending {
		t.Errorf("expected status pending, got %s", tx.Status)
	}
	if tx.ReferenceID != "ref-123" {
		t.Errorf("expected reference ID ref-123, got %s", tx.ReferenceID)
	}
	if tx.IdempotencyKey != "idem-key-456" {
		t.Errorf("expected idempotency key idem-key-456, got %s", tx.IdempotencyKey)
	}
}

func TestTransaction_Complete(t *testing.T) {
	tx := NewTransaction(
		uuid.New(),
		TransactionTypePayment,
		decimal.NewFromFloat(50.00),
		decimal.NewFromFloat(100.00),
		"",
		"",
		"Payment",
	)

	newBalance := decimal.NewFromFloat(50.00)
	tx.Complete(newBalance)

	if tx.Status != TransactionStatusCompleted {
		t.Errorf("expected status completed, got %s", tx.Status)
	}
	if !tx.BalanceAfter.Equal(newBalance) {
		t.Errorf("expected balance after %s, got %s", newBalance.String(), tx.BalanceAfter.String())
	}
	if !tx.IsCompleted() {
		t.Error("expected IsCompleted to return true")
	}
}

func TestTransaction_Fail(t *testing.T) {
	tx := NewTransaction(
		uuid.New(),
		TransactionTypePayment,
		decimal.NewFromFloat(50.00),
		decimal.NewFromFloat(100.00),
		"",
		"",
		"Payment",
	)

	tx.Fail()

	if tx.Status != TransactionStatusFailed {
		t.Errorf("expected status failed, got %s", tx.Status)
	}
	if tx.IsPending() {
		t.Error("expected IsPending to return false after failure")
	}
}

func TestTransaction_SetProvider(t *testing.T) {
	tx := NewTransaction(
		uuid.New(),
		TransactionTypePayment,
		decimal.NewFromFloat(50.00),
		decimal.NewFromFloat(100.00),
		"",
		"",
		"Payment",
	)

	providerID := uuid.New()
	tx.SetProvider(providerID)

	if tx.ProviderID == nil {
		t.Error("expected provider ID to be set")
	}
	if *tx.ProviderID != providerID {
		t.Errorf("expected provider ID %v, got %v", providerID, *tx.ProviderID)
	}
}

func TestTransaction_AddMetadata(t *testing.T) {
	tx := NewTransaction(
		uuid.New(),
		TransactionTypePayment,
		decimal.NewFromFloat(50.00),
		decimal.NewFromFloat(100.00),
		"",
		"",
		"Payment",
	)

	tx.AddMetadata("key1", "value1")
	tx.AddMetadata("key2", "value2")

	if tx.Metadata["key1"] != "value1" {
		t.Errorf("expected metadata key1=value1, got %s", tx.Metadata["key1"])
	}
	if tx.Metadata["key2"] != "value2" {
		t.Errorf("expected metadata key2=value2, got %s", tx.Metadata["key2"])
	}
}

func TestNewPaymentMethod(t *testing.T) {
	userID := uuid.New()

	pm := NewPaymentMethod(userID, "card", "visa", "tok_xxx")

	if pm.ID == uuid.Nil {
		t.Error("expected payment method ID to be set")
	}
	if pm.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, pm.UserID)
	}
	if pm.Type != "card" {
		t.Errorf("expected type card, got %s", pm.Type)
	}
	if pm.Provider != "visa" {
		t.Errorf("expected provider visa, got %s", pm.Provider)
	}
	if pm.Token != "tok_xxx" {
		t.Errorf("expected token tok_xxx, got %s", pm.Token)
	}
	if pm.IsDefault {
		t.Error("expected IsDefault to be false initially")
	}
}
