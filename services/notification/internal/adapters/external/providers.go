package external

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/notification/internal/ports"
)

// MockPushProvider simulates push notification delivery
type MockPushProvider struct{}

func NewMockPushProvider() *MockPushProvider {
	return &MockPushProvider{}
}

func (p *MockPushProvider) Send(ctx context.Context, req ports.PushRequest) (*ports.PushResponse, error) {
	log.Printf("[PUSH] to=%s title=%s body=%s", req.DeviceToken, req.Title, req.Body)
	return &ports.PushResponse{
		MessageID: uuid.New().String(),
		Success:   true,
	}, nil
}

// MockSMSProvider simulates SMS delivery
type MockSMSProvider struct{}

func NewMockSMSProvider() *MockSMSProvider {
	return &MockSMSProvider{}
}

func (p *MockSMSProvider) Send(ctx context.Context, req ports.SMSRequest) (*ports.SMSResponse, error) {
	log.Printf("[SMS] to=%s message=%s", req.PhoneNumber, req.Message)
	return &ports.SMSResponse{
		MessageID: uuid.New().String(),
		Status:    "sent",
	}, nil
}

// MockEmailProvider simulates email delivery
type MockEmailProvider struct{}

func NewMockEmailProvider() *MockEmailProvider {
	return &MockEmailProvider{}
}

func (p *MockEmailProvider) Send(ctx context.Context, req ports.EmailRequest) (*ports.EmailResponse, error) {
	log.Printf("[EMAIL] to=%s subject=%s", req.To, req.Subject)
	return &ports.EmailResponse{
		MessageID: uuid.New().String(),
		Status:    "sent",
	}, nil
}
