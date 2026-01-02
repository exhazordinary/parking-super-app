package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/parking-super-app/services/wallet/internal/application"
	"github.com/parking-super-app/services/wallet/internal/domain"
)

type WalletHandler struct {
	walletService *application.WalletService
}

func NewWalletHandler(walletService *application.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
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
	case errors.Is(err, domain.ErrWalletNotFound):
		return http.StatusNotFound, "WALLET_NOT_FOUND", "Wallet not found"
	case errors.Is(err, domain.ErrWalletAlreadyExists):
		return http.StatusConflict, "WALLET_EXISTS", "Wallet already exists for this user"
	case errors.Is(err, domain.ErrInsufficientBalance):
		return http.StatusBadRequest, "INSUFFICIENT_BALANCE", "Insufficient balance"
	case errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusBadRequest, "INVALID_AMOUNT", "Amount must be positive"
	case errors.Is(err, domain.ErrWalletInactive):
		return http.StatusForbidden, "WALLET_INACTIVE", "Wallet is inactive"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "An internal error occurred"
	}
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req application.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	resp, err := h.walletService.CreateWallet(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.walletService.GetWallet(r.Context(), userID)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) TopUp(w http.ResponseWriter, r *http.Request) {
	var req application.TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey != "" {
		req.IdempotencyKey = idempotencyKey
	}

	resp, err := h.walletService.TopUp(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) Pay(w http.ResponseWriter, r *http.Request) {
	var req application.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey != "" {
		req.IdempotencyKey = idempotencyKey
	}

	resp, err := h.walletService.Pay(r.Context(), req)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	walletIDStr := r.URL.Query().Get("wallet_id")
	if walletIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_WALLET_ID", "wallet_id query parameter required")
		return
	}

	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_WALLET_ID", "Invalid wallet ID format")
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

	resp, err := h.walletService.GetTransactions(r.Context(), walletID, limit, offset)
	if err != nil {
		status, code, msg := mapDomainError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
