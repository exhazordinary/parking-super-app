package domain

import (
	"time"

	"github.com/google/uuid"
)

// Vehicle represents a registered vehicle for a user
type Vehicle struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Plate     string    `json:"plate"`
	Type      string    `json:"type"`
	Make      string    `json:"make,omitempty"`
	Model     string    `json:"model,omitempty"`
	Color     string    `json:"color,omitempty"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

// VehicleType constants
const (
	VehicleTypeCar        = "car"
	VehicleTypeMotorcycle = "motorcycle"
	VehicleTypeTruck      = "truck"
)

// NewVehicle creates a new vehicle record
func NewVehicle(userID uuid.UUID, plate, vehicleType string) *Vehicle {
	return &Vehicle{
		ID:        uuid.New(),
		UserID:    userID,
		Plate:     plate,
		Type:      vehicleType,
		IsDefault: false,
		CreatedAt: time.Now().UTC(),
	}
}

// SetDetails adds additional vehicle details
func (v *Vehicle) SetDetails(make, model, color string) {
	v.Make = make
	v.Model = model
	v.Color = color
}

// MakeDefault sets this vehicle as the default for the user
func (v *Vehicle) MakeDefault() {
	v.IsDefault = true
}
