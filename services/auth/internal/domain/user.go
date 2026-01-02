// Package domain contains the core business entities and rules.
//
// MICROSERVICES PATTERN: Domain Layer (Hexagonal Architecture)
// ============================================================
// The domain layer is the innermost layer of our hexagonal architecture.
// It contains:
// - Business entities (User, Token, etc.)
// - Business rules and validation
// - Domain errors
//
// IMPORTANT: This layer has NO external dependencies.
// It doesn't know about databases, HTTP, gRPC, or any infrastructure.
// This makes it highly testable and portable.
package domain

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Domain errors - these are business-level errors that can occur
// in our domain logic. Using errors.New() here keeps them simple
// and framework-independent.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPhone       = errors.New("invalid phone format")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrUserInactive       = errors.New("user account is inactive")
)

// UserStatus represents the possible states of a user account.
// Using a custom type with constants is more type-safe than raw strings.
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending" // Awaiting OTP verification
	UserStatusBanned   UserStatus = "banned"
)

// User represents a user in our parking super app.
//
// DESIGN DECISION: Why use a struct with exported fields?
// - Simple and idiomatic Go
// - Easy to serialize/deserialize (JSON, DB, etc.)
// - Validation is done through separate methods
//
// MICROSERVICES PATTERN: Entity
// This is an Entity - it has a unique identity (ID) that persists
// over time, even if other attributes change.
type User struct {
	ID           uuid.UUID  `json:"id"`
	Phone        string     `json:"phone"`         // Malaysian phone format: +60xxxxxxxxx
	Email        string     `json:"email"`         // Optional, can be empty
	PasswordHash string     `json:"-"`             // "-" means don't include in JSON
	FullName     string     `json:"full_name"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// NewUser creates a new User entity with validation.
//
// PATTERN: Factory Function
// Instead of letting anyone create a User{} directly, we provide
// a factory function that ensures the entity is always valid.
// This is called "protecting invariants" in DDD terms.
func NewUser(phone, email, fullName, passwordHash string) (*User, error) {
	// Validate phone number (Malaysian format)
	if !isValidMalaysianPhone(phone) {
		return nil, ErrInvalidPhone
	}

	// Validate email if provided
	if email != "" && !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	now := time.Now().UTC()

	return &User{
		ID:           uuid.New(),
		Phone:        phone,
		Email:        email,
		FullName:     fullName,
		PasswordHash: passwordHash,
		Status:       UserStatusPending, // New users start as pending (need OTP verification)
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Activate changes user status to active.
// This is typically called after OTP verification.
//
// PATTERN: Behavior on Entity
// Instead of exposing the Status field for direct modification,
// we provide methods that encapsulate business logic.
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now().UTC()
}

// Deactivate changes user status to inactive.
func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now().UTC()
}

// IsActive checks if the user can perform actions.
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// CanLogin checks if the user is allowed to login.
func (u *User) CanLogin() bool {
	return u.Status == UserStatusActive || u.Status == UserStatusPending
}

// UpdateProfile updates user's profile information.
func (u *User) UpdateProfile(fullName, email string) error {
	if email != "" && !isValidEmail(email) {
		return ErrInvalidEmail
	}

	u.FullName = fullName
	u.Email = email
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdatePassword updates the user's password hash.
// Note: Password hashing should be done in the application layer,
// not here. This method just stores the already-hashed password.
func (u *User) UpdatePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = time.Now().UTC()
}

// Validation helpers - these are pure functions with no external dependencies

// isValidMalaysianPhone validates Malaysian phone number format.
// Format: +60 followed by 9-10 digits
// Examples: +60123456789, +6011234567890
func isValidMalaysianPhone(phone string) bool {
	// Malaysian phone regex: starts with +60, followed by 9-10 digits
	pattern := `^\+60\d{9,10}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// isValidEmail performs basic email validation.
func isValidEmail(email string) bool {
	// Simple email regex - not RFC 5322 compliant but good enough
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ValidatePassword checks if a password meets our requirements.
// This is a standalone function because we might need it before
// the User entity exists (during registration).
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	// Add more rules as needed:
	// - Must contain uppercase
	// - Must contain number
	// - Must contain special character
	return nil
}
