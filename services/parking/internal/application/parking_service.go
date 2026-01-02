package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/domain"
	"github.com/parking-super-app/services/parking/internal/ports"
	"github.com/shopspring/decimal"
)

// ParkingService handles parking session use cases
type ParkingService struct {
	sessions   ports.SessionRepository
	vehicles   ports.VehicleRepository
	provider   ports.ProviderClient
	wallet     ports.WalletClient
	events     ports.EventPublisher
	logger     ports.Logger
}

func NewParkingService(
	sessions ports.SessionRepository,
	vehicles ports.VehicleRepository,
	provider ports.ProviderClient,
	wallet ports.WalletClient,
	events ports.EventPublisher,
	logger ports.Logger,
) *ParkingService {
	return &ParkingService{
		sessions: sessions,
		vehicles: vehicles,
		provider: provider,
		wallet:   wallet,
		events:   events,
		logger:   logger,
	}
}

// Request/Response DTOs

type StartSessionRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	ProviderID   uuid.UUID `json:"provider_id"`
	LocationID   uuid.UUID `json:"location_id"`
	VehiclePlate string    `json:"vehicle_plate"`
	VehicleType  string    `json:"vehicle_type"`
}

type SessionResponse struct {
	ID                uuid.UUID        `json:"id"`
	UserID            uuid.UUID        `json:"user_id"`
	ProviderID        uuid.UUID        `json:"provider_id"`
	LocationID        uuid.UUID        `json:"location_id"`
	ExternalSessionID string           `json:"external_session_id,omitempty"`
	VehiclePlate      string           `json:"vehicle_plate"`
	VehicleType       string           `json:"vehicle_type"`
	EntryTime         string           `json:"entry_time"`
	ExitTime          string           `json:"exit_time,omitempty"`
	Duration          int              `json:"duration_minutes"`
	Amount            decimal.Decimal  `json:"amount"`
	Status            string           `json:"status"`
}

type EndSessionRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	WalletID  uuid.UUID `json:"wallet_id"`
}

type EndSessionResponse struct {
	SessionID     uuid.UUID       `json:"session_id"`
	Duration      int             `json:"duration_minutes"`
	Amount        decimal.Decimal `json:"amount"`
	PaymentStatus string          `json:"payment_status"`
}

