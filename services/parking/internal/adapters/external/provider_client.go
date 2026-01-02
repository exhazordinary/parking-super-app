package external

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/ports"
	"github.com/shopspring/decimal"
)

// MockProviderClient simulates provider API calls for development
type MockProviderClient struct{}

func NewMockProviderClient() *MockProviderClient {
	return &MockProviderClient{}
}

func (c *MockProviderClient) StartSession(ctx context.Context, req ports.StartSessionRequest) (*ports.StartSessionResponse, error) {
	return &ports.StartSessionResponse{
		ExternalSessionID: uuid.New().String(),
		EntryTime:         time.Now().UTC().Format(time.RFC3339),
		Status:            "active",
	}, nil
}

func (c *MockProviderClient) EndSession(ctx context.Context, req ports.EndSessionRequest) (*ports.EndSessionResponse, error) {
	// Simulate calculating parking fee
	return &ports.EndSessionResponse{
		ExitTime: time.Now().UTC().Format(time.RFC3339),
		Duration: 60, // 1 hour
		Amount:   decimal.NewFromFloat(5.00),
		Currency: "MYR",
	}, nil
}

func (c *MockProviderClient) GetSessionStatus(ctx context.Context, providerID uuid.UUID, externalSessionID string) (*ports.SessionStatusResponse, error) {
	return &ports.SessionStatusResponse{
		Status:   "active",
		Duration: 30,
		Amount:   decimal.NewFromFloat(2.50),
	}, nil
}
