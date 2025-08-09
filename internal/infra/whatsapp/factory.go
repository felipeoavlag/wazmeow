package whatsapp

import (
	"context"
	"fmt"

	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

// ClientFactory é responsável por criar clientes WhatsApp
type ClientFactory struct {
	deviceStore *sqlstore.Container
	sessionRepo repository.SessionRepository
}

// NewClientFactory cria uma nova instância do factory
func NewClientFactory(deviceStore *sqlstore.Container, sessionRepo repository.SessionRepository) *ClientFactory {
	return &ClientFactory{
		deviceStore: deviceStore,
		sessionRepo: sessionRepo,
	}
}

// CreateClient cria um novo cliente WhatsApp para uma sessão
func (cf *ClientFactory) CreateClient(session *entity.Session) (*WhatsAppClient, error) {
	var deviceStore *store.Device
	var err error

	// Estratégia de recuperação do device store:
	// 1. Tentar recuperar por DeviceJID se disponível (JID completo salvo)
	// 2. Tentar recuperar por telefone se disponível (fallback)
	// 3. Criar novo device store como fallback

	// Primeiro, tentar usar o DeviceJID se disponível (mais preciso)
	if session.DeviceJID != "" {
		jid, ok := parseJID(session.DeviceJID)
		if ok {
			logger.Debug("Tentando recuperar device store para DeviceJID %s", jid.String())
			deviceStore, err = cf.deviceStore.GetDevice(context.Background(), jid)
			if err != nil {
				logger.Warn("Erro ao obter device store para DeviceJID %s: %v", jid.String(), err)
			} else {
				logger.Info("Device store recuperado com sucesso para DeviceJID %s", jid.String())
			}
		} else {
			logger.Warn("DeviceJID inválido %s", session.DeviceJID)
		}
	}

	// Se não conseguiu recuperar pelo DeviceJID, tentar pelo Phone
	if deviceStore == nil && session.Phone != "" {
		jid, ok := parseJID(session.Phone)
		if ok {
			logger.Debug("Tentando recuperar device store para Phone %s", jid.String())
			deviceStore, err = cf.deviceStore.GetDevice(context.Background(), jid)
			if err != nil {
				logger.Warn("Erro ao obter device store para Phone %s: %v", jid.String(), err)
			} else {
				logger.Info("Device store recuperado com sucesso para Phone %s", jid.String())
			}
		} else {
			logger.Warn("JID inválido para telefone %s", session.Phone)
		}
	}

	// Se ainda não conseguiu recuperar, criar novo device store
	if deviceStore == nil {
		logger.Info("Criando novo device store para sessão '%s'", session.Name)
		deviceStore = cf.deviceStore.NewDevice()
	}

	if deviceStore == nil {
		return nil, fmt.Errorf("erro ao criar device store para sessão %s", session.ID)
	}

	// Criar cliente whatsmeow nativo
	nativeClient := whatsmeow.NewClient(deviceStore, logger.ForWhatsApp())

	// Verificar se já está logado
	isLoggedIn := nativeClient.Store.ID != nil

	// Log do estado da sessão
	if isLoggedIn {
		logger.Info("Sessão '%s' já está logada (JID: %s)", session.Name, nativeClient.Store.ID.String())
	} else {
		logger.Info("Sessão '%s' precisa de autenticação (QR code)", session.Name)
	}

	// Criar wrapper do cliente WhatsApp
	client := NewWhatsAppClient(nativeClient, session.ID, cf.sessionRepo)

	return client, nil
}

// parseJID converte uma string em JID do WhatsApp
func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !containsAt(arg) {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			logger.Error("JID inválido: %v", err)
			return recipient, false
		} else if recipient.User == "" {
			logger.Error("JID inválido: servidor não especificado")
			return recipient, false
		}
		return recipient, true
	}
}

// containsAt verifica se a string contém o caractere '@'
func containsAt(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}

// ConnectOnStartup conecta sessões que estavam conectadas antes do shutdown
func (cf *ClientFactory) ConnectOnStartup(sessionManager *SessionManager) error {
	// Buscar sessões que estavam conectadas
	sessions, err := cf.sessionRepo.List()
	if err != nil {
		return fmt.Errorf("erro ao buscar sessões: %w", err)
	}

	for _, session := range sessions {
		// Conectar apenas sessões que estavam conectadas
		if session.Status == entity.StatusConnected {
			logger.Info("Reconectando sessão '%s' na inicialização", session.Name)

			// Criar cliente
			client, err := cf.CreateClient(session)
			if err != nil {
				logger.Error("Erro ao criar cliente para sessão '%s': %v", session.Name, err)
				continue
			}

			// Conectar cliente usando ConnectWithQR
			// NÃO usar timeout para permitir que o loop QR continue ativo
			if err := client.ConnectWithQR(context.Background()); err != nil {
				logger.Error("Erro ao conectar sessão '%s': %v", session.Name, err)
				// Atualizar status para desconectado
				session.Status = entity.StatusDisconnected
				cf.sessionRepo.Update(session)
				continue
			}

			// Armazenar cliente no gerenciador
			sessionManager.SetClient(session.ID, client)
			logger.Info("Sessão '%s' reconectada com sucesso", session.Name)
		}
	}

	return nil
}
