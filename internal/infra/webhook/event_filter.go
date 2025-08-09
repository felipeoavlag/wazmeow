package webhook

import (
	"strings"

	"wazmeow/internal/domain/entity"
	"wazmeow/pkg/logger"
)

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
		event = strings.TrimSpace(event)
		if event != "" {
			cleanEvents = append(cleanEvents, strings.ToLower(event))
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

// matchesEventGroup verifica se o evento pertence a um grupo específico
func (ef *EventFilter) matchesEventGroup(group, eventType string) bool {
	switch group {
	case "messages":
		return ef.isMessageEvent(eventType)
	case "calls":
		return ef.isCallEvent(eventType)
	case "groups":
		return ef.isGroupEvent(eventType)
	case "presence":
		return ef.isPresenceEvent(eventType)
	case "connection":
		return ef.isConnectionEvent(eventType)
	case "media":
		return ef.isMediaEvent(eventType)
	case "newsletters":
		return ef.isNewsletterEvent(eventType)
	default:
		return false
	}
}

// isMessageEvent verifica se é um evento de mensagem
func (ef *EventFilter) isMessageEvent(eventType string) bool {
	messageEvents := []string{
		"message",
		"receipt",
	}
	
	for _, msgEvent := range messageEvents {
		if eventType == msgEvent {
			return true
		}
	}
	
	return false
}

// isCallEvent verifica se é um evento de chamada
func (ef *EventFilter) isCallEvent(eventType string) bool {
	callEvents := []string{
		"calloffer",
		"callaccept",
		"callterminate",
	}
	
	for _, callEvent := range callEvents {
		if eventType == callEvent {
			return true
		}
	}
	
	return false
}

// isGroupEvent verifica se é um evento de grupo
func (ef *EventFilter) isGroupEvent(eventType string) bool {
	groupEvents := []string{
		"groupinfo",
		"joinedgroup",
	}
	
	for _, groupEvent := range groupEvents {
		if eventType == groupEvent {
			return true
		}
	}
	
	return false
}

// isPresenceEvent verifica se é um evento de presença
func (ef *EventFilter) isPresenceEvent(eventType string) bool {
	presenceEvents := []string{
		"presence",
		"chatpresence",
	}
	
	for _, presenceEvent := range presenceEvents {
		if eventType == presenceEvent {
			return true
		}
	}
	
	return false
}

// isConnectionEvent verifica se é um evento de conexão
func (ef *EventFilter) isConnectionEvent(eventType string) bool {
	connectionEvents := []string{
		"connected",
		"disconnected",
		"logged_out",
		"qr",
		"pair_success",
	}
	
	for _, connEvent := range connectionEvents {
		if eventType == connEvent {
			return true
		}
	}
	
	return false
}

// isMediaEvent verifica se é um evento de mídia
func (ef *EventFilter) isMediaEvent(eventType string) bool {
	mediaEvents := []string{
		"picture",
	}
	
	for _, mediaEvent := range mediaEvents {
		if eventType == mediaEvent {
			return true
		}
	}
	
	return false
}

// isNewsletterEvent verifica se é um evento de newsletter
func (ef *EventFilter) isNewsletterEvent(eventType string) bool {
	newsletterEvents := []string{
		"newsletterjoin",
		"newsletterleave",
		"newslettermutechange",
	}
	
	for _, newsletterEvent := range newsletterEvents {
		if eventType == newsletterEvent {
			return true
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
		"genericevent",
	}
}

// GetEventGroups retorna os grupos de eventos disponíveis
func (ef *EventFilter) GetEventGroups() map[string][]string {
	return map[string][]string{
		"messages":   {"message", "receipt"},
		"calls":      {"calloffer", "callaccept", "callterminate"},
		"groups":     {"groupinfo", "joinedgroup"},
		"presence":   {"presence", "chatpresence"},
		"connection": {"connected", "disconnected", "logged_out", "qr", "pair_success"},
		"media":      {"picture"},
		"newsletters": {"newsletterjoin", "newsletterleave", "newslettermutechange"},
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
