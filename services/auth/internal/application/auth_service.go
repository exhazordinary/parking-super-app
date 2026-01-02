// Package application contains the application services (use cases).
//
// MICROSERVICES PATTERN: Application Layer
// ========================================
// The application layer orchestrates the flow of data between:
// - The outside world (HTTP handlers, gRPC servers)
// - The domain layer (business entities and rules)
// - External systems (databases, message queues, external APIs)
//
// This layer contains the USE CASES of our application.
// Each method represents a specific action a user can perform.
//
// IMPORTANT RULES:
// - Application services should be stateless
// - They coordinate work but don't contain business logic (that's in domain)
// - They depend on ports (interfaces), not concrete implementations
package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/auth/internal/domain"
	"github.com/parking-super-app/services/auth/internal/ports"
)

// AuthService provides authentication and user management functionality.
//
// PATTERN: Application Service
// This service orchestrates the authentication flow without knowing
// about HTTP, databases, or any infrastructure concerns.
type AuthService struct {
	// Dependencies (injected via constructor)
	users          ports.UserRepository
	tokens         ports.RefreshTokenRepository
	otps           ports.OTPRepository
	passwordHasher ports.PasswordHasher
	tokenService   ports.TokenService
	smsService     ports.SMSService
	otpGenerator   ports.OTPGenerator
	events         ports.EventPublisher
	logger         ports.Logger
}

// NewAuthService creates a new AuthService with all dependencies.
//
// PATTERN: Dependency Injection
// All dependencies are passed in via the constructor. This makes the
// service testable (we can pass mocks) and flexible (we can swap
// implementations).
func NewAuthService(
	users ports.UserRepository,
	tokens ports.RefreshTokenRepository,
	otps ports.OTPRepository,
	passwordHasher ports.PasswordHasher,
	tokenService ports.TokenService,
	smsService ports.SMSService,
	otpGenerator ports.OTPGenerator,
	events ports.EventPublisher,
	logger ports.Logger,
) *AuthService {
	return &AuthService{
		users:          users,
		tokens:         tokens,
		otps:           otps,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
		smsService:     smsService,
		otpGenerator:   otpGenerator,
		events:         events,
		logger:         logger,
	}
}

// ---- Request/Response DTOs ----
// DTOs (Data Transfer Objects) define the input/output of our use cases.
// They are different from domain entities because they're shaped for the
// specific use case, not the business domain.

// RegisterRequest contains the data needed to register a new user.
type RegisterRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

// RegisterResponse is returned after successful registration.
type RegisterResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

// LoginRequest contains credentials for login.
type LoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse contains tokens returned after successful login.
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"` // Seconds until access token expires
	UserID       uuid.UUID `json:"user_id"`
}

// RefreshTokenRequest contains the refresh token to exchange.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// VerifyOTPRequest contains the OTP code to verify.
type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required,len=6"`
}

// RequestOTPRequest contains the phone number to send OTP to.
type RequestOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
}

// UserProfile represents the user's public profile information.
type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
	FullName  string    `json:"full_name"`
	Status    string    `json:"status"`
}

// ---- Use Cases ----

