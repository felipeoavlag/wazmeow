package events

import (
	"sync"

	"wazmeow/pkg/logger"
)

// Subscriber define um subscriber de eventos
type Subscriber func(sessionID, eventType string, data interface{})

// Dispatcher gerencia routing de eventos de forma otimizada
type Dispatcher struct {
	subscribers sync.Map // string -> []Subscriber (eventType -> subscribers)
	mu          sync.RWMutex
}

// NewDispatcher cria um novo dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

// Subscribe adiciona um subscriber para um tipo de evento
func (d *Dispatcher) Subscribe(eventType string, subscriber Subscriber) {
	if subscriber == nil {
		return
	}

	// Carregar ou criar slice de subscribers
	value, _ := d.subscribers.LoadOrStore(eventType, make([]Subscriber, 0))
	subscribers := value.([]Subscriber)

	// Adicionar subscriber (thread-safe)
	d.mu.Lock()
	subscribers = append(subscribers, subscriber)
	d.subscribers.Store(eventType, subscribers)
	d.mu.Unlock()

	logger.Debug().Str("eventType", eventType).Msg("Event subscriber added")
}

// Unsubscribe remove um subscriber (implementação simplificada)
func (d *Dispatcher) Unsubscribe(eventType string, subscriber Subscriber) {
	if subscriber == nil {
		return
	}

	value, ok := d.subscribers.Load(eventType)
	if !ok {
		return
	}

	_ = value.([]Subscriber)

	// Remover subscriber (implementação simplificada - remove todos)
	d.mu.Lock()
	d.subscribers.Store(eventType, make([]Subscriber, 0))
	d.mu.Unlock()

	logger.Debug().Str("eventType", eventType).Msg("Event subscribers cleared")
}

// Dispatch envia evento para todos os subscribers
func (d *Dispatcher) Dispatch(sessionID, eventType string, data interface{}) {
	// Carregar subscribers para o tipo de evento
	value, ok := d.subscribers.Load(eventType)
	if !ok {
		// Nenhum subscriber para este tipo
		return
	}

	subscribers := value.([]Subscriber)
	if len(subscribers) == 0 {
		return
	}

	// Dispatch para todos os subscribers em goroutines
	for _, subscriber := range subscribers {
		go func(sub Subscriber) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error().
						Str("sessionID", sessionID).
						Str("eventType", eventType).
						Interface("panic", r).
						Msg("Subscriber panic recovered")
				}
			}()

			sub(sessionID, eventType, data)
		}(subscriber)
	}

	logger.Debug().
		Str("sessionID", sessionID).
		Str("eventType", eventType).
		Int("subscribers", len(subscribers)).
		Msg("Event dispatched")
}

// GetSubscriberCount retorna número de subscribers para um tipo
func (d *Dispatcher) GetSubscriberCount(eventType string) int {
	value, ok := d.subscribers.Load(eventType)
	if !ok {
		return 0
	}

	subscribers := value.([]Subscriber)
	return len(subscribers)
}

// GetEventTypes retorna todos os tipos de eventos com subscribers
func (d *Dispatcher) GetEventTypes() []string {
	var eventTypes []string

	d.subscribers.Range(func(key, value interface{}) bool {
		eventType := key.(string)
		subscribers := value.([]Subscriber)
		if len(subscribers) > 0 {
			eventTypes = append(eventTypes, eventType)
		}
		return true
	})

	return eventTypes
}

// Clear remove todos os subscribers
func (d *Dispatcher) Clear() {
	d.subscribers.Range(func(key, value interface{}) bool {
		d.subscribers.Delete(key)
		return true
	})

	logger.Info().Msg("All event subscribers cleared")
}

// GetStats retorna estatísticas do dispatcher
func (d *Dispatcher) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	totalSubscribers := 0

	d.subscribers.Range(func(key, value interface{}) bool {
		eventType := key.(string)
		subscribers := value.([]Subscriber)
		count := len(subscribers)
		stats[eventType] = count
		totalSubscribers += count
		return true
	})

	return map[string]interface{}{
		"eventTypes":       stats,
		"totalSubscribers": totalSubscribers,
	}
}
