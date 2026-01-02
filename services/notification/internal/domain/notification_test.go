package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewNotification(t *testing.T) {
	userID := uuid.New()

	notif, err := NewNotification(userID, ChannelPush, "test", "Title", "Body", "device-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if notif.ID == uuid.Nil {
		t.Error("expected ID to be set")
	}
	if notif.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, notif.UserID)
	}
	if notif.Channel != ChannelPush {
		t.Errorf("expected channel push, got %s", notif.Channel)
	}
	if notif.Status != StatusPending {
		t.Errorf("expected status pending, got %s", notif.Status)
	}
	if notif.Priority != PriorityNormal {
		t.Errorf("expected priority normal, got %s", notif.Priority)
	}
}

func TestNewNotification_InvalidChannel(t *testing.T) {
	_, err := NewNotification(uuid.New(), "invalid", "test", "Title", "Body", "recipient")
	if err != ErrInvalidChannel {
		t.Errorf("expected ErrInvalidChannel, got %v", err)
	}
}

func TestNewNotification_InvalidRecipient(t *testing.T) {
	_, err := NewNotification(uuid.New(), ChannelPush, "test", "Title", "Body", "")
	if err != ErrInvalidRecipient {
		t.Errorf("expected ErrInvalidRecipient, got %v", err)
	}
}

func TestNotification_MarkSent(t *testing.T) {
	notif, _ := NewNotification(uuid.New(), ChannelSMS, "test", "Title", "Body", "+60123456789")

	notif.MarkSent("provider-123")

	if notif.Status != StatusSent {
		t.Errorf("expected status sent, got %s", notif.Status)
	}
	if notif.ProviderID != "provider-123" {
		t.Errorf("expected provider ID provider-123, got %s", notif.ProviderID)
	}
	if notif.SentAt == nil {
		t.Error("expected sent_at to be set")
	}
}

func TestNotification_MarkDelivered(t *testing.T) {
	notif, _ := NewNotification(uuid.New(), ChannelEmail, "test", "Title", "Body", "test@example.com")
	notif.MarkSent("provider-123")

	notif.MarkDelivered()

	if notif.Status != StatusDelivered {
		t.Errorf("expected status delivered, got %s", notif.Status)
	}
	if notif.DeliveredAt == nil {
		t.Error("expected delivered_at to be set")
	}
}

func TestNotification_MarkFailed(t *testing.T) {
	notif, _ := NewNotification(uuid.New(), ChannelPush, "test", "Title", "Body", "token")

	notif.MarkFailed("connection timeout")

	if notif.Status != StatusFailed {
		t.Errorf("expected status failed, got %s", notif.Status)
	}
	if notif.ErrorMsg != "connection timeout" {
		t.Errorf("expected error message, got %s", notif.ErrorMsg)
	}
	if notif.FailedAt == nil {
		t.Error("expected failed_at to be set")
	}
}

func TestNotification_IsReady(t *testing.T) {
	notif, _ := NewNotification(uuid.New(), ChannelPush, "test", "Title", "Body", "token")

	if !notif.IsReady() {
		t.Error("pending notification should be ready")
	}

	// Schedule for future
	future := time.Now().Add(1 * time.Hour)
	notif.Schedule(future)

	if notif.IsReady() {
		t.Error("future scheduled notification should not be ready")
	}

	// Schedule for past
	past := time.Now().Add(-1 * time.Hour)
	notif.Schedule(past)

	if !notif.IsReady() {
		t.Error("past scheduled notification should be ready")
	}
}

func TestNotification_AddData(t *testing.T) {
	notif, _ := NewNotification(uuid.New(), ChannelPush, "test", "Title", "Body", "token")

	notif.AddData("key1", "value1")
	notif.AddData("key2", "value2")

	if notif.Data["key1"] != "value1" {
		t.Errorf("expected data key1=value1, got %s", notif.Data["key1"])
	}
	if notif.Data["key2"] != "value2" {
		t.Errorf("expected data key2=value2, got %s", notif.Data["key2"])
	}
}

func TestIsValidChannel(t *testing.T) {
	tests := []struct {
		channel Channel
		valid   bool
	}{
		{ChannelPush, true},
		{ChannelSMS, true},
		{ChannelEmail, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.channel), func(t *testing.T) {
			result := isValidChannel(tt.channel)
			if result != tt.valid {
				t.Errorf("isValidChannel(%s) = %v, want %v", tt.channel, result, tt.valid)
			}
		})
	}
}
