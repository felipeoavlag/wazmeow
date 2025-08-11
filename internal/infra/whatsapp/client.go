package whatsapp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/webhook"
	"wazmeow/pkg/logger"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// WhatsAppClient representa um cliente WhatsApp com funcionalidades completas
type WhatsAppClient struct {
	client        *whatsmeow.Client
	sessionID     string
	sessionRepo   repository.SessionRepository
	eventHandlers map[uint32]func(interface{})
	mutex         sync.RWMutex

	// Campos para QR code management
	qrChannel    chan string                    // Canal interno para compatibilidade
	nativeQRChan <-chan whatsmeow.QRChannelItem // Canal nativo do whatsmeow
	qrCancelFunc context.CancelFunc             // Para cancelar QR loop
	isQRActive   bool                           // Flag se QR está ativo
	isConnecting bool

	// Webhook components
	webhookService  *webhook.WebhookService
	eventSerializer *webhook.EventSerializer
	eventFilter     *webhook.EventFilter
}

// NewWhatsAppClient cria uma nova instância do cliente WhatsApp
func NewWhatsAppClient(client *whatsmeow.Client, sessionID string, sessionRepo repository.SessionRepository) *WhatsAppClient {
	wac := &WhatsAppClient{
		client:        client,
		sessionID:     sessionID,
		sessionRepo:   sessionRepo,
		eventHandlers: make(map[uint32]func(interface{})),
		qrChannel:     make(chan string, 1),
		isConnecting:  false,
	}

	// Inicializar componentes de webhook
	wac.eventSerializer = webhook.NewEventSerializer()
	wac.eventFilter = webhook.NewEventFilter()

	// Configurar event handlers padrão
	wac.setupDefaultEventHandlers()

	return wac
}

// SetWebhookService define o serviço de webhook para o cliente
func (wac *WhatsAppClient) SetWebhookService(webhookService *webhook.WebhookService) {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()
	wac.webhookService = webhookService
}

// Connect estabelece conexão com o WhatsApp (método legado)
func (wac *WhatsAppClient) Connect() error {
	// Usar ConnectWithQR sem timeout para manter loop QR ativo
	return wac.ConnectWithQR(context.Background())
}

// ConnectWithQR estabelece conexão com o WhatsApp usando GetQRChannel
func (wac *WhatsAppClient) ConnectWithQR(ctx context.Context) error {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()

	if wac.isConnecting {
		return fmt.Errorf("cliente já está tentando conectar")
	}

	if wac.client.IsConnected() {
		return fmt.Errorf("cliente já está conectado")
	}

	wac.isConnecting = true
	defer func() { wac.isConnecting = false }()

	// Atualizar status da sessão para conectando
	if err := wac.updateSessionStatus(entity.StatusConnecting); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}

	// Verificar se já está logado
	if wac.client.Store.ID != nil {
		logger.Info("Sessão %s já está logada (JID: %s), conectando diretamente", wac.sessionID, wac.client.Store.ID.String())

		// Já logado, apenas conectar
		if err := wac.client.Connect(); err != nil {
			wac.updateSessionStatus(entity.StatusDisconnected)
			return fmt.Errorf("erro ao conectar cliente logado: %w", err)
		}

		return nil
	}

	// Não está logado, precisa de QR code
	logger.Info("Sessão %s precisa de autenticação, iniciando processo QR", wac.sessionID)

	// Obter canal QR ANTES de conectar (usar context.Background() como na referência)
	qrChan, err := wac.client.GetQRChannel(context.Background())
	if err != nil {
		// Verificar se é erro de já estar logado
		if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			wac.updateSessionStatus(entity.StatusDisconnected)
			return fmt.Errorf("erro ao obter canal QR: %w", err)
		}
		// Se já está logado, apenas conectar
		if err := wac.client.Connect(); err != nil {
			wac.updateSessionStatus(entity.StatusDisconnected)
			return fmt.Errorf("erro ao conectar cliente já logado: %w", err)
		}
		return nil
	}

	// Conectar cliente DEPOIS de obter QR channel (como na referência)
	if err := wac.client.Connect(); err != nil {
		wac.updateSessionStatus(entity.StatusDisconnected)
		return fmt.Errorf("erro ao conectar cliente: %w", err)
	}

	// Armazenar canal QR
	wac.nativeQRChan = qrChan
	wac.isQRActive = true

	// Iniciar loop QR em goroutine com context independente
	qrCtx, qrCancel := context.WithCancel(context.Background())
	wac.qrCancelFunc = qrCancel

	go wac.processQREvents(qrCtx, qrChan)

	return nil
}

