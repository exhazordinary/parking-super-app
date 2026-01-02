// Package external provides implementations for external service interfaces.
package external

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher implements ports.PasswordHasher using bcrypt.
//
// SECURITY: Why bcrypt?
// =====================
// - Designed specifically for password hashing
// - Has built-in salting
// - Configurable work factor (cost)
// - Resistant to rainbow table attacks
// - Slow by design (makes brute force attacks expensive)
//
// Alternative: argon2 is the winner of the Password Hashing Competition
// and is considered more modern, but bcrypt is still very secure.
type BcryptPasswordHasher struct {
	// cost is the bcrypt cost factor (default: 12)
	// Higher cost = more secure but slower
	// Each increment doubles the computation time
	cost int
}

// NewBcryptPasswordHasher creates a new password hasher.
// Cost of 12 is recommended for most applications.
// Cost of 14+ for high-security applications.
func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordHasher{cost: cost}
}

// Hash generates a bcrypt hash of the password.
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// Compare checks if a password matches a hash.
// Returns nil if they match, error otherwise.
func (h *BcryptPasswordHasher) Compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
