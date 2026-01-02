package external

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/parking-super-app/services/auth/internal/ports"
)

// JWTTokenService implements ports.TokenService using JWT.
//
// JWT STRUCTURE:
// =============
// A JWT has three parts: Header.Payload.Signature
// - Header: Algorithm and token type
// - Payload: Claims (user data)
// - Signature: Ensures the token hasn't been tampered with
//
// SECURITY CONSIDERATIONS:
// - Use strong secret key (at least 32 bytes)
// - Keep access tokens short-lived (15 min)
// - Store refresh tokens in database, not JWT
// - Never store sensitive data in JWT payload (it's only base64 encoded, not encrypted)
type JWTTokenService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
}

// NewJWTTokenService creates a new JWT token service.
//
// Parameters:
// - secretKey: Should be at least 32 bytes of random data
// - accessTokenTTL: How long access tokens are valid (recommend 15 min)
func NewJWTTokenService(secretKey string, accessTokenTTL time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secretKey:      []byte(secretKey),
		accessTokenTTL: accessTokenTTL,
	}
}

// jwtClaims represents the custom claims in our JWT.
type jwtClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"uid"`
	Phone  string    `json:"phone"`
}

// GenerateAccessToken creates a new JWT access token.
func (s *JWTTokenService) GenerateAccessToken(userID uuid.UUID, phone string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.accessTokenTTL)

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Subject is a standard claim for the user identifier
			Subject: userID.String(),
			// IssuedAt is when the token was created
			IssuedAt: jwt.NewNumericDate(now),
			// ExpiresAt is when the token expires
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			// Issuer identifies who created the token
			Issuer: "parking-super-app-auth",
		},
		UserID: userID,
		Phone:  phone,
	}

	// Create token with HS256 algorithm
	// HS256 = HMAC with SHA-256 (symmetric key)
	// For distributed microservices, consider RS256 (asymmetric) instead
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateAccessToken validates a JWT and returns the claims.
func (s *JWTTokenService) ValidateAccessToken(tokenString string) (*ports.AccessTokenClaims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &ports.AccessTokenClaims{
		UserID:    claims.UserID,
		Phone:     claims.Phone,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}, nil
}

// GenerateRefreshToken creates a cryptographically secure random token.
//
// WHY NOT JWT FOR REFRESH TOKENS?
// ==============================
// - Refresh tokens are stored in the database anyway
// - We need to be able to revoke them
// - Simple random strings are sufficient and harder to forge
// - Less overhead than JWT parsing
func (s *JWTTokenService) GenerateRefreshToken() (string, error) {
	// Generate 32 bytes (256 bits) of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert to hex string (64 characters)
	return hex.EncodeToString(bytes), nil
}

// HashRefreshToken creates a SHA-256 hash of a refresh token.
//
// WHY HASH REFRESH TOKENS?
// ========================
// If our database is compromised, attackers can't use the hashes
// to create valid refresh tokens. We only store the hash.
func (s *JWTTokenService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
