package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/domain"
)

// SessionRepository defines persistence operations for parking sessions
type SessionRepository interface {
	Create(ctx context.Context, session *domain.ParkingSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ParkingSession, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ParkingSession, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ParkingSession, error)
	GetByProviderID(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*domain.ParkingSession, error)
	Update(ctx context.Context, session *domain.ParkingSession) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int, error)
}

// VehicleRepository defines persistence operations for vehicles
type VehicleRepository interface {
	Create(ctx context.Context, vehicle *domain.Vehicle) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Vehicle, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Vehicle, error)
	GetByPlate(ctx context.Context, plate string) (*domain.Vehicle, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetDefault(ctx context.Context, userID, vehicleID uuid.UUID) error
}
