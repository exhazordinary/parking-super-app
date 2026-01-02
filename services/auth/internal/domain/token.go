package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Token-related domain errors
var (
	ErrTokenNotFound  = errors.New("token not found")
	ErrTokenExpired   = errors.New("token has expired")
	ErrTokenRevoked   = errors.New("token has been revoked")
	ErrInvalidToken   = errors.New("invalid token")
)

// RefreshToken represents a refresh token stored in the database.
//
// SECURITY PATTERN: Token Rotation
// ================================
// We store refresh tokens (not access tokens) because:
// 1. Access tokens are short-lived (15 min) - no need to store
// 2. Refresh tokens are long-lived (7 days) - need to track for revocation
// 3. This allows us to implement "logout everywhere" by revoking all tokens
//
// WHY STORE TOKEN HASH, NOT THE TOKEN?
// If our database is compromised, attackers can't use the hashes
// to create valid refresh tokens. We only store the hash.
type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"` // SHA-256 hash of the actual token
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`

	// Metadata for security tracking
	UserAgent string `json:"user_agent,omitempty"` // Browser/app info
	IPAddress string `json:"ip_address,omitempty"` // IP when token was created
}

// RefreshTokenDuration is how long refresh tokens are valid.
// SECURITY: Don't make this too long - 7 days is a good balance.
const RefreshTokenDuration = 7 * 24 * time.Hour

// NewRefreshToken creates a new RefreshToken entity.
//
// IMPORTANT: The tokenHash parameter should be a SHA-256 hash
// of the actual token string. The actual token is returned to
// the user, but we only store the hash.
func NewRefreshToken(userID uuid.UUID, tokenHash, userAgent, ipAddress string) *RefreshToken {
	now := time.Now().UTC()

	return &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(RefreshTokenDuration),
		Revoked:   false,
		CreatedAt: now,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
}

// IsValid checks if the token can be used.
func (rt *RefreshToken) IsValid() bool {
	if rt.Revoked {
		return false
	}
	if time.Now().UTC().After(rt.ExpiresAt) {
		return false
	}
	return true
}

// Revoke marks the token as revoked.
// This is called during logout or when rotating tokens.
func (rt *RefreshToken) Revoke() {
	now := time.Now().UTC()
	rt.Revoked = true
	rt.RevokedAt = &now
}

// Validate checks the token and returns an appropriate error.
func (rt *RefreshToken) Validate() error {
	if rt.Revoked {
		return ErrTokenRevoked
	}
	if time.Now().UTC().After(rt.ExpiresAt) {
		return ErrTokenExpired
	}
	return nil
}

// OTP represents a one-time password for phone verification.
//
// SECURITY PATTERN: Rate Limiting & Expiry
// We should limit:
// - How many OTPs can be requested per phone number
// - How many verification attempts per OTP
// These limits are enforced in the application layer.
type OTP struct {
	ID          uuid.UUID `json:"id"`
	Phone       string    `json:"phone"`
	Code        string    `json:"-"` // 6-digit code, don't expose in JSON
	ExpiresAt   time.Time `json:"expires_at"`
	Verified    bool      `json:"verified"`
	Attempts    int       `json:"attempts"` // Track failed verification attempts
	CreatedAt   time.Time `json:"created_at"`
}

// OTPDuration is how long an OTP is valid.
const OTPDuration = 5 * time.Minute

// MaxOTPAttempts is the maximum number of verification attempts.
const MaxOTPAttempts = 3

// NewOTP creates a new OTP entity.
func NewOTP(phone, code string) *OTP {
	now := time.Now().UTC()

	return &OTP{
		ID:        uuid.New(),
		Phone:     phone,
		Code:      code,
		ExpiresAt: now.Add(OTPDuration),
		Verified:  false,
		Attempts:  0,
		CreatedAt: now,
	}
}

// IsValid checks if the OTP can still be used.
func (o *OTP) IsValid() bool {
	if o.Verified {
		return false // Already used
	}
	if o.Attempts >= MaxOTPAttempts {
		return false // Too many failed attempts
	}
	if time.Now().UTC().After(o.ExpiresAt) {
		return false // Expired
	}
	return true
}

// Verify attempts to verify the OTP with the given code.
// Returns true if verification succeeds.
func (o *OTP) Verify(code string) bool {
	if !o.IsValid() {
		return false
	}

	o.Attempts++

	if o.Code == code {
		o.Verified = true
		return true
	}

	return false
}
