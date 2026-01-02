package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/parking/internal/domain"
)

type VehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(ctx context.Context, vehicle *domain.Vehicle) error {
	query := `
		INSERT INTO vehicles (id, user_id, plate, type, make, model, color, is_default, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		vehicle.ID, vehicle.UserID, vehicle.Plate, vehicle.Type,
		vehicle.Make, vehicle.Model, vehicle.Color, vehicle.IsDefault,
		vehicle.CreatedAt,
	)
	return err
}

func (r *VehicleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error) {
	query := `
		SELECT id, user_id, plate, type, make, model, color, is_default, created_at
		FROM vehicles WHERE id = $1
	`
	var v domain.Vehicle
	err := r.db.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.UserID, &v.Plate, &v.Type,
		&v.Make, &v.Model, &v.Color, &v.IsDefault, &v.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("vehicle not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *VehicleRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Vehicle, error) {
	query := `
		SELECT id, user_id, plate, type, make, model, color, is_default, created_at
		FROM vehicles WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []*domain.Vehicle
	for rows.Next() {
		var v domain.Vehicle
		err := rows.Scan(
			&v.ID, &v.UserID, &v.Plate, &v.Type,
			&v.Make, &v.Model, &v.Color, &v.IsDefault, &v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, &v)
	}
	return vehicles, rows.Err()
}

func (r *VehicleRepository) GetByPlate(ctx context.Context, plate string) (*domain.Vehicle, error) {
	query := `
		SELECT id, user_id, plate, type, make, model, color, is_default, created_at
		FROM vehicles WHERE plate = $1
	`
	var v domain.Vehicle
	err := r.db.QueryRow(ctx, query, plate).Scan(
		&v.ID, &v.UserID, &v.Plate, &v.Type,
		&v.Make, &v.Model, &v.Color, &v.IsDefault, &v.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("vehicle not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *VehicleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM vehicles WHERE id = $1`, id)
	return err
}

func (r *VehicleRepository) SetDefault(ctx context.Context, userID, vehicleID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Clear existing default
	_, err = tx.Exec(ctx, `UPDATE vehicles SET is_default = false WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Set new default
	_, err = tx.Exec(ctx, `UPDATE vehicles SET is_default = true WHERE id = $1 AND user_id = $2`, vehicleID, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
