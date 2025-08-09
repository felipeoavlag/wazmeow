package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"wazmeow/pkg/logger"
)

// WebhookDispatcher é responsável pelo envio HTTP de webhooks
type WebhookDispatcher struct {
	service *WebhookService
}

// NewWebhookDispatcher cria uma nova instância do WebhookDispatcher
func NewWebhookDispatcher(service *WebhookService) *WebhookDispatcher {
	return &WebhookDispatcher{
		service: service,
	}
}

// Send envia um webhook via HTTP POST
func (wd *WebhookDispatcher) Send(event *WebhookEvent) error {
	if event.URL == "" {
		return fmt.Errorf("URL do webhook não definida")
	}

	// Verificar circuit breaker
	if !wd.service.circuitBreaker.CanExecute(event.URL) {
		wd.service.circuitBreaker.RecordFailure(event.URL)
		return fmt.Errorf("circuit breaker aberto para URL %s", event.URL)
	}

	// Serializar evento para JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		wd.service.circuitBreaker.RecordFailure(event.URL)
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	// Criar request HTTP
	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		event.URL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		wd.service.circuitBreaker.RecordFailure(event.URL)
		return fmt.Errorf("erro ao criar request HTTP: %w", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WazMeow-Webhook/1.0")
	req.Header.Set("X-Webhook-Event", event.Type)
	req.Header.Set("X-Webhook-Session", event.SessionID)
	req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", event.Timestamp))
	req.Header.Set("X-Webhook-ID", event.ID)

	// Registrar início da tentativa
	startTime := time.Now()
	wd.service.metrics.IncrementAttempts()

	// Enviar request
	resp, err := wd.service.httpClient.Do(req)
	if err != nil {
		duration := time.Since(startTime)
		wd.service.metrics.RecordLatency(duration)
		wd.service.circuitBreaker.RecordFailure(event.URL)
		
		logger.Error("Erro ao enviar webhook para %s: %v", event.URL, err)
		return fmt.Errorf("erro ao enviar webhook: %w", err)
	}
	defer resp.Body.Close()

	// Registrar latência
	duration := time.Since(startTime)
	wd.service.metrics.RecordLatency(duration)

	// Verificar status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Ler corpo da resposta para log
		body, _ := io.ReadAll(resp.Body)
		
		wd.service.circuitBreaker.RecordFailure(event.URL)
		
		logger.Error("Webhook retornou status %d para %s: %s", 
			resp.StatusCode, event.URL, string(body))
		
		return fmt.Errorf("webhook retornou status %d", resp.StatusCode)
	}

	// Sucesso
	wd.service.circuitBreaker.RecordSuccess(event.URL)
	
	logger.Debug("Webhook enviado com sucesso para %s (status: %d, latência: %v)", 
		event.URL, resp.StatusCode, duration)

	return nil
}

// SendWithRetry envia um webhook com retry automático
func (wd *WebhookDispatcher) SendWithRetry(event *WebhookEvent) error {
	var lastErr error
	
	for attempt := 0; attempt <= wd.service.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Aguardar antes do retry
			delay := wd.calculateRetryDelay(attempt)
			logger.Debug("Aguardando %v antes do retry %d para webhook %s", 
				delay, attempt, event.ID)
			time.Sleep(delay)
		}

		err := wd.Send(event)
		if err == nil {
			if attempt > 0 {
				logger.Info("Webhook %s enviado com sucesso após %d tentativas", 
					event.ID, attempt+1)
			}
			return nil
		}

		lastErr = err
		logger.Warn("Tentativa %d falhou para webhook %s: %v", 
			attempt+1, event.ID, err)
	}

	logger.Error("Webhook %s falhou após %d tentativas: %v", 
		event.ID, wd.service.config.MaxRetries+1, lastErr)
	
	return fmt.Errorf("webhook falhou após %d tentativas: %w", 
		wd.service.config.MaxRetries+1, lastErr)
}

// calculateRetryDelay calcula o delay para retry com backoff exponencial
func (wd *WebhookDispatcher) calculateRetryDelay(attempt int) time.Duration {
	baseDelay := wd.service.config.RetryDelay
	
	// Backoff exponencial: delay * 2^attempt
	multiplier := 1 << uint(attempt-1) // 2^(attempt-1)
	if multiplier > 8 {
		multiplier = 8 // Limitar o multiplicador máximo
	}
	
	delay := time.Duration(multiplier) * baseDelay
	
	// Limitar delay máximo a 5 minutos
	maxDelay := 5 * time.Minute
	if delay > maxDelay {
		delay = maxDelay
	}
	
	return delay
}

// ValidateWebhookURL valida se uma URL de webhook é válida
func (wd *WebhookDispatcher) ValidateWebhookURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL não pode estar vazia")
	}

	// Criar request de teste
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return fmt.Errorf("URL inválida: %w", err)
	}

	// Configurar timeout menor para validação
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Tentar fazer request HEAD
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao conectar com webhook: %w", err)
	}
	defer resp.Body.Close()

	// Aceitar qualquer status code para validação
	// O importante é que a URL seja acessível
	logger.Debug("Webhook URL %s validada (status: %d)", url, resp.StatusCode)
	
	return nil
}

// TestWebhook envia um webhook de teste
func (wd *WebhookDispatcher) TestWebhook(url, sessionID string) error {
	testEvent := &WebhookEvent{
		ID:        fmt.Sprintf("test_%d", time.Now().UnixNano()),
		Type:      "test",
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
		URL:       url,
		Data: map[string]interface{}{
			"message": "Este é um webhook de teste do WazMeow",
			"test":    true,
		},
	}

	return wd.Send(testEvent)
}
