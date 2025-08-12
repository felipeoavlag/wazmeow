package client

import (
	"context"
	"fmt"
	"sync"

	"go.mau.fi/whatsmeow/store/sqlstore"

	"wazmeow/internal/config"
	"wazmeow/internal/domain/repositories"
	"wazmeow/internal/infra/whatsapp/events"
	"wazmeow/internal/infra/whatsapp/qr"
	"wazmeow/pkg/logger"
)

// Manager gerencia múltiplas sessões WhatsApp com performance otimizada
type Manager struct {
	clients      sync.Map // string -> *Wrapper (lock-free para leituras)
	factory      *Factory
	pool         *Pool
	eventHandler *events.Handler
	qrProcessor  *qr.Processor
	config       *config.WhatsAppConfig
	sessionRepo  repositories.SessionRepository
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewManager cria um novo gerenciador otimizado
func NewManager(
	container *sqlstore.Container,
	sessionRepo repositories.SessionRepository,
	cfg *config.WhatsAppConfig,
) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	factory := NewFactory(container, cfg)
	pool := NewPool(cfg.PoolSize, cfg.PoolMaxIdle, cfg.PoolMaxLifetime)

	return &Manager{
		factory:      factory,
		pool:         pool,
		config:       cfg,
		sessionRepo:  sessionRepo,
		ctx:          ctx,
		cancel:       cancel,
		eventHandler: events.NewHandler(sessionRepo),
		qrProcessor:  qr.NewProcessor(sessionRepo, cfg),
	}
}

// Get retorna um wrapper de cliente (thread-safe, lock-free)
func (m *Manager) Get(sessionID string) *Wrapper {
	if value, ok := m.clients.Load(sessionID); ok {
		return value.(*Wrapper)
	}
	return nil
}

// Has verifica se uma sessão existe (thread-safe, lock-free)
func (m *Manager) Has(sessionID string) bool {
	_, exists := m.clients.Load(sessionID)
	return exists
}

// Create cria uma nova sessão
func (m *Manager) Create(ctx context.Context, sessionID string) error {
	// Verificar limite de sessões
	if m.Count() >= m.config.MaxSessions {
		return fmt.Errorf("maximum sessions limit reached: %d", m.config.MaxSessions)
	}

	// Verificar se já existe
	if m.Has(sessionID) {
		return fmt.Errorf("session %s already exists", sessionID)
	}

	// Criar wrapper usando factory
	wrapper, err := m.factory.CreateWrapper(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to create wrapper: %w", err)
	}

	// TODO: Configurar event handlers (implementar depois)
	// m.eventHandler.Setup(wrapper.GetWrapperAdapter())

	// Armazenar no mapa
	m.clients.Store(sessionID, wrapper)

	// TODO: Iniciar lifecycle em goroutine (implementar depois)
	// m.wg.Add(1)
	// go func() {
	// 	defer m.wg.Done()
	// 	m.lifecycle.Start(ctx, wrapper)
	// }()

	logger.Info().Str("sessionID", sessionID).Msg("Session created successfully")
	return nil
}

// Remove remove uma sessão
func (m *Manager) Remove(sessionID string) error {
	wrapper := m.Get(sessionID)
	if wrapper == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Desconectar e limpar
	wrapper.Disconnect()
	m.clients.Delete(sessionID)

	logger.Info().Str("sessionID", sessionID).Msg("Session removed successfully")
	return nil
}

// Count retorna o número de sessões ativas
func (m *Manager) Count() int {
	count := 0
	m.clients.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// List retorna lista de IDs de sessões ativas
func (m *Manager) List() []string {
	var sessions []string
	m.clients.Range(func(key, _ interface{}) bool {
		sessions = append(sessions, key.(string))
		return true
	})
	return sessions
}

// LoadAll carrega todas as sessões do banco
func (m *Manager) LoadAll(ctx context.Context) error {
	sessions, err := m.sessionRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sessions from database: %w", err)
	}

	logger.Info().Int("count", len(sessions)).Msg("Loading sessions from database")

	// Carregar sessões em paralelo (limitado pelo pool)
	semaphore := make(chan struct{}, 10) // Máximo 10 sessões carregando simultaneamente
	var wg sync.WaitGroup

	for _, session := range sessions {
		wg.Add(1)
		go func(sessionID string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			if err := m.Create(ctx, sessionID); err != nil {
				logger.Error().Str("sessionID", sessionID).Err(err).Msg("Failed to load session")
			}
		}(session.ID)
	}

	wg.Wait()
	logger.Info().Int("loaded", m.Count()).Msg("Sessions loaded successfully")
	return nil
}

// Shutdown para o manager gracefully
func (m *Manager) Shutdown(ctx context.Context) error {
	logger.Info().Msg("Shutting down client manager")

	// Cancelar context principal
	m.cancel()

	// Desconectar todas as sessões
	m.clients.Range(func(key, value interface{}) bool {
		wrapper := value.(*Wrapper)
		wrapper.Disconnect()
		return true
	})

	// Aguardar goroutines terminarem com timeout
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info().Msg("All sessions disconnected successfully")
	case <-ctx.Done():
		logger.Warn().Msg("Shutdown timeout reached, forcing exit")
	}

	// Fechar pool
	m.pool.Close()

	return nil
}

// GetStats retorna estatísticas do manager
func (m *Manager) GetStats() map[string]interface{} {
	connected := 0
	loggedIn := 0

	m.clients.Range(func(_, value interface{}) bool {
		wrapper := value.(*Wrapper)
		if wrapper.IsConnected() {
			connected++
		}
		if wrapper.IsLoggedIn() {
			loggedIn++
		}
		return true
	})

	return map[string]interface{}{
		"total":       m.Count(),
		"connected":   connected,
		"loggedIn":    loggedIn,
		"maxSessions": m.config.MaxSessions,
		"poolStats":   m.pool.GetStats(),
	}
}
