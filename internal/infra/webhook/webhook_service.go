// Package webhook provides webhook delivery services with reliability features
// including circuit breakers, rate limiting, retries, and metrics.
package webhook

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"wazmeow/internal/config"
	"wazmeow/pkg/logger"
)

// =============================================================================
// WEBHOOK SERVICE AND WORKERS
// =============================================================================

// WebhookService gerencia o envio de webhooks
type WebhookService struct {
	config         *config.WebhookConfig
	httpClient     *http.Client
	eventQueue     chan *WebhookEvent
	workers        []*WebhookWorker
	dispatcher     *WebhookDispatcher
	retryManager   *RetryManager
	metrics        *WebhookMetrics
	circuitBreaker *CircuitBreaker
	rateLimit      *RateLimit
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	mu             sync.RWMutex
	running        bool
}

// WebhookWorker representa um worker para processar webhooks
type WebhookWorker struct {
	id      int
	service *WebhookService
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewWebhookService cria uma nova instância do WebhookService
func NewWebhookService(cfg *config.WebhookConfig) *WebhookService {
	ctx, cancel := context.WithCancel(context.Background())

	// Configurar HTTP client com timeout
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	service := &WebhookService{
		config:     cfg,
		httpClient: httpClient,
		eventQueue: make(chan *WebhookEvent, cfg.QueueSize),
		ctx:        ctx,
		cancel:     cancel,
		running:    false,
	}

	// Inicializar componentes
	service.dispatcher = NewWebhookDispatcher(service)
	service.retryManager = NewRetryManager(service)
	service.metrics = NewWebhookMetrics()
	service.circuitBreaker = NewCircuitBreaker(&cfg.CircuitBreaker)
	service.rateLimit = NewRateLimit(&cfg.RateLimit)

	return service
}

// Start inicia o serviço de webhook
func (ws *WebhookService) Start() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.running {
		return fmt.Errorf("webhook service já está rodando")
	}

	logger.Info("Iniciando WebhookService com %d workers", ws.config.Workers)

	// Iniciar workers
	ws.workers = make([]*WebhookWorker, ws.config.Workers)
	for i := 0; i < ws.config.Workers; i++ {
		worker := &WebhookWorker{
			id:      i,
			service: ws,
		}
		worker.ctx, worker.cancel = context.WithCancel(ws.ctx)
		ws.workers[i] = worker

		ws.wg.Add(1)
		go worker.run()
	}

	// Iniciar retry manager
	ws.wg.Add(1)
	go ws.retryManager.Start()

	// Iniciar rate limit cleanup
	ws.wg.Add(1)
	go ws.rateLimit.StartCleanup()

	ws.running = true
	logger.Info("WebhookService iniciado com sucesso")

	return nil
}

// Stop para o serviço de webhook
func (ws *WebhookService) Stop() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.running {
		return nil
	}

	logger.Info("Parando WebhookService...")

	// Cancelar contexto principal (para todos os workers)
	ws.cancel()

	// Parar retry manager
	if ws.retryManager != nil {
		ws.retryManager.Stop()
	}

	// Aguardar workers com timeout
	done := make(chan struct{}, 1)
	go func() {
		ws.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		// Workers terminaram
	case <-time.After(2 * time.Second):
		logger.Warn("Timeout aguardando workers")
	}

	// Fechar canal se ainda aberto
	select {
	case <-ws.eventQueue:
	default:
		close(ws.eventQueue)
	}

	ws.running = false
	logger.Info("WebhookService parado")
	return nil
}

// SendEvent envia um evento via webhook
func (ws *WebhookService) SendEvent(event *WebhookEvent) error {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	if !ws.running {
		return fmt.Errorf("webhook service não está rodando")
	}

	// Verificar rate limit
	if !ws.rateLimit.Allow(event.SessionID) {
		ws.metrics.IncrementRateLimited()
		return fmt.Errorf("rate limit excedido para sessão %s", event.SessionID)
	}

	// Verificar circuit breaker
	if !ws.circuitBreaker.CanExecute(event.URL) {
		ws.metrics.IncrementCircuitBreakerOpen()
		return fmt.Errorf("circuit breaker aberto para URL %s", event.URL)
	}

	// Adicionar timestamp se não definido
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().Unix()
	}

	// Tentar enviar para a fila
	select {
	case ws.eventQueue <- event:
		ws.metrics.IncrementQueued()
		return nil
	default:
		ws.metrics.IncrementQueueFull()
		return fmt.Errorf("fila de webhooks está cheia")
	}
}

