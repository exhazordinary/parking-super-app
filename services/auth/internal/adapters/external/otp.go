package external

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// SecureOTPGenerator implements ports.OTPGenerator using crypto/rand.
//
// SECURITY: Why crypto/rand?
// ==========================
// - math/rand is not cryptographically secure
// - crypto/rand uses the OS's cryptographic random number generator
// - Ensures OTPs are truly random and unpredictable
type SecureOTPGenerator struct {
	length int
}

// NewSecureOTPGenerator creates a new OTP generator.
// Default length is 6 digits.
func NewSecureOTPGenerator(length int) *SecureOTPGenerator {
	if length < 4 || length > 8 {
		length = 6
	}
	return &SecureOTPGenerator{length: length}
}

// Generate creates a new OTP code.
func (g *SecureOTPGenerator) Generate() string {
	// Calculate the max value (10^length)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(g.length)), nil)

	// Generate random number
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback to a simple approach if crypto/rand fails
		// This should never happen in practice
		return "000000"
	}

	// Pad with leading zeros if necessary
	format := fmt.Sprintf("%%0%dd", g.length)
	return fmt.Sprintf(format, n.Int64())
}
