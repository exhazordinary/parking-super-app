package external

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/ports"
	"github.com/shopspring/decimal"
)

// MockWalletClient simulates wallet service calls for development
type MockWalletClient struct{}

func NewMockWalletClient() *MockWalletClient {
	return &MockWalletClient{}
}

func (c *MockWalletClient) Pay(ctx context.Context, req ports.PaymentRequest) (*ports.PaymentResponse, error) {
	return &ports.PaymentResponse{
		TransactionID: uuid.New(),
		Status:        "completed",
	}, nil
}

func (c *MockWalletClient) GetWallet(ctx context.Context, userID uuid.UUID) (*ports.WalletInfo, error) {
	return &ports.WalletInfo{
		ID:       uuid.New(),
		UserID:   userID,
		Balance:  decimal.NewFromFloat(100.00),
		Currency: "MYR",
		Status:   "active",
	}, nil
}
