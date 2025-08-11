package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// MyClient wraps whatsmeow.Client with session-specific metadata
type MyClient struct {
	WAClient       *whatsmeow.Client
	UserID         string
	Token          string
	webhook        string
	subscriptions  []string
	killChannel    chan bool
	qrChannel      <-chan whatsmeow.QRChannelItem
	db             *sql.DB
	eventHandlerID uint32
	mutex          sync.RWMutex
}

// NewMyClient creates a new MyClient instance
func NewMyClient(waClient *whatsmeow.Client, userID, token string, db *sql.DB) *MyClient {
	myClient := &MyClient{
		WAClient:      waClient,
		UserID:        userID,
		Token:         token,
		subscriptions: []string{"Message", "Connected", "Disconnected", "QR", "PairSuccess", "LoggedOut"}, // Default subscriptions
		killChannel:   make(chan bool, 1),
		db:            db,
	}

	// Register the event handler
	myClient.eventHandlerID = waClient.AddEventHandler(myClient.eventHandler)

	return myClient
}

// SetWebhook sets the webhook URL for this client
func (mc *MyClient) SetWebhook(webhook string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.webhook = webhook
}

// GetWebhook returns the webhook URL for this client
func (mc *MyClient) GetWebhook() string {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.webhook
}

// SetSubscriptions sets the event subscriptions for this client
func (mc *MyClient) SetSubscriptions(subscriptions []string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.subscriptions = subscriptions
	logger.Info().Str("sessionID", mc.UserID).Strs("subscriptions", subscriptions).Msg("Updated event subscriptions")
}

// GetSubscriptions returns the event subscriptions for this client
func (mc *MyClient) GetSubscriptions() []string {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	// Return copy to prevent external modification
	result := make([]string, len(mc.subscriptions))
	copy(result, mc.subscriptions)
	return result
}

// AddSubscription adds a single event subscription
func (mc *MyClient) AddSubscription(eventType string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// Check if already subscribed
	for _, sub := range mc.subscriptions {
		if sub == eventType {
			return // Already subscribed
		}
	}

	mc.subscriptions = append(mc.subscriptions, eventType)
	logger.Info().Str("sessionID", mc.UserID).Str("eventType", eventType).Msg("Added event subscription")
}

// RemoveSubscription removes a single event subscription
func (mc *MyClient) RemoveSubscription(eventType string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	for i, sub := range mc.subscriptions {
		if sub == eventType {
			mc.subscriptions = append(mc.subscriptions[:i], mc.subscriptions[i+1:]...)
			logger.Info().Str("sessionID", mc.UserID).Str("eventType", eventType).Msg("Removed event subscription")
			return
		}
	}
}

// GetSupportedEventTypes returns list of supported event types
func (mc *MyClient) GetSupportedEventTypes() []string {
	return []string{
		"All",
		"Connected",
		"Disconnected",
		"Message",
		"PairSuccess",
		"LoggedOut",
		"ReadReceipt",
		"Presence",
		"ConnectFailure",
		"QR",
	}
}

// HasSubscription checks if client is subscribed to a specific event
func (mc *MyClient) HasSubscription(event string) bool {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	for _, sub := range mc.subscriptions {
		if sub == event {
			return true
		}
	}
	return false
}

// Kill sends a kill signal to terminate this client session
func (mc *MyClient) Kill() {
	select {
	case mc.killChannel <- true:
	default:
		// Channel already has a signal or is closed
	}
}

// KillChannel returns the kill channel for this client
func (mc *MyClient) KillChannel() <-chan bool {
	return mc.killChannel
}

// Connect connects the WhatsApp client
func (mc *MyClient) Connect() error {
	return mc.WAClient.Connect()
}

// Disconnect disconnects the WhatsApp client
func (mc *MyClient) Disconnect() {
	mc.WAClient.Disconnect()
}

// IsConnected returns true if the client is connected
func (mc *MyClient) IsConnected() bool {
	return mc.WAClient.IsConnected()
}

// IsLoggedIn returns true if the client is logged in
func (mc *MyClient) IsLoggedIn() bool {
	return mc.WAClient.IsLoggedIn()
}

// SetQRChannel sets the QR channel for this client
func (mc *MyClient) SetQRChannel(qrChannel <-chan whatsmeow.QRChannelItem) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.qrChannel = qrChannel
}

// GetQRChannel returns the QR channel for this client
func (mc *MyClient) GetQRChannel() <-chan whatsmeow.QRChannelItem {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.qrChannel
}

// HasQRChannel checks if QR channel is set
func (mc *MyClient) HasQRChannel() bool {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.qrChannel != nil
}

