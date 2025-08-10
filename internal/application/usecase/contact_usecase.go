package usecase

import (
	"context"
	"fmt"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// GetUserInfoUseCase representa o caso de uso para obter informações do usuário
type GetUserInfoUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewGetUserInfoUseCase cria uma nova instância do use case
func NewGetUserInfoUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *GetUserInfoUseCase {
	return &GetUserInfoUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa a obtenção de informações do usuário
func (uc *GetUserInfoUseCase) Execute(sessionID string, req *requests.GetUserInfoRequest) (*responses.UserInfoResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	var jids []types.JID
	for _, phone := range req.Phone {
		jid, err := parseJID(phone)
		if err != nil {
			return nil, fmt.Errorf("número de telefone inválido '%s': %w", phone, err)
		}
		jids = append(jids, jid)
	}

	userInfos, err := client.GetClient().GetUserInfo(jids)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter informações do usuário: %w", err)
	}

	logger.Info("Informações do usuário obtidas - Session: %s, Count: %d", sessionID, len(userInfos))

	return &responses.UserInfoResponse{
		Users: userInfos,
	}, nil
}

// CheckUserUseCase representa o caso de uso para verificar se usuário está no WhatsApp
type CheckUserUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewCheckUserUseCase cria uma nova instância do use case
func NewCheckUserUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *CheckUserUseCase {
	return &CheckUserUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa a verificação de usuário
func (uc *CheckUserUseCase) Execute(sessionID string, req *requests.CheckUserRequest) (*responses.CheckUserResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	var phoneNumbers []string
	for _, phone := range req.Phone {
		jid, err := parseContactJID(phone)
		if err != nil {
			return nil, fmt.Errorf("número de telefone inválido '%s': %w", phone, err)
		}
		phoneNumbers = append(phoneNumbers, jid.User)
	}

	results, err := client.GetClient().IsOnWhatsApp(phoneNumbers)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar usuários: %w", err)
	}

	var users []responses.UserInfo
	for _, result := range results {
		verifiedName := ""
		if result.VerifiedName != nil {
			verifiedName = result.VerifiedName.Details.GetVerifiedName()
		}
		users = append(users, responses.UserInfo{
			Query:        result.Query,
			IsInWhatsapp: result.IsIn,
			JID:          result.JID.String(),
			VerifiedName: verifiedName,
		})
	}

	logger.Info("Usuários verificados - Session: %s, Count: %d", sessionID, len(users))

	return &responses.CheckUserResponse{
		Users: users,
	}, nil
}

// GetAvatarUseCase representa o caso de uso para obter avatar do usuário
type GetAvatarUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewGetAvatarUseCase cria uma nova instância do use case
func NewGetAvatarUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *GetAvatarUseCase {
	return &GetAvatarUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa a obtenção do avatar do usuário
func (uc *GetAvatarUseCase) Execute(sessionID string, req *requests.GetAvatarRequest) (*responses.AvatarResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	jid, err := parseContactJID(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("número de telefone inválido: %w", err)
	}

	pic, err := client.GetClient().GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
		Preview: req.Preview,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao obter avatar: %w", err)
	}

	logger.Info("Avatar obtido - Session: %s, Phone: %s", sessionID, req.Phone)

	return &responses.AvatarResponse{
		URL: pic.URL,
		ID:  pic.ID,
	}, nil
}

// GetContactsUseCase representa o caso de uso para obter lista de contatos
type GetContactsUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
	sessionFinder  *SessionFinder
}

// NewGetContactsUseCase cria uma nova instância do use case
func NewGetContactsUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *GetContactsUseCase {
	return &GetContactsUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
		sessionFinder:  NewSessionFinder(sessionRepo),
	}
}

// Execute executa a obtenção da lista de contatos
func (uc *GetContactsUseCase) Execute(sessionID string) (*responses.ContactsResponse, error) {
	session, err := uc.sessionFinder.FindSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("sessão não encontrada: %w", err)
	}

	client, exists := uc.sessionManager.GetClient(session.ID)
	if !exists {
		return nil, fmt.Errorf("sessão '%s' não está conectada", session.Name)
	}

	if !client.IsConnected() || !client.IsLoggedIn() {
		return nil, fmt.Errorf("sessão '%s' não está ativa", session.Name)
	}

	contacts, err := client.GetClient().Store.Contacts.GetAllContacts(context.Background())
	if err != nil {
		return nil, fmt.Errorf("erro ao obter contatos: %w", err)
	}

	logger.Info("Contatos obtidos - Session: %s, Count: %d", sessionID, len(contacts))

	return &responses.ContactsResponse{
		Contacts: contacts,
	}, nil
}

// parseContactJID converte um número de telefone em JID
func parseContactJID(phone string) (types.JID, error) {
	if phone == "" {
		return types.JID{}, fmt.Errorf("número de telefone não pode estar vazio")
	}

	// Remove caracteres não numéricos
	cleanPhone := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			cleanPhone += string(char)
		}
	}

	if cleanPhone == "" {
		return types.JID{}, fmt.Errorf("número de telefone inválido")
	}

	// Criar JID para usuário individual
	return types.NewJID(cleanPhone, types.DefaultUserServer), nil
}
