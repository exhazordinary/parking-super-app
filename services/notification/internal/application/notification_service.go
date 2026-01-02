package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/notification/internal/domain"
	"github.com/parking-super-app/services/notification/internal/ports"
)

// NotificationService handles notification use cases
type NotificationService struct {
	notifications ports.NotificationRepository
	templates     ports.TemplateRepository
	preferences   ports.PreferenceRepository
	push          ports.PushProvider
	sms           ports.SMSProvider
	email         ports.EmailProvider
	logger        ports.Logger
}

func NewNotificationService(
	notifications ports.NotificationRepository,
	templates ports.TemplateRepository,
	preferences ports.PreferenceRepository,
	push ports.PushProvider,
	sms ports.SMSProvider,
	email ports.EmailProvider,
	logger ports.Logger,
) *NotificationService {
	return &NotificationService{
		notifications: notifications,
		templates:     templates,
		preferences:   preferences,
		push:          push,
		sms:           sms,
		email:         email,
		logger:        logger,
	}
}

// Request/Response DTOs

type SendNotificationRequest struct {
	UserID    uuid.UUID         `json:"user_id"`
	Channel   string            `json:"channel"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Recipient string            `json:"recipient"`
	Data      map[string]string `json:"data,omitempty"`
	Priority  string            `json:"priority,omitempty"`
}

type NotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Channel   string    `json:"channel"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
}

type NotificationListResponse struct {
	Notifications []*NotificationResponse `json:"notifications"`
	Total         int                     `json:"total"`
	Limit         int                     `json:"limit"`
	Offset        int                     `json:"offset"`
}

type SendFromTemplateRequest struct {
	UserID       uuid.UUID         `json:"user_id"`
	TemplateName string            `json:"template_name"`
	Recipient    string            `json:"recipient"`
	Variables    map[string]string `json:"variables"`
}

type UpdatePreferenceRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	PushEnabled  *bool     `json:"push_enabled,omitempty"`
	SMSEnabled   *bool     `json:"sms_enabled,omitempty"`
	EmailEnabled *bool     `json:"email_enabled,omitempty"`
}

type PreferenceResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	PushEnabled  bool      `json:"push_enabled"`
	SMSEnabled   bool      `json:"sms_enabled"`
	EmailEnabled bool      `json:"email_enabled"`
}

// SendNotification sends a notification to a user
func (s *NotificationService) SendNotification(ctx context.Context, req SendNotificationRequest) (*NotificationResponse, error) {
	s.logger.Info("sending notification",
		ports.String("user_id", req.UserID.String()),
		ports.String("channel", req.Channel),
	)

	channel := domain.Channel(req.Channel)

	// Check user preferences
	pref, err := s.preferences.GetByUserID(ctx, req.UserID)
	if err == nil && pref != nil {
		if !pref.IsChannelEnabled(channel) {
			s.logger.Info("notification blocked by user preference")
			return nil, fmt.Errorf("channel %s is disabled for user", req.Channel)
		}
		if pref.IsInQuietHours() && req.Priority != "high" {
			s.logger.Info("notification delayed due to quiet hours")
		}
	}

	notif, err := domain.NewNotification(
		req.UserID,
		channel,
		req.Type,
		req.Title,
		req.Body,
		req.Recipient,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Data {
		notif.AddData(k, v)
	}

	if req.Priority != "" {
		notif.SetPriority(domain.Priority(req.Priority))
	}

	// Save notification
	if err := s.notifications.Create(ctx, notif); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notification
	if err := s.send(ctx, notif); err != nil {
		notif.MarkFailed(err.Error())
		s.notifications.Update(ctx, notif)
		return nil, err
	}

	return s.toResponse(notif), nil
}

// SendFromTemplate sends notification using a template
func (s *NotificationService) SendFromTemplate(ctx context.Context, req SendFromTemplateRequest) (*NotificationResponse, error) {
	template, err := s.templates.GetByName(ctx, req.TemplateName)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	title, body := template.Render(req.Variables)

	return s.SendNotification(ctx, SendNotificationRequest{
		UserID:    req.UserID,
		Channel:   string(template.Channel),
		Type:      template.Type,
		Title:     title,
		Body:      body,
		Recipient: req.Recipient,
	})
}

// GetNotification retrieves a notification by ID
func (s *NotificationService) GetNotification(ctx context.Context, id uuid.UUID) (*NotificationResponse, error) {
	notif, err := s.notifications.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(notif), nil
}

// GetUserNotifications retrieves notifications for a user
func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) (*NotificationListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	notifications, err := s.notifications.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	total, err := s.notifications.CountByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count notifications: %w", err)
	}

	responses := make([]*NotificationResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = s.toResponse(n)
	}

	return &NotificationListResponse{
		Notifications: responses,
		Total:         total,
		Limit:         limit,
		Offset:        offset,
	}, nil
}