// GetMetrics retorna as métricas do webhook service
func (ws *WebhookService) GetMetrics() *WebhookMetrics {
	return ws.metrics
}

// IsRunning verifica se o serviço está rodando
func (ws *WebhookService) IsRunning() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.running
}

// run executa o loop principal do worker
func (w *WebhookWorker) run() {
	defer w.service.wg.Done()

	logger.Debug("Worker %d iniciado", w.id)
	defer logger.Debug("Worker %d finalizado", w.id)

	for {
		select {
		case <-w.ctx.Done():
			return
		case event, ok := <-w.service.eventQueue:
			if !ok {
				return
			}
			w.processEvent(event)
		}
	}
}

// processEvent processa um evento de webhook
func (w *WebhookWorker) processEvent(event *WebhookEvent) {
	logger.Debug("Worker %d processando evento %s para sessão %s", w.id, event.Type, event.SessionID)

	// Tentar enviar o webhook
	err := w.service.dispatcher.Send(event)
	if err != nil {
		logger.Error("Erro ao enviar webhook (worker %d): %v", w.id, err)

		// Adicionar para retry se não excedeu o limite
		if event.Retries < w.service.config.MaxRetries {
			event.Retries++
			w.service.retryManager.AddForRetry(event)
		} else {
			logger.Error("Webhook descartado após %d tentativas: %s", event.Retries, event.ID)
			w.service.metrics.IncrementFailed()
		}
	} else {
		logger.Debug("Webhook enviado com sucesso (worker %d): %s", w.id, event.ID)
		w.service.metrics.IncrementSuccess()
	}
}

// =============================================================================
// CIRCUIT BREAKER
// =============================================================================

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

	if !exists {
		// Criar circuit breaker se não existir
		cb.mu.Lock()
		breaker = &URLCircuitBreaker{
			url:   url,
			state: StateClosed,
		}
		cb.breakers[url] = breaker
		cb.mu.Unlock()
	}

	breaker.recordFailure(cb.config)
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
	stats["config"] = map[string]interface{}{
		"max_failures":        cb.config.MaxFailures,
		"reset_timeout":       cb.config.ResetTimeout,
		"half_open_max_calls": cb.config.HalfOpenMaxCalls,
	}

	return stats
}

// Reset reseta o circuit breaker para uma URL
func (cb *CircuitBreaker) Reset(url string) {
	cb.mu.RLock()
	breaker, exists := cb.breakers[url]
	cb.mu.RUnlock()

	if exists {
		breaker.mu.Lock()
		breaker.state = StateClosed
		breaker.failures = 0
		breaker.halfOpenCalls = 0
		breaker.mu.Unlock()
		logger.Info("Circuit breaker resetado para URL %s", url)
	}
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
func (ucb *URLCircuitBreaker) recordSuccess(_ *config.CircuitBreakerConfig) {
	ucb.mu.Lock()
	defer ucb.mu.Unlock()

	switch ucb.state {
	case StateHalfOpen:
		// Se conseguiu sucesso em half-open, voltar para closed
		ucb.state = StateClosed
		ucb.failures = 0
		ucb.halfOpenCalls = 0
		logger.Info("Circuit breaker para %s voltou para closed após sucesso", ucb.url)
	case StateClosed:
		// Reset failures counter on success
		ucb.failures = 0
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

// =============================================================================
// RATE LIMITING
// =============================================================================

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

// =============================================================================
// RETRY MANAGER
// =============================================================================

// RetryManager gerencia o sistema de retry de webhooks
type RetryManager struct {
	service    *WebhookService
	retryQueue chan *WebhookEvent
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
	running    bool
}

// NewRetryManager cria uma nova instância do RetryManager
func NewRetryManager(service *WebhookService) *RetryManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &RetryManager{
		service:    service,
		retryQueue: make(chan *WebhookEvent, service.config.QueueSize),
		ctx:        ctx,
		cancel:     cancel,
		running:    false,
	}
}

// Start inicia o retry manager
func (rm *RetryManager) Start() {
	defer rm.service.wg.Done()

	rm.mu.Lock()
	if rm.running {
		rm.mu.Unlock()
		return
	}
	rm.running = true
	rm.mu.Unlock()

	logger.Info("RetryManager iniciado")
	defer logger.Info("RetryManager finalizado")

	// Iniciar worker de retry
	rm.wg.Add(1)
	go rm.retryWorker()

	// Aguardar cancelamento
	<-rm.ctx.Done()

	// Parar workers
	rm.cancel()
	rm.wg.Wait()

	rm.mu.Lock()
	rm.running = false
	rm.mu.Unlock()
}

// Stop para o retry manager
func (rm *RetryManager) Stop() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if !rm.running {
		return
	}

	logger.Info("Parando RetryManager...")
	rm.cancel()

	// Aguardar workers terminarem
	rm.wg.Wait()

	// Fechar canal apenas se ainda estiver aberto
	select {
	case <-rm.retryQueue:
	default:
		close(rm.retryQueue)
	}

	rm.running = false
	logger.Info("RetryManager parado com sucesso")
}

// AddForRetry adiciona um evento para retry
func (rm *RetryManager) AddForRetry(event *WebhookEvent) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if !rm.running {
		logger.Warn("RetryManager não está rodando, descartando evento %s", event.ID)
		return
	}

	// Calcular delay baseado no número de tentativas
	delay := rm.calculateRetryDelay(event.Retries)

	logger.Debug("Agendando retry para evento %s em %v (tentativa %d)",
		event.ID, delay, event.Retries)

	// Agendar retry
	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-timer.C:
			select {
			case rm.retryQueue <- event:
				logger.Debug("Evento %s adicionado à fila de retry", event.ID)
			default:
				logger.Error("Fila de retry cheia, descartando evento %s", event.ID)
				rm.service.metrics.IncrementFailed()
			}
		case <-rm.ctx.Done():
			return
		}
	}()
}

