package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewVehicle(t *testing.T) {
	userID := uuid.New()

	vehicle := NewVehicle(userID, "WKL1234", VehicleTypeCar)

	if vehicle.ID == uuid.Nil {
		t.Error("expected vehicle ID to be set")
	}
	if vehicle.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, vehicle.UserID)
	}
	if vehicle.Plate != "WKL1234" {
		t.Errorf("expected plate WKL1234, got %s", vehicle.Plate)
	}
	if vehicle.Type != VehicleTypeCar {
		t.Errorf("expected type car, got %s", vehicle.Type)
	}
	if vehicle.IsDefault {
		t.Error("new vehicle should not be default")
	}
}

func TestVehicle_SetDetails(t *testing.T) {
	vehicle := NewVehicle(uuid.New(), "ABC123", VehicleTypeCar)

	vehicle.SetDetails("Toyota", "Camry", "White")

	if vehicle.Make != "Toyota" {
		t.Errorf("expected make Toyota, got %s", vehicle.Make)
	}
	if vehicle.Model != "Camry" {
		t.Errorf("expected model Camry, got %s", vehicle.Model)
	}
	if vehicle.Color != "White" {
		t.Errorf("expected color White, got %s", vehicle.Color)
	}
}

func TestVehicle_MakeDefault(t *testing.T) {
	vehicle := NewVehicle(uuid.New(), "ABC123", VehicleTypeCar)

	if vehicle.IsDefault {
		t.Error("new vehicle should not be default")
	}

	vehicle.MakeDefault()

	if !vehicle.IsDefault {
		t.Error("vehicle should be default after MakeDefault")
	}
}

func TestVehicleTypes(t *testing.T) {
	if VehicleTypeCar != "car" {
		t.Errorf("expected VehicleTypeCar to be 'car', got %s", VehicleTypeCar)
	}
	if VehicleTypeMotorcycle != "motorcycle" {
		t.Errorf("expected VehicleTypeMotorcycle to be 'motorcycle', got %s", VehicleTypeMotorcycle)
	}
	if VehicleTypeTruck != "truck" {
		t.Errorf("expected VehicleTypeTruck to be 'truck', got %s", VehicleTypeTruck)
	}
}
