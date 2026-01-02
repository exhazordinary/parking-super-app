package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/provider/internal/domain"
	"github.com/parking-super-app/services/provider/internal/ports"
)

// ProviderService handles provider-related use cases
type ProviderService struct {
	providers   ports.ProviderRepository
	credentials ports.CredentialsRepository
	locations   ports.LocationRepository
	events      ports.EventPublisher
	logger      ports.Logger
}

func NewProviderService(
	providers ports.ProviderRepository,
	credentials ports.CredentialsRepository,
	locations ports.LocationRepository,
	events ports.EventPublisher,
	logger ports.Logger,
) *ProviderService {
	return &ProviderService{
		providers:   providers,
		credentials: credentials,
		locations:   locations,
		events:      events,
		logger:      logger,
	}
}

// Request/Response DTOs

type RegisterProviderRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	MFEURL      string `json:"mfe_url"`
	APIBaseURL  string `json:"api_base_url"`
}

type ProviderResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Code        string               `json:"code"`
	Description string               `json:"description"`
	LogoURL     string               `json:"logo_url,omitempty"`
	Status      string               `json:"status"`
	MFEURL      string               `json:"mfe_url"`
	APIBaseURL  string               `json:"api_base_url"`
	Config      domain.ProviderConfig `json:"config"`
}

type CredentialsResponse struct {
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	Environment string `json:"environment"`
}

