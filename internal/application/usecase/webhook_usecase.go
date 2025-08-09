package usecase

import (
	"fmt"
	"strings"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/repository"
	"wazmeow/pkg/logger"
)

// SetWebhookUseCase representa o caso de uso para definir webhook
type SetWebhookUseCase struct {
	sessionRepo   repository.SessionRepository
	sessionFinder *SessionFinder
}

// NewSetWebhookUseCase cria uma nova instância do use case
func NewSetWebhookUseCase(sessionRepo repository.SessionRepository) *SetWebhookUseCase {
	return &SetWebhookUseCase{
		sessionRepo:   sessionRepo,
		sessionFinder: NewSessionFinder(sessionRepo),
	}
}

// Execute executa a definição de webhook
func (uc *SetWebhookUseCase) Execute(sessionID string, req *requests.SetWebhookRequest) (*responses.WebhookResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	// Atualizar webhook na sessão
	session.WebhookURL = req.WebhookURL
	if len(req.Events) > 0 {
		session.Events = strings.Join(req.Events, ",")
	}

	err = uc.sessionRepo.Update(session)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar webhook: %w", err)
	}

	logger.Info("Webhook definido - Session: %s, URL: %s", sessionID, req.WebhookURL)

	return &responses.WebhookResponse{
		Webhook: req.WebhookURL,
		Events:  req.Events,
	}, nil
}

// GetWebhookUseCase representa o caso de uso para obter webhook
type GetWebhookUseCase struct {
	sessionRepo   repository.SessionRepository
	sessionFinder *SessionFinder
}

// NewGetWebhookUseCase cria uma nova instância do use case
func NewGetWebhookUseCase(sessionRepo repository.SessionRepository) *GetWebhookUseCase {
	return &GetWebhookUseCase{
		sessionRepo:   sessionRepo,
		sessionFinder: NewSessionFinder(sessionRepo),
	}
}

// Execute executa a obtenção de webhook
func (uc *GetWebhookUseCase) Execute(sessionID string) (*responses.WebhookResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	var events []string
	if session.Events != "" {
		events = strings.Split(session.Events, ",")
	}

	return &responses.WebhookResponse{
		Webhook:   session.WebhookURL,
		Subscribe: events,
	}, nil
}

// UpdateWebhookUseCase representa o caso de uso para atualizar webhook
type UpdateWebhookUseCase struct {
	sessionRepo   repository.SessionRepository
	sessionFinder *SessionFinder
}

// NewUpdateWebhookUseCase cria uma nova instância do use case
func NewUpdateWebhookUseCase(sessionRepo repository.SessionRepository) *UpdateWebhookUseCase {
	return &UpdateWebhookUseCase{
		sessionRepo:   sessionRepo,
		sessionFinder: NewSessionFinder(sessionRepo),
	}
}

// Execute executa a atualização de webhook
func (uc *UpdateWebhookUseCase) Execute(sessionID string, req *requests.UpdateWebhookRequest) (*responses.WebhookResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	// Atualizar webhook na sessão
	if req.Active {
		session.WebhookURL = req.WebhookURL
		if len(req.Events) > 0 {
			session.Events = strings.Join(req.Events, ",")
		}
	} else {
		session.WebhookURL = ""
		session.Events = ""
	}

	err = uc.sessionRepo.Update(session)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar webhook: %w", err)
	}

	logger.Info("Webhook atualizado - Session: %s, Active: %v", sessionID, req.Active)

	return &responses.WebhookResponse{
		Webhook: session.WebhookURL,
		Events:  req.Events,
		Active:  req.Active,
	}, nil
}

// DeleteWebhookUseCase representa o caso de uso para deletar webhook
type DeleteWebhookUseCase struct {
	sessionRepo   repository.SessionRepository
	sessionFinder *SessionFinder
}

// NewDeleteWebhookUseCase cria uma nova instância do use case
func NewDeleteWebhookUseCase(sessionRepo repository.SessionRepository) *DeleteWebhookUseCase {
	return &DeleteWebhookUseCase{
		sessionRepo:   sessionRepo,
		sessionFinder: NewSessionFinder(sessionRepo),
	}
}

// Execute executa a exclusão de webhook
func (uc *DeleteWebhookUseCase) Execute(sessionID string) (*responses.WebhookDeleteResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	// Limpar webhook na sessão
	session.WebhookURL = ""
	session.Events = ""

	err = uc.sessionRepo.Update(session)
	if err != nil {
		return nil, fmt.Errorf("erro ao deletar webhook: %w", err)
	}

	logger.Info("Webhook deletado - Session: %s", sessionID)

	return &responses.WebhookDeleteResponse{
		Details: "Webhook and events deleted successfully",
	}, nil
}
