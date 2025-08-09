package webhook

import (
	"context"
	"sync"
	"time"

	"wazmeow/internal/config"
	"wazmeow/pkg/logger"
)

// RateLimit implementa rate limiting para webhooks
type RateLimit struct {
	config   *config.RateLimitConfig
	limiters map[string]*SessionLimiter
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// SessionLimiter representa um rate limiter para uma sessão específica
type SessionLimiter struct {
	sessionID  string
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimit cria uma nova instância do RateLimit
func NewRateLimit(config *config.RateLimitConfig) *RateLimit {
	ctx, cancel := context.WithCancel(context.Background())

	return &RateLimit{
		config:   config,
		limiters: make(map[string]*SessionLimiter),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Allow verifica se uma sessão pode enviar um webhook
func (rl *RateLimit) Allow(sessionID string) bool {
	rl.mu.RLock()
	limiter, exists := rl.limiters[sessionID]
	rl.mu.RUnlock()

	if !exists {
		// Primeira vez, criar limiter para esta sessão
		rl.mu.Lock()
		limiter = &SessionLimiter{
			sessionID:  sessionID,
			tokens:     rl.config.BurstSize,
			lastRefill: time.Now(),
		}
		rl.limiters[sessionID] = limiter
		rl.mu.Unlock()
	}

	return limiter.allow(rl.config)
}

// allow verifica se o limiter permite uma requisição
func (sl *SessionLimiter) allow(config *config.RateLimitConfig) bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	now := time.Now()

	// Calcular quantos tokens devem ser adicionados desde a última recarga
	elapsed := now.Sub(sl.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * config.RequestsPerSecond

	if tokensToAdd > 0 {
		sl.tokens += tokensToAdd
		if sl.tokens > config.BurstSize {
			sl.tokens = config.BurstSize
		}
		sl.lastRefill = now
	}

	// Verificar se há tokens disponíveis
	if sl.tokens > 0 {
		sl.tokens--
		return true
	}

	return false
}

// StartCleanup inicia o processo de limpeza de limiters inativos
func (rl *RateLimit) StartCleanup() {
	defer func() {
		// Este método é chamado como goroutine no WebhookService
		// Não precisamos chamar wg.Done() aqui pois é gerenciado pelo service
	}()

	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	logger.Debug("Rate limit cleanup iniciado")
	defer logger.Debug("Rate limit cleanup finalizado")

	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-ticker.C:
			rl.cleanup()
		}
	}
}

// cleanup remove limiters inativos
func (rl *RateLimit) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	inactiveThreshold := 1 * time.Hour // Remove limiters inativos por mais de 1 hora

	var toRemove []string

	for sessionID, limiter := range rl.limiters {
		limiter.mu.Lock()
		if now.Sub(limiter.lastRefill) > inactiveThreshold {
			toRemove = append(toRemove, sessionID)
		}
		limiter.mu.Unlock()
	}

	for _, sessionID := range toRemove {
		delete(rl.limiters, sessionID)
	}

	if len(toRemove) > 0 {
		logger.Debug("Rate limit cleanup removeu %d limiters inativos", len(toRemove))
	}
}

// GetStats retorna estatísticas do rate limiter
func (rl *RateLimit) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	stats := make(map[string]interface{})
	sessionStats := make(map[string]interface{})

	for sessionID, limiter := range rl.limiters {
		limiter.mu.Lock()
		sessionStats[sessionID] = map[string]interface{}{
			"tokens":      limiter.tokens,
			"last_refill": limiter.lastRefill,
		}
		limiter.mu.Unlock()
	}

	stats["sessions"] = sessionStats
	stats["total_sessions"] = len(rl.limiters)
	stats["config"] = map[string]interface{}{
		"requests_per_second": rl.config.RequestsPerSecond,
		"burst_size":          rl.config.BurstSize,
		"cleanup_interval":    rl.config.CleanupInterval,
	}

	return stats
}

// Reset reseta o rate limiter para uma sessão
func (rl *RateLimit) Reset(sessionID string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.limiters[sessionID]; exists {
		limiter.mu.Lock()
		limiter.tokens = rl.config.BurstSize
		limiter.lastRefill = time.Now()
		limiter.mu.Unlock()
		logger.Info("Rate limiter resetado para sessão %s", sessionID)
	}
}

// ResetAll reseta todos os rate limiters
func (rl *RateLimit) ResetAll() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for _, limiter := range rl.limiters {
		limiter.mu.Lock()
		limiter.tokens = rl.config.BurstSize
		limiter.lastRefill = now
		limiter.mu.Unlock()
	}

	logger.Info("Todos os rate limiters foram resetados")
}

// Stop para o rate limiter
func (rl *RateLimit) Stop() {
	rl.cancel()
}
