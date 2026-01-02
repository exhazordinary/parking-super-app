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

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notif *domain.Notification) error {
	dataJSON, _ := json.Marshal(notif.Data)
	query := `
		INSERT INTO notifications (
			id, user_id, channel, type, title, body, data, priority,
			status, recipient, provider_id, scheduled_at, sent_at,
			delivered_at, failed_at, error_msg, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.db.Exec(ctx, query,
		notif.ID, notif.UserID, notif.Channel, notif.Type, notif.Title,
		notif.Body, dataJSON, notif.Priority, notif.Status, notif.Recipient,
		notif.ProviderID, notif.ScheduledAt, notif.SentAt, notif.DeliveredAt,
		notif.FailedAt, notif.ErrorMsg, notif.CreatedAt,
	)
	return err
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	query := `
		SELECT id, user_id, channel, type, title, body, data, priority,
			status, recipient, provider_id, scheduled_at, sent_at,
			delivered_at, failed_at, error_msg, created_at
		FROM notifications WHERE id = $1
	`
	return r.scanNotification(r.db.QueryRow(ctx, query, id))
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Notification, error) {
	query := `
		SELECT id, user_id, channel, type, title, body, data, priority,
			status, recipient, provider_id, scheduled_at, sent_at,
			delivered_at, failed_at, error_msg, created_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *NotificationRepository) GetPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	query := `
		SELECT id, user_id, channel, type, title, body, data, priority,
			status, recipient, provider_id, scheduled_at, sent_at,
			delivered_at, failed_at, error_msg, created_at
		FROM notifications
		WHERE status = 'pending' AND (scheduled_at IS NULL OR scheduled_at <= NOW())
		ORDER BY priority DESC, created_at
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanNotifications(rows)
}

func (r *NotificationRepository) Update(ctx context.Context, notif *domain.Notification) error {
	query := `
		UPDATE notifications
		SET status = $2, provider_id = $3, sent_at = $4, delivered_at = $5,
			failed_at = $6, error_msg = $7
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query,
		notif.ID, notif.Status, notif.ProviderID, notif.SentAt,
		notif.DeliveredAt, notif.FailedAt, notif.ErrorMsg,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}

func (r *NotificationRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) scanNotification(row pgx.Row) (*domain.Notification, error) {
	var n domain.Notification
	var dataJSON []byte
	err := row.Scan(
		&n.ID, &n.UserID, &n.Channel, &n.Type, &n.Title, &n.Body,
		&dataJSON, &n.Priority, &n.Status, &n.Recipient, &n.ProviderID,
		&n.ScheduledAt, &n.SentAt, &n.DeliveredAt, &n.FailedAt,
		&n.ErrorMsg, &n.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, err
	}
	json.Unmarshal(dataJSON, &n.Data)
	return &n, nil
}

func (r *NotificationRepository) scanNotifications(rows pgx.Rows) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		var dataJSON []byte
		err := rows.Scan(
			&n.ID, &n.UserID, &n.Channel, &n.Type, &n.Title, &n.Body,
			&dataJSON, &n.Priority, &n.Status, &n.Recipient, &n.ProviderID,
			&n.ScheduledAt, &n.SentAt, &n.DeliveredAt, &n.FailedAt,
			&n.ErrorMsg, &n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(dataJSON, &n.Data)
		notifications = append(notifications, &n)
	}
	return notifications, rows.Err()
}
