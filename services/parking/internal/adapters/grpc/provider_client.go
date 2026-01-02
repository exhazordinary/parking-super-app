package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/ports"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProviderGRPCClient implements ports.ProviderClient using gRPC
type ProviderGRPCClient struct {
	conn    *grpc.ClientConn
	address string
}

// NewProviderGRPCClient creates a new gRPC client for the provider service
func NewProviderGRPCClient(address string) (*ProviderGRPCClient, error) {
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to provider service: %w", err)
	}

	return &ProviderGRPCClient{
		conn:    conn,
		address: address,
	}, nil
}

// StartSession initiates a parking session with the provider
func (c *ProviderGRPCClient) StartSession(ctx context.Context, req ports.StartSessionRequest) (*ports.StartSessionResponse, error) {
	// This is a simplified implementation
	// In production with generated proto code, this would use the generated client

	// The actual implementation would look like:
	// resp, err := c.client.StartSession(ctx, &providerv1.StartSessionRequest{
	//     ProviderId:   req.ProviderID.String(),
	//     LocationId:   req.LocationID.String(),
	//     VehiclePlate: req.VehiclePlate,
	//     VehicleType:  req.VehicleType,
	//     UserRef:      req.UserRef,
	// })

	// Simulated successful response
	return &ports.StartSessionResponse{
		ExternalSessionID: uuid.New().String(),
		EntryTime:         time.Now().UTC().Format(time.RFC3339),
		Status:            "active",
	}, nil
}

// EndSession terminates a parking session
func (c *ProviderGRPCClient) EndSession(ctx context.Context, req ports.EndSessionRequest) (*ports.EndSessionResponse, error) {
	// Simulated response - in production this would use the generated client
	return &ports.EndSessionResponse{
		ExitTime: time.Now().UTC().Format(time.RFC3339),
		Duration: 60, // 1 hour
		Amount:   decimal.NewFromFloat(5.00),
		Currency: "MYR",
	}, nil
}

// GetSessionStatus retrieves the current status of a session
func (c *ProviderGRPCClient) GetSessionStatus(ctx context.Context, providerID uuid.UUID, externalSessionID string) (*ports.SessionStatusResponse, error) {
	// Simulated response - in production this would use the generated client
	return &ports.SessionStatusResponse{
		Status:   "active",
		Duration: 30,
		Amount:   decimal.NewFromFloat(2.50),
	}, nil
}

// Close closes the gRPC connection
func (c *ProviderGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Ensure ProviderGRPCClient implements ports.ProviderClient
var _ ports.ProviderClient = (*ProviderGRPCClient)(nil)
