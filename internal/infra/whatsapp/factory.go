package whatsapp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"wazmeow/internal/domain/entity"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/domain/service"
	"wazmeow/internal/infra/webhook"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"golang.org/x/net/proxy"
)

// ClientFactory é responsável por criar clientes WhatsApp
type ClientFactory struct {
	deviceStore          *sqlstore.Container
	sessionRepo          repository.SessionRepository
	sessionDomainService *service.SessionDomainService
	webhookService       *webhook.WebhookService
}

// NewClientFactory cria uma nova instância do factory
func NewClientFactory(deviceStore *sqlstore.Container, sessionRepo repository.SessionRepository, sessionDomainService *service.SessionDomainService, webhookService *webhook.WebhookService) *ClientFactory {
	return &ClientFactory{
		deviceStore:          deviceStore,
		sessionRepo:          sessionRepo,
		sessionDomainService: sessionDomainService,
		webhookService:       webhookService,
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

	// Configurar proxy se disponível
	if session.ProxyConfig != nil {
		logger.Info("Configurando proxy para sessão '%s': %s://%s:%d",
			session.Name, session.ProxyConfig.Type, session.ProxyConfig.Host, session.ProxyConfig.Port)

		if err := cf.configureProxy(session.ProxyConfig); err != nil {
			logger.Error("Erro ao configurar proxy para sessão '%s': %v", session.Name, err)
			return nil, fmt.Errorf("erro ao configurar proxy: %w", err)
		}
	}

	// Criar cliente whatsmeow nativo com configurações otimizadas
	nativeClient := whatsmeow.NewClient(deviceStore, logger.ForWhatsApp())

	// Configurar cliente para reduzir warnings de mídia
	cf.configureClientForMediaOptimization(nativeClient)

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

	// Configurar webhook service no cliente
	if cf.webhookService != nil {
		client.SetWebhookService(cf.webhookService)
	}

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

// ConnectOnStartup conecta sessões que possuem DeviceJID válido (já foram autenticadas)
func (cf *ClientFactory) ConnectOnStartup(sessionManager *SessionManager) error {
	// Buscar todas as sessões
	sessions, err := cf.sessionRepo.List()
	if err != nil {
		return fmt.Errorf("erro ao buscar sessões: %w", err)
	}

	for _, session := range sessions {
		// Usar domain service para determinar se deve reconectar automaticamente
		if cf.sessionDomainService.ShouldAutoReconnectOnStartup(session) {
			logger.Info("Reconectando sessão '%s' (DeviceJID: %s) na inicialização", session.Name, session.DeviceJID)

			// Atualizar status para connecting
			session.Status = entity.StatusConnecting
			session.UpdatedAt = time.Now()
			if err := cf.sessionRepo.Update(session); err != nil {
				logger.Error("Erro ao atualizar status da sessão '%s': %v", session.Name, err)
			}

			// Criar cliente
			client, err := cf.CreateClient(session)
			if err != nil {
				logger.Error("Erro ao criar cliente para sessão '%s': %v", session.Name, err)
				// Atualizar status para desconectado
				session.Status = entity.StatusDisconnected
				session.UpdatedAt = time.Now()
				cf.sessionRepo.Update(session)
				continue
			}

			// Para sessões já autenticadas, usar conexão direta (sem QR)
			if err := client.ConnectDirect(); err != nil {
				logger.Error("Erro ao reconectar sessão '%s': %v", session.Name, err)
				// Atualizar status para desconectado
				session.Status = entity.StatusDisconnected
				session.UpdatedAt = time.Now()
				cf.sessionRepo.Update(session)
				continue
			}

			// Armazenar cliente no gerenciador
			sessionManager.SetClient(session.ID, client)
			logger.Info("Sessão '%s' reconectada automaticamente com sucesso", session.Name)
		} else {
			logger.Debug("Sessão '%s' não possui DeviceJID, pulando reconexão automática", session.Name)
		}
	}

	return nil
}

// configureProxy configura o proxy para o cliente WhatsApp
func (cf *ClientFactory) configureProxy(proxyConfig *entity.ProxyConfig) error {
	if proxyConfig == nil {
		return nil
	}

	// Validar configuração de proxy
	if err := cf.sessionDomainService.ValidateProxyConfig(proxyConfig); err != nil {
		return fmt.Errorf("configuração de proxy inválida: %w", err)
	}

	// Construir URL do proxy
	proxyURL := cf.buildProxyURL(proxyConfig)

	// Configurar proxy baseado no tipo
	switch proxyConfig.Type {
	case "http":
		return cf.configureHTTPProxy(proxyURL)
	case "socks5":
		return cf.configureSOCKS5Proxy(proxyURL, proxyConfig)
	default:
		return fmt.Errorf("tipo de proxy não suportado: %s", proxyConfig.Type)
	}
}

// buildProxyURL constrói a URL do proxy
func (cf *ClientFactory) buildProxyURL(proxyConfig *entity.ProxyConfig) string {
	if proxyConfig.Username != "" && proxyConfig.Password != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%d",
			proxyConfig.Type,
			url.QueryEscape(proxyConfig.Username),
			url.QueryEscape(proxyConfig.Password),
			proxyConfig.Host,
			proxyConfig.Port)
	}

	return fmt.Sprintf("%s://%s:%d",
		proxyConfig.Type,
		proxyConfig.Host,
		proxyConfig.Port)
}

