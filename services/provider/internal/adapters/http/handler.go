package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/parking-super-app/services/provider/internal/application"
	"github.com/parking-super-app/services/provider/internal/domain"
)

type ProviderHandler struct {
	providerService *application.ProviderService
}

func NewProviderHandler(providerService *application.ProviderService) *ProviderHandler {
	return &ProviderHandler{providerService: providerService}
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
	case errors.Is(err, domain.ErrProviderNotFound):
		return http.StatusNotFound, "PROVIDER_NOT_FOUND", "Provider not found"
	case errors.Is(err, domain.ErrProviderAlreadyExists):
		return http.StatusConflict, "PROVIDER_EXISTS", "Provider with this code already exists"
	case errors.Is(err, domain.ErrInvalidProviderCode):
		return http.StatusBadRequest, "INVALID_CODE", "Provider code must be alphanumeric"
	case errors.Is(err, domain.ErrInvalidMFEURL):
		return http.StatusBadRequest, "INVALID_MFE_URL", "Invalid MFE URL"
	case errors.Is(err, domain.ErrProviderInactive):
		return http.StatusForbidden, "PROVIDER_INACTIVE", "Provider is not active"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred"
	}
}

func (h *ProviderHandler) RegisterProvider(w http.ResponseWriter, r *http.Request) {
	var req application.RegisterProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.providerService.RegisterProvider(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *ProviderHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid provider ID format")
		return
	}

	resp, err := h.providerService.GetProvider(r.Context(), id)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProviderHandler) GetProviderByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "MISSING_CODE", "Provider code is required")
		return
	}

	resp, err := h.providerService.GetProviderByCode(r.Context(), code)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProviderHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	resp, err := h.providerService.ListProviders(r.Context(), activeOnly)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProviderHandler) ActivateProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid provider ID format")
		return
	}

	if err := h.providerService.ActivateProvider(r.Context(), id); err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "activated"})
}

func (h *ProviderHandler) DeactivateProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid provider ID format")
		return
	}

	if err := h.providerService.DeactivateProvider(r.Context(), id); err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deactivated"})
}

type GenerateCredentialsRequest struct {
	Environment string `json:"environment"`
}

func (h *ProviderHandler) GenerateCredentials(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid provider ID format")
		return
	}

	var req GenerateCredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Environment = "sandbox"
	}

	env := domain.EnvironmentSandbox
	if req.Environment == "production" {
		env = domain.EnvironmentProduction
	}

	resp, err := h.providerService.GenerateCredentials(r.Context(), id, env)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *ProviderHandler) AddLocation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	providerID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid provider ID format")
		return
	}

	var req application.AddLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}
	req.ProviderID = providerID

	resp, err := h.providerService.AddLocation(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *ProviderHandler) GetProviderLocations(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid provider ID format")
		return
	}

	resp, err := h.providerService.GetProviderLocations(r.Context(), id)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
