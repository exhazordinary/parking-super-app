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

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, wallet_id, type, amount, balance_before, balance_after,
			reference_id, provider_id, status, description, idempotency_key,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		tx.ID, tx.WalletID, tx.Type, tx.Amount, tx.BalanceBefore, tx.BalanceAfter,
		tx.ReferenceID, tx.ProviderID, tx.Status, tx.Description, tx.IdempotencyKey,
		tx.CreatedAt, tx.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrDuplicateTransaction
		}
		return err
	}
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	query := `
		SELECT id, wallet_id, type, amount, balance_before, balance_after,
			reference_id, provider_id, status, description, idempotency_key,
			created_at, updated_at
		FROM transactions WHERE id = $1
	`
	return r.scanTransaction(r.db.QueryRow(ctx, query, id))
}

func (r *TransactionRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	if key == "" {
		return nil, domain.ErrTransactionNotFound
	}
	query := `
		SELECT id, wallet_id, type, amount, balance_before, balance_after,
			reference_id, provider_id, status, description, idempotency_key,
			created_at, updated_at
		FROM transactions WHERE idempotency_key = $1
	`
	return r.scanTransaction(r.db.QueryRow(ctx, query, key))
}

func (r *TransactionRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, wallet_id, type, amount, balance_before, balance_after,
			reference_id, provider_id, status, description, idempotency_key,
			created_at, updated_at
		FROM transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, walletID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx, err := r.scanTransactionRow(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, rows.Err()
}

func (r *TransactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET status = $2, balance_after = $3, updated_at = $4
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, tx.ID, tx.Status, tx.BalanceAfter, tx.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrTransactionNotFound
	}
	return nil
}

func (r *TransactionRepository) CountByWalletID(ctx context.Context, walletID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE wallet_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, walletID).Scan(&count)
	return count, err
}

func (r *TransactionRepository) scanTransaction(row pgx.Row) (*domain.Transaction, error) {
	tx := &domain.Transaction{}
	var amount, balanceBefore, balanceAfter decimal.Decimal
	err := row.Scan(
		&tx.ID, &tx.WalletID, &tx.Type, &amount, &balanceBefore, &balanceAfter,
		&tx.ReferenceID, &tx.ProviderID, &tx.Status, &tx.Description, &tx.IdempotencyKey,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, err
	}
	tx.Amount = amount
	tx.BalanceBefore = balanceBefore
	tx.BalanceAfter = balanceAfter
	return tx, nil
}

func (r *TransactionRepository) scanTransactionRow(rows pgx.Rows) (*domain.Transaction, error) {
	tx := &domain.Transaction{}
	var amount, balanceBefore, balanceAfter decimal.Decimal
	err := rows.Scan(
		&tx.ID, &tx.WalletID, &tx.Type, &amount, &balanceBefore, &balanceAfter,
		&tx.ReferenceID, &tx.ProviderID, &tx.Status, &tx.Description, &tx.IdempotencyKey,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	tx.Amount = amount
	tx.BalanceBefore = balanceBefore
	tx.BalanceAfter = balanceAfter
	return tx, nil
}
