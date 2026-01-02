package domain

import (
	"testing"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name       string
		provName   string
		code       string
		mfeURL     string
		apiBaseURL string
		wantErr    error
	}{
		{
			name:       "valid provider",
			provName:   "Test Provider",
			code:       "test-prov",
			mfeURL:     "https://mfe.example.com",
			apiBaseURL: "https://api.example.com",
			wantErr:    nil,
		},
		{
			name:       "invalid code - too short",
			provName:   "Test",
			code:       "t",
			mfeURL:     "https://mfe.example.com",
			apiBaseURL: "https://api.example.com",
			wantErr:    ErrInvalidProviderCode,
		},
		{
			name:       "invalid code - special characters",
			provName:   "Test",
			code:       "test@provider",
			mfeURL:     "https://mfe.example.com",
			apiBaseURL: "https://api.example.com",
			wantErr:    ErrInvalidProviderCode,
		},
		{
			name:       "invalid MFE URL",
			provName:   "Test",
			code:       "test-prov",
			mfeURL:     "not-a-url",
			apiBaseURL: "https://api.example.com",
			wantErr:    ErrInvalidMFEURL,
		},
		{
			name:       "empty MFE URL",
			provName:   "Test",
			code:       "test-prov",
			mfeURL:     "",
			apiBaseURL: "https://api.example.com",
			wantErr:    ErrInvalidMFEURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(tt.provName, tt.code, tt.mfeURL, tt.apiBaseURL)

			if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if provider.Name != tt.provName {
					t.Errorf("expected name %s, got %s", tt.provName, provider.Name)
				}
				if provider.Code != tt.code {
					t.Errorf("expected code %s, got %s", tt.code, provider.Code)
				}
				if provider.Status != ProviderStatusPending {
					t.Errorf("expected status pending, got %s", provider.Status)
				}
			}
		})
	}
}

func TestProvider_Activate(t *testing.T) {
	provider, _ := NewProvider("Test", "test", "https://mfe.example.com", "https://api.example.com")

	if provider.IsActive() {
		t.Error("new provider should not be active")
	}

	provider.Activate()

	if !provider.IsActive() {
		t.Error("provider should be active after activation")
	}
	if provider.Status != ProviderStatusActive {
		t.Errorf("expected status active, got %s", provider.Status)
	}
}

func TestProvider_Deactivate(t *testing.T) {
	provider, _ := NewProvider("Test", "test", "https://mfe.example.com", "https://api.example.com")
	provider.Activate()

	provider.Deactivate()

	if provider.IsActive() {
		t.Error("provider should not be active after deactivation")
	}
	if provider.Status != ProviderStatusInactive {
		t.Errorf("expected status inactive, got %s", provider.Status)
	}
}

func TestProvider_UpdateMFEURL(t *testing.T) {
	provider, _ := NewProvider("Test", "test", "https://mfe.example.com", "https://api.example.com")

	err := provider.UpdateMFEURL("https://new-mfe.example.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if provider.MFEURL != "https://new-mfe.example.com" {
		t.Errorf("expected new MFE URL, got %s", provider.MFEURL)
	}

	err = provider.UpdateMFEURL("invalid-url")
	if err != ErrInvalidMFEURL {
		t.Errorf("expected ErrInvalidMFEURL, got %v", err)
	}
}

func TestProvider_AddFeature(t *testing.T) {
	provider, _ := NewProvider("Test", "test", "https://mfe.example.com", "https://api.example.com")

	if provider.HasFeature("ev_charging") {
		t.Error("should not have feature before adding")
	}

	provider.AddFeature("ev_charging")

	if !provider.HasFeature("ev_charging") {
		t.Error("should have feature after adding")
	}
}

func TestIsValidCode(t *testing.T) {
	tests := []struct {
		code  string
		valid bool
	}{
		{"abc", true},
		{"ABC", true},
		{"abc123", true},
		{"test-code", true},
		{"test_code", true},
		{"a", false},
		{"", false},
		{"test@code", false},
		{"test code", false},
		{"abcdefghijklmnopqrstuvwxyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := isValidCode(tt.code)
			if result != tt.valid {
				t.Errorf("isValidCode(%s) = %v, want %v", tt.code, result, tt.valid)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url   string
		valid bool
	}{
		{"https://example.com", true},
		{"http://example.com", true},
		{"https://example.com/path", true},
		{"", false},
		{"not-a-url", false},
		{"ftp://example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := isValidURL(tt.url)
			if result != tt.valid {
				t.Errorf("isValidURL(%s) = %v, want %v", tt.url, result, tt.valid)
			}
		})
	}
}
