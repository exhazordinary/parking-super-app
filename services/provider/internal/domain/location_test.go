package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewLocation(t *testing.T) {
	providerID := uuid.New()

	location := NewLocation(
		providerID,
		"KLCC Parking",
		"Suria KLCC",
		"Kuala Lumpur",
		"Wilayah Persekutuan",
		3.1579,
		101.7116,
	)

	if location.ID == uuid.Nil {
		t.Error("expected location ID to be set")
	}
	if location.ProviderID != providerID {
		t.Errorf("expected provider ID %v, got %v", providerID, location.ProviderID)
	}
	if location.Name != "KLCC Parking" {
		t.Errorf("expected name KLCC Parking, got %s", location.Name)
	}
	if location.City != "Kuala Lumpur" {
		t.Errorf("expected city Kuala Lumpur, got %s", location.City)
	}
	if !location.IsActive {
		t.Error("new location should be active")
	}
	if location.Pricing.Currency != "MYR" {
		t.Errorf("expected currency MYR, got %s", location.Pricing.Currency)
	}
	if location.Pricing.GracePeriodMin != 15 {
		t.Errorf("expected grace period 15, got %d", location.Pricing.GracePeriodMin)
	}
}

func TestLocation_SetPricing(t *testing.T) {
	location := NewLocation(uuid.New(), "Test", "Address", "City", "State", 0, 0)

	location.SetPricing(5.00, 50.00)

	if location.Pricing.HourlyRate != 5.00 {
		t.Errorf("expected hourly rate 5.00, got %f", location.Pricing.HourlyRate)
	}
	if location.Pricing.DailyMax != 50.00 {
		t.Errorf("expected daily max 50.00, got %f", location.Pricing.DailyMax)
	}
}

func TestLocation_AddAmenity(t *testing.T) {
	location := NewLocation(uuid.New(), "Test", "Address", "City", "State", 0, 0)

	if len(location.Amenities) != 0 {
		t.Error("new location should have no amenities")
	}

	location.AddAmenity("EV Charging")
	location.AddAmenity("Covered Parking")

	if len(location.Amenities) != 2 {
		t.Errorf("expected 2 amenities, got %d", len(location.Amenities))
	}
	if location.Amenities[0] != "EV Charging" {
		t.Errorf("expected first amenity EV Charging, got %s", location.Amenities[0])
	}
}

func TestLocation_Deactivate(t *testing.T) {
	location := NewLocation(uuid.New(), "Test", "Address", "City", "State", 0, 0)

	if !location.IsActive {
		t.Error("new location should be active")
	}

	location.Deactivate()

	if location.IsActive {
		t.Error("location should be inactive after deactivation")
	}
}
