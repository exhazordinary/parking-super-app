package external

import (
	"testing"
)

func TestSecureOTPGenerator_Generate(t *testing.T) {
	generator := NewSecureOTPGenerator(6)

	otp1 := generator.Generate()
	otp2 := generator.Generate()

	if len(otp1) != 6 {
		t.Errorf("OTP length = %d, want 6", len(otp1))
	}

	// Check all characters are digits
	for _, c := range otp1 {
		if c < '0' || c > '9' {
			t.Errorf("OTP contains non-digit character: %c", c)
		}
	}

	// OTPs should be different (statistically)
	// This test could theoretically fail with 1 in 1,000,000 chance
	if otp1 == otp2 {
		t.Log("Warning: two consecutive OTPs were the same (unlikely but possible)")
	}
}

func TestSecureOTPGenerator_LengthVariants(t *testing.T) {
	tests := []struct {
		length   int
		expected int
	}{
		{4, 4},
		{6, 6},
		{8, 8},
		{3, 6},  // Too short, should default to 6
		{10, 6}, // Too long, should default to 6
	}

	for _, tt := range tests {
		generator := NewSecureOTPGenerator(tt.length)
		otp := generator.Generate()

		if len(otp) != tt.expected {
			t.Errorf("Generator(%d) produced OTP of length %d, want %d", tt.length, len(otp), tt.expected)
		}
	}
}

func TestSecureOTPGenerator_LeadingZeros(t *testing.T) {
	generator := NewSecureOTPGenerator(6)

	// Generate many OTPs and check that short numbers are zero-padded
	hasLeadingZero := false
	for i := 0; i < 100; i++ {
		otp := generator.Generate()
		if otp[0] == '0' {
			hasLeadingZero = true
			break
		}
	}

	// We should see at least one OTP starting with 0 in 100 tries
	// (probability of not seeing one is about 0.9^100 which is tiny)
	if !hasLeadingZero {
		t.Log("Note: No OTP with leading zero found in 100 tries (unlikely but possible)")
	}
}
