// Package postgres provides PostgreSQL implementations of repository interfaces.
//
// MICROSERVICES PATTERN: Adapter Layer (Hexagonal Architecture)
// =============================================================
// Adapters are concrete implementations of ports (interfaces).
// This package implements the UserRepository interface using PostgreSQL.
//
// WHY SEPARATE PACKAGE?
// - Clear separation of concerns
// - Easy to swap implementations (e.g., postgres to mysql)
// - Can be tested independently with integration tests
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parking-super-app/services/auth/internal/domain"
)

// UserRepository implements ports.UserRepository using PostgreSQL.
//
// PATTERN: Repository Implementation
// This struct wraps a database connection pool and provides methods
// that translate between domain objects and database rows.
type UserRepository struct {
	// db is a connection pool, not a single connection.
	// This allows concurrent database operations.
	// pgxpool is preferred over database/sql for PostgreSQL because:
	// - Native PostgreSQL types support
	// - Better performance
	// - Connection pooling built-in
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database.
//
// LEARNING: SQL Query Explanation
// ===============================
// - We use $1, $2, etc. for parameterized queries (prevents SQL injection)
// - RETURNING clause returns the inserted values (useful for auto-generated fields)
// - ON CONFLICT DO NOTHING could be used to handle duplicates gracefully
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, phone, email, password_hash, full_name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Phone,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation (duplicate phone)
		// PostgreSQL error code 23505 = unique_violation
		if isUniqueViolation(err) {
			return domain.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID.
//
// LEARNING: Scanning Rows
// =======================
// pgx's Scan method maps database columns to struct fields in order.
// Make sure the SELECT columns match the Scan arguments exactly.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, phone, email, password_hash, full_name, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetByPhone retrieves a user by their phone number.
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query := `
		SELECT id, phone, email, password_hash, full_name, status, created_at, updated_at
		FROM users
		WHERE phone = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, phone, email, password_hash, full_name, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// Update saves changes to an existing user.
//
// LEARNING: Optimistic Locking
// ===========================
// In production, you might want to add a version column for optimistic locking:
// UPDATE users SET ... WHERE id = $1 AND version = $2
// This prevents lost updates when two requests modify the same user.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET phone = $2, email = $3, password_hash = $4, full_name = $5, status = $6, updated_at = $7
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		user.ID,
		user.Phone,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Status,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Check if any row was actually updated
	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete removes a user. This is a soft delete (sets status to inactive).
//
// BEST PRACTICE: Soft Delete
// ==========================
// We don't actually delete the row - we set status to inactive.
// Benefits:
// - Maintain data integrity (foreign keys)
// - Audit trail
// - Easy to restore if needed
// - Avoid orphaned records
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id, domain.UserStatusInactive)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// ExistsByPhone checks if a user with the given phone exists.
//
// PERFORMANCE: EXISTS is more efficient than SELECT *
// because it returns as soon as it finds one match.
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, phone).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// isUniqueViolation checks if the error is a PostgreSQL unique constraint violation.
// PostgreSQL error code 23505 = unique_violation
func isUniqueViolation(err error) bool {
	// Check if it's a pgx error with code 23505
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}

// ---- Helper for database/sql compatibility (optional) ----
// If you need to use standard database/sql instead of pgx:

// UserRepositorySQL implements ports.UserRepository using database/sql.
// Use this if you prefer the standard library or need driver flexibility.
type UserRepositorySQL struct {
	db *sql.DB
}

// NewUserRepositorySQL creates a new SQL-based UserRepository.
func NewUserRepositorySQL(db *sql.DB) *UserRepositorySQL {
	return &UserRepositorySQL{db: db}
}

// GetByID retrieves a user by ID using database/sql.
// Implementation similar to above, but uses sql.Row instead of pgx.Row.
func (r *UserRepositorySQL) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, phone, email, password_hash, full_name, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}