// UpdatePreferences updates user notification preferences
func (s *NotificationService) UpdatePreferences(ctx context.Context, req UpdatePreferenceRequest) (*PreferenceResponse, error) {
	pref, err := s.preferences.GetByUserID(ctx, req.UserID)
	if err != nil {
		// Create new preferences
		pref = domain.NewUserPreference(req.UserID)
	}

	if req.PushEnabled != nil {
		pref.SetChannelEnabled(domain.ChannelPush, *req.PushEnabled)
	}
	if req.SMSEnabled != nil {
		pref.SetChannelEnabled(domain.ChannelSMS, *req.SMSEnabled)
	}
	if req.EmailEnabled != nil {
		pref.SetChannelEnabled(domain.ChannelEmail, *req.EmailEnabled)
	}

	if err := s.preferences.Upsert(ctx, pref); err != nil {
		return nil, fmt.Errorf("failed to update preferences: %w", err)
	}

	return &PreferenceResponse{
		UserID:       pref.UserID,
		PushEnabled:  pref.PushEnabled,
		SMSEnabled:   pref.SMSEnabled,
		EmailEnabled: pref.EmailEnabled,
	}, nil
}

// GetPreferences retrieves user notification preferences
func (s *NotificationService) GetPreferences(ctx context.Context, userID uuid.UUID) (*PreferenceResponse, error) {
	pref, err := s.preferences.GetByUserID(ctx, userID)
	if err != nil {
		// Return default preferences
		pref = domain.NewUserPreference(userID)
	}

	return &PreferenceResponse{
		UserID:       pref.UserID,
		PushEnabled:  pref.PushEnabled,
		SMSEnabled:   pref.SMSEnabled,
		EmailEnabled: pref.EmailEnabled,
	}, nil
}

func (s *NotificationService) send(ctx context.Context, notif *domain.Notification) error {
	var providerID string
	var err error

	switch notif.Channel {
	case domain.ChannelPush:
		resp, sendErr := s.push.Send(ctx, ports.PushRequest{
			DeviceToken: notif.Recipient,
			Title:       notif.Title,
			Body:        notif.Body,
			Data:        notif.Data,
			Priority:    string(notif.Priority),
		})
		if sendErr != nil {
			return sendErr
		}
		providerID = resp.MessageID

	case domain.ChannelSMS:
		resp, sendErr := s.sms.Send(ctx, ports.SMSRequest{
			PhoneNumber: notif.Recipient,
			Message:     notif.Body,
		})
		if sendErr != nil {
			return sendErr
		}
		providerID = resp.MessageID

	case domain.ChannelEmail:
		resp, sendErr := s.email.Send(ctx, ports.EmailRequest{
			To:      notif.Recipient,
			Subject: notif.Title,
			Body:    notif.Body,
			IsHTML:  false,
		})
		if sendErr != nil {
			return sendErr
		}
		providerID = resp.MessageID

	default:
		return domain.ErrInvalidChannel
	}

	notif.MarkSent(providerID)
	err = s.notifications.Update(ctx, notif)

	return err
}

func (s *NotificationService) toResponse(n *domain.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Channel:   string(n.Channel),
		Type:      n.Type,
		Title:     n.Title,
		Body:      n.Body,
		Status:    string(n.Status),
		CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
