package webhook

import (
	"context"
	"sync"
	"time"

	"wazmeow/pkg/logger"
)

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
