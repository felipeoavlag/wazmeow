package webhook

import (
	"sync"
	"time"
)

// WebhookMetrics coleta métricas do sistema de webhooks
type WebhookMetrics struct {
	// Contadores
	totalAttempts     int64
	totalSuccess      int64
	totalFailed       int64
	totalQueued       int64
	totalQueueFull    int64
	totalRateLimited  int64
	totalCBOpen       int64

	// Latência
	totalLatency    time.Duration
	minLatency      time.Duration
	maxLatency      time.Duration
	latencyCount    int64

	// Timestamps
	startTime       time.Time
	lastEventTime   time.Time

	mu sync.RWMutex
}

// NewWebhookMetrics cria uma nova instância de WebhookMetrics
func NewWebhookMetrics() *WebhookMetrics {
	now := time.Now()
	return &WebhookMetrics{
		startTime:     now,
		lastEventTime: now,
		minLatency:    time.Duration(0),
		maxLatency:    time.Duration(0),
	}
}

// IncrementAttempts incrementa o contador de tentativas
func (wm *WebhookMetrics) IncrementAttempts() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.totalAttempts++
	wm.lastEventTime = time.Now()
}

// IncrementSuccess incrementa o contador de sucessos
func (wm *WebhookMetrics) IncrementSuccess() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.totalSuccess++
	wm.lastEventTime = time.Now()
}

// IncrementFailed incrementa o contador de falhas
func (wm *WebhookMetrics) IncrementFailed() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.totalFailed++
	wm.lastEventTime = time.Now()
}

// IncrementQueued incrementa o contador de eventos enfileirados
func (wm *WebhookMetrics) IncrementQueued() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.totalQueued++
	wm.lastEventTime = time.Now()
}

// IncrementQueueFull incrementa o contador de fila cheia
func (wm *WebhookMetrics) IncrementQueueFull() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.totalQueueFull++
	wm.lastEventTime = time.Now()
}

// IncrementRateLimited incrementa o contador de rate limit
func (wm *WebhookMetrics) IncrementRateLimited() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.totalRateLimited++
	wm.lastEventTime = time.Now()
}

// IncrementCircuitBreakerOpen incrementa o contador de circuit breaker aberto
func (wm *WebhookMetrics) IncrementCircuitBreakerOpen() {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.totalCBOpen++
	wm.lastEventTime = time.Now()
}

// RecordLatency registra a latência de uma requisição
func (wm *WebhookMetrics) RecordLatency(latency time.Duration) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	
	wm.totalLatency += latency
	wm.latencyCount++
	
	if wm.minLatency == 0 || latency < wm.minLatency {
		wm.minLatency = latency
	}
	
	if latency > wm.maxLatency {
		wm.maxLatency = latency
	}
	
	wm.lastEventTime = time.Now()
}

// GetStats retorna todas as estatísticas
func (wm *WebhookMetrics) GetStats() map[string]interface{} {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	uptime := time.Since(wm.startTime)
	
	stats := map[string]interface{}{
		"uptime": map[string]interface{}{
			"start_time":      wm.startTime,
			"uptime_seconds":  uptime.Seconds(),
			"last_event_time": wm.lastEventTime,
		},
		"counters": map[string]interface{}{
			"total_attempts":      wm.totalAttempts,
			"total_success":       wm.totalSuccess,
			"total_failed":        wm.totalFailed,
			"total_queued":        wm.totalQueued,
			"total_queue_full":    wm.totalQueueFull,
			"total_rate_limited":  wm.totalRateLimited,
			"total_cb_open":       wm.totalCBOpen,
		},
		"rates": map[string]interface{}{
			"success_rate":     wm.calculateSuccessRate(),
			"failure_rate":     wm.calculateFailureRate(),
			"events_per_second": wm.calculateEventsPerSecond(uptime),
		},
		"latency": map[string]interface{}{
			"min_latency_ms":     wm.minLatency.Milliseconds(),
			"max_latency_ms":     wm.maxLatency.Milliseconds(),
			"avg_latency_ms":     wm.calculateAverageLatency().Milliseconds(),
			"total_latency_ms":   wm.totalLatency.Milliseconds(),
			"latency_count":      wm.latencyCount,
		},
	}

	return stats
}

// calculateSuccessRate calcula a taxa de sucesso
func (wm *WebhookMetrics) calculateSuccessRate() float64 {
	if wm.totalAttempts == 0 {
		return 0.0
	}
	return float64(wm.totalSuccess) / float64(wm.totalAttempts) * 100.0
}

// calculateFailureRate calcula a taxa de falha
func (wm *WebhookMetrics) calculateFailureRate() float64 {
	if wm.totalAttempts == 0 {
		return 0.0
	}
	return float64(wm.totalFailed) / float64(wm.totalAttempts) * 100.0
}

// calculateEventsPerSecond calcula eventos por segundo
func (wm *WebhookMetrics) calculateEventsPerSecond(uptime time.Duration) float64 {
	if uptime.Seconds() == 0 {
		return 0.0
	}
	return float64(wm.totalAttempts) / uptime.Seconds()
}

// calculateAverageLatency calcula a latência média
func (wm *WebhookMetrics) calculateAverageLatency() time.Duration {
	if wm.latencyCount == 0 {
		return 0
	}
	return time.Duration(int64(wm.totalLatency) / wm.latencyCount)
}

// Reset reseta todas as métricas
func (wm *WebhookMetrics) Reset() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	now := time.Now()
	wm.totalAttempts = 0
	wm.totalSuccess = 0
	wm.totalFailed = 0
	wm.totalQueued = 0
	wm.totalQueueFull = 0
	wm.totalRateLimited = 0
	wm.totalCBOpen = 0
	wm.totalLatency = 0
	wm.minLatency = 0
	wm.maxLatency = 0
	wm.latencyCount = 0
	wm.startTime = now
	wm.lastEventTime = now
}

// GetSummary retorna um resumo das métricas
func (wm *WebhookMetrics) GetSummary() map[string]interface{} {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	return map[string]interface{}{
		"total_attempts":   wm.totalAttempts,
		"total_success":    wm.totalSuccess,
		"total_failed":     wm.totalFailed,
		"success_rate":     wm.calculateSuccessRate(),
		"avg_latency_ms":   wm.calculateAverageLatency().Milliseconds(),
		"uptime_seconds":   time.Since(wm.startTime).Seconds(),
	}
}
