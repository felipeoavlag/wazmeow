package usecase

import (
	"fmt"

	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow/types"
)

// ListNewsletterUseCase representa o caso de uso para listar newsletters
type ListNewsletterUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewListNewsletterUseCase cria uma nova instância do use case
func NewListNewsletterUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *ListNewsletterUseCase {
	return &ListNewsletterUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a listagem de newsletters
func (uc *ListNewsletterUseCase) Execute(sessionID string) (*responses.NewsletterListResponse, error) {
	session, err := uc.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(sessionID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	// Obter newsletters subscritas
	newsletters, err := client.GetClient().GetSubscribedNewsletters()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar newsletters: %w", err)
	}

	// Converter para slice de NewsletterMetadata
	var newsletterList []types.NewsletterMetadata
	for _, newsletter := range newsletters {
		newsletterList = append(newsletterList, *newsletter)
	}

	logger.Info("Newsletters listadas - Session: %s, Count: %d", sessionID, len(newsletterList))

	return &responses.NewsletterListResponse{
		Newsletters: newsletterList,
	}, nil
}