// retryWorker processa eventos da fila de retry
func (rm *RetryManager) retryWorker() {
	defer rm.wg.Done()

	logger.Debug("Retry worker iniciado")
	defer logger.Debug("Retry worker finalizado")

	for {
		select {
		case <-rm.ctx.Done():
			return
		case event, ok := <-rm.retryQueue:
			if !ok {
				return
			}
			rm.processRetry(event)
		}
	}
}

// processRetry processa um evento de retry
func (rm *RetryManager) processRetry(event *WebhookEvent) {
	logger.Debug("Processando retry para evento %s (tentativa %d)",
		event.ID, event.Retries)

	// Verificar se ainda não excedeu o limite de retries
	if event.Retries > rm.service.config.MaxRetries {
		logger.Error("Evento %s excedeu limite de retries (%d), descartando",
			event.ID, rm.service.config.MaxRetries)
		rm.service.metrics.IncrementFailed()
		return
	}

	// Tentar enviar novamente
	err := rm.service.dispatcher.Send(event)
	if err != nil {
		logger.Error("Retry falhou para evento %s: %v", event.ID, err)

		// Incrementar contador de retries e tentar novamente se não excedeu limite
		if event.Retries < rm.service.config.MaxRetries {
			event.Retries++
			rm.AddForRetry(event)
		} else {
			logger.Error("Evento %s descartado após %d tentativas",
				event.ID, event.Retries)
			rm.service.metrics.IncrementFailed()
		}
	} else {
		logger.Info("Retry bem-sucedido para evento %s após %d tentativas",
			event.ID, event.Retries)
		rm.service.metrics.IncrementSuccess()
	}
}

// calculateRetryDelay calcula o delay para retry com backoff exponencial
func (rm *RetryManager) calculateRetryDelay(attempt int) time.Duration {
	baseDelay := rm.service.config.RetryDelay

	// Backoff exponencial: delay * 2^attempt
	multiplier := 1 << uint(attempt-1) // 2^(attempt-1)
	if multiplier > 16 {
		multiplier = 16 // Limitar o multiplicador máximo
	}

	delay := time.Duration(multiplier) * baseDelay

	// Limitar delay máximo a 10 minutos
	maxDelay := 10 * time.Minute
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// IsRunning verifica se o retry manager está rodando
func (rm *RetryManager) IsRunning() bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.running
}

// GetQueueSize retorna o tamanho atual da fila de retry
func (rm *RetryManager) GetQueueSize() int {
	return len(rm.retryQueue)
}