type AddLocationRequest struct {
	ProviderID uuid.UUID `json:"provider_id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	HourlyRate float64   `json:"hourly_rate"`
	DailyMax   float64   `json:"daily_max"`
}

type LocationResponse struct {
	ID          uuid.UUID              `json:"id"`
	ProviderID  uuid.UUID              `json:"provider_id"`
	Name        string                 `json:"name"`
	Address     string                 `json:"address"`
	City        string                 `json:"city"`
	Latitude    float64                `json:"latitude"`
	Longitude   float64                `json:"longitude"`
	TotalSpaces int                    `json:"total_spaces"`
	Pricing     domain.LocationPricing `json:"pricing"`
}

// RegisterProvider creates a new parking provider
func (s *ProviderService) RegisterProvider(ctx context.Context, req RegisterProviderRequest) (*ProviderResponse, error) {
	s.logger.Info("registering provider", ports.String("code", req.Code))

	// Check if provider code already exists
	existing, err := s.providers.GetByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, domain.ErrProviderAlreadyExists
	}

	provider, err := domain.NewProvider(req.Name, req.Code, req.MFEURL, req.APIBaseURL)
	if err != nil {
		return nil, err
	}
	provider.Description = req.Description
	provider.LogoURL = req.LogoURL

	if err := s.providers.Create(ctx, provider); err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	// Publish event asynchronously
	go func() {
		event := ports.Event{
			Type: ports.EventProviderCreated,
			Payload: map[string]interface{}{
				"provider_id": provider.ID.String(),
				"code":        provider.Code,
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return s.toProviderResponse(provider), nil
}

// GetProvider retrieves a provider by ID
func (s *ProviderService) GetProvider(ctx context.Context, id uuid.UUID) (*ProviderResponse, error) {
	provider, err := s.providers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toProviderResponse(provider), nil
}

// GetProviderByCode retrieves a provider by code
func (s *ProviderService) GetProviderByCode(ctx context.Context, code string) (*ProviderResponse, error) {
	provider, err := s.providers.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return s.toProviderResponse(provider), nil
}

// ListProviders retrieves all providers
func (s *ProviderService) ListProviders(ctx context.Context, activeOnly bool) ([]*ProviderResponse, error) {
	providers, err := s.providers.GetAll(ctx, activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}

	responses := make([]*ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = s.toProviderResponse(p)
	}
	return responses, nil
}

// ActivateProvider activates a pending or inactive provider
func (s *ProviderService) ActivateProvider(ctx context.Context, id uuid.UUID) error {
	provider, err := s.providers.GetByID(ctx, id)
	if err != nil {
		return err
	}

	provider.Activate()
	if err := s.providers.Update(ctx, provider); err != nil {
		return fmt.Errorf("failed to activate provider: %w", err)
	}

	go func() {
		event := ports.Event{
			Type: ports.EventProviderActivated,
			Payload: map[string]interface{}{
				"provider_id": provider.ID.String(),
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return nil
}

// DeactivateProvider deactivates a provider
func (s *ProviderService) DeactivateProvider(ctx context.Context, id uuid.UUID) error {
	provider, err := s.providers.GetByID(ctx, id)
	if err != nil {
		return err
	}

	provider.Deactivate()
	if err := s.providers.Update(ctx, provider); err != nil {
		return fmt.Errorf("failed to deactivate provider: %w", err)
	}

	return nil
}

// GenerateCredentials creates API credentials for a provider
func (s *ProviderService) GenerateCredentials(ctx context.Context, providerID uuid.UUID, env domain.Environment) (*CredentialsResponse, error) {
	// Verify provider exists
	_, err := s.providers.GetByID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	creds, err := domain.NewProviderCredentials(providerID, env)
	if err != nil {
		return nil, fmt.Errorf("failed to generate credentials: %w", err)
	}

	if err := s.credentials.Create(ctx, creds); err != nil {
		return nil, fmt.Errorf("failed to store credentials: %w", err)
	}

	// Return credentials with secret visible only once
	return &CredentialsResponse{
		APIKey:      creds.APIKey,
		APISecret:   creds.APISecret,
		Environment: string(creds.Environment),
	}, nil
}

// AddLocation adds a parking location for a provider
func (s *ProviderService) AddLocation(ctx context.Context, req AddLocationRequest) (*LocationResponse, error) {
	// Verify provider exists and is active
	provider, err := s.providers.GetByID(ctx, req.ProviderID)
	if err != nil {
		return nil, err
	}
	if !provider.IsActive() {
		return nil, domain.ErrProviderInactive
	}

	location := domain.NewLocation(
		req.ProviderID,
		req.Name,
		req.Address,
		req.City,
		req.State,
		req.Latitude,
		req.Longitude,
	)
	location.PostalCode = req.PostalCode
	location.SetPricing(req.HourlyRate, req.DailyMax)

	if err := s.locations.Create(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	go func() {
		event := ports.Event{
			Type: ports.EventLocationAdded,
			Payload: map[string]interface{}{
				"location_id": location.ID.String(),
				"provider_id": provider.ID.String(),
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return s.toLocationResponse(location), nil
}

// GetProviderLocations retrieves all locations for a provider
func (s *ProviderService) GetProviderLocations(ctx context.Context, providerID uuid.UUID) ([]*LocationResponse, error) {
	locations, err := s.locations.GetByProviderID(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	responses := make([]*LocationResponse, len(locations))
	for i, loc := range locations {
		responses[i] = s.toLocationResponse(loc)
	}
	return responses, nil
}

// GetNearbyLocations finds parking locations near coordinates
func (s *ProviderService) GetNearbyLocations(ctx context.Context, lat, lng, radiusKm float64) ([]*LocationResponse, error) {
	locations, err := s.locations.GetNearby(ctx, lat, lng, radiusKm)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby locations: %w", err)
	}

	responses := make([]*LocationResponse, len(locations))
	for i, loc := range locations {
		responses[i] = s.toLocationResponse(loc)
	}
	return responses, nil
}

func (s *ProviderService) toProviderResponse(p *domain.Provider) *ProviderResponse {
	return &ProviderResponse{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		LogoURL:     p.LogoURL,
		Status:      string(p.Status),
		MFEURL:      p.MFEURL,
		APIBaseURL:  p.APIBaseURL,
		Config:      p.Config,
	}
}

func (s *ProviderService) toLocationResponse(l *domain.Location) *LocationResponse {
	return &LocationResponse{
		ID:          l.ID,
		ProviderID:  l.ProviderID,
		Name:        l.Name,
		Address:     l.Address,
		City:        l.City,
		Latitude:    l.Latitude,
		Longitude:   l.Longitude,
		TotalSpaces: l.TotalSpaces,
		Pricing:     l.Pricing,
	}
}
