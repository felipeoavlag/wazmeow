package client

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"

	"wazmeow/internal/config"
	"wazmeow/pkg/logger"
)

// Factory cria clientes e devices de forma otimizada
type Factory struct {
	container *sqlstore.Container
	config    *config.WhatsAppConfig
}

// NewFactory cria uma nova factory
func NewFactory(container *sqlstore.Container, config *config.WhatsAppConfig) *Factory {
	return &Factory{
		container: container,
		config:    config,
	}
}

// CreateWrapper cria um novo wrapper com cliente WhatsApp
func (f *Factory) CreateWrapper(ctx context.Context, sessionID string) (*Wrapper, error) {
	// Criar device
	device := f.CreateDevice(nil)
	if device == nil {
		return nil, fmt.Errorf("failed to create device for session %s", sessionID)
	}

	// Criar cliente WhatsApp com logger específico
	clientLog := logger.NewWALogger(fmt.Sprintf("Client-%s", sessionID))
	client := whatsmeow.NewClient(device, clientLog)

	// Criar context com timeout
	_, cancel := context.WithTimeout(ctx, f.config.ConnectionTimeout)

	// Criar wrapper
	wrapper := NewWrapper(client, sessionID, cancel)

	logger.Info().Str("sessionID", sessionID).Msg("Wrapper created successfully")
	return wrapper, nil
}

// CreateWrapperFromDevice cria wrapper a partir de device existente
func (f *Factory) CreateWrapperFromDevice(ctx context.Context, sessionID string, device *store.Device) (*Wrapper, error) {
	if device == nil {
		return nil, fmt.Errorf("device cannot be nil")
	}

	// Criar cliente WhatsApp
	clientLog := logger.NewWALogger(fmt.Sprintf("Client-%s", sessionID))
	client := whatsmeow.NewClient(device, clientLog)

	// Criar context com timeout
	_, cancel := context.WithTimeout(ctx, f.config.ConnectionTimeout)

	// Criar wrapper
	wrapper := NewWrapper(client, sessionID, cancel)

	// Definir JID se disponível
	if device.ID != nil {
		wrapper.SetJID(*device.ID)
	}

	logger.Info().Str("sessionID", sessionID).Msg("Wrapper created from existing device")
	return wrapper, nil
}

// CreateDevice cria um novo device ou retorna existente
func (f *Factory) CreateDevice(jid *types.JID) *store.Device {
	if jid != nil {
		// Tentar obter device existente
		if device, err := f.container.GetDevice(context.Background(), *jid); err == nil && device != nil {
			logger.Debug().Str("jid", jid.String()).Msg("Using existing device")
			return device
		}
	}

	// Criar novo device
	device := f.container.NewDevice()
	logger.Debug().Msg("Created new device")
	return device
}

// GetDeviceByJID obtém device por JID
func (f *Factory) GetDeviceByJID(jid types.JID) (*store.Device, error) {
	device, err := f.container.GetDevice(context.Background(), jid)
	if err != nil {
		return nil, fmt.Errorf("failed to get device for JID %s: %w", jid.String(), err)
	}
	return device, nil
}

// GetAllDevices retorna todos os devices
func (f *Factory) GetAllDevices() ([]*store.Device, error) {
	devices, err := f.container.GetAllDevices(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}
	return devices, nil
}
