package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/wallet/internal/domain"
	"github.com/parking-super-app/services/wallet/internal/ports"
	"github.com/shopspring/decimal"
)

type WalletService struct {
	wallets      ports.WalletRepository
	transactions ports.TransactionRepository
	uow          ports.UnitOfWork
	gateway      ports.PaymentGateway
	events       ports.EventPublisher
	logger       ports.Logger
}

func NewWalletService(
	wallets ports.WalletRepository,
	transactions ports.TransactionRepository,
	uow ports.UnitOfWork,
	gateway ports.PaymentGateway,
	events ports.EventPublisher,
	logger ports.Logger,
) *WalletService {
	return &WalletService{
		wallets:      wallets,
		transactions: transactions,
		uow:          uow,
		gateway:      gateway,
		events:       events,
		logger:       logger,
	}
}

type CreateWalletRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Currency string    `json:"currency"`
}

type WalletResponse struct {
	ID       uuid.UUID       `json:"id"`
	UserID   uuid.UUID       `json:"user_id"`
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"currency"`
	Status   string          `json:"status"`
}

type TopUpRequest struct {
	WalletID       uuid.UUID       `json:"wallet_id"`
	Amount         decimal.Decimal `json:"amount"`
	PaymentMethod  string          `json:"payment_method"`
	IdempotencyKey string          `json:"idempotency_key"`
}

type PaymentRequest struct {
	WalletID       uuid.UUID       `json:"wallet_id"`
	Amount         decimal.Decimal `json:"amount"`
	ProviderID     uuid.UUID       `json:"provider_id"`
	ReferenceID    string          `json:"reference_id"`
	Description    string          `json:"description"`
	IdempotencyKey string          `json:"idempotency_key"`
}

type TransactionResponse struct {
	ID            uuid.UUID       `json:"id"`
	Type          string          `json:"type"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceBefore decimal.Decimal `json:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	Status        string          `json:"status"`
	Description   string          `json:"description"`
	CreatedAt     string          `json:"created_at"`
}

type TransactionListResponse struct {
	Transactions []*TransactionResponse `json:"transactions"`
	Total        int                    `json:"total"`
	Limit        int                    `json:"limit"`
	Offset       int                    `json:"offset"`
}

func (s *WalletService) CreateWallet(ctx context.Context, req CreateWalletRequest) (*WalletResponse, error) {
	s.logger.Info("creating wallet", ports.String("user_id", req.UserID.String()))

	exists, err := s.wallets.ExistsByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check wallet existence: %w", err)
	}
	if exists {
		return nil, domain.ErrWalletAlreadyExists
	}

	currency := req.Currency
	if currency == "" {
		currency = "MYR"
	}

	wallet := domain.NewWallet(req.UserID, currency)
	if err := s.wallets.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	go func() {
		event := ports.Event{
			Type: ports.EventWalletCreated,
			Payload: map[string]interface{}{
				"wallet_id": wallet.ID.String(),
				"user_id":   wallet.UserID.String(),
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return &WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
		Status:   string(wallet.Status),
	}, nil
}

func (s *WalletService) GetWallet(ctx context.Context, userID uuid.UUID) (*WalletResponse, error) {
	wallet, err := s.wallets.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
		Status:   string(wallet.Status),
	}, nil
}

func (s *WalletService) TopUp(ctx context.Context, req TopUpRequest) (*TransactionResponse, error) {
	s.logger.Info("processing topup",
		ports.String("wallet_id", req.WalletID.String()),
		ports.String("amount", req.Amount.String()),
	)

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}

	existingTx, err := s.transactions.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil && existingTx != nil {
		return s.toTransactionResponse(existingTx), nil
	}

	wallet, err := s.wallets.GetByID(ctx, req.WalletID)
	if err != nil {
		return nil, err
	}

	if !wallet.CanTransact() {
		return nil, domain.ErrWalletInactive
	}

	tx := domain.NewTransaction(
		wallet.ID,
		domain.TransactionTypeTopUp,
		req.Amount,
		wallet.Balance,
		"",
		req.IdempotencyKey,
		"Wallet top-up",
	)

	if err := s.transactions.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := wallet.Credit(req.Amount); err != nil {
		tx.Fail()
		s.transactions.Update(ctx, tx)
		return nil, err
	}

	if err := s.wallets.Update(ctx, wallet); err != nil {
		tx.Fail()
		s.transactions.Update(ctx, tx)
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	tx.Complete(wallet.Balance)
	if err := s.transactions.Update(ctx, tx); err != nil {
		s.logger.Error("failed to update transaction status", ports.Err(err))
	}

	go func() {
		event := ports.Event{
			Type: ports.EventTopUpCompleted,
			Payload: map[string]interface{}{
				"transaction_id": tx.ID.String(),
				"wallet_id":      wallet.ID.String(),
				"amount":         req.Amount.String(),
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return s.toTransactionResponse(tx), nil
}

func (s *WalletService) Pay(ctx context.Context, req PaymentRequest) (*TransactionResponse, error) {
	s.logger.Info("processing payment",
		ports.String("wallet_id", req.WalletID.String()),
		ports.String("amount", req.Amount.String()),
	)

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidAmount
	}

	existingTx, err := s.transactions.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil && existingTx != nil {
		return s.toTransactionResponse(existingTx), nil
	}

	wallet, err := s.wallets.GetByID(ctx, req.WalletID)
	if err != nil {
		return nil, err
	}

	if !wallet.CanTransact() {
		return nil, domain.ErrWalletInactive
	}

	if !wallet.HasSufficientBalance(req.Amount) {
		return nil, domain.ErrInsufficientBalance
	}

	tx := domain.NewTransaction(
		wallet.ID,
		domain.TransactionTypePayment,
		req.Amount,
		wallet.Balance,
		req.ReferenceID,
		req.IdempotencyKey,
		req.Description,
	)
	tx.SetProvider(req.ProviderID)

	if err := s.transactions.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := wallet.Debit(req.Amount); err != nil {
		tx.Fail()
		s.transactions.Update(ctx, tx)
		return nil, err
	}

	if err := s.wallets.Update(ctx, wallet); err != nil {
		tx.Fail()
		s.transactions.Update(ctx, tx)
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	tx.Complete(wallet.Balance)
	if err := s.transactions.Update(ctx, tx); err != nil {
		s.logger.Error("failed to update transaction status", ports.Err(err))
	}

	go func() {
		event := ports.Event{
			Type: ports.EventPaymentCompleted,
			Payload: map[string]interface{}{
				"transaction_id": tx.ID.String(),
				"wallet_id":      wallet.ID.String(),
				"provider_id":    req.ProviderID.String(),
				"amount":         req.Amount.String(),
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return s.toTransactionResponse(tx), nil
}

func (s *WalletService) GetTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) (*TransactionListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	transactions, err := s.transactions.GetByWalletID(ctx, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	total, err := s.transactions.CountByWalletID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}

	var txResponses []*TransactionResponse
	for _, tx := range transactions {
		txResponses = append(txResponses, s.toTransactionResponse(tx))
	}

	return &TransactionListResponse{
		Transactions: txResponses,
		Total:        total,
		Limit:        limit,
		Offset:       offset,
	}, nil
}

func (s *WalletService) toTransactionResponse(tx *domain.Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:            tx.ID,
		Type:          string(tx.Type),
		Amount:        tx.Amount,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		Status:        string(tx.Status),
		Description:   tx.Description,
		CreatedAt:     tx.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
