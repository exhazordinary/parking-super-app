package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/wallet/internal/domain"
	"github.com/shopspring/decimal"
)

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, balance, currency, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency,
		wallet.Status, wallet.CreatedAt, wallet.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrWalletAlreadyExists
		}
		return err
	}
	return nil
}

func (r *WalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	query := `
		SELECT id, user_id, balance, currency, status, created_at, updated_at
		FROM wallets WHERE id = $1
	`
	wallet := &domain.Wallet{}
	var balance decimal.Decimal
	err := r.db.QueryRow(ctx, query, id).Scan(
		&wallet.ID, &wallet.UserID, &balance, &wallet.Currency,
		&wallet.Status, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}
	wallet.Balance = balance
	return wallet, nil
}

func (r *WalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	query := `
		SELECT id, user_id, balance, currency, status, created_at, updated_at
		FROM wallets WHERE user_id = $1
	`
	wallet := &domain.Wallet{}
	var balance decimal.Decimal
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&wallet.ID, &wallet.UserID, &balance, &wallet.Currency,
		&wallet.Status, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}
	wallet.Balance = balance
	return wallet, nil
}

func (r *WalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		UPDATE wallets
		SET balance = $2, status = $3, updated_at = $4
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query,
		wallet.ID, wallet.Balance, wallet.Status, wallet.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrWalletNotFound
	}
	return nil
}

func (r *WalletRepository) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM wallets WHERE user_id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