type SessionListResponse struct {
	Sessions []*SessionResponse `json:"sessions"`
	Total    int                `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}

type RegisterVehicleRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Plate  string    `json:"plate"`
	Type   string    `json:"type"`
	Make   string    `json:"make,omitempty"`
	Model  string    `json:"model,omitempty"`
	Color  string    `json:"color,omitempty"`
}

type VehicleResponse struct {
	ID        uuid.UUID `json:"id"`
	Plate     string    `json:"plate"`
	Type      string    `json:"type"`
	Make      string    `json:"make,omitempty"`
	Model     string    `json:"model,omitempty"`
	Color     string    `json:"color,omitempty"`
	IsDefault bool      `json:"is_default"`
}

// StartSession initiates a new parking session
func (s *ParkingService) StartSession(ctx context.Context, req StartSessionRequest) (*SessionResponse, error) {
	s.logger.Info("starting parking session",
		ports.String("user_id", req.UserID.String()),
		ports.String("provider_id", req.ProviderID.String()),
	)

	// Create session in our system first
	session, err := domain.NewParkingSession(
		req.UserID,
		req.ProviderID,
		req.LocationID,
		req.VehiclePlate,
		req.VehicleType,
	)
	if err != nil {
		return nil, err
	}

	// Call provider API to start session
	providerResp, err := s.provider.StartSession(ctx, ports.StartSessionRequest{
		ProviderID:   req.ProviderID,
		LocationID:   req.LocationID,
		VehiclePlate: req.VehiclePlate,
		VehicleType:  req.VehicleType,
		UserRef:      session.ID.String(),
	})
	if err != nil {
		s.logger.Error("failed to start session with provider", ports.Err(err))
		return nil, fmt.Errorf("failed to start session with provider: %w", err)
	}

	session.SetExternalSessionID(providerResp.ExternalSessionID)

	// Persist session
	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	// Publish event
	go func() {
		event := ports.Event{
			Type: ports.EventSessionStarted,
			Payload: map[string]interface{}{
				"session_id":  session.ID.String(),
				"user_id":     session.UserID.String(),
				"provider_id": session.ProviderID.String(),
				"plate":       session.VehiclePlate,
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return s.toSessionResponse(session), nil
}

// EndSession completes a parking session and processes payment
func (s *ParkingService) EndSession(ctx context.Context, req EndSessionRequest) (*EndSessionResponse, error) {
	s.logger.Info("ending parking session", ports.String("session_id", req.SessionID.String()))

	session, err := s.sessions.GetByID(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

	if !session.IsActive() {
		return nil, domain.ErrSessionAlreadyEnded
	}

	// Get final amount from provider
	providerResp, err := s.provider.EndSession(ctx, ports.EndSessionRequest{
		ProviderID:        session.ProviderID,
		ExternalSessionID: session.ExternalSessionID,
	})
	if err != nil {
		s.logger.Error("failed to end session with provider", ports.Err(err))
		return nil, fmt.Errorf("failed to end session with provider: %w", err)
	}

	// End session with the calculated amount
	if err := session.End(providerResp.Amount); err != nil {
		return nil, err
	}

	// Process payment through wallet
	paymentResp, err := s.wallet.Pay(ctx, ports.PaymentRequest{
		WalletID:       req.WalletID,
		Amount:         session.Amount,
		ProviderID:     session.ProviderID,
		ReferenceID:    session.ID.String(),
		Description:    fmt.Sprintf("Parking at location %s", session.LocationID),
		IdempotencyKey: fmt.Sprintf("parking-%s", session.ID),
	})
	if err != nil {
		s.logger.Error("payment failed", ports.Err(err))
		// Session ended but payment failed - needs handling
		session.Status = domain.SessionStatusFailed
		s.sessions.Update(ctx, session)
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	session.MarkPaid(paymentResp.TransactionID)

	// Update session
	if err := s.sessions.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Publish event
	go func() {
		event := ports.Event{
			Type: ports.EventSessionEnded,
			Payload: map[string]interface{}{
				"session_id": session.ID.String(),
				"user_id":    session.UserID.String(),
				"amount":     session.Amount.String(),
				"duration":   session.Duration,
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return &EndSessionResponse{
		SessionID:     session.ID,
		Duration:      session.Duration,
		Amount:        session.Amount,
		PaymentStatus: paymentResp.Status,
	}, nil
}

// GetSession retrieves a parking session by ID
func (s *ParkingService) GetSession(ctx context.Context, id uuid.UUID) (*SessionResponse, error) {
	session, err := s.sessions.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toSessionResponse(session), nil
}

// GetUserSessions retrieves parking sessions for a user
func (s *ParkingService) GetUserSessions(ctx context.Context, userID uuid.UUID, limit, offset int) (*SessionListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	sessions, err := s.sessions.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	total, err := s.sessions.CountByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
	}

	responses := make([]*SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = s.toSessionResponse(session)
	}

	return &SessionListResponse{
		Sessions: responses,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

// GetActiveSessions retrieves active parking sessions for a user
func (s *ParkingService) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]*SessionResponse, error) {
	sessions, err := s.sessions.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	responses := make([]*SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = s.toSessionResponse(session)
	}

	return responses, nil
}

// CancelSession cancels an active session
func (s *ParkingService) CancelSession(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if err := session.Cancel(); err != nil {
		return err
	}

	if err := s.sessions.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	go func() {
		event := ports.Event{
			Type: ports.EventSessionCancelled,
			Payload: map[string]interface{}{
				"session_id": session.ID.String(),
				"user_id":    session.UserID.String(),
			},
		}
		s.events.Publish(context.Background(), event)
	}()

	return nil
}

// RegisterVehicle adds a new vehicle for a user
func (s *ParkingService) RegisterVehicle(ctx context.Context, req RegisterVehicleRequest) (*VehicleResponse, error) {
	vehicle := domain.NewVehicle(req.UserID, req.Plate, req.Type)
	vehicle.SetDetails(req.Make, req.Model, req.Color)

	if err := s.vehicles.Create(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to register vehicle: %w", err)
	}

	return s.toVehicleResponse(vehicle), nil
}

// GetUserVehicles retrieves all vehicles for a user
func (s *ParkingService) GetUserVehicles(ctx context.Context, userID uuid.UUID) ([]*VehicleResponse, error) {
	vehicles, err := s.vehicles.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}

	responses := make([]*VehicleResponse, len(vehicles))
	for i, v := range vehicles {
		responses[i] = s.toVehicleResponse(v)
	}

	return responses, nil
}

func (s *ParkingService) toSessionResponse(session *domain.ParkingSession) *SessionResponse {
	resp := &SessionResponse{
		ID:                session.ID,
		UserID:            session.UserID,
		ProviderID:        session.ProviderID,
		LocationID:        session.LocationID,
		ExternalSessionID: session.ExternalSessionID,
		VehiclePlate:      session.VehiclePlate,
		VehicleType:       session.VehicleType,
		EntryTime:         session.EntryTime.Format("2006-01-02T15:04:05Z"),
		Duration:          session.CalculateDuration(),
		Amount:            session.Amount,
		Status:            string(session.Status),
	}
	if session.ExitTime != nil {
		resp.ExitTime = session.ExitTime.Format("2006-01-02T15:04:05Z")
		resp.Duration = session.Duration
	}
	return resp
}

func (s *ParkingService) toVehicleResponse(v *domain.Vehicle) *VehicleResponse {
	return &VehicleResponse{
		ID:        v.ID,
		Plate:     v.Plate,
		Type:      v.Type,
		Make:      v.Make,
		Model:     v.Model,
		Color:     v.Color,
		IsDefault: v.IsDefault,
	}
}
