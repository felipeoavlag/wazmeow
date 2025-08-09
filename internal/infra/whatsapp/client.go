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

	// Configurar event handlers padrão
	wac.setupDefaultEventHandlers()

	return wac
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
		case *events.Message:
			wac.handleMessage(e)
		case *events.Receipt:
			wac.handleReceipt(e)
		case *events.Presence:
			wac.handlePresence(e)
		case *events.ChatPresence:
			wac.handleChatPresence(e)
		default:
			logger.Debug("Evento não tratado: %T", evt)
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

	// Enviar presença disponível
	if err := wac.client.SendPresence(types.PresenceAvailable); err != nil {
		logger.Warn("Erro ao enviar presença disponível: %v", err)
	}
}

// handleDisconnected trata eventos de desconexão
func (wac *WhatsAppClient) handleDisconnected(evt *events.Disconnected) {
	logger.Info("Sessão %s desconectada do WhatsApp", wac.sessionID)

	// Atualizar status da sessão
	if err := wac.updateSessionStatus(entity.StatusDisconnected); err != nil {
		logger.Error("Erro ao atualizar status da sessão: %v", err)
	}
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

	// Criar mapa de dados da mensagem
	messageData := map[string]interface{}{
		"type":      "Message",
		"sessionId": wac.sessionID,
		"messageId": evt.Info.ID,
		"from":      evt.Info.Sender.String(),
		"chat":      evt.Info.Chat.String(),
		"timestamp": evt.Info.Timestamp.Unix(),
		"pushName":  evt.Info.PushName,
		"isFromMe":  evt.Info.IsFromMe,
		"isGroup":   evt.Info.IsGroup,
		"event":     evt,
	}

	// Processar mídia se presente
	wac.processMessageMedia(evt, messageData)

	// Aqui você pode implementar webhook ou outras integrações
	wac.sendWebhook(messageData)
}

// handleReceipt trata eventos de confirmação de leitura
func (wac *WhatsAppClient) handleReceipt(evt *events.Receipt) {
	logger.Debug("Confirmação de leitura recebida na sessão %s", wac.sessionID)

	receiptData := map[string]interface{}{
		"type":        "ReadReceipt",
		"sessionId":   wac.sessionID,
		"messageIds":  evt.MessageIDs,
		"from":        evt.SourceString(),
		"timestamp":   evt.Timestamp.Unix(),
		"receiptType": string(evt.Type),
	}

	// Aqui você pode implementar webhook ou outras integrações
	wac.sendWebhook(receiptData)
}

// handlePresence trata eventos de presença
func (wac *WhatsAppClient) handlePresence(evt *events.Presence) {
	logger.Debug("Presença recebida na sessão %s de %s", wac.sessionID, evt.From.String())

	presenceData := map[string]interface{}{
		"type":        "Presence",
		"sessionId":   wac.sessionID,
		"from":        evt.From.String(),
		"unavailable": evt.Unavailable,
		"lastSeen":    evt.LastSeen.Unix(),
	}

	// Aqui você pode implementar webhook ou outras integrações
	wac.sendWebhook(presenceData)
}

// handleChatPresence trata eventos de presença em chat
func (wac *WhatsAppClient) handleChatPresence(evt *events.ChatPresence) {
	logger.Debug("Presença de chat recebida na sessão %s", wac.sessionID)

	chatPresenceData := map[string]interface{}{
		"type":      "ChatPresence",
		"sessionId": wac.sessionID,
		"chat":      evt.MessageSource.Chat.String(),
		"sender":    evt.MessageSource.Sender.String(),
		"state":     string(evt.State),
		"media":     string(evt.Media),
	}

	// Aqui você pode implementar webhook ou outras integrações
	wac.sendWebhook(chatPresenceData)
}

// processMessageMedia processa mídia de mensagens
func (wac *WhatsAppClient) processMessageMedia(evt *events.Message, messageData map[string]interface{}) {
	// TODO: Implementar processamento de mídia (imagens, áudios, vídeos, documentos)
	// Por enquanto, apenas log
	if evt.Message.GetImageMessage() != nil {
		logger.Debug("Mensagem contém imagem")
		messageData["mediaType"] = "image"
	} else if evt.Message.GetAudioMessage() != nil {
		logger.Debug("Mensagem contém áudio")
		messageData["mediaType"] = "audio"
	} else if evt.Message.GetVideoMessage() != nil {
		logger.Debug("Mensagem contém vídeo")
		messageData["mediaType"] = "video"
	} else if evt.Message.GetDocumentMessage() != nil {
		logger.Debug("Mensagem contém documento")
		messageData["mediaType"] = "document"
	} else {
		messageData["mediaType"] = "text"
	}
}

// sendWebhook envia dados para webhook (implementação placeholder)
func (wac *WhatsAppClient) sendWebhook(data map[string]interface{}) {
	// TODO: Implementar envio de webhook
	// Por enquanto, apenas log
	logger.Debug("Webhook data para sessão %s: %+v", wac.sessionID, data)
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
func (wac *WhatsAppClient) saveQRCodeToDB(code string) {
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