// ConnectDirect conecta diretamente ao WhatsApp (para sessões já autenticadas)
func (wac *WhatsAppClient) ConnectDirect() error {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()

	if wac.isConnecting {
		return fmt.Errorf("cliente já está tentando conectar")
	}

	if wac.client.IsConnected() {
		return fmt.Errorf("cliente já está conectado")
	}

	// Verificar se já está logado (tem DeviceJID)
	if wac.client.Store.ID == nil {
		return fmt.Errorf("sessão não está autenticada, use ConnectWithQR para autenticar")
	}

	wac.isConnecting = true
	defer func() { wac.isConnecting = false }()

	logger.Info("Conectando diretamente sessão %s (JID: %s)", wac.sessionID, wac.client.Store.ID.String())

	// Atualizar status da sessão para conectando
	if err := wac.updateSessionStatus(entity.StatusConnecting); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}

	// Conectar diretamente (sem QR code)
	if err := wac.client.Connect(); err != nil {
		wac.updateSessionStatus(entity.StatusDisconnected)
		return fmt.Errorf("erro ao conectar cliente já autenticado: %w", err)
	}

	logger.Info("Sessão %s conectada diretamente com sucesso", wac.sessionID)
	return nil
}

// Disconnect desconecta do WhatsApp
func (wac *WhatsAppClient) Disconnect() {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()

	if wac.client != nil {
		wac.client.Disconnect()
	}

	// Atualizar status da sessão
	if err := wac.updateSessionStatus(entity.StatusDisconnected); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}
}

// IsConnected verifica se está conectado
func (wac *WhatsAppClient) IsConnected() bool {
	wac.mutex.RLock()
	defer wac.mutex.RUnlock()

	if wac.client == nil {
		return false
	}
	return wac.client.IsConnected()
}

// IsLoggedIn verifica se está logado
func (wac *WhatsAppClient) IsLoggedIn() bool {
	wac.mutex.RLock()
	defer wac.mutex.RUnlock()

	if wac.client == nil {
		return false
	}
	return wac.client.IsLoggedIn()
}

// Logout faz logout da sessão
func (wac *WhatsAppClient) Logout(ctx context.Context) error {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()

	if wac.client == nil {
		return fmt.Errorf("cliente não inicializado")
	}

	if err := wac.client.Logout(ctx); err != nil {
		return fmt.Errorf("erro ao fazer logout: %w", err)
	}

	// Atualizar status da sessão para desconectado (logout = desconectado)
	if err := wac.updateSessionStatus(entity.StatusDisconnected); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}

	return nil
}

