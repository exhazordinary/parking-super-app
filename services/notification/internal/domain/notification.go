package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidChannel       = errors.New("invalid notification channel")
	ErrInvalidRecipient     = errors.New("invalid recipient")
	ErrNotificationFailed   = errors.New("notification delivery failed")
)

// Channel represents a notification delivery channel
type Channel string

const (
	ChannelPush  Channel = "push"
	ChannelSMS   Channel = "sms"
	ChannelEmail Channel = "email"
)

// Status represents notification delivery status
type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusDelivered Status = "delivered"
	StatusFailed    Status = "failed"
)

// Priority represents notification urgency
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
)

// Notification represents a notification to be sent to a user
type Notification struct {
	ID          uuid.UUID         `json:"id"`
	UserID      uuid.UUID         `json:"user_id"`
	Channel     Channel           `json:"channel"`
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data,omitempty"`
	Priority    Priority          `json:"priority"`
	Status      Status            `json:"status"`
	Recipient   string            `json:"recipient"`
	ProviderID  string            `json:"provider_id,omitempty"`
	ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
	SentAt      *time.Time        `json:"sent_at,omitempty"`
	DeliveredAt *time.Time        `json:"delivered_at,omitempty"`
	FailedAt    *time.Time        `json:"failed_at,omitempty"`
	ErrorMsg    string            `json:"error_msg,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// NewNotification creates a new notification
func NewNotification(
	userID uuid.UUID,
	channel Channel,
	notifType, title, body, recipient string,
) (*Notification, error) {
	if !isValidChannel(channel) {
		return nil, ErrInvalidChannel
	}
	if recipient == "" {
		return nil, ErrInvalidRecipient
	}

	return &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Channel:   channel,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Data:      make(map[string]string),
		Priority:  PriorityNormal,
		Status:    StatusPending,
		Recipient: recipient,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// SetPriority sets the notification priority
func (n *Notification) SetPriority(priority Priority) {
	n.Priority = priority
}

// Schedule sets a future delivery time
func (n *Notification) Schedule(at time.Time) {
	n.ScheduledAt = &at
}

// AddData adds custom data to the notification
func (n *Notification) AddData(key, value string) {
	if n.Data == nil {
		n.Data = make(map[string]string)
	}
	n.Data[key] = value
}

// MarkSent updates status to sent
func (n *Notification) MarkSent(providerID string) {
	now := time.Now().UTC()
	n.Status = StatusSent
	n.ProviderID = providerID
	n.SentAt = &now
}

// MarkDelivered updates status to delivered
func (n *Notification) MarkDelivered() {
	now := time.Now().UTC()
	n.Status = StatusDelivered
	n.DeliveredAt = &now
}

// MarkFailed records delivery failure
func (n *Notification) MarkFailed(errMsg string) {
	now := time.Now().UTC()
	n.Status = StatusFailed
	n.FailedAt = &now
	n.ErrorMsg = errMsg
}

// IsReady checks if notification is ready to send
func (n *Notification) IsReady() bool {
	if n.Status != StatusPending {
		return false
	}
	if n.ScheduledAt != nil && n.ScheduledAt.After(time.Now()) {
		return false
	}
	return true
}

// IsPending returns true if notification has not been sent
func (n *Notification) IsPending() bool {
	return n.Status == StatusPending
}

func isValidChannel(c Channel) bool {
	return c == ChannelPush || c == ChannelSMS || c == ChannelEmail
}
