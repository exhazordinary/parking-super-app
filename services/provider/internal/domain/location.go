package domain

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a parking location operated by a provider
type Location struct {
	ID          uuid.UUID       `json:"id"`
	ProviderID  uuid.UUID       `json:"provider_id"`
	Name        string          `json:"name"`
	Address     string          `json:"address"`
	City        string          `json:"city"`
	State       string          `json:"state"`
	PostalCode  string          `json:"postal_code"`
	Latitude    float64         `json:"latitude"`
	Longitude   float64         `json:"longitude"`
	TotalSpaces int             `json:"total_spaces"`
	Amenities   []string        `json:"amenities"`
	Pricing     LocationPricing `json:"pricing"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// LocationPricing defines the pricing structure for a location
type LocationPricing struct {
	HourlyRate     float64 `json:"hourly_rate"`
	DailyMax       float64 `json:"daily_max"`
	Currency       string  `json:"currency"`
	GracePeriodMin int     `json:"grace_period_min"`
}

// NewLocation creates a new parking location
func NewLocation(providerID uuid.UUID, name, address, city, state string, lat, lng float64) *Location {
	now := time.Now().UTC()
	return &Location{
		ID:         uuid.New(),
		ProviderID: providerID,
		Name:       name,
		Address:    address,
		City:       city,
		State:      state,
		Latitude:   lat,
		Longitude:  lng,
		Amenities:  []string{},
		Pricing: LocationPricing{
			Currency:       "MYR",
			GracePeriodMin: 15,
		},
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetPricing updates the location pricing
func (l *Location) SetPricing(hourlyRate, dailyMax float64) {
	l.Pricing.HourlyRate = hourlyRate
	l.Pricing.DailyMax = dailyMax
	l.UpdatedAt = time.Now().UTC()
}

// AddAmenity adds an amenity to the location
func (l *Location) AddAmenity(amenity string) {
	l.Amenities = append(l.Amenities, amenity)
	l.UpdatedAt = time.Now().UTC()
}

// Deactivate disables the location
func (l *Location) Deactivate() {
	l.IsActive = false
	l.UpdatedAt = time.Now().UTC()
}
