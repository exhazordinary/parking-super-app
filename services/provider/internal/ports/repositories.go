package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/provider/internal/domain"
)

// ProviderRepository defines the interface for provider persistence
type ProviderRepository interface {
	Create(ctx context.Context, provider *domain.Provider) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Provider, error)
	GetByCode(ctx context.Context, code string) (*domain.Provider, error)
	GetAll(ctx context.Context, activeOnly bool) ([]*domain.Provider, error)
	Update(ctx context.Context, provider *domain.Provider) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// CredentialsRepository defines the interface for credential persistence
type CredentialsRepository interface {
	Create(ctx context.Context, creds *domain.ProviderCredentials) error
	GetByAPIKey(ctx context.Context, apiKey string) (*domain.ProviderCredentials, error)
	GetByProviderID(ctx context.Context, providerID uuid.UUID, env domain.Environment) (*domain.ProviderCredentials, error)
	Update(ctx context.Context, creds *domain.ProviderCredentials) error
	Revoke(ctx context.Context, id uuid.UUID) error
}

// LocationRepository defines the interface for location persistence
type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	GetByProviderID(ctx context.Context, providerID uuid.UUID) ([]*domain.Location, error)
	GetNearby(ctx context.Context, lat, lng float64, radiusKm float64) ([]*domain.Location, error)
	Update(ctx context.Context, location *domain.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
}
