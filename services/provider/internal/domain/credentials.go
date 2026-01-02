package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// Environment represents the deployment environment for credentials
type Environment string

const (
	EnvironmentSandbox    Environment = "sandbox"
	EnvironmentProduction Environment = "production"
)

// ProviderCredentials stores API credentials for a provider
// These are used to authenticate requests from the super app to the provider
type ProviderCredentials struct {
	ID          uuid.UUID   `json:"id"`
	ProviderID  uuid.UUID   `json:"provider_id"`
	APIKey      string      `json:"api_key"`
	APISecret   string      `json:"-"`
	Environment Environment `json:"environment"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
}

// NewProviderCredentials creates new credentials for a provider
func NewProviderCredentials(providerID uuid.UUID, env Environment) (*ProviderCredentials, error) {
	apiKey, err := generateSecureKey(32)
	if err != nil {
		return nil, err
	}
	apiSecret, err := generateSecureKey(64)
	if err != nil {
		return nil, err
	}

	return &ProviderCredentials{
		ID:          uuid.New(),
		ProviderID:  providerID,
		APIKey:      apiKey,
		APISecret:   apiSecret,
		Environment: env,
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// IsExpired checks if credentials have expired
func (c *ProviderCredentials) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// IsValid checks if credentials are active and not expired
func (c *ProviderCredentials) IsValid() bool {
	return c.IsActive && !c.IsExpired()
}

// Revoke invalidates the credentials
func (c *ProviderCredentials) Revoke() {
	c.IsActive = false
}

// SetExpiration sets an expiration date for the credentials
func (c *ProviderCredentials) SetExpiration(expiresAt time.Time) {
	c.ExpiresAt = &expiresAt
}

// generateSecureKey generates a cryptographically secure random key
func generateSecureKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
