package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware_Authenticate(t *testing.T) {
	secret := "test-secret-key"
	authMw := NewAuthMiddleware(secret)

	// Create a valid token
	validToken := createTestToken(t, secret, "user-123", time.Now().Add(time.Hour))

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectUserID   bool
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectUserID:   true,
		},
		{
			name:           "missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
		{
			name:           "invalid format - no Bearer prefix",
			authHeader:     validToken,
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
		{
			name:           "expired token",
			authHeader:     "Bearer " + createTestToken(t, secret, "user-123", time.Now().Add(-time.Hour)),
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
		{
			name:           "wrong secret",
			authHeader:     "Bearer " + createTestToken(t, "wrong-secret", "user-123", time.Now().Add(time.Hour)),
			expectedStatus: http.StatusUnauthorized,
			expectUserID:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUserID string

			handler := authMw.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUserID = r.Header.Get("X-User-ID")
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectUserID && capturedUserID != "user-123" {
				t.Errorf("expected user ID 'user-123', got '%s'", capturedUserID)
			}
		})
	}
}

func TestAuthMiddleware_OptionalAuth(t *testing.T) {
	secret := "test-secret-key"
	authMw := NewAuthMiddleware(secret)

	validToken := createTestToken(t, secret, "user-456", time.Now().Add(time.Hour))

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUserID string
	}{
		{
			name:           "with valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedUserID: "user-456",
		},
		{
			name:           "without token - still passes",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			expectedUserID: "",
		},
		{
			name:           "with invalid token - still passes but no user",
			authHeader:     "Bearer invalid.token",
			expectedStatus: http.StatusOK,
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUserID string

			handler := authMw.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUserID = r.Header.Get("X-User-ID")
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if capturedUserID != tt.expectedUserID {
				t.Errorf("expected user ID '%s', got '%s'", tt.expectedUserID, capturedUserID)
			}
		})
	}
}

func createTestToken(t *testing.T, secret, userID string, expiresAt time.Time) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

	return tokenString
}