// Register creates a new user account.
//
// Flow:
// 1. Validate input
// 2. Check if phone already exists
// 3. Hash password
// 4. Create user (with pending status)
// 5. Generate and send OTP for verification
// 6. Publish user.registered event
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	s.logger.Info("registering new user", ports.String("phone", req.Phone))

	// Check if password meets requirements
	if err := domain.ValidatePassword(req.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Check if user already exists
	exists, err := s.users.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		s.logger.Error("failed to check user existence", ports.Err(err))
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash the password
	passwordHash, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", ports.Err(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user entity
	user, err := domain.NewUser(req.Phone, req.Email, req.FullName, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	// Persist the user
	if err := s.users.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user", ports.Err(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate and send OTP
	otp := domain.NewOTP(req.Phone, s.otpGenerator.Generate())
	if err := s.otps.Create(ctx, otp); err != nil {
		s.logger.Error("failed to create OTP", ports.Err(err))
		// Continue - user is created, they can request OTP again
	} else {
		// Send OTP via SMS (don't fail registration if SMS fails)
		go func() {
			if err := s.smsService.SendOTP(context.Background(), req.Phone, otp.Code); err != nil {
				s.logger.Error("failed to send OTP", ports.Err(err), ports.String("phone", req.Phone))
			}
		}()
	}

	// Publish event (async)
	go func() {
		event := ports.Event{
			Type: ports.EventUserRegistered,
			Payload: map[string]interface{}{
				"user_id": user.ID.String(),
				"phone":   user.Phone,
			},
		}
		if err := s.events.Publish(context.Background(), event); err != nil {
			s.logger.Error("failed to publish event", ports.Err(err))
		}
	}()

	s.logger.Info("user registered successfully", ports.String("user_id", user.ID.String()))

	return &RegisterResponse{
		UserID:  user.ID,
		Message: "Registration successful. Please verify your phone number with the OTP sent.",
	}, nil
}

// Login authenticates a user and returns tokens.
//
// Flow:
// 1. Find user by phone
// 2. Verify password
// 3. Check if user can login (status check)
// 4. Generate access token and refresh token
// 5. Store refresh token hash
// 6. Publish user.logged_in event
func (s *AuthService) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*LoginResponse, error) {
	s.logger.Info("user attempting login", ports.String("phone", req.Phone))

	// Find user
	user, err := s.users.GetByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials // Don't reveal if user exists
		}
		s.logger.Error("failed to get user", ports.Err(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := s.passwordHasher.Compare(req.Password, user.PasswordHash); err != nil {
		s.logger.Warn("invalid password attempt", ports.String("phone", req.Phone))
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user can login
	if !user.CanLogin() {
		return nil, domain.ErrUserInactive
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Phone)
	if err != nil {
		s.logger.Error("failed to generate access token", ports.Err(err))
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		s.logger.Error("failed to generate refresh token", ports.Err(err))
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash and store refresh token
	tokenHash := s.tokenService.HashRefreshToken(refreshToken)
	rt := domain.NewRefreshToken(user.ID, tokenHash, userAgent, ipAddress)
	if err := s.tokens.Create(ctx, rt); err != nil {
		s.logger.Error("failed to store refresh token", ports.Err(err))
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Publish event (async)
	go func() {
		event := ports.Event{
			Type: ports.EventUserLoggedIn,
			Payload: map[string]interface{}{
				"user_id":    user.ID.String(),
				"ip_address": ipAddress,
			},
		}
		if err := s.events.Publish(context.Background(), event); err != nil {
			s.logger.Error("failed to publish event", ports.Err(err))
		}
	}()

	s.logger.Info("user logged in successfully", ports.String("user_id", user.ID.String()))

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
		UserID:       user.ID,
	}, nil
}

// RefreshToken exchanges a refresh token for new access and refresh tokens.
//
// SECURITY: Token Rotation
// When a refresh token is used, we:
// 1. Revoke the old refresh token
// 2. Issue a new refresh token
// This limits the window of opportunity if a token is stolen.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, userAgent, ipAddress string) (*LoginResponse, error) {
	// Hash the provided token to look it up
	tokenHash := s.tokenService.HashRefreshToken(refreshToken)

	// Find the token
	storedToken, err := s.tokens.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return nil, domain.ErrInvalidToken
		}
		s.logger.Error("failed to get refresh token", ports.Err(err))
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Validate the token
	if err := storedToken.Validate(); err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.users.GetByID(ctx, storedToken.UserID)
	if err != nil {
		s.logger.Error("failed to get user for token refresh", ports.Err(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is still active
	if !user.CanLogin() {
		return nil, domain.ErrUserInactive
	}

	// Revoke the old token (token rotation)
	if err := s.tokens.Revoke(ctx, storedToken.ID); err != nil {
		s.logger.Error("failed to revoke old token", ports.Err(err))
		// Continue anyway - don't block the user
	}

	// Generate new access token
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store new refresh token
	newTokenHash := s.tokenService.HashRefreshToken(newRefreshToken)
	newRT := domain.NewRefreshToken(user.ID, newTokenHash, userAgent, ipAddress)
	if err := s.tokens.Create(ctx, newRT); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    900,
		UserID:       user.ID,
	}, nil
}

// Logout revokes the user's refresh token.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := s.tokenService.HashRefreshToken(refreshToken)

	storedToken, err := s.tokens.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return nil // Token doesn't exist, consider it logged out
		}
		return fmt.Errorf("failed to get token: %w", err)
	}

	if err := s.tokens.Revoke(ctx, storedToken.ID); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// LogoutAllDevices revokes all refresh tokens for a user.
func (s *AuthService) LogoutAllDevices(ctx context.Context, userID uuid.UUID) error {
	if err := s.tokens.RevokeAllForUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke all tokens: %w", err)
	}
	return nil
}

// RequestOTP generates and sends a new OTP to the user's phone.
func (s *AuthService) RequestOTP(ctx context.Context, req RequestOTPRequest) error {
	// Check if user exists
	_, err := s.users.GetByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Don't reveal if user exists - just pretend we sent OTP
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Generate OTP
	otp := domain.NewOTP(req.Phone, s.otpGenerator.Generate())
	if err := s.otps.Create(ctx, otp); err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	// Send OTP
	if err := s.smsService.SendOTP(ctx, req.Phone, otp.Code); err != nil {
		s.logger.Error("failed to send OTP", ports.Err(err))
		return fmt.Errorf("failed to send OTP: %w", err)
	}

	return nil
}

// VerifyOTP verifies an OTP code and activates the user if pending.
func (s *AuthService) VerifyOTP(ctx context.Context, req VerifyOTPRequest) error {
	// Get the latest OTP
	otp, err := s.otps.GetLatestByPhone(ctx, req.Phone)
	if err != nil {
		return domain.ErrInvalidToken
	}

	// Verify the code
	if !otp.Verify(req.Code) {
		// Update the OTP to record the failed attempt
		if err := s.otps.Update(ctx, otp); err != nil {
			s.logger.Error("failed to update OTP attempts", ports.Err(err))
		}
		return domain.ErrInvalidToken
	}

	// OTP verified - activate user if pending
	user, err := s.users.GetByPhone(ctx, req.Phone)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Status == domain.UserStatusPending {
		user.Activate()
		if err := s.users.Update(ctx, user); err != nil {
			return fmt.Errorf("failed to activate user: %w", err)
		}
	}

	// Clean up OTPs for this phone
	if err := s.otps.DeleteByPhone(ctx, req.Phone); err != nil {
		s.logger.Error("failed to delete OTPs", ports.Err(err))
	}

	return nil
}

// GetProfile returns the user's profile.
func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		ID:       user.ID,
		Phone:    user.Phone,
		Email:    user.Email,
		FullName: user.FullName,
		Status:   string(user.Status),
	}, nil
}
