package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserPreference stores user notification preferences
type UserPreference struct {
	ID              uuid.UUID         `json:"id"`
	UserID          uuid.UUID         `json:"user_id"`
	PushEnabled     bool              `json:"push_enabled"`
	SMSEnabled      bool              `json:"sms_enabled"`
	EmailEnabled    bool              `json:"email_enabled"`
	QuietHoursStart *int              `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd   *int              `json:"quiet_hours_end,omitempty"`
	TypePreferences map[string]bool   `json:"type_preferences"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// NewUserPreference creates default preferences for a user
func NewUserPreference(userID uuid.UUID) *UserPreference {
	now := time.Now().UTC()
	return &UserPreference{
		ID:              uuid.New(),
		UserID:          userID,
		PushEnabled:     true,
		SMSEnabled:      true,
		EmailEnabled:    true,
		TypePreferences: make(map[string]bool),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// IsChannelEnabled checks if a channel is enabled
func (p *UserPreference) IsChannelEnabled(channel Channel) bool {
	switch channel {
	case ChannelPush:
		return p.PushEnabled
	case ChannelSMS:
		return p.SMSEnabled
	case ChannelEmail:
		return p.EmailEnabled
	default:
		return false
	}
}

// IsTypeEnabled checks if a notification type is enabled
func (p *UserPreference) IsTypeEnabled(notifType string) bool {
	// If no specific preference, default to enabled
	enabled, exists := p.TypePreferences[notifType]
	if !exists {
		return true
	}
	return enabled
}

// SetChannelEnabled enables/disables a channel
func (p *UserPreference) SetChannelEnabled(channel Channel, enabled bool) {
	switch channel {
	case ChannelPush:
		p.PushEnabled = enabled
	case ChannelSMS:
		p.SMSEnabled = enabled
	case ChannelEmail:
		p.EmailEnabled = enabled
	}
	p.UpdatedAt = time.Now().UTC()
}

// SetTypeEnabled enables/disables a notification type
func (p *UserPreference) SetTypeEnabled(notifType string, enabled bool) {
	if p.TypePreferences == nil {
		p.TypePreferences = make(map[string]bool)
	}
	p.TypePreferences[notifType] = enabled
	p.UpdatedAt = time.Now().UTC()
}

// SetQuietHours sets the quiet hours window (24-hour format)
func (p *UserPreference) SetQuietHours(start, end int) {
	p.QuietHoursStart = &start
	p.QuietHoursEnd = &end
	p.UpdatedAt = time.Now().UTC()
}

// IsInQuietHours checks if current time is in quiet hours
func (p *UserPreference) IsInQuietHours() bool {
	if p.QuietHoursStart == nil || p.QuietHoursEnd == nil {
		return false
	}

	currentHour := time.Now().Hour()
	start := *p.QuietHoursStart
	end := *p.QuietHoursEnd

	if start < end {
		return currentHour >= start && currentHour < end
	}
	// Quiet hours span midnight
	return currentHour >= start || currentHour < end
}
