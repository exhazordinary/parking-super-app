package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/notification/internal/domain"
)

type PreferenceRepository struct {
	db *pgxpool.Pool
}

func NewPreferenceRepository(db *pgxpool.Pool) *PreferenceRepository {
	return &PreferenceRepository{db: db}
}

func (r *PreferenceRepository) Create(ctx context.Context, pref *domain.UserPreference) error {
	typePrefsJSON, _ := json.Marshal(pref.TypePreferences)
	query := `
		INSERT INTO user_preferences (
			id, user_id, push_enabled, sms_enabled, email_enabled,
			quiet_hours_start, quiet_hours_end, type_preferences,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		pref.ID, pref.UserID, pref.PushEnabled, pref.SMSEnabled,
		pref.EmailEnabled, pref.QuietHoursStart, pref.QuietHoursEnd,
		typePrefsJSON, pref.CreatedAt, pref.UpdatedAt,
	)
	return err
}

func (r *PreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserPreference, error) {
	query := `
		SELECT id, user_id, push_enabled, sms_enabled, email_enabled,
			quiet_hours_start, quiet_hours_end, type_preferences,
			created_at, updated_at
		FROM user_preferences WHERE user_id = $1
	`
	var p domain.UserPreference
	var typePrefsJSON []byte
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.PushEnabled, &p.SMSEnabled, &p.EmailEnabled,
		&p.QuietHoursStart, &p.QuietHoursEnd, &typePrefsJSON,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("preferences not found")
		}
		return nil, err
	}
	json.Unmarshal(typePrefsJSON, &p.TypePreferences)
	return &p, nil
}

func (r *PreferenceRepository) Update(ctx context.Context, pref *domain.UserPreference) error {
	typePrefsJSON, _ := json.Marshal(pref.TypePreferences)
	query := `
		UPDATE user_preferences
		SET push_enabled = $2, sms_enabled = $3, email_enabled = $4,
			quiet_hours_start = $5, quiet_hours_end = $6,
			type_preferences = $7, updated_at = $8
		WHERE user_id = $1
	`
	_, err := r.db.Exec(ctx, query,
		pref.UserID, pref.PushEnabled, pref.SMSEnabled, pref.EmailEnabled,
		pref.QuietHoursStart, pref.QuietHoursEnd, typePrefsJSON, pref.UpdatedAt,
	)
	return err
}

func (r *PreferenceRepository) Upsert(ctx context.Context, pref *domain.UserPreference) error {
	typePrefsJSON, _ := json.Marshal(pref.TypePreferences)
	query := `
		INSERT INTO user_preferences (
			id, user_id, push_enabled, sms_enabled, email_enabled,
			quiet_hours_start, quiet_hours_end, type_preferences,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id) DO UPDATE SET
			push_enabled = EXCLUDED.push_enabled,
			sms_enabled = EXCLUDED.sms_enabled,
			email_enabled = EXCLUDED.email_enabled,
			quiet_hours_start = EXCLUDED.quiet_hours_start,
			quiet_hours_end = EXCLUDED.quiet_hours_end,
			type_preferences = EXCLUDED.type_preferences,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, query,
		pref.ID, pref.UserID, pref.PushEnabled, pref.SMSEnabled,
		pref.EmailEnabled, pref.QuietHoursStart, pref.QuietHoursEnd,
		typePrefsJSON, pref.CreatedAt, pref.UpdatedAt,
	)
	return err
}
