package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/auth/internal/application"
	"github.com/parking-super-app/services/auth/internal/domain"
	"github.com/parking-super-app/services/auth/internal/ports"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// UserIDKey is the context key for the authenticated user's ID.
	UserIDKey contextKey = "user_id"
)

// AuthHandler handles HTTP requests for authentication endpoints.
type AuthHandler struct {
	authService  *application.AuthService
	tokenService ports.TokenService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SetTokenService sets the token service for JWT validation.
// This is called after the handler is created because of circular dependencies.
func (h *AuthHandler) SetTokenService(ts ports.TokenService) {
	h.tokenService = ts
}

// ---- Response Helpers ----
// These functions help create consistent JSON responses.

// APIResponse is the standard response format.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error in the API response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	})
}

// writeError writes an error response.
//
// BEST PRACTICE: Error Responses
// ==============================
// - Use consistent error format
// - Include error codes for programmatic handling
// - Don't expose internal errors to clients
// - Log detailed errors server-side
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// mapDomainError maps domain errors to HTTP status codes and error codes.
func mapDomainError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "USER_NOT_FOUND", "User not found"
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, "USER_EXISTS", "A user with this phone number already exists"
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid phone number or password"
	case errors.Is(err, domain.ErrInvalidEmail):
		return http.StatusBadRequest, "INVALID_EMAIL", "Invalid email format"
	case errors.Is(err, domain.ErrInvalidPhone):
		return http.StatusBadRequest, "INVALID_PHONE", "Invalid phone number format. Use +60xxxxxxxxx"
	case errors.Is(err, domain.ErrWeakPassword):
		return http.StatusBadRequest, "WEAK_PASSWORD", "Password must be at least 8 characters"
	case errors.Is(err, domain.ErrUserInactive):
		return http.StatusForbidden, "USER_INACTIVE", "Your account is inactive"
	case errors.Is(err, domain.ErrTokenExpired):
		return http.StatusUnauthorized, "TOKEN_EXPIRED", "Token has expired"
	case errors.Is(err, domain.ErrTokenRevoked):
		return http.StatusUnauthorized, "TOKEN_REVOKED", "Token has been revoked"
	case errors.Is(err, domain.ErrInvalidToken):
		return http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred"
	}
}

// ---- Handlers ----

// Register handles user registration.
//
// POST /api/v1/auth/register
// Request: { "phone": "+60123456789", "password": "...", "full_name": "...", "email": "..." }
// Response: { "success": true, "data": { "user_id": "...", "message": "..." } }
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req application.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.authService.Register(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// Login handles user login.
//
// POST /api/v1/auth/login
// Request: { "phone": "+60123456789", "password": "..." }
// Response: { "success": true, "data": { "access_token": "...", "refresh_token": "...", "expires_in": 900 } }
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req application.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Get client info for security tracking
	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.RemoteAddr

	resp, err := h.authService.Login(r.Context(), req, userAgent, ipAddress)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// RefreshToken handles token refresh.
//
// POST /api/v1/auth/refresh
// Request: { "refresh_token": "..." }
// Response: { "success": true, "data": { "access_token": "...", "refresh_token": "...", "expires_in": 900 } }
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req application.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	userAgent := r.Header.Get("User-Agent")
	ipAddress := r.RemoteAddr

	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken, userAgent, ipAddress)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// RequestOTP handles OTP request.
//
// POST /api/v1/auth/otp/request
// Request: { "phone": "+60123456789" }
// Response: { "success": true, "data": { "message": "OTP sent" } }
func (h *AuthHandler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req application.RequestOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Always return success to prevent phone enumeration attacks
	_ = h.authService.RequestOTP(r.Context(), req)

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "If the phone number is registered, an OTP has been sent",
	})
}

// VerifyOTP handles OTP verification.
//
// POST /api/v1/auth/otp/verify
// Request: { "phone": "+60123456789", "code": "123456" }
// Response: { "success": true, "data": { "message": "Phone verified" } }
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req application.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.authService.VerifyOTP(r.Context(), req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_OTP", "Invalid or expired OTP")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Phone number verified successfully",
	})
}

// Logout handles user logout.
//
// POST /api/v1/auth/logout (requires authentication)
// Request: { "refresh_token": "..." }
// Response: { "success": true }
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		// Log but don't fail - user should be logged out regardless
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// LogoutAllDevices handles logout from all devices.
//
// POST /api/v1/auth/logout/all (requires authentication)
// Response: { "success": true }
func (h *AuthHandler) LogoutAllDevices(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)

	if err := h.authService.LogoutAllDevices(r.Context(), userID); err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out from all devices",
	})
}

// GetProfile handles getting user profile.
//
// GET /api/v1/auth/me (requires authentication)
// Response: { "success": true, "data": { "id": "...", "phone": "...", ... } }
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)

	profile, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// ---- Middleware ----

// AuthMiddleware validates JWT access tokens and sets user ID in context.
//
// PATTERN: Middleware Authentication
// ==================================
// 1. Extract token from Authorization header
// 2. Validate token (signature, expiration)
// 3. Set user ID in context for downstream handlers
// 4. Continue to next handler or return 401
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header required")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid authorization format")
			return
		}
		token := parts[1]

		// Validate token
		if h.tokenService == nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Token service not configured")
			return
		}

		claims, err := h.tokenService.ValidateAccessToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
