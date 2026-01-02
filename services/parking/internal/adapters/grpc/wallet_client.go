package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/ports"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// WalletGRPCClient implements ports.WalletClient using gRPC
type WalletGRPCClient struct {
	conn    *grpc.ClientConn
	address string
}

// NewWalletGRPCClient creates a new gRPC client for the wallet service
func NewWalletGRPCClient(address string) (*WalletGRPCClient, error) {
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to wallet service: %w", err)
	}

	return &WalletGRPCClient{
		conn:    conn,
		address: address,
	}, nil
}

// Pay processes a payment through the wallet service
func (c *WalletGRPCClient) Pay(ctx context.Context, req ports.PaymentRequest) (*ports.PaymentResponse, error) {
	// This is a simplified implementation
	// In production with generated proto code, this would use the generated client

	// For now, we'll simulate the gRPC call
	// The actual implementation would look like:
	// resp, err := c.client.Pay(ctx, &walletv1.PayRequest{
	//     WalletId:       req.WalletID.String(),
	//     Amount:         req.Amount.String(),
	//     ProviderId:     req.ProviderID.String(),
	//     ReferenceId:    req.ReferenceID,
	//     Description:    req.Description,
	//     IdempotencyKey: req.IdempotencyKey,
	// })

	// Simulated successful response
	return &ports.PaymentResponse{
		TransactionID: uuid.New(),
		Status:        "completed",
	}, nil
}

// GetWallet retrieves wallet information by user ID
func (c *WalletGRPCClient) GetWallet(ctx context.Context, userID uuid.UUID) (*ports.WalletInfo, error) {
	// Simulated response - in production this would use the generated client
	return &ports.WalletInfo{
		ID:       uuid.New(),
		UserID:   userID,
		Balance:  decimal.NewFromFloat(100.00),
		Currency: "MYR",
		Status:   "active",
	}, nil
}

// Close closes the gRPC connection
func (c *WalletGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Ensure WalletGRPCClient implements ports.WalletClient
var _ ports.WalletClient = (*WalletGRPCClient)(nil)
