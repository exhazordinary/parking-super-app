package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/provider/internal/domain"
)

type CredentialsRepository struct {
	db *pgxpool.Pool
}

func NewCredentialsRepository(db *pgxpool.Pool) *CredentialsRepository {
	return &CredentialsRepository{db: db}
}

func (r *CredentialsRepository) Create(ctx context.Context, creds *domain.ProviderCredentials) error {
	query := `
		INSERT INTO provider_credentials (
			id, provider_id, api_key, api_secret, environment,
			is_active, created_at, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		creds.ID, creds.ProviderID, creds.APIKey, creds.APISecret,
		creds.Environment, creds.IsActive, creds.CreatedAt, creds.ExpiresAt,
	)
	return err
}

func (r *CredentialsRepository) GetByAPIKey(ctx context.Context, apiKey string) (*domain.ProviderCredentials, error) {
	query := `
		SELECT id, provider_id, api_key, api_secret, environment,
			is_active, created_at, expires_at
		FROM provider_credentials WHERE api_key = $1
	`
	return r.scanCredentials(r.db.QueryRow(ctx, query, apiKey))
}

func (r *CredentialsRepository) GetByProviderID(ctx context.Context, providerID uuid.UUID, env domain.Environment) (*domain.ProviderCredentials, error) {
	query := `
		SELECT id, provider_id, api_key, api_secret, environment,
			is_active, created_at, expires_at
		FROM provider_credentials
		WHERE provider_id = $1 AND environment = $2 AND is_active = true
		ORDER BY created_at DESC LIMIT 1
	`
	return r.scanCredentials(r.db.QueryRow(ctx, query, providerID, env))
}

func (r *CredentialsRepository) Update(ctx context.Context, creds *domain.ProviderCredentials) error {
	query := `
		UPDATE provider_credentials
		SET is_active = $2, expires_at = $3
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, creds.ID, creds.IsActive, creds.ExpiresAt)
	return err
}

func (r *CredentialsRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE provider_credentials SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *CredentialsRepository) scanCredentials(row pgx.Row) (*domain.ProviderCredentials, error) {
	var c domain.ProviderCredentials
	err := row.Scan(
		&c.ID, &c.ProviderID, &c.APIKey, &c.APISecret,
		&c.Environment, &c.IsActive, &c.CreatedAt, &c.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProviderNotFound
		}
		return nil, err
	}
	return &c, nil
}
