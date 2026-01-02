package external

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTTokenService_GenerateAccessToken(t *testing.T) {
	service := NewJWTTokenService("test-secret-key-32-chars-long!!", 15*time.Minute)
	userID := uuid.New()
	phone := "+60123456789"

	token, err := service.GenerateAccessToken(userID, phone)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Error("token should not be empty")
	}
}

func TestJWTTokenService_ValidateAccessToken(t *testing.T) {
	service := NewJWTTokenService("test-secret-key-32-chars-long!!", 15*time.Minute)
	userID := uuid.New()
	phone := "+60123456789"

	token, _ := service.GenerateAccessToken(userID, phone)

	claims, err := service.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("claims.UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Phone != phone {
		t.Errorf("claims.Phone = %v, want %v", claims.Phone, phone)
	}
}

func TestJWTTokenService_ValidateExpiredToken(t *testing.T) {
	// Create service with very short TTL
	service := NewJWTTokenService("test-secret-key-32-chars-long!!", 1*time.Millisecond)
	userID := uuid.New()

	token, _ := service.GenerateAccessToken(userID, "+60123456789")

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	_, err := service.ValidateAccessToken(token)
	if err == nil {
		t.Error("expired token should fail validation")
	}
}

func TestJWTTokenService_ValidateInvalidToken(t *testing.T) {
	service := NewJWTTokenService("test-secret-key-32-chars-long!!", 15*time.Minute)

	_, err := service.ValidateAccessToken("invalid.token.here")
	if err == nil {
		t.Error("invalid token should fail validation")
	}
}

func TestJWTTokenService_ValidateTokenWithWrongSecret(t *testing.T) {
	service1 := NewJWTTokenService("secret-key-one-32-chars-long!!!", 15*time.Minute)
	service2 := NewJWTTokenService("secret-key-two-32-chars-long!!!", 15*time.Minute)

	token, _ := service1.GenerateAccessToken(uuid.New(), "+60123456789")

	_, err := service2.ValidateAccessToken(token)
	if err == nil {
		t.Error("token signed with different secret should fail validation")
	}
}

func TestJWTTokenService_GenerateRefreshToken(t *testing.T) {
	service := NewJWTTokenService("test-secret-key-32-chars-long!!", 15*time.Minute)

	token1, err := service.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	token2, _ := service.GenerateRefreshToken()

	if token1 == "" {
		t.Error("refresh token should not be empty")
	}

	if token1 == token2 {
		t.Error("refresh tokens should be unique")
	}

	// Should be 64 hex characters (32 bytes)
	if len(token1) != 64 {
		t.Errorf("refresh token length = %d, want 64", len(token1))
	}
}

func TestJWTTokenService_HashRefreshToken(t *testing.T) {
	service := NewJWTTokenService("test-secret-key-32-chars-long!!", 15*time.Minute)

	token := "test-refresh-token"
	hash1 := service.HashRefreshToken(token)
	hash2 := service.HashRefreshToken(token)

	if hash1 == "" {
		t.Error("hash should not be empty")
	}

	if hash1 != hash2 {
		t.Error("same token should produce same hash")
	}

	if hash1 == token {
		t.Error("hash should not equal original token")
	}

	// SHA-256 produces 64 hex characters
	if len(hash1) != 64 {
		t.Errorf("hash length = %d, want 64", len(hash1))
	}
}
