package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wazmeow/internal/application/dto"
	"wazmeow/internal/application/usecases/session"
	"wazmeow/internal/domain/entities"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// SessionHandler handles HTTP requests for session management
type SessionHandler struct {
	createUseCase   *session.CreateSessionUseCase
	listUseCase     *session.ListSessionsUseCase
	connectUseCase  *session.ConnectSessionUseCase
	whatsappService *whatsapp.Service
}

// NewSessionHandler creates a new SessionHandler
func NewSessionHandler(
	createUseCase *session.CreateSessionUseCase,
	listUseCase *session.ListSessionsUseCase,
	connectUseCase *session.ConnectSessionUseCase,
	whatsappService *whatsapp.Service,
) *SessionHandler {
	return &SessionHandler{
		createUseCase:   createUseCase,
		listUseCase:     listUseCase,
		connectUseCase:  connectUseCase,
		whatsappService: whatsappService,
	}
}

// CreateSession handles POST /sessions/add
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to decode create session request")
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.createUseCase.Execute(r.Context(), req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create session")
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Session created successfully", response)
}

// ListSessions handles GET /sessions/list
func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	response, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list sessions")
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, "Sessions retrieved successfully", response)
}

// ConnectSession handles POST /sessions/{sessionID}/connect
func (h *SessionHandler) ConnectSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	var req dto.ConnectSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to decode connect session request")
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.connectUseCase.Execute(r.Context(), sessionID, req)
	if err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to connect session")
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// GetSessionInfo handles GET /sessions/{sessionID}/info
func (h *SessionHandler) GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	// Get session info from WhatsApp service
	info, err := h.whatsappService.GetSessionInfo(sessionID)
	if err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to get session info")
		h.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get session info: %v", err))
		return
	}

	h.respondSuccess(w, http.StatusOK, "Session info retrieved successfully", info)
}

// DeleteSession handles DELETE /sessions/{sessionID}
func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "sessionID") // TODO: Use sessionID when implementing

	// TODO: Implement delete session logic
	h.respondError(w, http.StatusNotImplemented, "Not implemented yet")
}

// LogoutSession handles POST /sessions/{sessionID}/logout
func (h *SessionHandler) LogoutSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	// Logout from WhatsApp service
	if err := h.whatsappService.Logout(r.Context(), sessionID); err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to logout session")
		h.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to logout: %v", err))
		return
	}

	h.respondSuccess(w, http.StatusOK, "Session logged out successfully", map[string]interface{}{
		"sessionId": sessionID,
		"status":    "disconnected",
		"message":   "Session has been logged out from WhatsApp",
	})
}

// GetQRCode handles GET /sessions/{sessionID}/qr
func (h *SessionHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	// Get QR code from WhatsApp service
	qrCode, err := h.whatsappService.GetQRCode(r.Context(), sessionID)
	if err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to get QR code")
		h.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get QR code: %v", err))
		return
	}

	h.respondSuccess(w, http.StatusOK, "QR code retrieved successfully", map[string]interface{}{
		"sessionId": sessionID,
		"qrCode":    qrCode,
		"message":   "Scan this QR code with WhatsApp to authenticate",
	})
}

// PairPhone handles POST /sessions/{sessionID}/pairphone
func (h *SessionHandler) PairPhone(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	var req dto.PairPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to decode pair phone request")
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Pair phone with WhatsApp service
	linkingCode, err := h.whatsappService.PairPhone(r.Context(), sessionID, req.Phone)
	if err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Str("phone", req.Phone).Msg("Failed to pair phone")
		h.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to pair phone: %v", err))
		return
	}

	h.respondSuccess(w, http.StatusOK, "Phone pairing initiated", map[string]interface{}{
		"sessionId":   sessionID,
		"phone":       req.Phone,
		"linkingCode": linkingCode,
		"message":     "Enter the linking code in WhatsApp to complete pairing",
	})
}

// SetProxy handles POST /sessions/{sessionID}/proxy/set
func (h *SessionHandler) SetProxy(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")

	var req dto.SetProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to decode set proxy request")
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Set proxy configuration
	proxyConfig := &entities.ProxyConfig{
		Enabled:  req.Enabled,
		ProxyURL: req.ProxyURL,
	}

	if err := h.whatsappService.SetProxy(sessionID, proxyConfig); err != nil {
		logger.Error().Err(err).Str("sessionId", sessionID).Msg("Failed to set proxy")
		h.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to set proxy: %v", err))
		return
	}

	h.respondSuccess(w, http.StatusOK, "Proxy configuration updated", map[string]interface{}{
		"sessionId":   sessionID,
		"proxyConfig": proxyConfig,
		"message":     "Proxy configuration has been updated",
	})
}

// respondSuccess sends a successful response
func (h *SessionHandler) respondSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	response := dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	h.respondJSON(w, status, response)
}

// respondError sends an error response
func (h *SessionHandler) respondError(w http.ResponseWriter, status int, message string) {
	response := dto.APIResponse{
		Success: false,
		Error:   message,
	}
	h.respondJSON(w, status, response)
}

// respondJSON sends a JSON response
func (h *SessionHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error().Err(err).Msg("Failed to encode JSON response")
	}
}
