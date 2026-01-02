package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		phone       string
		email       string
		fullName    string
		passHash    string
		wantErr     error
		description string
	}{
		{
			name:        "valid user with email",
			phone:       "+60123456789",
			email:       "test@example.com",
			fullName:    "John Doe",
			passHash:    "hashedpassword",
			wantErr:     nil,
			description: "should create user with valid Malaysian phone and email",
		},
		{
			name:        "valid user without email",
			phone:       "+60198765432",
			email:       "",
			fullName:    "Jane Doe",
			passHash:    "hashedpassword",
			wantErr:     nil,
			description: "should create user without email",
		},
		{
			name:        "invalid phone format",
			phone:       "0123456789",
			email:       "",
			fullName:    "Test User",
			passHash:    "hashedpassword",
			wantErr:     ErrInvalidPhone,
			description: "should reject phone without +60 prefix",
		},
		{
			name:        "phone too short",
			phone:       "+6012345",
			email:       "",
			fullName:    "Test User",
			passHash:    "hashedpassword",
			wantErr:     ErrInvalidPhone,
			description: "should reject phone with insufficient digits",
		},
		{
			name:        "invalid email format",
			phone:       "+60123456789",
			email:       "invalid-email",
			fullName:    "Test User",
			passHash:    "hashedpassword",
			wantErr:     ErrInvalidEmail,
			description: "should reject invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.phone, tt.email, tt.fullName, tt.passHash)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}

			if user.Phone != tt.phone {
				t.Errorf("user.Phone = %v, want %v", user.Phone, tt.phone)
			}
			if user.Email != tt.email {
				t.Errorf("user.Email = %v, want %v", user.Email, tt.email)
			}
			if user.FullName != tt.fullName {
				t.Errorf("user.FullName = %v, want %v", user.FullName, tt.fullName)
			}
			if user.Status != UserStatusPending {
				t.Errorf("user.Status = %v, want %v", user.Status, UserStatusPending)
			}
			if user.ID.String() == "" {
				t.Error("user.ID should not be empty")
			}
		})
	}
}

func TestUser_Activate(t *testing.T) {
	user, _ := NewUser("+60123456789", "", "Test", "hash")

	if user.Status != UserStatusPending {
		t.Errorf("initial status = %v, want %v", user.Status, UserStatusPending)
	}

	user.Activate()

	if user.Status != UserStatusActive {
		t.Errorf("status after Activate() = %v, want %v", user.Status, UserStatusActive)
	}
}

func TestUser_IsActive(t *testing.T) {
	user, _ := NewUser("+60123456789", "", "Test", "hash")

	if user.IsActive() {
		t.Error("new user should not be active")
	}

	user.Activate()

	if !user.IsActive() {
		t.Error("activated user should be active")
	}

	user.Deactivate()

	if user.IsActive() {
		t.Error("deactivated user should not be active")
	}
}

func TestUser_CanLogin(t *testing.T) {
	user, _ := NewUser("+60123456789", "", "Test", "hash")

	// Pending users can login (to complete verification)
	if !user.CanLogin() {
		t.Error("pending user should be able to login")
	}

	user.Activate()
	if !user.CanLogin() {
		t.Error("active user should be able to login")
	}

	user.Status = UserStatusBanned
	if user.CanLogin() {
		t.Error("banned user should not be able to login")
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "password123", false},
		{"exactly 8 chars", "12345678", false},
		{"too short", "1234567", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
