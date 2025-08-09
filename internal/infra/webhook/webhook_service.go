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
