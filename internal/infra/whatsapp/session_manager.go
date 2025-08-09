package whatsapp

import (
	"sync"

	"go.mau.fi/whatsmeow"
)

// SessionManager gerencia os clientes WhatsApp ativos
type SessionManager struct {
	clients map[string]*WhatsAppClient
	mutex   sync.RWMutex
}

// NewSessionManager cria uma nova instância do gerenciador de sessões
func NewSessionManager() *SessionManager {
	return &SessionManager{
		clients: make(map[string]*WhatsAppClient),
	}
}

// SetClient define o cliente WhatsApp para uma sessão
func (sm *SessionManager) SetClient(sessionID string, client *WhatsAppClient) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.clients[sessionID] = client
}

// GetClient retorna o cliente WhatsApp para uma sessão
func (sm *SessionManager) GetClient(sessionID string) (*WhatsAppClient, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	client, exists := sm.clients[sessionID]
	return client, exists
}

// GetWhatsmeowClient retorna o cliente whatsmeow nativo para uma sessão
func (sm *SessionManager) GetWhatsmeowClient(sessionID string) (*whatsmeow.Client, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	client, exists := sm.clients[sessionID]
	if !exists {
		return nil, false
	}
	return client.client, true
}

// RemoveClient remove o cliente WhatsApp de uma sessão
func (sm *SessionManager) RemoveClient(sessionID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if client, exists := sm.clients[sessionID]; exists {
		client.Disconnect()
		delete(sm.clients, sessionID)
	}
}

// IsConnected verifica se uma sessão está conectada
func (sm *SessionManager) IsConnected(sessionID string) bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	client, exists := sm.clients[sessionID]
	if !exists {
		return false
	}
	return client.IsConnected()
}

// IsLoggedIn verifica se uma sessão está logada
func (sm *SessionManager) IsLoggedIn(sessionID string) bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	client, exists := sm.clients[sessionID]
	if !exists {
		return false
	}
	return client.IsLoggedIn()
}

// GetAllClients retorna todos os clientes ativos
func (sm *SessionManager) GetAllClients() map[string]*WhatsAppClient {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	result := make(map[string]*WhatsAppClient)
	for sessionID, client := range sm.clients {
		result[sessionID] = client
	}
	return result
}

// DisconnectAll desconecta todos os clientes
func (sm *SessionManager) DisconnectAll() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for sessionID, client := range sm.clients {
		client.Disconnect()
		delete(sm.clients, sessionID)
	}
}
