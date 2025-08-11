package whatsapp

import (
	"sync"
)

// ClientManager manages multiple WhatsApp sessions using MyClient
type ClientManager struct {
	sync.RWMutex
	myClients map[string]*MyClient
}

// NewClientManager creates a new client manager
func NewClientManager() *ClientManager {
	return &ClientManager{
		myClients: make(map[string]*MyClient),
	}
}

// SetMyClient stores a MyClient for a user
func (cm *ClientManager) SetMyClient(userID string, client *MyClient) {
	cm.Lock()
	defer cm.Unlock()
	cm.myClients[userID] = client
}

// GetMyClient retrieves a MyClient for a user
func (cm *ClientManager) GetMyClient(userID string) *MyClient {
	cm.RLock()
	defer cm.RUnlock()
	return cm.myClients[userID]
}

// DeleteMyClient removes a MyClient for a user
func (cm *ClientManager) DeleteMyClient(userID string) {
	cm.Lock()
	defer cm.Unlock()
	if client, exists := cm.myClients[userID]; exists {
		client.Cleanup()
		delete(cm.myClients, userID)
	}
}

// HasSession checks if a session exists
func (cm *ClientManager) HasSession(userID string) bool {
	cm.RLock()
	defer cm.RUnlock()
	_, exists := cm.myClients[userID]
	return exists
}

// UpdateMyClientSubscriptions updates event subscriptions for a client
func (cm *ClientManager) UpdateMyClientSubscriptions(userID string, subscriptions []string) {
	cm.Lock()
	defer cm.Unlock()
	if client, exists := cm.myClients[userID]; exists {
		client.SetSubscriptions(subscriptions)
	}
}

// GetAllSessions returns all active session IDs
func (cm *ClientManager) GetAllSessions() []string {
	cm.RLock()
	defer cm.RUnlock()
	sessions := make([]string, 0, len(cm.myClients))
	for userID := range cm.myClients {
		sessions = append(sessions, userID)
	}
	return sessions
}

// CleanupSession removes and cleans up a client for a user session
func (cm *ClientManager) CleanupSession(userID string) {
	cm.Lock()
	defer cm.Unlock()
	if client, exists := cm.myClients[userID]; exists {
		client.Cleanup()
		delete(cm.myClients, userID)
	}
}

// CleanupAll cleans up all sessions
func (cm *ClientManager) CleanupAll() {
	cm.Lock()
	defer cm.Unlock()
	for _, client := range cm.myClients {
		client.Cleanup()
	}
	cm.myClients = make(map[string]*MyClient)
}
