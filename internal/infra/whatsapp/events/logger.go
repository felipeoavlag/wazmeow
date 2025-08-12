package events

import (
	"encoding/json"

	"go.mau.fi/whatsmeow/types/events"

	"wazmeow/pkg/logger"
)

// Logger especializado para eventos WhatsApp com JSON otimizado
type Logger struct{}

// NewLogger cria um novo logger de eventos
func NewLogger() *Logger {
	return &Logger{}
}

// LogEvent faz log estruturado de eventos WhatsApp
func (l *Logger) LogEvent(sessionID string, evt interface{}) {
	switch e := evt.(type) {
	case *events.Message:
		l.logMessage(sessionID, e)
	case *events.Connected:
		l.logConnected(sessionID, e)
	case *events.Disconnected:
		l.logDisconnected(sessionID, e)
	case *events.QR:
		l.logQR(sessionID, e)
	case *events.PairSuccess:
		l.logPairSuccess(sessionID, e)
	case *events.LoggedOut:
		l.logLoggedOut(sessionID, e)
	case *events.Receipt:
		l.logReceipt(sessionID, e)
	case *events.Presence:
		l.logPresence(sessionID, e)
	case *events.PushName:
		l.logPushName(sessionID, e)
	default:
		l.logGeneric(sessionID, evt)
	}
}

// logMessage faz log otimizado de mensagens
func (l *Logger) logMessage(sessionID string, evt *events.Message) {
	logger.Info().
		Str("sessionID", sessionID).
		Str("chat", evt.Info.Chat.String()).
		Str("sender", evt.Info.Sender.String()).
		Bool("fromMe", evt.Info.IsFromMe).
		Str("type", evt.Info.Type).
		Time("timestamp", evt.Info.Timestamp).
		Msg("📨 Message received")

	// Log do payload completo com JSON limpo
	if payloadBytes, err := json.MarshalIndent(evt, "", "  "); err == nil {
		logger.Info().
			Str("sessionID", sessionID).
			RawJSON("payload", payloadBytes).
			Msg("📨 PAYLOAD")
	} else {
		logger.Error().
			Str("sessionID", sessionID).
			Err(err).
			Msg("📨 Failed to marshal message payload")
	}
}

// logConnected faz log de conexão
func (l *Logger) logConnected(sessionID string, evt *events.Connected) {
	logger.Info().
		Str("sessionID", sessionID).
		Msg("✅ Connected")
}

// logDisconnected faz log de desconexão
func (l *Logger) logDisconnected(sessionID string, evt *events.Disconnected) {
	logger.Warn().
		Str("sessionID", sessionID).
		Msg("❌ Disconnected")
}

// logQR faz log de eventos QR
func (l *Logger) logQR(sessionID string, evt *events.QR) {
	logger.Info().
		Str("sessionID", sessionID).
		Msg("📱 QR")
}

// logPairSuccess faz log de pareamento bem-sucedido
func (l *Logger) logPairSuccess(sessionID string, evt *events.PairSuccess) {
	logger.Info().
		Str("sessionID", sessionID).
		Str("jid", evt.ID.String()).
		Msg("🔗 Pair Success")
}

// logLoggedOut faz log de logout
func (l *Logger) logLoggedOut(sessionID string, evt *events.LoggedOut) {
	logger.Info().
		Str("sessionID", sessionID).
		Str("reason", evt.Reason.String()).
		Msg("🚪 Logged Out")
}

// logReceipt faz log de confirmações
func (l *Logger) logReceipt(sessionID string, evt *events.Receipt) {
	logger.Debug().
		Str("sessionID", sessionID).
		Str("from", evt.Chat.String()).
		Str("type", string(evt.Type)).
		Msg("✅ Receipt")
}

// logPresence faz log de presença
func (l *Logger) logPresence(sessionID string, evt *events.Presence) {
	logger.Debug().
		Str("sessionID", sessionID).
		Str("from", evt.From.String()).
		Msg("👁️ Presence")
}

// logPushName faz log de nomes
func (l *Logger) logPushName(sessionID string, evt *events.PushName) {
	logger.Debug().
		Str("sessionID", sessionID).
		Str("jid", evt.JID.String()).
		Msg("👤 Push Name")
}

// logGeneric faz log de eventos genéricos
func (l *Logger) logGeneric(sessionID string, evt interface{}) {
	// Log básico para eventos não específicos
	logger.Debug().
		Str("sessionID", sessionID).
		Str("type", getEventType(evt)).
		Msg("🔄 Generic Event")

	// Payload completo apenas para debug
	if payloadBytes, err := json.MarshalIndent(evt, "", "  "); err == nil {
		logger.Debug().
			Str("sessionID", sessionID).
			RawJSON("payload", payloadBytes).
			Msg("🔄 GENERIC PAYLOAD")
	}
}

// getEventType retorna o tipo do evento
func getEventType(evt interface{}) string {
	switch evt.(type) {
	case *events.Message:
		return "Message"
	case *events.Connected:
		return "Connected"
	case *events.Disconnected:
		return "Disconnected"
	case *events.QR:
		return "QR"
	case *events.PairSuccess:
		return "PairSuccess"
	case *events.LoggedOut:
		return "LoggedOut"
	case *events.Receipt:
		return "Receipt"
	case *events.Presence:
		return "Presence"
	case *events.PushName:
		return "PushName"
	default:
		return "Unknown"
	}
}

// LogStats faz log de estatísticas de eventos
func (l *Logger) LogStats(sessionID string, stats map[string]interface{}) {
	logger.Info().
		Str("sessionID", sessionID).
		Interface("stats", stats).
		Msg("📊 Event Stats")
}

// LogError faz log de erros relacionados a eventos
func (l *Logger) LogError(sessionID string, err error, context string) {
	logger.Error().
		Str("sessionID", sessionID).
		Str("context", context).
		Err(err).
		Msg("❌ Event Error")
}
