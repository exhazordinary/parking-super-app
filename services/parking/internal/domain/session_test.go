package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestNewParkingSession(t *testing.T) {
	userID := uuid.New()
	providerID := uuid.New()
	locationID := uuid.New()

	session, err := NewParkingSession(userID, providerID, locationID, "WKL1234", "car")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.ID == uuid.Nil {
		t.Error("expected session ID to be set")
	}
	if session.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, session.UserID)
	}
	if session.VehiclePlate != "WKL1234" {
		t.Errorf("expected plate WKL1234, got %s", session.VehiclePlate)
	}
	if session.Status != SessionStatusActive {
		t.Errorf("expected status active, got %s", session.Status)
	}
	if !session.IsActive() {
		t.Error("new session should be active")
	}
}

func TestNewParkingSession_InvalidPlate(t *testing.T) {
	_, err := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "X", "car")
	if err != ErrInvalidVehiclePlate {
		t.Errorf("expected ErrInvalidVehiclePlate, got %v", err)
	}
}

func TestParkingSession_End(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")
	amount := decimal.NewFromFloat(10.00)

	err := session.End(amount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.Status != SessionStatusCompleted {
		t.Errorf("expected status completed, got %s", session.Status)
	}
	if session.ExitTime == nil {
		t.Error("expected exit time to be set")
	}
	if !session.Amount.Equal(amount) {
		t.Errorf("expected amount %s, got %s", amount.String(), session.Amount.String())
	}
	if !session.IsCompleted() {
		t.Error("session should be completed")
	}
}

func TestParkingSession_EndTwice(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")
	session.End(decimal.NewFromFloat(10.00))

	err := session.End(decimal.NewFromFloat(20.00))
	if err != ErrSessionAlreadyEnded {
		t.Errorf("expected ErrSessionAlreadyEnded, got %v", err)
	}
}

func TestParkingSession_Cancel(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")

	err := session.Cancel()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.Status != SessionStatusCancelled {
		t.Errorf("expected status cancelled, got %s", session.Status)
	}
	if session.ExitTime == nil {
		t.Error("expected exit time to be set")
	}
}

func TestParkingSession_CancelEnded(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")
	session.End(decimal.NewFromFloat(10.00))

	err := session.Cancel()
	if err != ErrSessionAlreadyEnded {
		t.Errorf("expected ErrSessionAlreadyEnded, got %v", err)
	}
}

func TestParkingSession_MarkPaid(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")
	paymentID := uuid.New()

	session.MarkPaid(paymentID)

	if session.PaymentID == nil {
		t.Error("expected payment ID to be set")
	}
	if *session.PaymentID != paymentID {
		t.Errorf("expected payment ID %v, got %v", paymentID, *session.PaymentID)
	}
}

func TestParkingSession_CalculateDuration(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")
	session.EntryTime = time.Now().Add(-30 * time.Minute)

	duration := session.CalculateDuration()
	if duration < 29 || duration > 31 {
		t.Errorf("expected duration around 30 minutes, got %d", duration)
	}
}

func TestParkingSession_CalculateAmount(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")
	session.EntryTime = time.Now().Add(-90 * time.Minute) // 1.5 hours

	hourlyRate := decimal.NewFromFloat(5.00)
	dailyMax := decimal.NewFromFloat(50.00)

	amount := session.CalculateAmount(hourlyRate, dailyMax)

	// Should be 2 hours (rounded up) * 5 = 10
	expected := decimal.NewFromFloat(10.00)
	if !amount.Equal(expected) {
		t.Errorf("expected amount %s, got %s", expected.String(), amount.String())
	}
}

func TestParkingSession_CalculateAmount_DailyCap(t *testing.T) {
	session, _ := NewParkingSession(uuid.New(), uuid.New(), uuid.New(), "ABC123", "car")
	session.EntryTime = time.Now().Add(-12 * time.Hour)

	hourlyRate := decimal.NewFromFloat(10.00)
	dailyMax := decimal.NewFromFloat(50.00)

	amount := session.CalculateAmount(hourlyRate, dailyMax)

	// 12 hours * 10 = 120, but capped at 50
	if !amount.Equal(dailyMax) {
		t.Errorf("expected amount capped at %s, got %s", dailyMax.String(), amount.String())
	}
}

func TestIsValidPlate(t *testing.T) {
	tests := []struct {
		plate string
		valid bool
	}{
		{"WKL1234", true},
		{"ABC123", true},
		{"JJ", true},
		{"X", false},
		{"", false},
		{"ABCDEFGHIJK", false},
	}

	for _, tt := range tests {
		t.Run(tt.plate, func(t *testing.T) {
			result := isValidPlate(tt.plate)
			if result != tt.valid {
				t.Errorf("isValidPlate(%s) = %v, want %v", tt.plate, result, tt.valid)
			}
		})
	}
}