// PairPhone pairs a phone number with this session
func (mc *MyClient) PairPhone(ctx context.Context, phone string) (string, error) {
	if mc.WAClient.IsLoggedIn() {
		return "", fmt.Errorf("session %s already logged in", mc.UserID)
	}

	linkingCode, err := mc.WAClient.PairPhone(ctx, phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	return linkingCode, nil
}

// Logout logs out from WhatsApp
func (mc *MyClient) Logout(ctx context.Context) error {
	if !mc.WAClient.IsLoggedIn() {
		return fmt.Errorf("session %s not logged in", mc.UserID)
	}

	return mc.WAClient.Logout(ctx)
}

// eventHandler handles WhatsApp events for this specific session
func (mc *MyClient) eventHandler(rawEvt interface{}) {
	// Only process events if subscribed
	if !mc.shouldProcessEvent(rawEvt) {
		return
	}

	// Create event data map
	postmap := make(map[string]interface{})
	postmap["event"] = rawEvt
	postmap["sessionID"] = mc.UserID

	// Process different event types
	switch evt := rawEvt.(type) {
	case *events.Connected:
		postmap["type"] = "Connected"
		logger.Info().Str("sessionID", mc.UserID).Msg("WhatsApp connected")

	case *events.Disconnected:
		postmap["type"] = "Disconnected"
		logger.Info().Str("sessionID", mc.UserID).Msg("WhatsApp disconnected")

	case *events.Message:
		postmap["type"] = "Message"
		logger.Info().Str("sessionID", mc.UserID).Str("id", evt.Info.ID).Str("source", evt.Info.SourceString()).Msg("Message received")

	case *events.PairSuccess:
		postmap["type"] = "PairSuccess"
		logger.Info().Str("sessionID", mc.UserID).Str("jid", evt.ID.String()).Msg("Pairing successful")

	case *events.LoggedOut:
		postmap["type"] = "LoggedOut"
		logger.Info().Str("sessionID", mc.UserID).Str("reason", evt.Reason.String()).Msg("Logged out")
		mc.Kill() // Trigger session termination

	case *events.Receipt:
		postmap["type"] = "ReadReceipt"
		switch evt.Type {
		case types.ReceiptTypeRead:
			postmap["state"] = "Read"
		case types.ReceiptTypeDelivered:
			postmap["state"] = "Delivered"
		default:
			return // Skip unhandled receipt types
		}

	case *events.Presence:
		postmap["type"] = "Presence"
		if evt.Unavailable {
			postmap["state"] = "offline"
		} else {
			postmap["state"] = "online"
		}

	case *events.ConnectFailure:
		postmap["type"] = "ConnectFailure"
		logger.Error().Str("sessionID", mc.UserID).Interface("reason", evt).Msg("Connection failed")

	default:
		logger.Debug().Str("sessionID", mc.UserID).Str("event", fmt.Sprintf("%T", evt)).Msg("Unhandled event type")
		return
	}

	// Here you could send to webhook or process further
	mc.processEvent(postmap)
}

// shouldProcessEvent checks if this event should be processed based on subscriptions
func (mc *MyClient) shouldProcessEvent(rawEvt interface{}) bool {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// If subscribed to "All", process everything
	if mc.hasSubscription("All") {
		return true
	}

	// Check specific event type subscriptions
	switch rawEvt.(type) {
	case *events.Connected:
		return mc.hasSubscription("Connected")
	case *events.Disconnected:
		return mc.hasSubscription("Disconnected")
	case *events.Message:
		return mc.hasSubscription("Message")
	case *events.PairSuccess:
		return mc.hasSubscription("PairSuccess")
	case *events.LoggedOut:
		return mc.hasSubscription("LoggedOut")
	case *events.Receipt:
		return mc.hasSubscription("ReadReceipt")
	case *events.Presence:
		return mc.hasSubscription("Presence")
	case *events.ConnectFailure:
		return mc.hasSubscription("ConnectFailure")
	default:
		return false
	}
}

// hasSubscription checks if subscribed to a specific event type (internal, no mutex)
func (mc *MyClient) hasSubscription(eventType string) bool {
	for _, sub := range mc.subscriptions {
		if sub == eventType {
			return true
		}
	}
	return false
}

// processEvent processes the event data (placeholder for webhook/further processing)
func (mc *MyClient) processEvent(eventData map[string]interface{}) {
	// Here you could:
	// 1. Send to webhook if configured
	// 2. Store in database
	// 3. Send to message queue
	// 4. Emit to WebSocket clients

	logger.Debug().
		Str("sessionID", mc.UserID).
		Interface("eventData", eventData).
		Msg("Event processed")
}

// Cleanup performs cleanup operations for this client
func (mc *MyClient) Cleanup() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// Remove event handler if registered
	if mc.eventHandlerID != 0 && mc.WAClient != nil {
		mc.WAClient.RemoveEventHandler(mc.eventHandlerID)
		mc.eventHandlerID = 0
	}

	if mc.WAClient != nil && mc.WAClient.IsConnected() {
		mc.WAClient.Disconnect()
	}

	// Close kill channel safely
	select {
	case <-mc.killChannel:
		// Already closed
	default:
		close(mc.killChannel)
	}
}
