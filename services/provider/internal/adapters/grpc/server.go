package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/provider/internal/application"
	"github.com/parking-super-app/services/provider/internal/domain"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProviderServiceServer implements the gRPC ProviderService
type ProviderServiceServer struct {
	providerService *application.ProviderService
}

// NewProviderServiceServer creates a new gRPC server for the provider service
func NewProviderServiceServer(ps *application.ProviderService) *ProviderServiceServer {
	return &ProviderServiceServer{
		providerService: ps,
	}
}

// Request/Response types for gRPC

type StartSessionRequest struct {
	ProviderID   string
	LocationID   string
	VehiclePlate string
	VehicleType  string
	UserRef      string
}

type StartSessionResponse struct {
	ExternalSessionID string
	EntryTime         string
	Status            string
	ErrorMessage      string
}

type EndSessionRequest struct {
	ProviderID        string
	ExternalSessionID string
}

type EndSessionResponse struct {
	ExitTime        string
	DurationMinutes int32
	Amount          string
	Currency        string
	ErrorMessage    string
}

type GetSessionStatusRequest struct {
	ProviderID        string
	ExternalSessionID string
}

type SessionStatusResponse struct {
	Status          string
	DurationMinutes int32
	CurrentAmount   string
	Currency        string
	EntryTime       string
}

type GetProviderRequest struct {
	ID string
}

type ProviderResponse struct {
	ID         string
	Name       string
	Code       string
	Status     string
	MFEURL     string
	APIBaseURL string
	LogoURL    string
	CreatedAt  string
	UpdatedAt  string
}

type ListProvidersRequest struct {
	ActiveOnly bool
	Limit      int32
	Offset     int32
}

type ListProvidersResponse struct {
	Providers []*ProviderResponse
	Total     int32
}

// StartSession initiates a parking session with the provider
// This simulates the provider's API - in production this would call the actual provider
func (s *ProviderServiceServer) StartSession(ctx context.Context, req *StartSessionRequest) (*StartSessionResponse, error) {
	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provider_id")
	}

	// Verify provider exists and is active
	provider, err := s.providerService.GetProvider(ctx, providerID)
	if err != nil {
		if err == domain.ErrProviderNotFound {
			return nil, status.Error(codes.NotFound, "provider not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if provider.Status != "active" {
		return nil, status.Error(codes.FailedPrecondition, "provider is not active")
	}

	// Generate external session ID (simulating provider's system)
	externalSessionID := uuid.New().String()
	entryTime := time.Now().UTC()

	return &StartSessionResponse{
		ExternalSessionID: externalSessionID,
		EntryTime:         entryTime.Format(time.RFC3339),
		Status:            "active",
	}, nil
}

// EndSession terminates a parking session
func (s *ProviderServiceServer) EndSession(ctx context.Context, req *EndSessionRequest) (*EndSessionResponse, error) {
	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provider_id")
	}

	// Verify provider exists
	_, err = s.providerService.GetProvider(ctx, providerID)
	if err != nil {
		if err == domain.ErrProviderNotFound {
			return nil, status.Error(codes.NotFound, "provider not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Simulate ending the session
	// In production, this would call the provider's API
	exitTime := time.Now().UTC()
	duration := 60 // Simulated 60 minutes
	amount := decimal.NewFromFloat(5.00)

	return &EndSessionResponse{
		ExitTime:        exitTime.Format(time.RFC3339),
		DurationMinutes: int32(duration),
		Amount:          amount.String(),
		Currency:        "MYR",
	}, nil
}

// GetSessionStatus retrieves the current status of a session
func (s *ProviderServiceServer) GetSessionStatus(ctx context.Context, req *GetSessionStatusRequest) (*SessionStatusResponse, error) {
	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provider_id")
	}

	// Verify provider exists
	_, err = s.providerService.GetProvider(ctx, providerID)
	if err != nil {
		if err == domain.ErrProviderNotFound {
			return nil, status.Error(codes.NotFound, "provider not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Simulated session status
	return &SessionStatusResponse{
		Status:          "active",
		DurationMinutes: 30,
		CurrentAmount:   "2.50",
		Currency:        "MYR",
		EntryTime:       time.Now().Add(-30 * time.Minute).UTC().Format(time.RFC3339),
	}, nil
}

// GetProvider retrieves provider information by ID
func (s *ProviderServiceServer) GetProvider(ctx context.Context, req *GetProviderRequest) (*ProviderResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provider id")
	}

	provider, err := s.providerService.GetProvider(ctx, id)
	if err != nil {
		if err == domain.ErrProviderNotFound {
			return nil, status.Error(codes.NotFound, "provider not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ProviderResponse{
		ID:         provider.ID.String(),
		Name:       provider.Name,
		Code:       provider.Code,
		Status:     provider.Status,
		MFEURL:     provider.MFEURL,
		APIBaseURL: provider.APIBaseURL,
		LogoURL:    provider.LogoURL,
	}, nil
}

// ListProviders lists all available providers
func (s *ProviderServiceServer) ListProviders(ctx context.Context, req *ListProvidersRequest) (*ListProvidersResponse, error) {
	providers, err := s.providerService.ListProviders(ctx, req.ActiveOnly)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responses := make([]*ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = &ProviderResponse{
			ID:         p.ID.String(),
			Name:       p.Name,
			Code:       p.Code,
			Status:     p.Status,
			MFEURL:     p.MFEURL,
			APIBaseURL: p.APIBaseURL,
			LogoURL:    p.LogoURL,
		}
	}

	return &ListProvidersResponse{
		Providers: responses,
		Total:     int32(len(responses)),
	}, nil
}