// PairPhone emparelha um telefone
func (wac *WhatsAppClient) PairPhone(ctx context.Context, phone string, showPushNotification bool, clientType string, clientName string) (string, error) {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()

	if wac.client == nil {
		return "", fmt.Errorf("cliente não inicializado")
	}

	// Limpar caracteres especiais do telefone
	phone = strings.ReplaceAll(phone, "+", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Emparelhar telefone
	code, err := wac.client.PairPhone(ctx, phone, showPushNotification, whatsmeow.PairClientChrome, clientName)
	if err != nil {
		return "", fmt.Errorf("erro ao emparelhar telefone: %w", err)
	}

	return code, nil
}

// AddEventHandler adiciona um handler de eventos
func (wac *WhatsAppClient) AddEventHandler(handler func(interface{})) string {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()

	handlerID := wac.client.AddEventHandler(handler)
	wac.eventHandlers[handlerID] = handler

	return strconv.FormatUint(uint64(handlerID), 10)
}

// RemoveEventHandler remove um handler de eventos
func (wac *WhatsAppClient) RemoveEventHandler(handlerID string) {
	wac.mutex.Lock()
	defer wac.mutex.Unlock()

	// Converter string para uint32
	id, err := strconv.ParseUint(handlerID, 10, 32)
	if err != nil {
		logger.Error("Erro ao converter handler ID: %v", err)
		return
	}

	handlerIDUint := uint32(id)
	wac.client.RemoveEventHandler(handlerIDUint)
	delete(wac.eventHandlers, handlerIDUint)
}

// GetQRChannel retorna o canal para receber QR codes
func (wac *WhatsAppClient) GetQRChannel() <-chan string {
	return wac.qrChannel
}

// GetClient retorna o cliente whatsmeow nativo
func (wac *WhatsAppClient) GetClient() *whatsmeow.Client {
	wac.mutex.RLock()
	defer wac.mutex.RUnlock()
	return wac.client
}

// IsQRActive verifica se o QR está ativo
func (wac *WhatsAppClient) IsQRActive() bool {
	wac.mutex.RLock()
	defer wac.mutex.RUnlock()
	return wac.isQRActive
}

// processQREvents processa eventos do canal QR em loop contínuo
func (wac *WhatsAppClient) processQREvents(ctx context.Context, qrChan <-chan whatsmeow.QRChannelItem) {
	logger.Info("Iniciando loop de processamento QR para sessão %s", wac.sessionID)

	defer func() {
		wac.mutex.Lock()
		wac.isQRActive = false
		wac.qrCancelFunc = nil
		wac.mutex.Unlock()
		logger.Info("Loop QR finalizado para sessão %s", wac.sessionID)
	}()

	for {
		select {
		case evt := <-qrChan:
			switch evt.Event {
			case "code":
				logger.Info("Evento QR 'code' recebido para sessão %s", wac.sessionID)
				wac.handleQRCode(evt.Code)

			case "timeout":
				logger.Info("Evento QR 'timeout' recebido para sessão %s", wac.sessionID)
				wac.handleQRTimeout()

			case "success":
				logger.Info("Evento QR 'success' recebido para sessão %s", wac.sessionID)
				wac.handleQRSuccess()
				return // Sair do loop - autenticação bem-sucedida

			case "error":
				logger.Error("Evento QR 'error' recebido para sessão %s: %v", wac.sessionID, evt.Error)
				wac.handleQRError(evt.Error)
				return // Sair do loop - erro

			default:
				logger.Info("Evento QR desconhecido '%s' para sessão %s", evt.Event, wac.sessionID)
			}

		case <-ctx.Done():
			logger.Info("Context cancelado, finalizando loop QR para sessão %s", wac.sessionID)
			return // Context cancelado
		}
	}
}

// updateSessionStatus atualiza o status da sessão no banco de dados
func (wac *WhatsAppClient) updateSessionStatus(status entity.SessionStatus) error {
	session, err := wac.sessionRepo.GetByID(wac.sessionID)
	if err != nil {
		return err
	}

	session.Status = status
	session.UpdatedAt = time.Now()

	return wac.sessionRepo.Update(session)
}

// setupDefaultEventHandlers configura os event handlers padrão
func (wac *WhatsAppClient) setupDefaultEventHandlers() {
	wac.client.AddEventHandler(func(evt interface{}) {
		switch e := evt.(type) {
		// Eventos de conectividade
		case *events.Connected:
			wac.handleConnected(e)
		case *events.Disconnected:
			wac.handleDisconnected(e)
		case *events.LoggedOut:
			wac.handleLoggedOut(e)
		case *events.QR:
			wac.handleQR(e)
		case *events.PairSuccess:
			wac.handlePairSuccess(e)

		// Eventos de mensagem
		case *events.Message:
			wac.handleMessage(e)
		case *events.Receipt:
			wac.handleReceipt(e)

		// Eventos de presença
		case *events.Presence:
			wac.handlePresence(e)
		case *events.ChatPresence:
			wac.handleChatPresence(e)

		// Eventos de grupo
		case *events.GroupInfo:
			wac.handleGroupInfo(e)

		// Eventos de mídia e perfil
		case *events.Picture:
			wac.handlePicture(e)

		// Eventos de histórico
		case *events.HistorySync:
			wac.handleHistorySync(e)

		// Eventos de chamada
		case *events.CallOffer:
			wac.handleCallOffer(e)
		case *events.CallAccept:
			wac.handleCallAccept(e)
		case *events.CallTerminate:
			wac.handleCallTerminate(e)

		// Eventos de newsletter
		case *events.NewsletterJoin:
			wac.handleNewsletterJoin(e)
		case *events.NewsletterLeave:
			wac.handleNewsletterLeave(e)
		case *events.NewsletterMuteChange:
			wac.handleNewsletterMuteChange(e)

		// Outros eventos importantes
		case *events.BlocklistChange:
			wac.handleBlocklistChange(e)
		case *events.PushName:
			wac.handlePushName(e)
		case *events.BusinessName:
			wac.handleBusinessName(e)
		case *events.JoinedGroup:
			wac.handleJoinedGroup(e)
		case *events.Contact:
			wac.handleContact(e)

		default:
			logger.Debug("Evento não tratado: %T", evt)
			// Enviar evento genérico via webhook
			wac.handleGenericEvent(evt)
		}
	})
}

// handleConnected trata eventos de conexão
func (wac *WhatsAppClient) handleConnected(evt *events.Connected) {
	logger.Info("Sessão %s conectada ao WhatsApp", wac.sessionID)

	// Log de monitoramento do estado
	logger.Info("Estado da sessão %s: conectado=%v, logado=%v, QR_ativo=%v",
		wac.sessionID, true, wac.client.IsLoggedIn(), wac.IsQRActive())

	// Atualizar status da sessão
	if err := wac.updateSessionStatus(entity.StatusConnected); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}

	// Enviar presença disponível apenas se o PushName estiver definido
	if wac.client.Store.PushName != "" {
		if err := wac.client.SendPresence(types.PresenceAvailable); err != nil {
			logger.Warn("Erro ao enviar presença disponível: %v", err)
		}
	} else {
		logger.Debug("PushName não definido, aguardando para enviar presença")
	}

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleDisconnected trata eventos de desconexão
func (wac *WhatsAppClient) handleDisconnected(evt *events.Disconnected) {
	logger.Info("Sessão %s desconectada do WhatsApp", wac.sessionID)

	// Atualizar status da sessão
	if err := wac.updateSessionStatus(entity.StatusDisconnected); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleLoggedOut trata eventos de logout
func (wac *WhatsAppClient) handleLoggedOut(evt *events.LoggedOut) {
	logger.Info("Sessão %s fez logout do WhatsApp. Motivo: %s", wac.sessionID, evt.Reason.String())

	// Atualizar status da sessão para desconectado (logout = desconectado)
	if err := wac.updateSessionStatus(entity.StatusDisconnected); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}
}

// handleQR trata eventos de QR code
func (wac *WhatsAppClient) handleQR(evt *events.QR) {
	logger.Info("QR code gerado para sessão %s", wac.sessionID)

	if len(evt.Codes) > 0 {
		qrCode := evt.Codes[0]

		// Exibir QR code no terminal (útil para desenvolvimento/teste)
		logger.Info("=== QR CODE PARA SESSÃO %s ===", wac.sessionID)
		qrterminal.GenerateHalfBlock(qrCode, qrterminal.L, os.Stdout)
		fmt.Printf("\n📱 Escaneie o QR code acima com o WhatsApp\n")
		fmt.Printf("🔗 Código QR: %s\n", qrCode)
		fmt.Printf("⏰ Sessão: %s\n", wac.sessionID)
		fmt.Printf("=======================================\n\n")

		// Atualizar sessão com QR code
		session, err := wac.sessionRepo.GetByID(wac.sessionID)
		if err != nil {
			logger.Error("Erro ao buscar sessão: %v", err)
			return
		}

		// Atualizar QR code na sessão
		session.UpdatedAt = time.Now()
		if err := wac.sessionRepo.Update(session); err != nil {
			logger.Error("Erro ao atualizar sessão com QR code: %v", err)
		}

		// Enviar QR code pelo canal
		select {
		case wac.qrChannel <- qrCode:
		default:
			// Canal cheio, ignorar
		}

		logger.Info("QR code atualizado para sessão %s", wac.sessionID)
	}
}

// handlePairSuccess trata eventos de emparelhamento bem-sucedido
func (wac *WhatsAppClient) handlePairSuccess(evt *events.PairSuccess) {
	logger.Info("🎉 EMPARELHAMENTO BEM-SUCEDIDO! 🎉")
	logger.Info("Sessão: %s", wac.sessionID)
	logger.Info("JID: %s", evt.ID.String())
	logger.Info("Plataforma: %s", evt.Platform)
	if evt.BusinessName != "" {
		logger.Info("Business: %s", evt.BusinessName)
	}

	// Atualizar sessão com JID do dispositivo
	session, err := wac.sessionRepo.GetByID(wac.sessionID)
	if err != nil {
		logger.Error("Erro ao buscar sessão: %v", err)
		return
	}

	// Atualizar JID do dispositivo e STATUS na sessão
	session.DeviceJID = evt.ID.String()     // JID completo (ex: 554988989314:12@s.whatsapp.net)
	session.Phone = evt.ID.User             // Apenas o número do telefone (ex: 554988989314)
	session.Status = entity.StatusConnected // IMPORTANTE: Atualizar status para connected
	session.UpdatedAt = time.Now()

	if err := wac.sessionRepo.Update(session); err != nil {
		logger.Error("Erro ao atualizar sessão com JID: %v", err)
	} else {
		logger.Info("Sessão %s atualizada com JID: %s (Phone: %s) - Status: connected", wac.sessionID, evt.ID.String(), evt.ID.User)
	}

	fmt.Printf("\n✅ WhatsApp conectado com sucesso!\n")
	fmt.Printf("📱 Sessão: %s\n", wac.sessionID)
	fmt.Printf("🆔 JID: %s\n", evt.ID.String())
	fmt.Printf("=======================================\n\n")
}

// handleMessage trata eventos de mensagem
func (wac *WhatsAppClient) handleMessage(evt *events.Message) {
	logger.Info("Mensagem recebida na sessão %s de %s", wac.sessionID, evt.Info.SourceString())

	// Enviar webhook com evento bruto do whatsmeow
	wac.sendWebhookForEvent(evt)
}

// handleReceipt trata eventos de confirmação de leitura
func (wac *WhatsAppClient) handleReceipt(evt *events.Receipt) {
	logger.Debug("Confirmação de leitura recebida na sessão %s", wac.sessionID)

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handlePresence trata eventos de presença
func (wac *WhatsAppClient) handlePresence(evt *events.Presence) {
	logger.Debug("Presença recebida na sessão %s de %s", wac.sessionID, evt.From.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleChatPresence trata eventos de presença em chat
func (wac *WhatsAppClient) handleChatPresence(evt *events.ChatPresence) {
	logger.Debug("Presença de chat recebida na sessão %s", wac.sessionID)

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// sendWebhook envia dados para webhook usando o sistema completo
func (wac *WhatsAppClient) sendWebhook(data map[string]interface{}) {
	wac.mutex.RLock()
	webhookService := wac.webhookService
	wac.mutex.RUnlock()

	// Se não há serviço de webhook configurado, apenas log
	if webhookService == nil {
		logger.Debug("Webhook service não configurado para sessão %s", wac.sessionID)
		return
	}

	// Buscar configuração da sessão
	session, err := wac.sessionRepo.GetByID(wac.sessionID)
	if err != nil {
		logger.Error("Erro ao buscar sessão para webhook: %v", err)
		return
	}

	// Verificar se há pelo menos um webhook configurado
	if session.WebhookURL == "" && session.Webhook == "" {
		logger.Debug("Nenhum webhook configurado para sessão %s", wac.sessionID)
		return
	}

	// Extrair tipo do evento
	eventType, ok := data["type"].(string)
	if !ok {
		logger.Error("Tipo de evento não encontrado nos dados do webhook")
		return
	}

	// Verificar se deve enviar este evento
	if !wac.eventFilter.ShouldSendEvent(session, eventType) {
		logger.Debug("Evento %s filtrado para sessão %s", eventType, wac.sessionID)
		return
	}

	// Enviar para webhook customizado (se configurado)
	if session.WebhookURL != "" {
		// Verificar se deve enviar este evento para webhook customizado
		if wac.eventFilter.ShouldSendEvent(session, eventType) {
			webhookEvent := &webhook.WebhookEvent{
				ID:        fmt.Sprintf("evt_custom_%s_%d", wac.sessionID, time.Now().UnixNano()),
				Type:      eventType,
				SessionID: wac.sessionID,
				Timestamp: time.Now().Unix(),
				Data:      data,
				URL:       session.WebhookURL,
				Retries:   0,
			}

			err = webhookService.SendEvent(webhookEvent)
			if err != nil {
				logger.Error("Erro ao enviar webhook customizado para sessão %s: %v", wac.sessionID, err)
			} else {
				logger.Debug("Webhook customizado enviado para sessão %s: %s", wac.sessionID, eventType)
			}
		}
	}

	// Enviar para webhook padrão (se configurado) - sempre envia todos os eventos
	if session.Webhook != "" {
		webhookEvent := &webhook.WebhookEvent{
			ID:        fmt.Sprintf("evt_default_%s_%d", wac.sessionID, time.Now().UnixNano()),
			Type:      eventType,
			SessionID: wac.sessionID,
			Timestamp: time.Now().Unix(),
			Data:      data,
			URL:       session.Webhook,
			Retries:   0,
		}

		err = webhookService.SendEvent(webhookEvent)
		if err != nil {
			logger.Error("Erro ao enviar webhook padrão para sessão %s: %v", wac.sessionID, err)
		} else {
			logger.Debug("Webhook padrão enviado para sessão %s: %s", wac.sessionID, eventType)
		}
	}
}

// sendWebhookForEvent envia evento bruto do whatsmeow via webhook
func (wac *WhatsAppClient) sendWebhookForEvent(evt interface{}) {
	wac.mutex.RLock()
	webhookService := wac.webhookService
	eventSerializer := wac.eventSerializer
	wac.mutex.RUnlock()

	logger.Debug("🔍 Tentando enviar webhook para evento %T na sessão %s", evt, wac.sessionID)

	// Se não há serviço de webhook configurado, apenas log
	if webhookService == nil {
		logger.Error("❌ Webhook service não configurado para sessão %s", wac.sessionID)
		return
	}

	logger.Debug("✅ Webhook service encontrado para sessão %s", wac.sessionID)

	// Buscar configuração da sessão
	session, err := wac.sessionRepo.GetByID(wac.sessionID)
	if err != nil {
		logger.Error("❌ Erro ao buscar sessão para webhook: %v", err)
		return
	}

	logger.Debug("✅ Sessão encontrada: %s, WebhookURL: %s", wac.sessionID, session.WebhookURL)

	// Verificar se há pelo menos um webhook configurado
	if session.WebhookURL == "" && session.Webhook == "" {
		logger.Debug("❌ Nenhum webhook configurado para sessão %s", wac.sessionID)
		return
	}

	logger.Debug("✅ Webhooks configurados - Custom: %s, Padrão: %s", session.WebhookURL, session.Webhook)

	// Serializar evento (payload bruto)
	payload, err := eventSerializer.SerializeEvent(wac.sessionID, evt)
	if err != nil {
		logger.Error("❌ Erro ao serializar evento para webhook: %v", err)
		return
	}

	logger.Debug("✅ Evento serializado: %s", payload.Event)

	// Enviar para webhook customizado (se configurado e aprovado pelo filtro)
	if session.WebhookURL != "" {
		if wac.eventFilter.ShouldSendEvent(session, payload.Event) {
			webhookEvent := &webhook.WebhookEvent{
				ID:        fmt.Sprintf("custom_%s", payload.Metadata.EventID),
				Type:      payload.Event,
				SessionID: wac.sessionID,
				Timestamp: payload.Timestamp,
				Data:      payload.Data,
				URL:       session.WebhookURL,
				Retries:   0,
			}

			logger.Debug("🚀 Enviando webhook customizado: ID=%s, Type=%s, URL=%s", webhookEvent.ID, webhookEvent.Type, webhookEvent.URL)

			err = webhookService.SendEvent(webhookEvent)
			if err != nil {
				logger.Error("❌ Erro ao enviar webhook customizado para sessão %s: %v", wac.sessionID, err)
			} else {
				logger.Info("✅ Webhook customizado enviado com sucesso para sessão %s: %s", wac.sessionID, payload.Event)
			}
		} else {
			logger.Debug("🔧 Evento %s filtrado para webhook customizado da sessão %s (eventos configurados: %s)", payload.Event, wac.sessionID, session.Events)
		}
	}

	// Enviar para webhook padrão (se configurado) - sempre envia todos os eventos
	if session.Webhook != "" {
		webhookEvent := &webhook.WebhookEvent{
			ID:        fmt.Sprintf("default_%s", payload.Metadata.EventID),
			Type:      payload.Event,
			SessionID: wac.sessionID,
			Timestamp: payload.Timestamp,
			Data:      payload.Data,
			URL:       session.Webhook,
			Retries:   0,
		}

		logger.Debug("🚀 Enviando webhook padrão: ID=%s, Type=%s, URL=%s", webhookEvent.ID, webhookEvent.Type, webhookEvent.URL)

		err = webhookService.SendEvent(webhookEvent)
		if err != nil {
			logger.Error("❌ Erro ao enviar webhook padrão para sessão %s: %v", wac.sessionID, err)
		} else {
			logger.Info("✅ Webhook padrão enviado com sucesso para sessão %s: %s", wac.sessionID, payload.Event)
		}
	}
}

// handleQRCode trata evento de novo código QR
func (wac *WhatsAppClient) handleQRCode(code string) {
	logger.Info("=== NOVO QR CODE PARA SESSÃO %s ===", wac.sessionID)

	// Exibir QR code no terminal
	qrterminal.GenerateHalfBlock(code, qrterminal.L, os.Stdout)
	fmt.Printf("\n📱 Escaneie o QR code acima com o WhatsApp\n")
	fmt.Printf("🔗 Código QR: %s\n", code)
	fmt.Printf("⏰ Sessão: %s\n", wac.sessionID)
	fmt.Printf("⏱️  Expira em: ~20 segundos\n")
	fmt.Printf("=======================================\n\n")

	// Salvar QR code no banco
	wac.saveQRCodeToDB(code)

	// Enviar QR code pelo canal interno (compatibilidade)
	select {
	case wac.qrChannel <- code:
	default:
		// Canal cheio, ignorar
	}

	// Enviar webhook
	wac.sendQRWebhook(code, "code")
}

// handleQRTimeout trata evento de timeout do QR code
func (wac *WhatsAppClient) handleQRTimeout() {
	logger.Warn("QR code expirou para sessão %s - aguardando novo...", wac.sessionID)

	// Limpar QR code do banco
	wac.clearQRCodeFromDB()

	// Enviar webhook
	wac.sendQRWebhook("", "timeout")

	fmt.Printf("\n⏰ QR code expirou - aguardando novo...\n")
	fmt.Printf("📱 Sessão: %s\n", wac.sessionID)
	fmt.Printf("🔄 Novo QR code será gerado automaticamente\n")
	fmt.Printf("=======================================\n\n")

	// NÃO terminar o loop - aguardar novo QR code
	// O whatsmeow automaticamente gerará um novo QR code
}

// handleQRSuccess trata evento de sucesso na autenticação
func (wac *WhatsAppClient) handleQRSuccess() {
	logger.Info("🎉 QR code escaneado com sucesso para sessão %s!", wac.sessionID)

	// Limpar QR code do banco
	wac.clearQRCodeFromDB()

	// Atualizar status da sessão
	wac.updateSessionStatus(entity.StatusConnected)

	// Enviar webhook
	wac.sendQRWebhook("", "success")

	fmt.Printf("\n✅ QR code escaneado com sucesso!\n")
	fmt.Printf("📱 Sessão: %s\n", wac.sessionID)
	fmt.Printf("🎉 WhatsApp conectado!\n")
	fmt.Printf("=======================================\n\n")
}

// handleQRError trata evento de erro na autenticação
func (wac *WhatsAppClient) handleQRError(err error) {
	logger.Error("Erro no processo QR para sessão %s: %v", wac.sessionID, err)

	// Limpar QR code do banco
	wac.clearQRCodeFromDB()

	// Atualizar status da sessão
	wac.updateSessionStatus(entity.StatusDisconnected)

	// Enviar webhook
	wac.sendQRWebhook("", "error")

	fmt.Printf("\n❌ Erro no processo QR!\n")
	fmt.Printf("📱 Sessão: %s\n", wac.sessionID)
	fmt.Printf("🚨 Erro: %v\n", err)
	fmt.Printf("=======================================\n\n")
}

// saveQRCodeToDB salva o QR code no banco de dados
func (wac *WhatsAppClient) saveQRCodeToDB(_ string) {
	session, err := wac.sessionRepo.GetByID(wac.sessionID)
	if err != nil {
		logger.Error("Erro ao buscar sessão para salvar QR: %v", err)
		return
	}

	// Atualizar QR code e timestamp
	session.UpdatedAt = time.Now()
	if err := wac.sessionRepo.Update(session); err != nil {
		logger.Error("Erro ao salvar QR code no banco: %v", err)
	} else {
		logger.Debug("QR code salvo no banco para sessão %s", wac.sessionID)
	}
}

// clearQRCodeFromDB limpa o QR code do banco de dados
func (wac *WhatsAppClient) clearQRCodeFromDB() {
	session, err := wac.sessionRepo.GetByID(wac.sessionID)
	if err != nil {
		logger.Error("Erro ao buscar sessão para limpar QR: %v", err)
		return
	}

	// Limpar QR code e atualizar timestamp
	session.UpdatedAt = time.Now()
	if err := wac.sessionRepo.Update(session); err != nil {
		logger.Error("Erro ao limpar QR code do banco: %v", err)
	} else {
		logger.Debug("QR code limpo do banco para sessão %s", wac.sessionID)
	}
}

// sendQRWebhook envia eventos QR via webhook
func (wac *WhatsAppClient) sendQRWebhook(code, event string) {
	webhookData := map[string]interface{}{
		"type":      "QREvent",
		"sessionId": wac.sessionID,
		"event":     event,
		"timestamp": time.Now().Unix(),
	}

	if code != "" {
		webhookData["qrCode"] = code
	}

	// Usar método existente de webhook
	wac.sendWebhook(webhookData)
}

// handleGroupInfo trata eventos de informações de grupo
func (wac *WhatsAppClient) handleGroupInfo(evt *events.GroupInfo) {
	logger.Debug("Informações de grupo recebidas na sessão %s para %s", wac.sessionID, evt.JID.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handlePicture trata eventos de mudança de foto
func (wac *WhatsAppClient) handlePicture(evt *events.Picture) {
	logger.Debug("Mudança de foto recebida na sessão %s para %s", wac.sessionID, evt.JID.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleHistorySync trata eventos de sincronização de histórico
func (wac *WhatsAppClient) handleHistorySync(evt *events.HistorySync) {
	logger.Debug("Sincronização de histórico recebida na sessão %s", wac.sessionID)

	// Configurar client para ser mais tolerante a falhas de download de mídia
	if wac.client != nil {
		// Reduzir timeout de download para evitar travamentos longos
		wac.client.AutoTrustIdentity = false
	}

	// Enviar webhook com evento bruto (sem dados de mídia pesados)
	wac.sendWebhookForEvent(evt)
}

// handleCallOffer trata eventos de oferta de chamada
func (wac *WhatsAppClient) handleCallOffer(evt *events.CallOffer) {
	logger.Info("Oferta de chamada recebida na sessão %s de %s", wac.sessionID, evt.From.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleCallAccept trata eventos de aceitação de chamada
func (wac *WhatsAppClient) handleCallAccept(evt *events.CallAccept) {
	logger.Info("Chamada aceita na sessão %s de %s", wac.sessionID, evt.From.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleCallTerminate trata eventos de término de chamada
func (wac *WhatsAppClient) handleCallTerminate(evt *events.CallTerminate) {
	logger.Info("Chamada terminada na sessão %s de %s", wac.sessionID, evt.From.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleNewsletterJoin trata eventos de entrada em newsletter
func (wac *WhatsAppClient) handleNewsletterJoin(evt *events.NewsletterJoin) {
	logger.Debug("Entrada em newsletter na sessão %s", wac.sessionID)

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleNewsletterLeave trata eventos de saída de newsletter
func (wac *WhatsAppClient) handleNewsletterLeave(evt *events.NewsletterLeave) {
	logger.Debug("Saída de newsletter na sessão %s", wac.sessionID)

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleNewsletterMuteChange trata eventos de mudança de mute em newsletter
func (wac *WhatsAppClient) handleNewsletterMuteChange(evt *events.NewsletterMuteChange) {
	logger.Debug("Mudança de mute em newsletter na sessão %s", wac.sessionID)

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleBlocklistChange trata eventos de mudança na lista de bloqueados
func (wac *WhatsAppClient) handleBlocklistChange(evt *events.BlocklistChange) {
	logger.Debug("Mudança na lista de bloqueados na sessão %s", wac.sessionID)

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handlePushName trata eventos de mudança de nome de exibição
func (wac *WhatsAppClient) handlePushName(evt *events.PushName) {
	logger.Info("Nome de exibição definido na sessão %s para %s", wac.sessionID, evt.JID.String())

	// Agora que temos PushName, podemos enviar presença disponível
	if wac.client.IsConnected() {
		if err := wac.client.SendPresence(types.PresenceAvailable); err != nil {
			logger.Warn("Erro ao enviar presença disponível após PushName: %v", err)
		} else {
			logger.Debug("Presença disponível enviada após definição do PushName")
		}
	}

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleBusinessName trata eventos de mudança de nome comercial
func (wac *WhatsAppClient) handleBusinessName(evt *events.BusinessName) {
	logger.Debug("Mudança de nome comercial na sessão %s para %s", wac.sessionID, evt.JID.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleJoinedGroup trata eventos de entrada em grupo
func (wac *WhatsAppClient) handleJoinedGroup(evt *events.JoinedGroup) {
	logger.Info("Entrada em grupo na sessão %s: %s", wac.sessionID, evt.JID.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleContact trata eventos de mudança de contato
func (wac *WhatsAppClient) handleContact(evt *events.Contact) {
	logger.Info("Mudança de contato na sessão %s para %s", wac.sessionID, evt.JID.String())

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}

// handleGenericEvent trata eventos genéricos não mapeados
func (wac *WhatsAppClient) handleGenericEvent(evt interface{}) {
	logger.Debug("Evento genérico na sessão %s: %T", wac.sessionID, evt)

	// Enviar webhook com evento bruto
	wac.sendWebhookForEvent(evt)
}
