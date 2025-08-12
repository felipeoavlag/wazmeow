package client

import (
	"context"
	"sync/atomic"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"

	"wazmeow/internal/domain/entities"
)

// Wrapper encapsula um cliente WhatsApp com estado thread-safe otimizado
type Wrapper struct {
	client    *whatsmeow.Client
	sessionID string
	jid       types.JID
	state     *State
	cancel    context.CancelFunc
}

// State gerencia estados com operações atômicas para performance
type State struct {
	connected int32 // 0 = false, 1 = true
	loggedIn  int32 // 0 = false, 1 = true
	status    int32 // entities.SessionStatus como int32
}

// NewWrapper cria um novo wrapper otimizado
func NewWrapper(client *whatsmeow.Client, sessionID string, cancel context.CancelFunc) *Wrapper {
	return &Wrapper{
		client:    client,
		sessionID: sessionID,
		state:     &State{},
		cancel:    cancel,
	}
}

// Client retorna o cliente WhatsApp
func (w *Wrapper) Client() *whatsmeow.Client {
	return w.client
}

// SessionID retorna o ID da sessão
func (w *Wrapper) SessionID() string {
	return w.sessionID
}

// JID retorna o JID do dispositivo
func (w *Wrapper) JID() types.JID {
	return w.jid
}

// SetJID define o JID do dispositivo
func (w *Wrapper) SetJID(jid types.JID) {
	w.jid = jid
}

// IsConnected verifica se está conectado (thread-safe)
func (w *Wrapper) IsConnected() bool {
	return atomic.LoadInt32(&w.state.connected) == 1
}

// SetConnected define estado de conexão (thread-safe)
func (w *Wrapper) SetConnected(connected bool) {
	var val int32
	if connected {
		val = 1
	}
	atomic.StoreInt32(&w.state.connected, val)
}

// IsLoggedIn verifica se está logado (thread-safe)
func (w *Wrapper) IsLoggedIn() bool {
	return atomic.LoadInt32(&w.state.loggedIn) == 1
}

// SetLoggedIn define estado de login (thread-safe)
func (w *Wrapper) SetLoggedIn(loggedIn bool) {
	var val int32
	if loggedIn {
		val = 1
	}
	atomic.StoreInt32(&w.state.loggedIn, val)
}

// GetStatus retorna o status atual (thread-safe)
func (w *Wrapper) GetStatus() entities.SessionStatus {
	status := atomic.LoadInt32(&w.state.status)
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
func (w *Wrapper) SetStatus(status entities.SessionStatus) {
	var val int32
	switch status {
	case entities.StatusConnecting:
		val = 1
	case entities.StatusConnected:
		val = 2
	default:
		val = 0 // StatusDisconnected
	}
	atomic.StoreInt32(&w.state.status, val)
}

// Disconnect desconecta o cliente e cancela o context
func (w *Wrapper) Disconnect() {
	if w.client != nil && w.client.IsConnected() {
		w.client.Disconnect()
	}
	w.SetConnected(false)
	w.SetStatus(entities.StatusDisconnected)
	if w.cancel != nil {
		w.cancel()
	}
}

// ClientAdapter adapta o cliente para a interface de eventos
type ClientAdapter struct {
	client *whatsmeow.Client
}

// AddEventHandler implementa a interface ClientInterface
func (ca *ClientAdapter) AddEventHandler(handler func(interface{})) {
	ca.client.AddEventHandler(handler)
}

// GetClientAdapter retorna um adapter para eventos
func (w *Wrapper) GetClientAdapter() *ClientAdapter {
	return &ClientAdapter{client: w.client}
}

// WrapperAdapter adapta o wrapper para a interface de eventos
type WrapperAdapter struct {
	wrapper *Wrapper
}

// SessionID implementa WrapperInterface
func (wa *WrapperAdapter) SessionID() string {
	return wa.wrapper.SessionID()
}

// Client implementa WrapperInterface
func (wa *WrapperAdapter) Client() *ClientAdapter {
	return wa.wrapper.GetClientAdapter()
}

// GetWrapperAdapter retorna um adapter para eventos
func (w *Wrapper) GetWrapperAdapter() *WrapperAdapter {
	return &WrapperAdapter{wrapper: w}
}
