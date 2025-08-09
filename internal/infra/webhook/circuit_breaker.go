package webhook

import (
	"sync"
	"time"

	"wazmeow/internal/config"
	"wazmeow/pkg/logger"
)

// CircuitBreakerState representa o estado do circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implementa o padrão Circuit Breaker para webhooks
type CircuitBreaker struct {
	config   *config.CircuitBreakerConfig
	breakers map[string]*URLCircuitBreaker
	mu       sync.RWMutex
}

// URLCircuitBreaker representa um circuit breaker para uma URL específica
type URLCircuitBreaker struct {
	url           string
	state         CircuitBreakerState
	failures      int
	lastFailTime  time.Time
	halfOpenCalls int
	mu            sync.RWMutex
}

// NewCircuitBreaker cria uma nova instância do CircuitBreaker
func NewCircuitBreaker(config *config.CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config:   config,
		breakers: make(map[string]*URLCircuitBreaker),
	}
}

// CanExecute verifica se uma URL pode ser executada
func (cb *CircuitBreaker) CanExecute(url string) bool {
	cb.mu.RLock()
	breaker, exists := cb.breakers[url]
	cb.mu.RUnlock()

	if !exists {
		// Primeira vez, criar circuit breaker para esta URL
		cb.mu.Lock()
		breaker = &URLCircuitBreaker{
			url:   url,
			state: StateClosed,
		}
		cb.breakers[url] = breaker
		cb.mu.Unlock()
		return true
	}

	return breaker.canExecute(cb.config)
}

// RecordSuccess registra um sucesso para uma URL
func (cb *CircuitBreaker) RecordSuccess(url string) {
	cb.mu.RLock()
	breaker, exists := cb.breakers[url]
	cb.mu.RUnlock()

	if exists {
		breaker.recordSuccess(cb.config)
	}
}

// RecordFailure registra uma falha para uma URL
func (cb *CircuitBreaker) RecordFailure(url string) {
	cb.mu.RLock()
	breaker, exists := cb.breakers[url]
	cb.mu.RUnlock()

	if exists {
		breaker.recordFailure(cb.config)
	}
}

// GetState retorna o estado atual do circuit breaker para uma URL
func (cb *CircuitBreaker) GetState(url string) CircuitBreakerState {
	cb.mu.RLock()
	breaker, exists := cb.breakers[url]
	cb.mu.RUnlock()

	if !exists {
		return StateClosed
	}

	breaker.mu.RLock()
	defer breaker.mu.RUnlock()
	return breaker.state
}

// GetStats retorna estatísticas do circuit breaker
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	stats := make(map[string]interface{})
	urlStats := make(map[string]interface{})

	for url, breaker := range cb.breakers {
		breaker.mu.RLock()
		urlStats[url] = map[string]interface{}{
			"state":           breaker.state,
			"failures":        breaker.failures,
			"last_fail_time":  breaker.lastFailTime,
			"half_open_calls": breaker.halfOpenCalls,
		}
		breaker.mu.RUnlock()
	}

	stats["urls"] = urlStats
	stats["total_urls"] = len(cb.breakers)

	return stats
}

// canExecute verifica se o circuit breaker permite execução
func (ucb *URLCircuitBreaker) canExecute(config *config.CircuitBreakerConfig) bool {
	ucb.mu.Lock()
	defer ucb.mu.Unlock()

	switch ucb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Verificar se é hora de tentar half-open
		if time.Since(ucb.lastFailTime) >= config.ResetTimeout {
			ucb.state = StateHalfOpen
			ucb.halfOpenCalls = 0
			logger.Info("Circuit breaker para %s mudou para half-open", ucb.url)
			return true
		}
		return false
	case StateHalfOpen:
		// Permitir apenas um número limitado de chamadas em half-open
		if ucb.halfOpenCalls < config.HalfOpenMaxCalls {
			ucb.halfOpenCalls++
			return true
		}
		return false
	default:
		return false
	}
}

// recordSuccess registra um sucesso
func (ucb *URLCircuitBreaker) recordSuccess(config *config.CircuitBreakerConfig) {
	ucb.mu.Lock()
	defer ucb.mu.Unlock()

	switch ucb.state {
	case StateClosed:
		// Reset failures counter on success
		ucb.failures = 0
	case StateHalfOpen:
		// Se todas as chamadas half-open foram bem-sucedidas, fechar o circuit
		if ucb.halfOpenCalls >= config.HalfOpenMaxCalls {
			ucb.state = StateClosed
			ucb.failures = 0
			ucb.halfOpenCalls = 0
			logger.Info("Circuit breaker para %s fechado após sucesso em half-open", ucb.url)
		}
	}
}

// recordFailure registra uma falha
func (ucb *URLCircuitBreaker) recordFailure(config *config.CircuitBreakerConfig) {
	ucb.mu.Lock()
	defer ucb.mu.Unlock()

	ucb.failures++
	ucb.lastFailTime = time.Now()

	switch ucb.state {
	case StateClosed:
		if ucb.failures >= config.MaxFailures {
			ucb.state = StateOpen
			logger.Warn("Circuit breaker para %s aberto após %d falhas", ucb.url, ucb.failures)
		}
	case StateHalfOpen:
		// Qualquer falha em half-open volta para open
		ucb.state = StateOpen
		ucb.halfOpenCalls = 0
		logger.Warn("Circuit breaker para %s voltou para open após falha em half-open", ucb.url)
	}
}

// Reset reseta o circuit breaker para uma URL
func (cb *CircuitBreaker) Reset(url string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if breaker, exists := cb.breakers[url]; exists {
		breaker.mu.Lock()
		breaker.state = StateClosed
		breaker.failures = 0
		breaker.halfOpenCalls = 0
		breaker.mu.Unlock()
		logger.Info("Circuit breaker para %s resetado", url)
	}
}

// ResetAll reseta todos os circuit breakers
func (cb *CircuitBreaker) ResetAll() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	for _, breaker := range cb.breakers {
		breaker.mu.Lock()
		breaker.state = StateClosed
		breaker.failures = 0
		breaker.halfOpenCalls = 0
		breaker.mu.Unlock()
	}

	logger.Info("Todos os circuit breakers foram resetados")
}
