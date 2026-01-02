package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewRefreshToken(t *testing.T) {
	userID := uuid.New()
	tokenHash := "somehash"
	userAgent := "Mozilla/5.0"
	ipAddress := "192.168.1.1"

	token := NewRefreshToken(userID, tokenHash, userAgent, ipAddress)

	if token.UserID != userID {
		t.Errorf("UserID = %v, want %v", token.UserID, userID)
	}
	if token.TokenHash != tokenHash {
		t.Errorf("TokenHash = %v, want %v", token.TokenHash, tokenHash)
	}
	if token.Revoked {
		t.Error("new token should not be revoked")
	}
	if token.ExpiresAt.Before(time.Now()) {
		t.Error("token should not be expired immediately")
	}
}

func TestRefreshToken_IsValid(t *testing.T) {
	userID := uuid.New()

	t.Run("valid token", func(t *testing.T) {
		token := NewRefreshToken(userID, "hash", "", "")
		if !token.IsValid() {
			t.Error("new token should be valid")
		}
	})

	t.Run("revoked token", func(t *testing.T) {
		token := NewRefreshToken(userID, "hash", "", "")
		token.Revoke()
		if token.IsValid() {
			t.Error("revoked token should not be valid")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token := NewRefreshToken(userID, "hash", "", "")
		token.ExpiresAt = time.Now().Add(-time.Hour)
		if token.IsValid() {
			t.Error("expired token should not be valid")
		}
	})
}

func TestRefreshToken_Validate(t *testing.T) {
	userID := uuid.New()

	t.Run("valid token returns nil", func(t *testing.T) {
		token := NewRefreshToken(userID, "hash", "", "")
		if err := token.Validate(); err != nil {
			t.Errorf("valid token should return nil, got %v", err)
		}
	})

	t.Run("revoked token returns error", func(t *testing.T) {
		token := NewRefreshToken(userID, "hash", "", "")
		token.Revoke()
		if err := token.Validate(); err != ErrTokenRevoked {
			t.Errorf("revoked token should return ErrTokenRevoked, got %v", err)
		}
	})

	t.Run("expired token returns error", func(t *testing.T) {
		token := NewRefreshToken(userID, "hash", "", "")
		token.ExpiresAt = time.Now().Add(-time.Hour)
		if err := token.Validate(); err != ErrTokenExpired {
			t.Errorf("expired token should return ErrTokenExpired, got %v", err)
		}
	})
}

func TestNewOTP(t *testing.T) {
	phone := "+60123456789"
	code := "123456"

	otp := NewOTP(phone, code)

	if otp.Phone != phone {
		t.Errorf("Phone = %v, want %v", otp.Phone, phone)
	}
	if otp.Code != code {
		t.Errorf("Code = %v, want %v", otp.Code, code)
	}
	if otp.Verified {
		t.Error("new OTP should not be verified")
	}
	if otp.Attempts != 0 {
		t.Error("new OTP should have 0 attempts")
	}
}

func TestOTP_IsValid(t *testing.T) {
	t.Run("new OTP is valid", func(t *testing.T) {
		otp := NewOTP("+60123456789", "123456")
		if !otp.IsValid() {
			t.Error("new OTP should be valid")
		}
	})

	t.Run("verified OTP is invalid", func(t *testing.T) {
		otp := NewOTP("+60123456789", "123456")
		otp.Verified = true
		if otp.IsValid() {
			t.Error("verified OTP should be invalid")
		}
	})

	t.Run("expired OTP is invalid", func(t *testing.T) {
		otp := NewOTP("+60123456789", "123456")
		otp.ExpiresAt = time.Now().Add(-time.Minute)
		if otp.IsValid() {
			t.Error("expired OTP should be invalid")
		}
	})

	t.Run("max attempts reached is invalid", func(t *testing.T) {
		otp := NewOTP("+60123456789", "123456")
		otp.Attempts = MaxOTPAttempts
		if otp.IsValid() {
			t.Error("OTP with max attempts should be invalid")
		}
	})
}

func TestOTP_Verify(t *testing.T) {
	t.Run("correct code verifies", func(t *testing.T) {
		otp := NewOTP("+60123456789", "123456")
		if !otp.Verify("123456") {
			t.Error("correct code should verify")
		}
		if !otp.Verified {
			t.Error("OTP should be marked as verified")
		}
	})

	t.Run("wrong code increments attempts", func(t *testing.T) {
		otp := NewOTP("+60123456789", "123456")
		if otp.Verify("000000") {
			t.Error("wrong code should not verify")
		}
		if otp.Attempts != 1 {
			t.Errorf("attempts = %d, want 1", otp.Attempts)
		}
	})

	t.Run("cannot verify after max attempts", func(t *testing.T) {
		otp := NewOTP("+60123456789", "123456")
		otp.Attempts = MaxOTPAttempts
		if otp.Verify("123456") {
			t.Error("should not verify after max attempts")
		}
	})
}
