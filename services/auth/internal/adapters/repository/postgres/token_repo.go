package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/auth/internal/domain"
)

// RefreshTokenRepository implements ports.RefreshTokenRepository using PostgreSQL.
type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository.
func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create stores a new refresh token.
func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked, created_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.Revoked,
		token.CreatedAt,
		token.UserAgent,
		token.IPAddress,
	)

	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetByTokenHash retrieves a token by its hash.
//
// SECURITY NOTE: We lookup by hash, not by the raw token.
// The client sends the raw token, we hash it, then look it up.
func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked, created_at, revoked_at, user_agent, ip_address
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	token := &domain.RefreshToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.Revoked,
		&token.CreatedAt,
		&token.RevokedAt,
		&token.UserAgent,
		&token.IPAddress,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return token, nil
}

// GetByUserID retrieves all tokens for a user.
// Useful for showing active sessions or implementing "logout everywhere".
func (r *RefreshTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked, created_at, revoked_at, user_agent, ip_address
		FROM refresh_tokens
		WHERE user_id = $1 AND revoked = false AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*domain.RefreshToken
	for rows.Next() {
		token := &domain.RefreshToken{}
		if err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.TokenHash,
			&token.ExpiresAt,
			&token.Revoked,
			&token.CreatedAt,
			&token.RevokedAt,
			&token.UserAgent,
			&token.IPAddress,
		); err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tokens: %w", err)
	}

	return tokens, nil
}

// Revoke marks a specific token as revoked.
func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, revoked_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTokenNotFound
	}

	return nil
}

// RevokeAllForUser revokes all tokens for a user.
// This is the "logout everywhere" functionality.
func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, revoked_at = NOW()
		WHERE user_id = $1 AND revoked = false
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all tokens: %w", err)
	}

	return nil
}

// DeleteExpired removes tokens that have expired.
// This should be called periodically by a cleanup job (e.g., cron).
//
// BEST PRACTICE: Batch Deletion
// =============================
// In production with millions of tokens, consider:
// - Deleting in batches (LIMIT 1000)
// - Using RETURNING to log deleted tokens
// - Running during low-traffic periods
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW()
	`

	result, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	// Log how many were deleted (in a real app, use proper logging)
	_ = result.RowsAffected()

	return nil
}
