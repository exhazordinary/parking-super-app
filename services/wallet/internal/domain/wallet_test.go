package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestNewWallet(t *testing.T) {
	userID := uuid.New()
	currency := "MYR"

	wallet := NewWallet(userID, currency)

	if wallet.ID == uuid.Nil {
		t.Error("expected wallet ID to be set")
	}
	if wallet.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, wallet.UserID)
	}
	if wallet.Currency != currency {
		t.Errorf("expected currency %s, got %s", currency, wallet.Currency)
	}
	if !wallet.Balance.Equal(decimal.Zero) {
		t.Errorf("expected zero balance, got %s", wallet.Balance.String())
	}
	if wallet.Status != WalletStatusActive {
		t.Errorf("expected status active, got %s", wallet.Status)
	}
}

func TestWallet_Credit(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  decimal.Decimal
		amount          decimal.Decimal
		status          WalletStatus
		expectedBalance decimal.Decimal
		expectedErr     error
	}{
		{
			name:            "successful credit",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.NewFromFloat(50.00),
			status:          WalletStatusActive,
			expectedBalance: decimal.NewFromFloat(150.00),
			expectedErr:     nil,
		},
		{
			name:            "credit zero amount fails",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.Zero,
			status:          WalletStatusActive,
			expectedBalance: decimal.NewFromFloat(100.00),
			expectedErr:     ErrInvalidAmount,
		},
		{
			name:            "credit negative amount fails",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.NewFromFloat(-50.00),
			status:          WalletStatusActive,
			expectedBalance: decimal.NewFromFloat(100.00),
			expectedErr:     ErrInvalidAmount,
		},
		{
			name:            "credit to inactive wallet fails",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.NewFromFloat(50.00),
			status:          WalletStatusInactive,
			expectedBalance: decimal.NewFromFloat(100.00),
			expectedErr:     ErrWalletInactive,
		},
		{
			name:            "credit to frozen wallet fails",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.NewFromFloat(50.00),
			status:          WalletStatusFrozen,
			expectedBalance: decimal.NewFromFloat(100.00),
			expectedErr:     ErrWalletInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet := &Wallet{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				Balance:  tt.initialBalance,
				Currency: "MYR",
				Status:   tt.status,
			}

			err := wallet.Credit(tt.amount)

			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if !wallet.Balance.Equal(tt.expectedBalance) {
				t.Errorf("expected balance %s, got %s", tt.expectedBalance.String(), wallet.Balance.String())
			}
		})
	}
}

func TestWallet_Debit(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  decimal.Decimal
		amount          decimal.Decimal
		status          WalletStatus
		expectedBalance decimal.Decimal
		expectedErr     error
	}{
		{
			name:            "successful debit",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.NewFromFloat(50.00),
			status:          WalletStatusActive,
			expectedBalance: decimal.NewFromFloat(50.00),
			expectedErr:     nil,
		},
		{
			name:            "debit exact balance succeeds",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.NewFromFloat(100.00),
			status:          WalletStatusActive,
			expectedBalance: decimal.Zero,
			expectedErr:     nil,
		},
		{
			name:            "debit more than balance fails",
			initialBalance:  decimal.NewFromFloat(50.00),
			amount:          decimal.NewFromFloat(100.00),
			status:          WalletStatusActive,
			expectedBalance: decimal.NewFromFloat(50.00),
			expectedErr:     ErrInsufficientBalance,
		},
		{
			name:            "debit zero amount fails",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.Zero,
			status:          WalletStatusActive,
			expectedBalance: decimal.NewFromFloat(100.00),
			expectedErr:     ErrInvalidAmount,
		},
		{
			name:            "debit from frozen wallet fails",
			initialBalance:  decimal.NewFromFloat(100.00),
			amount:          decimal.NewFromFloat(50.00),
			status:          WalletStatusFrozen,
			expectedBalance: decimal.NewFromFloat(100.00),
			expectedErr:     ErrWalletInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet := &Wallet{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				Balance:  tt.initialBalance,
				Currency: "MYR",
				Status:   tt.status,
			}

			err := wallet.Debit(tt.amount)

			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if !wallet.Balance.Equal(tt.expectedBalance) {
				t.Errorf("expected balance %s, got %s", tt.expectedBalance.String(), wallet.Balance.String())
			}
		})
	}
}

func TestWallet_HasSufficientBalance(t *testing.T) {
	wallet := &Wallet{
		Balance: decimal.NewFromFloat(100.00),
	}

	if !wallet.HasSufficientBalance(decimal.NewFromFloat(50.00)) {
		t.Error("expected sufficient balance for 50.00")
	}
	if !wallet.HasSufficientBalance(decimal.NewFromFloat(100.00)) {
		t.Error("expected sufficient balance for 100.00")
	}
	if wallet.HasSufficientBalance(decimal.NewFromFloat(150.00)) {
		t.Error("expected insufficient balance for 150.00")
	}
}

func TestWallet_Freeze(t *testing.T) {
	wallet := NewWallet(uuid.New(), "MYR")

	wallet.Freeze()

	if wallet.Status != WalletStatusFrozen {
		t.Errorf("expected status frozen, got %s", wallet.Status)
	}
	if !wallet.CanTransact() == false {
		t.Error("frozen wallet should not be able to transact")
	}
}

func TestWallet_Activate(t *testing.T) {
	wallet := &Wallet{
		ID:     uuid.New(),
		Status: WalletStatusFrozen,
	}

	wallet.Activate()

	if wallet.Status != WalletStatusActive {
		t.Errorf("expected status active, got %s", wallet.Status)
	}
	if !wallet.CanTransact() {
		t.Error("active wallet should be able to transact")
	}
}
