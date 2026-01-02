package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/parking/internal/domain"
	"github.com/shopspring/decimal"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.ParkingSession) error {
	query := `
		INSERT INTO parking_sessions (
			id, user_id, provider_id, location_id, external_session_id,
			vehicle_plate, vehicle_type, entry_time, exit_time,
			duration_minutes, amount, currency, status, payment_id,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID, session.UserID, session.ProviderID, session.LocationID,
		session.ExternalSessionID, session.VehiclePlate, session.VehicleType,
		session.EntryTime, session.ExitTime, session.Duration,
		session.Amount, session.Currency, session.Status, session.PaymentID,
		session.CreatedAt, session.UpdatedAt,
	)
	return err
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ParkingSession, error) {
	query := `
		SELECT id, user_id, provider_id, location_id, external_session_id,
			vehicle_plate, vehicle_type, entry_time, exit_time,
			duration_minutes, amount, currency, status, payment_id,
			created_at, updated_at
		FROM parking_sessions WHERE id = $1
	`
	return r.scanSession(r.db.QueryRow(ctx, query, id))
}

func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ParkingSession, error) {
	query := `
		SELECT id, user_id, provider_id, location_id, external_session_id,
			vehicle_plate, vehicle_type, entry_time, exit_time,
			duration_minutes, amount, currency, status, payment_id,
			created_at, updated_at
		FROM parking_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSessions(rows)
}

func (r *SessionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ParkingSession, error) {
	query := `
		SELECT id, user_id, provider_id, location_id, external_session_id,
			vehicle_plate, vehicle_type, entry_time, exit_time,
			duration_minutes, amount, currency, status, payment_id,
			created_at, updated_at
		FROM parking_sessions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY entry_time DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSessions(rows)
}

func (r *SessionRepository) GetByProviderID(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*domain.ParkingSession, error) {
	query := `
		SELECT id, user_id, provider_id, location_id, external_session_id,
			vehicle_plate, vehicle_type, entry_time, exit_time,
			duration_minutes, amount, currency, status, payment_id,
			created_at, updated_at
		FROM parking_sessions
		WHERE provider_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, providerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSessions(rows)
}

func (r *SessionRepository) Update(ctx context.Context, session *domain.ParkingSession) error {
	query := `
		UPDATE parking_sessions
		SET external_session_id = $2, exit_time = $3, duration_minutes = $4,
			amount = $5, status = $6, payment_id = $7, updated_at = $8
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query,
		session.ID, session.ExternalSessionID, session.ExitTime,
		session.Duration, session.Amount, session.Status,
		session.PaymentID, session.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrSessionNotFound
	}
	return nil
}

func (r *SessionRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM parking_sessions WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *SessionRepository) scanSession(row pgx.Row) (*domain.ParkingSession, error) {
	var s domain.ParkingSession
	var amount decimal.Decimal
	err := row.Scan(
		&s.ID, &s.UserID, &s.ProviderID, &s.LocationID, &s.ExternalSessionID,
		&s.VehiclePlate, &s.VehicleType, &s.EntryTime, &s.ExitTime,
		&s.Duration, &amount, &s.Currency, &s.Status, &s.PaymentID,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}
	s.Amount = amount
	return &s, nil
}

func (r *SessionRepository) scanSessions(rows pgx.Rows) ([]*domain.ParkingSession, error) {
	var sessions []*domain.ParkingSession
	for rows.Next() {
		var s domain.ParkingSession
		var amount decimal.Decimal
		err := rows.Scan(
			&s.ID, &s.UserID, &s.ProviderID, &s.LocationID, &s.ExternalSessionID,
			&s.VehiclePlate, &s.VehicleType, &s.EntryTime, &s.ExitTime,
			&s.Duration, &amount, &s.Currency, &s.Status, &s.PaymentID,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		s.Amount = amount
		sessions = append(sessions, &s)
	}
	return sessions, rows.Err()
}
