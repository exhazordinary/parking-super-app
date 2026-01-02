package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/provider/internal/domain"
)

type ProviderRepository struct {
	db *pgxpool.Pool
}

func NewProviderRepository(db *pgxpool.Pool) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) Create(ctx context.Context, provider *domain.Provider) error {
	configJSON, err := json.Marshal(provider.Config)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO providers (
			id, name, code, description, logo_url, status,
			mfe_url, api_base_url, webhook_secret, config,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = r.db.Exec(ctx, query,
		provider.ID, provider.Name, provider.Code, provider.Description,
		provider.LogoURL, provider.Status, provider.MFEURL, provider.APIBaseURL,
		provider.WebhookSecret, configJSON, provider.CreatedAt, provider.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrProviderAlreadyExists
		}
		return err
	}
	return nil
}

func (r *ProviderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Provider, error) {
	query := `
		SELECT id, name, code, description, logo_url, status,
			mfe_url, api_base_url, webhook_secret, config,
			created_at, updated_at
		FROM providers WHERE id = $1
	`
	return r.scanProvider(r.db.QueryRow(ctx, query, id))
}

func (r *ProviderRepository) GetByCode(ctx context.Context, code string) (*domain.Provider, error) {
	query := `
		SELECT id, name, code, description, logo_url, status,
			mfe_url, api_base_url, webhook_secret, config,
			created_at, updated_at
		FROM providers WHERE code = $1
	`
	return r.scanProvider(r.db.QueryRow(ctx, query, code))
}

func (r *ProviderRepository) GetAll(ctx context.Context, activeOnly bool) ([]*domain.Provider, error) {
	query := `
		SELECT id, name, code, description, logo_url, status,
			mfe_url, api_base_url, webhook_secret, config,
			created_at, updated_at
		FROM providers
	`
	if activeOnly {
		query += ` WHERE status = 'active'`
	}
	query += ` ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*domain.Provider
	for rows.Next() {
		p, err := r.scanProviderRow(rows)
		if err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	return providers, rows.Err()
}

func (r *ProviderRepository) Update(ctx context.Context, provider *domain.Provider) error {
	configJSON, err := json.Marshal(provider.Config)
	if err != nil {
		return err
	}

	query := `
		UPDATE providers
		SET name = $2, description = $3, logo_url = $4, status = $5,
			mfe_url = $6, api_base_url = $7, webhook_secret = $8,
			config = $9, updated_at = $10
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query,
		provider.ID, provider.Name, provider.Description, provider.LogoURL,
		provider.Status, provider.MFEURL, provider.APIBaseURL,
		provider.WebhookSecret, configJSON, provider.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrProviderNotFound
	}
	return nil
}

func (r *ProviderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM providers WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrProviderNotFound
	}
	return nil
}

func (r *ProviderRepository) scanProvider(row pgx.Row) (*domain.Provider, error) {
	var p domain.Provider
	var configJSON []byte
	err := row.Scan(
		&p.ID, &p.Name, &p.Code, &p.Description, &p.LogoURL, &p.Status,
		&p.MFEURL, &p.APIBaseURL, &p.WebhookSecret, &configJSON,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProviderNotFound
		}
		return nil, err
	}
	if err := json.Unmarshal(configJSON, &p.Config); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProviderRepository) scanProviderRow(rows pgx.Rows) (*domain.Provider, error) {
	var p domain.Provider
	var configJSON []byte
	err := rows.Scan(
		&p.ID, &p.Name, &p.Code, &p.Description, &p.LogoURL, &p.Status,
		&p.MFEURL, &p.APIBaseURL, &p.WebhookSecret, &configJSON,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configJSON, &p.Config); err != nil {
		return nil, err
	}
	return &p, nil
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