// configureHTTPProxy configura proxy HTTP/HTTPS com bypass para CDNs do WhatsApp
func (cf *ClientFactory) configureHTTPProxy(proxyURL string) error {
	// Configurar variáveis de ambiente para proxy HTTP
	// Esta é a abordagem mais compatível com whatsmeow
	os.Setenv("HTTP_PROXY", proxyURL)
	os.Setenv("HTTPS_PROXY", proxyURL)

	// Também configurar o transport padrão
	proxyURLParsed, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("erro ao parsear URL do proxy: %w", err)
	}

	// Configurar transport customizado com bypass para CDNs do WhatsApp
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			// Verificar se o host deve fazer bypass
			if cf.shouldBypassProxy(req.URL.Host) {
				logger.Debug("Bypass do proxy para host: %s", req.URL.Host)
				return nil, nil // Conexão direta
			}

			// Usar proxy para outros hosts
			return proxyURLParsed, nil
		},
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	// Aplicar ao cliente HTTP padrão (usado pelo whatsmeow)
	http.DefaultTransport = transport

	logger.Info("Proxy HTTP configurado com bypass para CDNs: %s", proxyURL)
	return nil
}

// configureSOCKS5Proxy configura proxy SOCKS5 com bypass para CDNs do WhatsApp
func (cf *ClientFactory) configureSOCKS5Proxy(proxyURL string, proxyConfig *entity.ProxyConfig) error {
	// Para SOCKS5, precisamos usar uma abordagem diferente
	proxyAddr := fmt.Sprintf("%s:%d", proxyConfig.Host, proxyConfig.Port)

	// Criar dialer SOCKS5
	var auth *proxy.Auth
	if proxyConfig.Username != "" && proxyConfig.Password != "" {
		auth = &proxy.Auth{
			User:     proxyConfig.Username,
			Password: proxyConfig.Password,
		}
	}

	dialer, err := proxy.SOCKS5("tcp", proxyAddr, auth, proxy.Direct)
	if err != nil {
		return fmt.Errorf("erro ao criar dialer SOCKS5: %w", err)
	}

	// Configurar transport com dialer SOCKS5 e bypass para CDNs
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Extrair host do endereço
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
			}

			// Verificar se o host deve fazer bypass
			if cf.shouldBypassProxy(host) {
				logger.Debug("Bypass do proxy SOCKS5 para host: %s", host)
				// Usar dialer direto
				directDialer := &net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}
				return directDialer.DialContext(ctx, network, addr)
			}

			// Usar proxy SOCKS5 para outros hosts
			return dialer.Dial(network, addr)
		},
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	// Aplicar ao cliente HTTP padrão
	http.DefaultTransport = transport

	// Também configurar variáveis de ambiente para compatibilidade
	os.Setenv("HTTP_PROXY", proxyURL)
	os.Setenv("HTTPS_PROXY", proxyURL)

	logger.Info("Proxy SOCKS5 configurado com bypass para CDNs: %s", proxyAddr)
	return nil
}

// getWhatsAppBypassHosts retorna a lista de hosts do WhatsApp que devem fazer bypass do proxy
func (cf *ClientFactory) getWhatsAppBypassHosts() []string {
	return []string{
		"cdn.whatsapp.net",
		"media-lga3-1.cdn.whatsapp.net",
		"media-lga3-2.cdn.whatsapp.net",
		"media-iad3-1.cdn.whatsapp.net",
		"media-iad3-2.cdn.whatsapp.net",
		"media-sjc3-1.cdn.whatsapp.net",
		"media-sjc3-2.cdn.whatsapp.net",
		"media-dfw5-1.cdn.whatsapp.net",
		"media-dfw5-2.cdn.whatsapp.net",
		"mmg.whatsapp.net",
		"pps.whatsapp.net",
		"web.whatsapp.com",
		"static.whatsapp.net",
	}
}

// shouldBypassProxy verifica se um host deve fazer bypass do proxy
func (cf *ClientFactory) shouldBypassProxy(host string) bool {
	bypassHosts := cf.getWhatsAppBypassHosts()
	for _, bypassHost := range bypassHosts {
		if strings.Contains(host, bypassHost) {
			return true
		}
	}
	return false
}

// configureClientForMediaOptimization configura o cliente para otimizar downloads de mídia
func (cf *ClientFactory) configureClientForMediaOptimization(_ *whatsmeow.Client) {
	// O whatsmeow usa o http.DefaultTransport que já foi configurado com bypass
	// Apenas log informativo de que a configuração foi aplicada
	logger.Debug("Cliente WhatsApp configurado com proxy bypass para CDNs de mídia")
}
