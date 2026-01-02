package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewProviderCredentials(t *testing.T) {
	providerID := uuid.New()

	creds, err := NewProviderCredentials(providerID, EnvironmentSandbox)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if creds.ID == uuid.Nil {
		t.Error("expected credentials ID to be set")
	}
	if creds.ProviderID != providerID {
		t.Errorf("expected provider ID %v, got %v", providerID, creds.ProviderID)
	}
	if creds.APIKey == "" {
		t.Error("expected API key to be set")
	}
	if len(creds.APIKey) != 64 {
		t.Errorf("expected API key length 64, got %d", len(creds.APIKey))
	}
	if creds.APISecret == "" {
		t.Error("expected API secret to be set")
	}
	if len(creds.APISecret) != 128 {
		t.Errorf("expected API secret length 128, got %d", len(creds.APISecret))
	}
	if creds.Environment != EnvironmentSandbox {
		t.Errorf("expected environment sandbox, got %s", creds.Environment)
	}
	if !creds.IsActive {
		t.Error("expected credentials to be active")
	}
}

func TestProviderCredentials_IsExpired(t *testing.T) {
	creds, _ := NewProviderCredentials(uuid.New(), EnvironmentSandbox)

	if creds.IsExpired() {
		t.Error("credentials without expiry should not be expired")
	}

	// Set expiration in the past
	past := time.Now().Add(-1 * time.Hour)
	creds.SetExpiration(past)

	if !creds.IsExpired() {
		t.Error("credentials with past expiry should be expired")
	}

	// Set expiration in the future
	future := time.Now().Add(1 * time.Hour)
	creds.SetExpiration(future)

	if creds.IsExpired() {
		t.Error("credentials with future expiry should not be expired")
	}
}

func TestProviderCredentials_IsValid(t *testing.T) {
	creds, _ := NewProviderCredentials(uuid.New(), EnvironmentSandbox)

	if !creds.IsValid() {
		t.Error("new credentials should be valid")
	}

	creds.Revoke()

	if creds.IsValid() {
		t.Error("revoked credentials should not be valid")
	}
}

func TestProviderCredentials_Revoke(t *testing.T) {
	creds, _ := NewProviderCredentials(uuid.New(), EnvironmentProduction)

	if !creds.IsActive {
		t.Error("new credentials should be active")
	}

	creds.Revoke()

	if creds.IsActive {
		t.Error("credentials should be inactive after revoke")
	}
}
