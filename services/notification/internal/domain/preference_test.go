package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUserPreference(t *testing.T) {
	userID := uuid.New()

	pref := NewUserPreference(userID)

	if pref.ID == uuid.Nil {
		t.Error("expected ID to be set")
	}
	if pref.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, pref.UserID)
	}
	if !pref.PushEnabled {
		t.Error("expected push to be enabled by default")
	}
	if !pref.SMSEnabled {
		t.Error("expected sms to be enabled by default")
	}
	if !pref.EmailEnabled {
		t.Error("expected email to be enabled by default")
	}
}

func TestUserPreference_IsChannelEnabled(t *testing.T) {
	pref := NewUserPreference(uuid.New())

	if !pref.IsChannelEnabled(ChannelPush) {
		t.Error("push should be enabled by default")
	}

	pref.SetChannelEnabled(ChannelPush, false)

	if pref.IsChannelEnabled(ChannelPush) {
		t.Error("push should be disabled")
	}
}

func TestUserPreference_IsTypeEnabled(t *testing.T) {
	pref := NewUserPreference(uuid.New())

	// Should be enabled by default
	if !pref.IsTypeEnabled("payment.success") {
		t.Error("type should be enabled by default")
	}

	pref.SetTypeEnabled("payment.success", false)

	if pref.IsTypeEnabled("payment.success") {
		t.Error("type should be disabled")
	}
}

func TestUserPreference_QuietHours(t *testing.T) {
	pref := NewUserPreference(uuid.New())

	if pref.IsInQuietHours() {
		t.Error("should not be in quiet hours without setting")
	}

	// Set quiet hours that span current time
	pref.SetQuietHours(0, 24)

	if !pref.IsInQuietHours() {
		t.Error("should be in quiet hours")
	}
}

func TestUserPreference_SetChannelEnabled(t *testing.T) {
	pref := NewUserPreference(uuid.New())

	pref.SetChannelEnabled(ChannelSMS, false)
	if pref.SMSEnabled {
		t.Error("sms should be disabled")
	}

	pref.SetChannelEnabled(ChannelEmail, false)
	if pref.EmailEnabled {
		t.Error("email should be disabled")
	}
}
