package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/wallet/internal/application"
	"github.com/parking-super-app/services/wallet/internal/domain"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WalletServiceServer implements the gRPC WalletService
// This is a manual implementation until proto files are generated
type WalletServiceServer struct {
	walletService *application.WalletService
}

// NewWalletServiceServer creates a new gRPC server for the wallet service
func NewWalletServiceServer(ws *application.WalletService) *WalletServiceServer {
	return &WalletServiceServer{
		walletService: ws,
	}
}

// PayRequest represents a payment request
type PayRequest struct {
	WalletID       string
	Amount         string
	Currency       string
	ProviderID     string
	ReferenceID    string
	Description    string
	IdempotencyKey string
}

// PayResponse represents a payment response
type PayResponse struct {
	TransactionID string
	Status        string
	BalanceAfter  string
	ErrorMessage  string
}

// GetWalletRequest represents a get wallet request
type GetWalletRequest struct {
	UserID string
}

// GetWalletByIDRequest represents a get wallet by ID request
type GetWalletByIDRequest struct {
	WalletID string
}

// GetWalletResponse represents a wallet response
type GetWalletResponse struct {
	ID        string
	UserID    string
	Balance   string
	Currency  string
	Status    string
	CreatedAt string
	UpdatedAt string
}

// Pay processes a payment from a wallet
func (s *WalletServiceServer) Pay(ctx context.Context, req *PayRequest) (*PayResponse, error) {
	walletID, err := uuid.Parse(req.WalletID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid wallet_id")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}

	providerID := uuid.Nil
	if req.ProviderID != "" {
		providerID, _ = uuid.Parse(req.ProviderID)
	}

	resp, err := s.walletService.Pay(ctx, application.PaymentRequest{
		WalletID:       walletID,
		Amount:         amount,
		ProviderID:     providerID,
		ReferenceID:    req.ReferenceID,
		Description:    req.Description,
		IdempotencyKey: req.IdempotencyKey,
	})

	if err != nil {
		switch err {
		case domain.ErrWalletNotFound:
			return nil, status.Error(codes.NotFound, "wallet not found")
		case domain.ErrInsufficientBalance:
			return nil, status.Error(codes.FailedPrecondition, "insufficient balance")
		case domain.ErrWalletInactive:
			return nil, status.Error(codes.FailedPrecondition, "wallet is inactive")
		case domain.ErrInvalidAmount:
			return nil, status.Error(codes.InvalidArgument, "invalid amount")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &PayResponse{
		TransactionID: resp.ID.String(),
		Status:        resp.Status,
		BalanceAfter:  resp.BalanceAfter.String(),
	}, nil
}

// GetWallet retrieves wallet information by user ID
func (s *WalletServiceServer) GetWallet(ctx context.Context, req *GetWalletRequest) (*GetWalletResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	wallet, err := s.walletService.GetWallet(ctx, userID)
	if err != nil {
		if err == domain.ErrWalletNotFound {
			return nil, status.Error(codes.NotFound, "wallet not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetWalletResponse{
		ID:       wallet.ID.String(),
		UserID:   wallet.UserID.String(),
		Balance:  wallet.Balance.String(),
		Currency: wallet.Currency,
		Status:   wallet.Status,
	}, nil
}

// GetWalletByID retrieves wallet information by wallet ID
func (s *WalletServiceServer) GetWalletByID(ctx context.Context, req *GetWalletByIDRequest) (*GetWalletResponse, error) {
	walletID, err := uuid.Parse(req.WalletID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid wallet_id")
	}

	wallet, err := s.walletService.GetWalletByID(ctx, walletID)
	if err != nil {
		if err == domain.ErrWalletNotFound {
			return nil, status.Error(codes.NotFound, "wallet not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetWalletResponse{
		ID:       wallet.ID.String(),
		UserID:   wallet.UserID.String(),
		Balance:  wallet.Balance.String(),
		Currency: wallet.Currency,
		Status:   wallet.Status,
	}, nil
}
