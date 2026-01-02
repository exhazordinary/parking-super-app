package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/parking-super-app/services/provider/internal/domain"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(ctx context.Context, location *domain.Location) error {
	query := `
		INSERT INTO locations (
			id, provider_id, name, address, city, state, postal_code,
			latitude, longitude, total_spaces, amenities,
			hourly_rate, daily_max, currency, grace_period_min,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	_, err := r.db.Exec(ctx, query,
		location.ID, location.ProviderID, location.Name, location.Address,
		location.City, location.State, location.PostalCode,
		location.Latitude, location.Longitude, location.TotalSpaces,
		pq.Array(location.Amenities),
		location.Pricing.HourlyRate, location.Pricing.DailyMax,
		location.Pricing.Currency, location.Pricing.GracePeriodMin,
		location.IsActive, location.CreatedAt, location.UpdatedAt,
	)
	return err
}

func (r *LocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	query := `
		SELECT id, provider_id, name, address, city, state, postal_code,
			latitude, longitude, total_spaces, amenities,
			hourly_rate, daily_max, currency, grace_period_min,
			is_active, created_at, updated_at
		FROM locations WHERE id = $1
	`
	return r.scanLocation(r.db.QueryRow(ctx, query, id))
}

func (r *LocationRepository) GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*domain.Location, error) {
	query := `
		SELECT id, provider_id, name, address, city, state, postal_code,
			latitude, longitude, total_spaces, amenities,
			hourly_rate, daily_max, currency, grace_period_min,
			is_active, created_at, updated_at
		FROM locations WHERE provider_id = $1 AND is_active = true
		ORDER BY name
	`
	rows, err := r.db.Query(ctx, query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*domain.Location
	for rows.Next() {
		loc, err := r.scanLocationRow(rows)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, rows.Err()
}

func (r *LocationRepository) GetNearby(ctx context.Context, lat, lng float64, radiusKm float64) ([]*domain.Location, error) {
	// Using Haversine formula for distance calculation
	// This is approximate but works well for short distances
	query := `
		SELECT id, provider_id, name, address, city, state, postal_code,
			latitude, longitude, total_spaces, amenities,
			hourly_rate, daily_max, currency, grace_period_min,
			is_active, created_at, updated_at,
			(6371 * acos(cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))) AS distance
		FROM locations
		WHERE is_active = true
		HAVING distance < $3
		ORDER BY distance
		LIMIT 50
	`
	rows, err := r.db.Query(ctx, query, lat, lng, radiusKm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*domain.Location
	for rows.Next() {
		loc, err := r.scanLocationRowWithDistance(rows)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, rows.Err()
}

func (r *LocationRepository) Update(ctx context.Context, location *domain.Location) error {
	query := `
		UPDATE locations
		SET name = $2, address = $3, city = $4, state = $5, postal_code = $6,
			latitude = $7, longitude = $8, total_spaces = $9, amenities = $10,
			hourly_rate = $11, daily_max = $12, is_active = $13, updated_at = $14
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query,
		location.ID, location.Name, location.Address, location.City,
		location.State, location.PostalCode, location.Latitude, location.Longitude,
		location.TotalSpaces, pq.Array(location.Amenities),
		location.Pricing.HourlyRate, location.Pricing.DailyMax,
		location.IsActive, location.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrProviderNotFound
	}
	return nil
}

func (r *LocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM locations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrProviderNotFound
	}
	return nil
}

func (r *LocationRepository) scanLocation(row pgx.Row) (*domain.Location, error) {
	var loc domain.Location
	var amenities []string
	err := row.Scan(
		&loc.ID, &loc.ProviderID, &loc.Name, &loc.Address, &loc.City,
		&loc.State, &loc.PostalCode, &loc.Latitude, &loc.Longitude,
		&loc.TotalSpaces, pq.Array(&amenities),
		&loc.Pricing.HourlyRate, &loc.Pricing.DailyMax,
		&loc.Pricing.Currency, &loc.Pricing.GracePeriodMin,
		&loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProviderNotFound
		}
		return nil, err
	}
	loc.Amenities = amenities
	return &loc, nil
}

func (r *LocationRepository) scanLocationRow(rows pgx.Rows) (*domain.Location, error) {
	var loc domain.Location
	var amenities []string
	err := rows.Scan(
		&loc.ID, &loc.ProviderID, &loc.Name, &loc.Address, &loc.City,
		&loc.State, &loc.PostalCode, &loc.Latitude, &loc.Longitude,
		&loc.TotalSpaces, pq.Array(&amenities),
		&loc.Pricing.HourlyRate, &loc.Pricing.DailyMax,
		&loc.Pricing.Currency, &loc.Pricing.GracePeriodMin,
		&loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	loc.Amenities = amenities
	return &loc, nil
}

func (r *LocationRepository) scanLocationRowWithDistance(rows pgx.Rows) (*domain.Location, error) {
	var loc domain.Location
	var amenities []string
	var distance float64
	err := rows.Scan(
		&loc.ID, &loc.ProviderID, &loc.Name, &loc.Address, &loc.City,
		&loc.State, &loc.PostalCode, &loc.Latitude, &loc.Longitude,
		&loc.TotalSpaces, pq.Array(&amenities),
		&loc.Pricing.HourlyRate, &loc.Pricing.DailyMax,
		&loc.Pricing.Currency, &loc.Pricing.GracePeriodMin,
		&loc.IsActive, &loc.CreatedAt, &loc.UpdatedAt,
		&distance,
	)
	if err != nil {
		return nil, err
	}
	loc.Amenities = amenities
	return &loc, nil
}
