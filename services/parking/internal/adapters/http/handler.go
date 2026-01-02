package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/parking-super-app/services/parking/internal/application"
	"github.com/parking-super-app/services/parking/internal/domain"
)

type ParkingHandler struct {
	parkingService *application.ParkingService
}

func NewParkingHandler(parkingService *application.ParkingService) *ParkingHandler {
	return &ParkingHandler{parkingService: parkingService}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   &APIError{Code: code, Message: message},
	})
}

func mapDomainError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrSessionNotFound):
		return http.StatusNotFound, "SESSION_NOT_FOUND", "Parking session not found"
	case errors.Is(err, domain.ErrSessionAlreadyEnded):
		return http.StatusBadRequest, "SESSION_ENDED", "Session has already ended"
	case errors.Is(err, domain.ErrInvalidVehiclePlate):
		return http.StatusBadRequest, "INVALID_PLATE", "Invalid vehicle plate number"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred"
	}
}

func (h *ParkingHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	var req application.StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.parkingService.StartSession(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *ParkingHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid session ID format")
		return
	}

	var req struct {
		WalletID uuid.UUID `json:"wallet_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.parkingService.EndSession(r.Context(), application.EndSessionRequest{
		SessionID: sessionID,
		WalletID:  req.WalletID,
	})
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ParkingHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid session ID format")
		return
	}

	resp, err := h.parkingService.GetSession(r.Context(), sessionID)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ParkingHandler) GetUserSessions(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_USER_ID", "X-User-ID header required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format")
		return
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	resp, err := h.parkingService.GetUserSessions(r.Context(), userID, limit, offset)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ParkingHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_USER_ID", "X-User-ID header required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format")
		return
	}

	resp, err := h.parkingService.GetActiveSessions(r.Context(), userID)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ParkingHandler) CancelSession(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid session ID format")
		return
	}

	if err := h.parkingService.CancelSession(r.Context(), sessionID); err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *ParkingHandler) RegisterVehicle(w http.ResponseWriter, r *http.Request) {
	var req application.RegisterVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.parkingService.RegisterVehicle(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register vehicle")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *ParkingHandler) GetUserVehicles(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_USER_ID", "X-User-ID header required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format")
		return
	}

	resp, err := h.parkingService.GetUserVehicles(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get vehicles")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
