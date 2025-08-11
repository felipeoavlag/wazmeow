package webhook

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"wazmeow/internal/domain/entity"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow/types/events"
)

// =============================================================================
// EVENT FILTERING AND VALIDATION
// =============================================================================

// EventFilter filtra eventos baseado na configuração da sessão
type EventFilter struct{}

// NewEventFilter cria uma nova instância do EventFilter
func NewEventFilter() *EventFilter {
	return &EventFilter{}
}

// ShouldSendEvent verifica se um evento deve ser enviado baseado na configuração da sessão
func (ef *EventFilter) ShouldSendEvent(session *entity.Session, eventType string) bool {
	// Se não há webhook configurado, não enviar
	if session.WebhookURL == "" {
		return false
	}

	// Se não há eventos configurados, enviar todos
	if session.Events == "" {
		return true
	}

	// Verificar se o evento está na lista de eventos configurados
	configuredEvents := ef.parseEvents(session.Events)

	// Verificar se o evento específico está configurado
	if ef.isEventConfigured(configuredEvents, eventType) {
		return true
	}

	// Verificar wildcards e grupos de eventos
	return ef.matchesWildcard(configuredEvents, eventType)
}

// parseEvents converte a string de eventos em um slice
func (ef *EventFilter) parseEvents(eventsStr string) []string {
	if eventsStr == "" {
		return []string{}
	}

	events := strings.Split(eventsStr, ",")
	var cleanEvents []string

	for _, event := range events {
		cleanEvent := strings.TrimSpace(strings.ToLower(event))
		if cleanEvent != "" {
			cleanEvents = append(cleanEvents, cleanEvent)
		}
	}

	return cleanEvents
}

// isEventConfigured verifica se um evento específico está configurado
func (ef *EventFilter) isEventConfigured(configuredEvents []string, eventType string) bool {
	eventTypeLower := strings.ToLower(eventType)

	for _, configuredEvent := range configuredEvents {
		if configuredEvent == eventTypeLower {
			return true
		}
	}

	return false
}

// matchesWildcard verifica se o evento corresponde a algum wildcard ou grupo
func (ef *EventFilter) matchesWildcard(configuredEvents []string, eventType string) bool {
	eventTypeLower := strings.ToLower(eventType)

	for _, configuredEvent := range configuredEvents {
		// Wildcard completo
		if configuredEvent == "*" || configuredEvent == "all" {
			return true
		}

		// Grupos de eventos
		if ef.matchesEventGroup(configuredEvent, eventTypeLower) {
			return true
		}

		// Wildcards com prefixo (ex: "message.*", "call.*")
		if strings.HasSuffix(configuredEvent, "*") {
			prefix := strings.TrimSuffix(configuredEvent, "*")
			if strings.HasPrefix(eventTypeLower, prefix) {
				return true
			}
		}
	}

	return false
}

// matchesEventGroup verifica se o evento pertence a um grupo configurado
func (ef *EventFilter) matchesEventGroup(groupName, eventType string) bool {
	eventGroups := ef.GetEventGroups()

	if events, exists := eventGroups[groupName]; exists {
		for _, event := range events {
			if event == eventType {
				return true
			}
		}
	}

	return false
}

// GetSupportedEvents retorna a lista de eventos suportados
func (ef *EventFilter) GetSupportedEvents() []string {
	return []string{
		// Eventos de conexão
		"connected",
		"disconnected",
		"logged_out",
		"qr",
		"pair_success",

		// Eventos de mensagem
		"message",
		"receipt",

		// Eventos de presença
		"presence",
		"chatpresence",

		// Eventos de grupo
		"groupinfo",
		"joinedgroup",

		// Eventos de mídia
		"picture",

		// Eventos de histórico
		"historysync",

		// Eventos de chamada
		"calloffer",
		"callaccept",
		"callterminate",

		// Eventos de newsletter
		"newsletterjoin",
		"newsletterleave",
		"newslettermutechange",

		// Outros eventos
		"blocklistchange",
		"pushname",
		"businessname",
		"contact",
		"genericevent",
	}
}

// GetEventGroups retorna os grupos de eventos disponíveis
func (ef *EventFilter) GetEventGroups() map[string][]string {
	return map[string][]string{
		"messages":    {"message", "receipt"},
		"calls":       {"calloffer", "callaccept", "callterminate"},
		"groups":      {"groupinfo", "joinedgroup"},
		"presence":    {"presence", "chatpresence"},
		"connection":  {"connected", "disconnected", "logged_out", "qr", "pair_success"},
		"media":       {"picture"},
		"newsletters": {"newsletterjoin", "newsletterleave", "newslettermutechange"},
		"contacts":    {"contact", "pushname", "businessname"},
	}
}

// ValidateEvents valida se os eventos configurados são válidos
func (ef *EventFilter) ValidateEvents(eventsStr string) error {
	if eventsStr == "" {
		return nil // Vazio é válido (significa todos os eventos)
	}

	configuredEvents := ef.parseEvents(eventsStr)
	supportedEvents := ef.GetSupportedEvents()
	eventGroups := ef.GetEventGroups()

	for _, configuredEvent := range configuredEvents {
		// Verificar wildcards
		if configuredEvent == "*" || configuredEvent == "all" {
			continue
		}

		// Verificar grupos
		if _, exists := eventGroups[configuredEvent]; exists {
			continue
		}

		// Verificar wildcards com prefixo
		if strings.HasSuffix(configuredEvent, "*") {
			continue // Assumir que wildcards são válidos
		}

		// Verificar eventos específicos
		found := false
		for _, supportedEvent := range supportedEvents {
			if configuredEvent == supportedEvent {
				found = true
				break
			}
		}

		if !found {
			logger.Warn("Evento não suportado configurado: %s", configuredEvent)
		}
	}

	return nil
}

// =============================================================================
// WEBHOOK PAYLOADS AND SERIALIZATION
// =============================================================================

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

// WebhookEvent representa um evento a ser enviado via webhook
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	SessionID string                 `json:"session_id"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	URL       string                 `json:"-"` // URL do webhook (não serializada)
	Retries   int                    `json:"-"` // Número de tentativas (não serializada)
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
	case *events.JoinedGroup:
		payload.Event = "joined_group"
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
	case *events.NewsletterJoin:
		payload.Event = "newsletter_join"
		payload.Data = map[string]interface{}{"event": e}
	case *events.NewsletterLeave:
		payload.Event = "newsletter_leave"
		payload.Data = map[string]interface{}{"event": e}
	case *events.NewsletterMuteChange:
		payload.Event = "newsletter_mute_change"
		payload.Data = map[string]interface{}{"event": e}
	case *events.BlocklistChange:
		payload.Event = "blocklist_change"
		payload.Data = map[string]interface{}{"event": e}
	case *events.PushName:
		payload.Event = "push_name"
		payload.Data = map[string]interface{}{"event": e}
	case *events.BusinessName:
		payload.Event = "business_name"
		payload.Data = map[string]interface{}{"event": e}
	case *events.QR:
		payload.Event = "qr"
		payload.Data = map[string]interface{}{"event": e}
	case *events.PairSuccess:
		payload.Event = "pair_success"
		payload.Data = map[string]interface{}{"event": e}
	case *events.Contact:
		payload.Event = "contact"
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
