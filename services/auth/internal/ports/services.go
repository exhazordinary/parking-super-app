package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PasswordHasher defines the contract for password hashing operations.
//
// WHY AN INTERFACE FOR PASSWORD HASHING?
// =====================================
// - We might want to upgrade algorithms (bcrypt -> argon2)
// - Testing: We can mock this to avoid slow hashing in tests
// - Flexibility: Different environments might use different implementations
type PasswordHasher interface {
	// Hash generates a secure hash of the password.
	// The implementation should use a modern algorithm like bcrypt or argon2.
	Hash(password string) (string, error)

	// Compare checks if a password matches a hash.
	// Returns nil if they match, error otherwise.
	Compare(password, hash string) error
}

// SMSService defines the contract for sending SMS messages.
//
// MICROSERVICES PATTERN: External Service Interface
// =================================================
// By defining an interface, our application layer doesn't know
// (or care) whether we're using Twilio, AWS SNS, or a local
// Malaysian SMS provider. This is the Dependency Inversion Principle.
type SMSService interface {
	// SendOTP sends an OTP code to the given phone number.
	// Returns an error if the SMS couldn't be sent.
	SendOTP(ctx context.Context, phone, code string) error

	// SendMessage sends a generic message to the given phone number.
	// Use this for notifications, welcome messages, etc.
	SendMessage(ctx context.Context, phone, message string) error
}

// TokenService defines the contract for JWT token operations.
//
// This handles the creation and validation of JWT access tokens.
// Refresh tokens are handled separately by the RefreshTokenRepository.
type TokenService interface {
	// GenerateAccessToken creates a new JWT access token for the user.
	// The token contains claims like user ID, phone, and expiration.
	GenerateAccessToken(userID uuid.UUID, phone string) (string, error)

	// ValidateAccessToken validates a JWT and returns the claims.
	// Returns an error if the token is invalid or expired.
	ValidateAccessToken(token string) (*AccessTokenClaims, error)

	// GenerateRefreshToken creates a cryptographically secure random token.
	// This is NOT a JWT - it's a random string that gets hashed and stored.
	GenerateRefreshToken() (string, error)

	// HashRefreshToken creates a SHA-256 hash of a refresh token.
	// We store the hash, not the token itself.
	HashRefreshToken(token string) string
}

// AccessTokenClaims represents the claims extracted from a JWT access token.
type AccessTokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Phone     string    `json:"phone"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
}

// OTPGenerator defines the contract for generating OTP codes.
//
// Why an interface? In tests, we might want predictable OTPs.
// In production, we want cryptographically secure random codes.
type OTPGenerator interface {
	// Generate creates a new OTP code.
	// Typically a 6-digit numeric code.
	Generate() string
}

// EventPublisher defines the contract for publishing domain events.
//
// MICROSERVICES PATTERN: Event-Driven Architecture
// ================================================
// When something important happens (user registered, password changed),
// we publish events. Other services can subscribe to these events.
//
// Benefits:
// - Loose coupling: Auth service doesn't need to know about other services
// - Async processing: Slow operations (send welcome email) don't block the request
// - Audit trail: Events can be logged/stored for debugging
//
// In production, this would publish to Kafka or NATS.
// For simplicity, we might start with an in-memory implementation.
type EventPublisher interface {
	// Publish sends an event to the message broker.
	// The event is sent asynchronously - this method returns immediately.
	Publish(ctx context.Context, event Event) error
}

// Event represents a domain event.
type Event struct {
	Type      string                 `json:"type"`      // e.g., "user.registered", "user.password_changed"
	Payload   map[string]interface{} `json:"payload"`   // Event-specific data
	Timestamp time.Time              `json:"timestamp"`
	TraceID   string                 `json:"trace_id,omitempty"` // For distributed tracing
}

// Common event types
const (
	EventUserRegistered     = "user.registered"
	EventUserActivated      = "user.activated"
	EventUserLoggedIn       = "user.logged_in"
	EventUserLoggedOut      = "user.logged_out"
	EventPasswordChanged    = "user.password_changed"
	EventPasswordReset      = "user.password_reset"
	EventTokenRefreshed     = "user.token_refreshed"
	EventOTPRequested       = "user.otp_requested"
	EventOTPVerified        = "user.otp_verified"
)

// Logger defines the contract for structured logging.
//
// We use an interface instead of a concrete logger so we can:
// - Switch implementations (zap, zerolog, logrus)
// - Mock in tests
// - Add middleware (add trace IDs, user IDs to all logs)
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)

	// WithFields returns a new logger with the given fields attached.
	// All subsequent logs will include these fields.
	WithFields(fields ...Field) Logger
}

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// Helper functions for creating fields
func String(key, value string) Field { return Field{Key: key, Value: value} }
func Int(key string, value int) Field { return Field{Key: key, Value: value} }
func Bool(key string, value bool) Field { return Field{Key: key, Value: value} }
func Err(err error) Field { return Field{Key: "error", Value: err} }
func Any(key string, value interface{}) Field { return Field{Key: key, Value: value} }
