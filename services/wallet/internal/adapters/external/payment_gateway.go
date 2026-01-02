package external

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/wallet/internal/ports"
)

// MockPaymentGateway simulates payment gateway operations for development.
// In production, replace with actual payment gateway integrations (FPX, iPay88, etc.)
type MockPaymentGateway struct{}

func NewMockPaymentGateway() *MockPaymentGateway {
	return &MockPaymentGateway{}
}

func (g *MockPaymentGateway) ProcessTopUp(ctx context.Context, req ports.TopUpRequest) (*ports.TopUpResponse, error) {
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	return &ports.TopUpResponse{
		TransactionID: uuid.New().String(),
		Status:        "success",
		Message:       "Top-up processed successfully",
	}, nil
}

func (g *MockPaymentGateway) ProcessPayment(ctx context.Context, req ports.PaymentRequest) (*ports.PaymentResponse, error) {
	time.Sleep(100 * time.Millisecond)

	return &ports.PaymentResponse{
		TransactionID: uuid.New().String(),
		Status:        "success",
		Message:       "Payment processed successfully",
	}, nil
}

func (g *MockPaymentGateway) ProcessRefund(ctx context.Context, req ports.RefundRequest) (*ports.RefundResponse, error) {
	time.Sleep(100 * time.Millisecond)

	return &ports.RefundResponse{
		RefundID: uuid.New().String(),
		Status:   "success",
		Message:  "Refund processed successfully",
	}, nil
}
