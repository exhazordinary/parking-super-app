package external

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/ports"
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
