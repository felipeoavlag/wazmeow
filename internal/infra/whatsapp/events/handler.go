package events

import (
	"context"

	"go.mau.fi/whatsmeow/types/events"

	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/pkg/logger"
)

// Handler gerencia eventos WhatsApp de forma otimizada
type Handler struct {
	dispatcher  *Dispatcher
	logger      *Logger
	sessionRepo repositories.SessionRepository
}

// NewHandler cria um novo handler de eventos
func NewHandler(sessionRepo repositories.SessionRepository) *Handler {
	return &Handler{
		dispatcher:  NewDispatcher(),
		logger:      NewLogger(),
		sessionRepo: sessionRepo,
	}
}

// Setup configura event handlers para um wrapper
func (h *Handler) Setup(wrapper WrapperInterface) {
	client := wrapper.Client()
	sessionID := wrapper.SessionID()

	// Registrar handler principal
	client.AddEventHandler(func(evt interface{}) {
		h.handleEvent(sessionID, evt)
	})

	logger.Info().Str("sessionID", sessionID).Msg("Event handlers configured")
}

// handleEvent processa eventos de forma otimizada
func (h *Handler) handleEvent(sessionID string, evt interface{}) {
	// Log estruturado do evento
	h.logger.LogEvent(sessionID, evt)

	// Dispatch por tipo para handlers específicos
	switch e := evt.(type) {
	case *events.Connected:
		h.handleConnected(sessionID, e)
	case *events.Disconnected:
		h.handleDisconnected(sessionID, e)
	case *events.QR:
		h.handleQR(sessionID, e)
	case *events.PairSuccess:
		h.handlePairSuccess(sessionID, e)
	case *events.LoggedOut:
		h.handleLoggedOut(sessionID, e)
	case *events.Message:
		h.handleMessage(sessionID, e)
	case *events.Receipt:
		h.handleReceipt(sessionID, e)
	case *events.Presence:
		h.handlePresence(sessionID, e)
	case *events.PushName:
		h.handlePushName(sessionID, e)
	default:
		// Eventos não tratados especificamente
		h.dispatcher.Dispatch(sessionID, "unknown", evt)
	}
}

// handleConnected processa evento de conexão
func (h *Handler) handleConnected(sessionID string, evt *events.Connected) {
	logger.Info().Str("sessionID", sessionID).Msg("✅ Session connected")

	// Atualizar status no banco
	h.updateSessionStatus(sessionID, entities.StatusConnected)

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "connected", evt)
}

// handleDisconnected processa evento de desconexão
func (h *Handler) handleDisconnected(sessionID string, evt *events.Disconnected) {
	logger.Warn().Str("sessionID", sessionID).Msg("❌ Session disconnected")

	// Atualizar status no banco
	h.updateSessionStatus(sessionID, entities.StatusDisconnected)

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "disconnected", evt)
}

// handleQR processa evento de QR code
func (h *Handler) handleQR(sessionID string, evt *events.QR) {
	logger.Info().Str("sessionID", sessionID).Msg("📱 QR event")

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "qr", evt)
}

// handlePairSuccess processa sucesso de pareamento
func (h *Handler) handlePairSuccess(sessionID string, evt *events.PairSuccess) {
	logger.Info().Str("sessionID", sessionID).Str("jid", evt.ID.String()).Msg("🔗 Pair success")

	// Atualizar JID no banco
	h.updateSessionJID(sessionID, evt.ID.String())

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "pair_success", evt)
}

// handleLoggedOut processa logout
func (h *Handler) handleLoggedOut(sessionID string, evt *events.LoggedOut) {
	logger.Info().Str("sessionID", sessionID).Msg("🚪 Session logged out")

	// Atualizar status no banco
	h.updateSessionStatus(sessionID, entities.StatusDisconnected)

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "logged_out", evt)
}

// handleMessage processa mensagens
func (h *Handler) handleMessage(sessionID string, evt *events.Message) {
	logger.Info().
		Str("sessionID", sessionID).
		Str("chat", evt.Info.Chat.String()).
		Str("sender", evt.Info.Sender.String()).
		Bool("fromMe", evt.Info.IsFromMe).
		Msg("📨 Message received")

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "message", evt)
}

// handleReceipt processa confirmações de leitura
func (h *Handler) handleReceipt(sessionID string, evt *events.Receipt) {
	logger.Debug().
		Str("sessionID", sessionID).
		Str("from", evt.Chat.String()).
		Msg("✅ Receipt received")

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "receipt", evt)
}

// handlePresence processa eventos de presença
func (h *Handler) handlePresence(sessionID string, evt *events.Presence) {
	logger.Debug().
		Str("sessionID", sessionID).
		Str("from", evt.From.String()).
		Msg("👁️ Presence update")

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "presence", evt)
}

// handlePushName processa nomes de contatos
func (h *Handler) handlePushName(sessionID string, evt *events.PushName) {
	logger.Debug().
		Str("sessionID", sessionID).
		Str("jid", evt.JID.String()).
		Msg("👤 Push name update")

	// Dispatch para subscribers
	h.dispatcher.Dispatch(sessionID, "push_name", evt)
}

// updateSessionStatus atualiza status da sessão no banco
func (h *Handler) updateSessionStatus(sessionID string, status entities.SessionStatus) {
	ctx := context.Background()
	if err := h.sessionRepo.UpdateStatus(ctx, sessionID, status); err != nil {
		logger.Error().
			Str("sessionID", sessionID).
			Str("status", string(status)).
			Err(err).
			Msg("Failed to update session status")
	}
}

// updateSessionJID atualiza JID da sessão no banco
func (h *Handler) updateSessionJID(sessionID, jid string) {
	ctx := context.Background()
	session, err := h.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.Error().Str("sessionID", sessionID).Err(err).Msg("Failed to get session for JID update")
		return
	}

	session.SetDeviceJID(jid)
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		logger.Error().
			Str("sessionID", sessionID).
			Str("jid", jid).
			Err(err).
			Msg("Failed to update session JID")
	}
}

// WrapperInterface define interface mínima para wrapper
type WrapperInterface interface {
	SessionID() string
	Client() ClientInterface
}

// ClientInterface define interface mínima para cliente
type ClientInterface interface {
	AddEventHandler(handler func(interface{}))
}
