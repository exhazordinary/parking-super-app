package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/notification/internal/domain"
)

// NotificationRepository defines persistence for notifications
type NotificationRepository interface {
	Create(ctx context.Context, notif *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Notification, error)
	GetPending(ctx context.Context, limit int) ([]*domain.Notification, error)
	Update(ctx context.Context, notif *domain.Notification) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
}

// TemplateRepository defines persistence for notification templates
type TemplateRepository interface {
	Create(ctx context.Context, template *domain.Template) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error)
	GetByName(ctx context.Context, name string) (*domain.Template, error)
	GetByType(ctx context.Context, notifType string, channel domain.Channel) (*domain.Template, error)
	GetAll(ctx context.Context) ([]*domain.Template, error)
	Update(ctx context.Context, template *domain.Template) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PreferenceRepository defines persistence for user preferences
type PreferenceRepository interface {
	Create(ctx context.Context, pref *domain.UserPreference) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserPreference, error)
	Update(ctx context.Context, pref *domain.UserPreference) error
	Upsert(ctx context.Context, pref *domain.UserPreference) error
}
