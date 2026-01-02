package domain

import (
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProviderNotFound      = errors.New("provider not found")
	ErrProviderAlreadyExists = errors.New("provider already exists")
	ErrInvalidProviderCode   = errors.New("provider code must be alphanumeric")
	ErrInvalidWebhookURL     = errors.New("invalid webhook URL")
	ErrInvalidMFEURL         = errors.New("invalid MFE URL")
	ErrProviderInactive      = errors.New("provider is inactive")
)

// ProviderStatus represents the operational status of a parking provider
type ProviderStatus string

const (
	ProviderStatusActive   ProviderStatus = "active"
	ProviderStatusInactive ProviderStatus = "inactive"
	ProviderStatusPending  ProviderStatus = "pending"
)

// Provider represents a parking provider that integrates with the super app.
// Each provider operates their own parking infrastructure and exposes it via MFE.
type Provider struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Code          string         `json:"code"`
	Description   string         `json:"description"`
	LogoURL       string         `json:"logo_url,omitempty"`
	Status        ProviderStatus `json:"status"`
	MFEURL        string         `json:"mfe_url"`
	APIBaseURL    string         `json:"api_base_url"`
	WebhookSecret string         `json:"-"`
	Config        ProviderConfig `json:"config"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// ProviderConfig holds provider-specific configuration
type ProviderConfig struct {
	SupportedPaymentMethods []string          `json:"supported_payment_methods"`
	MaxSessionDuration      int               `json:"max_session_duration_hours"`
	RequiresPlateValidation bool              `json:"requires_plate_validation"`
	Features                map[string]bool   `json:"features"`
	CustomSettings          map[string]string `json:"custom_settings"`
}

// NewProvider creates a new provider with default values
func NewProvider(name, code, mfeURL, apiBaseURL string) (*Provider, error) {
	if !isValidCode(code) {
		return nil, ErrInvalidProviderCode
	}
	if !isValidURL(mfeURL) {
		return nil, ErrInvalidMFEURL
	}

	now := time.Now().UTC()
	return &Provider{
		ID:         uuid.New(),
		Name:       name,
		Code:       code,
		Status:     ProviderStatusPending,
		MFEURL:     mfeURL,
		APIBaseURL: apiBaseURL,
		Config: ProviderConfig{
			SupportedPaymentMethods: []string{"wallet"},
			MaxSessionDuration:      24,
			RequiresPlateValidation: false,
			Features:                make(map[string]bool),
			CustomSettings:          make(map[string]string),
		},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// IsActive checks if provider is available for use
func (p *Provider) IsActive() bool {
	return p.Status == ProviderStatusActive
}

// Activate sets provider status to active
func (p *Provider) Activate() {
	p.Status = ProviderStatusActive
	p.UpdatedAt = time.Now().UTC()
}

// Deactivate sets provider status to inactive
func (p *Provider) Deactivate() {
	p.Status = ProviderStatusInactive
	p.UpdatedAt = time.Now().UTC()
}

// SetWebhookSecret sets the webhook secret for signature verification
func (p *Provider) SetWebhookSecret(secret string) {
	p.WebhookSecret = secret
	p.UpdatedAt = time.Now().UTC()
}

// UpdateMFEURL updates the MFE URL after validation
func (p *Provider) UpdateMFEURL(mfeURL string) error {
	if !isValidURL(mfeURL) {
		return ErrInvalidMFEURL
	}
	p.MFEURL = mfeURL
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// AddFeature enables a feature for this provider
func (p *Provider) AddFeature(feature string) {
	if p.Config.Features == nil {
		p.Config.Features = make(map[string]bool)
	}
	p.Config.Features[feature] = true
	p.UpdatedAt = time.Now().UTC()
}

// HasFeature checks if a feature is enabled
func (p *Provider) HasFeature(feature string) bool {
	if p.Config.Features == nil {
		return false
	}
	return p.Config.Features[feature]
}

// isValidCode checks if provider code is alphanumeric and reasonable length
func isValidCode(code string) bool {
	if len(code) < 2 || len(code) > 20 {
		return false
	}
	for _, c := range code {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

// isValidURL checks if a URL is valid
func isValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
