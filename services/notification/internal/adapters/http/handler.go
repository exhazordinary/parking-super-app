package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/parking-super-app/services/notification/internal/application"
	"github.com/parking-super-app/services/notification/internal/domain"
)

type NotificationHandler struct {
	service *application.NotificationService
}

func NewNotificationHandler(service *application.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
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
	case errors.Is(err, domain.ErrNotificationNotFound):
		return http.StatusNotFound, "NOTIFICATION_NOT_FOUND", "Notification not found"
	case errors.Is(err, domain.ErrInvalidChannel):
		return http.StatusBadRequest, "INVALID_CHANNEL", "Invalid notification channel"
	case errors.Is(err, domain.ErrInvalidRecipient):
		return http.StatusBadRequest, "INVALID_RECIPIENT", "Invalid recipient"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred"
	}
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req application.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.service.SendNotification(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *NotificationHandler) SendFromTemplate(w http.ResponseWriter, r *http.Request) {
	var req application.SendFromTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.service.SendFromTemplate(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "Invalid notification ID")
		return
	}

	resp, err := h.service.GetNotification(r.Context(), id)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_USER_ID", "X-User-ID header required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
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

	resp, err := h.service.GetUserNotifications(r.Context(), userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *NotificationHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_USER_ID", "X-User-ID header required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
		return
	}

	resp, err := h.service.GetPreferences(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *NotificationHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_USER_ID", "X-User-ID header required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
		return
	}

	var req application.UpdatePreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}
	req.UserID = userID

	resp, err := h.service.UpdatePreferences(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
