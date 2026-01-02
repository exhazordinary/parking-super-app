// Package ports defines the interfaces (ports) that our application uses.
//
// MICROSERVICES PATTERN: Ports (Hexagonal Architecture)
// =====================================================
// Ports are interfaces that define how the outside world interacts
// with our application. There are two types:
//
// 1. PRIMARY PORTS (Driving) - How external actors call our app
//    Example: HTTP handlers, gRPC servers, CLI commands
//    These CALL our application layer.
//
// 2. SECONDARY PORTS (Driven) - How our app calls external systems
//    Example: Database repositories, external API clients
//    These are CALLED BY our application layer.
//
// This file contains SECONDARY PORTS - interfaces that our
// application layer uses to interact with external systems.
//
// WHY INTERFACES?
// ===============
// - Dependency Inversion: High-level modules don't depend on low-level modules
// - Testability: We can mock these interfaces in tests
// - Flexibility: We can swap implementations (e.g., PostgreSQL to MongoDB)
package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/auth/internal/domain"
)

// UserRepository defines the contract for user persistence operations.
//
// PATTERN: Repository
// The repository pattern mediates between the domain and data mapping layers.
// It provides a collection-like interface for accessing domain objects.
//
// All methods take a context.Context as first parameter. This is idiomatic Go
// and allows for:
// - Request cancellation
// - Timeouts
// - Passing request-scoped values (like trace IDs)
type UserRepository interface {
	// Create stores a new user in the database.
	// Returns ErrUserAlreadyExists if phone number already exists.
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by their unique ID.
	// Returns ErrUserNotFound if user doesn't exist.
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByPhone retrieves a user by their phone number.
	// Returns ErrUserNotFound if user doesn't exist.
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)

	// GetByEmail retrieves a user by their email address.
	// Returns ErrUserNotFound if user doesn't exist.
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update saves changes to an existing user.
	// Returns ErrUserNotFound if user doesn't exist.
	Update(ctx context.Context, user *domain.User) error

	// Delete removes a user from the database.
	// This is typically a soft delete (sets status to inactive).
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByPhone checks if a user with the given phone exists.
	// This is more efficient than GetByPhone when we just need to check existence.
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

// RefreshTokenRepository defines the contract for refresh token persistence.
//
// SECURITY NOTE: We store hashed tokens, not the actual tokens.
// The actual token is sent to the client, and we hash it before storing.
type RefreshTokenRepository interface {
	// Create stores a new refresh token.
	Create(ctx context.Context, token *domain.RefreshToken) error

	// GetByTokenHash retrieves a token by its hash.
	// This is used during token refresh to validate the token.
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)

	// GetByUserID retrieves all tokens for a user.
	// Useful for showing active sessions or "logout everywhere".
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error)

	// Revoke marks a specific token as revoked.
	Revoke(ctx context.Context, id uuid.UUID) error

	// RevokeAllForUser revokes all tokens for a user.
	// Used for "logout everywhere" functionality.
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired removes tokens that have expired.
	// This should be called periodically by a cleanup job.
	DeleteExpired(ctx context.Context) error
}

// OTPRepository defines the contract for OTP persistence.
//
// OTPs are temporary and should be cleaned up after expiration.
// Consider using Redis for OTP storage in production for better
// performance and automatic expiration.
type OTPRepository interface {
	// Create stores a new OTP.
	// Any existing OTPs for the same phone should be invalidated.
	Create(ctx context.Context, otp *domain.OTP) error

	// GetLatestByPhone retrieves the most recent valid OTP for a phone.
	GetLatestByPhone(ctx context.Context, phone string) (*domain.OTP, error)

	// Update saves changes to an OTP (e.g., incrementing attempts).
	Update(ctx context.Context, otp *domain.OTP) error

	// DeleteByPhone removes all OTPs for a phone number.
	// Called after successful verification.
	DeleteByPhone(ctx context.Context, phone string) error

	// DeleteExpired removes expired OTPs.
	// This should be called periodically by a cleanup job.
	DeleteExpired(ctx context.Context) error
}

// UnitOfWork provides transaction management across repositories.
//
// PATTERN: Unit of Work
// The Unit of Work pattern maintains a list of objects affected by a
// business transaction and coordinates the writing out of changes.
//
// This is useful when you need to perform multiple repository operations
// that should either all succeed or all fail (atomicity).
//
// Example usage:
//
//	err := uow.Execute(ctx, func(tx Transaction) error {
//	    // All operations here are in the same transaction
//	    if err := tx.Users().Create(ctx, user); err != nil {
//	        return err // Transaction will be rolled back
//	    }
//	    if err := tx.Tokens().Create(ctx, token); err != nil {
//	        return err // Transaction will be rolled back
//	    }
//	    return nil // Transaction will be committed
//	})
type UnitOfWork interface {
	// Execute runs the given function within a transaction.
	// If the function returns an error, the transaction is rolled back.
	// If the function returns nil, the transaction is committed.
	Execute(ctx context.Context, fn func(tx Transaction) error) error
}

// Transaction provides access to repositories within a transaction.
type Transaction interface {
	Users() UserRepository
	Tokens() RefreshTokenRepository
	OTPs() OTPRepository
}
