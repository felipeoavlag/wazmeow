package qr

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"

	"wazmeow/internal/config"
	"wazmeow/internal/domain/entities"
	"wazmeow/internal/domain/repositories"
	"wazmeow/pkg/logger"
)

// Processor gerencia processamento de QR codes de forma otimizada
type Processor struct {
	generator   *Generator
	timeout     time.Duration
	sessionRepo repositories.SessionRepository
}

// NewProcessor cria um novo processador de QR
func NewProcessor(sessionRepo repositories.SessionRepository, config *config.WhatsAppConfig) *Processor {
	return &Processor{
		generator:   NewGenerator(),
		timeout:     config.QRTimeout,
		sessionRepo: sessionRepo,
	}
}

// Process processa QR code para uma sessão
func (p *Processor) Process(ctx context.Context, client *whatsmeow.Client, sessionID string) error {
	logger.Info().Str("sessionID", sessionID).Msg("Starting QR process")

	// Context com timeout
	qrCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// Canal para receber QR codes
	qrChan, err := client.GetQRChannel(qrCtx)
	if err != nil {
		return fmt.Errorf("failed to get QR channel: %w", err)
	}

	// Conectar para iniciar processo de QR
	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect for QR: %w", err)
	}

	// Processar eventos QR
	return p.handleQREvents(qrCtx, sessionID, qrChan)
}

// handleQREvents processa eventos QR com timeout
func (p *Processor) handleQREvents(ctx context.Context, sessionID string, qrChan <-chan whatsmeow.QRChannelItem) error {
	for {
		select {
		case <-ctx.Done():
			logger.Error().Str("sessionID", sessionID).Msg("QR process timeout")
			p.updateSessionStatus(sessionID, entities.StatusDisconnected)
			return fmt.Errorf("QR process timeout")

		case evt, ok := <-qrChan:
			if !ok {
				logger.Error().Str("sessionID", sessionID).Msg("QR channel closed")
				return fmt.Errorf("QR channel closed")
			}

			if err := p.processQREvent(sessionID, evt); err != nil {
				return err
			}

			// Se foi sucesso, sair do loop
			if evt.Event == "success" {
				logger.Info().Str("sessionID", sessionID).Msg("QR authentication successful")
				return nil
			}
		}
	}
}

// processQREvent processa um evento QR específico
func (p *Processor) processQREvent(sessionID string, evt whatsmeow.QRChannelItem) error {
	switch evt.Event {
	case "code":
		return p.handleQRCode(sessionID, evt.Code)
	case "timeout":
		logger.Error().Str("sessionID", sessionID).Msg("QR code timeout")
		p.updateSessionStatus(sessionID, entities.StatusDisconnected)
		return fmt.Errorf("QR code timeout")
	case "success":
		logger.Info().Str("sessionID", sessionID).Msg("QR authentication successful")
		p.clearQRCode(sessionID)
		p.updateSessionStatus(sessionID, entities.StatusConnected)
		return nil
	default:
		logger.Debug().Str("sessionID", sessionID).Str("event", evt.Event).Msg("QR event")
		return nil
	}
}

// handleQRCode processa novo QR code
func (p *Processor) handleQRCode(sessionID, code string) error {
	logger.Info().Str("sessionID", sessionID).Msg("New QR code generated")

	// Gerar PNG base64
	base64PNG, err := p.generator.GenerateBase64PNG(code)
	if err != nil {
		logger.Error().Str("sessionID", sessionID).Err(err).Msg("Failed to generate QR PNG")
		return fmt.Errorf("failed to generate QR PNG: %w", err)
	}

	// Salvar no banco
	if err := p.saveQRCode(sessionID, base64PNG); err != nil {
		logger.Error().Str("sessionID", sessionID).Err(err).Msg("Failed to save QR code")
		return fmt.Errorf("failed to save QR code: %w", err)
	}

	// Display no terminal (opcional)
	p.generator.DisplayTerminal(code)

	logger.Info().Str("sessionID", sessionID).Msg("QR code saved and displayed")
	return nil
}

// saveQRCode salva QR code no banco
func (p *Processor) saveQRCode(sessionID, base64PNG string) error {
	ctx := context.Background()
	session, err := p.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.QRCode = base64PNG
	session.UpdatedAt = time.Now()

	if err := p.sessionRepo.Update(ctx, session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// clearQRCode limpa QR code do banco
func (p *Processor) clearQRCode(sessionID string) {
	ctx := context.Background()
	session, err := p.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.Error().Str("sessionID", sessionID).Err(err).Msg("Failed to get session for QR clear")
		return
	}

	session.QRCode = ""
	session.UpdatedAt = time.Now()

	if err := p.sessionRepo.Update(ctx, session); err != nil {
		logger.Error().Str("sessionID", sessionID).Err(err).Msg("Failed to clear QR code")
	}
}

// updateSessionStatus atualiza status da sessão
func (p *Processor) updateSessionStatus(sessionID string, status entities.SessionStatus) {
	ctx := context.Background()
	if err := p.sessionRepo.UpdateStatus(ctx, sessionID, status); err != nil {
		logger.Error().
			Str("sessionID", sessionID).
			Str("status", string(status)).
			Err(err).
			Msg("Failed to update session status")
	}
}

// GetQRCode retorna QR code salvo
func (p *Processor) GetQRCode(sessionID string) (string, error) {
	ctx := context.Background()
	session, err := p.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	if session.QRCode == "" {
		return "", fmt.Errorf("no QR code available")
	}

	return session.QRCode, nil
}
