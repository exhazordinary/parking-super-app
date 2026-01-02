package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrSessionNotFound       = errors.New("parking session not found")
	ErrSessionAlreadyEnded   = errors.New("session has already ended")
	ErrSessionStillActive    = errors.New("session is still active")
	ErrInvalidVehiclePlate   = errors.New("invalid vehicle plate number")
	ErrInvalidSessionDuration = errors.New("invalid session duration")
)

// SessionStatus represents the current state of a parking session
type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
	SessionStatusFailed    SessionStatus = "failed"
)

// ParkingSession represents a single parking session from entry to exit.
// This is the core domain entity for the parking service.
type ParkingSession struct {
	ID                uuid.UUID       `json:"id"`
	UserID            uuid.UUID       `json:"user_id"`
	ProviderID        uuid.UUID       `json:"provider_id"`
	LocationID        uuid.UUID       `json:"location_id"`
	ExternalSessionID string          `json:"external_session_id"`
	VehiclePlate      string          `json:"vehicle_plate"`
	VehicleType       string          `json:"vehicle_type"`
	EntryTime         time.Time       `json:"entry_time"`
	ExitTime          *time.Time      `json:"exit_time,omitempty"`
	Duration          int             `json:"duration_minutes"`
	Amount            decimal.Decimal `json:"amount"`
	Currency          string          `json:"currency"`
	Status            SessionStatus   `json:"status"`
	PaymentID         *uuid.UUID      `json:"payment_id,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// NewParkingSession creates a new active parking session
func NewParkingSession(
	userID, providerID, locationID uuid.UUID,
	vehiclePlate, vehicleType string,
) (*ParkingSession, error) {
	if !isValidPlate(vehiclePlate) {
		return nil, ErrInvalidVehiclePlate
	}

	now := time.Now().UTC()
	return &ParkingSession{
		ID:           uuid.New(),
		UserID:       userID,
		ProviderID:   providerID,
		LocationID:   locationID,
		VehiclePlate: vehiclePlate,
		VehicleType:  vehicleType,
		EntryTime:    now,
		Amount:       decimal.Zero,
		Currency:     "MYR",
		Status:       SessionStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsActive returns true if the session is still ongoing
func (s *ParkingSession) IsActive() bool {
	return s.Status == SessionStatusActive
}

// IsCompleted returns true if the session has been completed
func (s *ParkingSession) IsCompleted() bool {
	return s.Status == SessionStatusCompleted
}

// SetExternalSessionID sets the session ID from the provider's system
func (s *ParkingSession) SetExternalSessionID(externalID string) {
	s.ExternalSessionID = externalID
	s.UpdatedAt = time.Now().UTC()
}

// End completes the parking session with the final amount
func (s *ParkingSession) End(amount decimal.Decimal) error {
	if !s.IsActive() {
		return ErrSessionAlreadyEnded
	}

	now := time.Now().UTC()
	s.ExitTime = &now
	s.Duration = int(now.Sub(s.EntryTime).Minutes())
	s.Amount = amount
	s.Status = SessionStatusCompleted
	s.UpdatedAt = now

	return nil
}

// Cancel cancels an active session
func (s *ParkingSession) Cancel() error {
	if !s.IsActive() {
		return ErrSessionAlreadyEnded
	}

	now := time.Now().UTC()
	s.ExitTime = &now
	s.Status = SessionStatusCancelled
	s.UpdatedAt = now

	return nil
}

// MarkPaid records the payment for this session
func (s *ParkingSession) MarkPaid(paymentID uuid.UUID) {
	s.PaymentID = &paymentID
	s.UpdatedAt = time.Now().UTC()
}

// CalculateDuration returns the duration of the session in minutes
func (s *ParkingSession) CalculateDuration() int {
	endTime := time.Now().UTC()
	if s.ExitTime != nil {
		endTime = *s.ExitTime
	}
	return int(endTime.Sub(s.EntryTime).Minutes())
}

// CalculateAmount calculates the parking fee based on hourly rate
func (s *ParkingSession) CalculateAmount(hourlyRate, dailyMax decimal.Decimal) decimal.Decimal {
	duration := s.CalculateDuration()
	hours := decimal.NewFromInt(int64(duration)).Div(decimal.NewFromInt(60))

	// Round up to nearest hour for billing
	if duration%60 > 0 {
		hours = hours.Ceil()
	}

	amount := hours.Mul(hourlyRate)

	// Cap at daily maximum
	if amount.GreaterThan(dailyMax) && dailyMax.GreaterThan(decimal.Zero) {
		amount = dailyMax
	}

	return amount.Round(2)
}

// isValidPlate validates Malaysian vehicle plate format (basic validation)
func isValidPlate(plate string) bool {
	if len(plate) < 2 || len(plate) > 10 {
		return false
	}
	return true
}
