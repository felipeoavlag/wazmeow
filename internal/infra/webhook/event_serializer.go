package webhook

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow/types/events"
)

// WebhookPayload representa o payload completo de um webhook
type WebhookPayload struct {
	Event     string                 `json:"event"`
	SessionID string                 `json:"session_id"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Metadata  *Metadata              `json:"metadata"`
}

// Metadata contém informações adicionais sobre o evento
type Metadata struct {
	Version   string `json:"version"`
	Source    string `json:"source"`
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

// EventSerializer serializa eventos do whatsmeow para JSON
type EventSerializer struct{}

// NewEventSerializer cria uma nova instância do EventSerializer
func NewEventSerializer() *EventSerializer {
	return &EventSerializer{}
}

// SerializeEvent serializa um evento do whatsmeow para WebhookPayload (payload bruto)
func (es *EventSerializer) SerializeEvent(sessionID string, evt interface{}) (*WebhookPayload, error) {
	payload := &WebhookPayload{
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
		Metadata: &Metadata{
			Version:   "1.0",
			Source:    "wazmeow",
			EventID:   generateEventID(),
			Timestamp: time.Now().Unix(),
		},
	}

	// Determinar o tipo de evento e enviar payload bruto
	switch e := evt.(type) {
	case *events.Connected:
		payload.Event = "connected"
		payload.Data = map[string]interface{}{"event": e}
	case *events.Disconnected:
		payload.Event = "disconnected"
		payload.Data = map[string]interface{}{"event": e}
	case *events.LoggedOut:
		payload.Event = "logged_out"
		payload.Data = map[string]interface{}{"event": e}
	case *events.Message:
		payload.Event = "message"
		payload.Data = map[string]interface{}{"event": e}
	case *events.Receipt:
		payload.Event = "receipt"
		payload.Data = map[string]interface{}{"event": e}
	case *events.Presence:
		payload.Event = "presence"
		payload.Data = map[string]interface{}{"event": e}
	case *events.ChatPresence:
		payload.Event = "chat_presence"
		payload.Data = map[string]interface{}{"event": e}
	case *events.GroupInfo:
		payload.Event = "group_info"
		payload.Data = map[string]interface{}{"event": e}
	case *events.Picture:
		payload.Event = "picture"
		payload.Data = map[string]interface{}{"event": e}
	case *events.HistorySync:
		payload.Event = "history_sync"
		payload.Data = map[string]interface{}{"event": e}
	case *events.CallOffer:
		payload.Event = "call_offer"
		payload.Data = map[string]interface{}{"event": e}
	case *events.CallAccept:
		payload.Event = "call_accept"
		payload.Data = map[string]interface{}{"event": e}
	case *events.CallTerminate:
		payload.Event = "call_terminate"
		payload.Data = map[string]interface{}{"event": e}
	case *events.QR:
		payload.Event = "qr"
		payload.Data = map[string]interface{}{"event": e}
	case *events.PairSuccess:
		payload.Event = "pair_success"
		payload.Data = map[string]interface{}{"event": e}
	default:
		// Para eventos não mapeados, enviar como genérico
		payload.Event = "generic"
		payload.Data = map[string]interface{}{
			"event_type": fmt.Sprintf("%T", evt),
			"event":      evt,
		}
	}

	return payload, nil
}

// generateEventID gera um ID único para o evento
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// ToJSON converte o payload para JSON
func (wp *WebhookPayload) ToJSON() ([]byte, error) {
	return json.Marshal(wp)
}
