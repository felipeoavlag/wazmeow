package session

import (
	"sync/atomic"
	"time"

	"wazmeow/internal/domain/entities"
)

// StateManager gerencia estado de sessões com operações atômicas
type StateManager struct {
	connected   int32 // 0 = false, 1 = true
	loggedIn    int32 // 0 = false, 1 = true
	status      int32 // entities.SessionStatus como int32
	lastUpdated int64 // Unix timestamp
}

// NewStateManager cria um novo gerenciador de estado
func NewStateManager() *StateManager {
	return &StateManager{
		lastUpdated: time.Now().Unix(),
	}
}

// IsConnected verifica se está conectado (thread-safe)
func (sm *StateManager) IsConnected() bool {
	return atomic.LoadInt32(&sm.connected) == 1
}

// SetConnected define estado de conexão (thread-safe)
func (sm *StateManager) SetConnected(connected bool) {
	var val int32
	if connected {
		val = 1
	}
	atomic.StoreInt32(&sm.connected, val)
	atomic.StoreInt64(&sm.lastUpdated, time.Now().Unix())
}

// IsLoggedIn verifica se está logado (thread-safe)
func (sm *StateManager) IsLoggedIn() bool {
	return atomic.LoadInt32(&sm.loggedIn) == 1
}

// SetLoggedIn define estado de login (thread-safe)
func (sm *StateManager) SetLoggedIn(loggedIn bool) {
	var val int32
	if loggedIn {
		val = 1
	}
	atomic.StoreInt32(&sm.loggedIn, val)
	atomic.StoreInt64(&sm.lastUpdated, time.Now().Unix())
}

// GetStatus retorna o status atual (thread-safe)
func (sm *StateManager) GetStatus() entities.SessionStatus {
	status := atomic.LoadInt32(&sm.status)
	switch status {
	case 1:
		return entities.StatusConnecting
	case 2:
		return entities.StatusConnected
	default:
		return entities.StatusDisconnected
	}
}

// SetStatus define o status (thread-safe)
func (sm *StateManager) SetStatus(status entities.SessionStatus) {
	var val int32
	switch status {
	case entities.StatusConnecting:
		val = 1
	case entities.StatusConnected:
		val = 2
	default:
		val = 0 // StatusDisconnected
	}
	atomic.StoreInt32(&sm.status, val)
	atomic.StoreInt64(&sm.lastUpdated, time.Now().Unix())
}

// GetLastUpdated retorna timestamp da última atualização
func (sm *StateManager) GetLastUpdated() time.Time {
	timestamp := atomic.LoadInt64(&sm.lastUpdated)
	return time.Unix(timestamp, 0)
}

// Reset reseta todos os estados para valores padrão
func (sm *StateManager) Reset() {
	atomic.StoreInt32(&sm.connected, 0)
	atomic.StoreInt32(&sm.loggedIn, 0)
	atomic.StoreInt32(&sm.status, 0)
	atomic.StoreInt64(&sm.lastUpdated, time.Now().Unix())
}
