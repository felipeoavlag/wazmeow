package usecase

import (
	"fmt"

	"wazmeow/internal/application/dto/requests"
	"wazmeow/internal/application/dto/responses"
	"wazmeow/internal/domain/repository"
	"wazmeow/internal/infra/whatsapp"
	"wazmeow/pkg/logger"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// CreateGroupUseCase representa o caso de uso para criar grupo
type CreateGroupUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewCreateGroupUseCase cria uma nova instância do use case
func NewCreateGroupUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *CreateGroupUseCase {
	return &CreateGroupUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a criação do grupo
func (uc *CreateGroupUseCase) Execute(sessionID string, req *requests.CreateGroupRequest) (*responses.GroupResponse, error) {
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

	// Converter participantes para JIDs
	var participants []types.JID
	for _, phone := range req.Participants {
		jid, err := parseJID(phone)
		if err != nil {
			return nil, fmt.Errorf("número de telefone inválido '%s': %w", phone, err)
		}
		participants = append(participants, jid)
	}

	// Criar o grupo
	groupInfo, err := client.GetClient().CreateGroup(whatsmeow.ReqCreateGroup{
		Name:         req.Name,
		Participants: participants,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar grupo: %w", err)
	}

	logger.Info("Grupo criado - Session: %s, GroupID: %s, Name: %s", sessionID, groupInfo.JID.String(), req.Name)

	return &responses.GroupResponse{
		GroupID:      groupInfo.JID.String(),
		Name:         req.Name,
		Participants: req.Participants,
		Details:      "Group created successfully",
	}, nil
}

// SetGroupPhotoUseCase representa o caso de uso para definir foto do grupo
type SetGroupPhotoUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewSetGroupPhotoUseCase cria uma nova instância do use case
func NewSetGroupPhotoUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SetGroupPhotoUseCase {
	return &SetGroupPhotoUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a definição da foto do grupo
func (uc *SetGroupPhotoUseCase) Execute(sessionID string, req *requests.SetGroupPhotoRequest) (*responses.GroupResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	// Decodificar a imagem base64
	imageData, err := decodeMediaData(req.Image)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar imagem: %w", err)
	}

	// Definir foto do grupo
	pictureID, err := client.GetClient().SetGroupPhoto(groupJID, imageData)
	if err != nil {
		return nil, fmt.Errorf("erro ao definir foto do grupo: %w", err)
	}

	logger.Info("Foto do grupo definida - Session: %s, GroupID: %s, PictureID: %s", sessionID, req.GroupID, pictureID)

	return &responses.GroupResponse{
		GroupID: req.GroupID,
		Details: "Group photo set successfully",
	}, nil
}

// UpdateGroupParticipantsUseCase representa o caso de uso para atualizar participantes do grupo
type UpdateGroupParticipantsUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewUpdateGroupParticipantsUseCase cria uma nova instância do use case
func NewUpdateGroupParticipantsUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *UpdateGroupParticipantsUseCase {
	return &UpdateGroupParticipantsUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a atualização dos participantes do grupo
func (uc *UpdateGroupParticipantsUseCase) Execute(sessionID string, req *requests.UpdateGroupParticipantsRequest) (*responses.GroupResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	// Converter participantes para JIDs
	var participants []types.JID
	for _, phone := range req.Participants {
		jid, err := parseJID(phone)
		if err != nil {
			return nil, fmt.Errorf("número de telefone inválido '%s': %w", phone, err)
		}
		participants = append(participants, jid)
	}

	var result []types.GroupParticipant
	switch req.Action {
	case "add":
		result, err = client.GetClient().UpdateGroupParticipants(groupJID, participants, whatsmeow.ParticipantChangeAdd)
	case "remove":
		result, err = client.GetClient().UpdateGroupParticipants(groupJID, participants, whatsmeow.ParticipantChangeRemove)
	case "promote":
		result, err = client.GetClient().UpdateGroupParticipants(groupJID, participants, whatsmeow.ParticipantChangePromote)
	case "demote":
		result, err = client.GetClient().UpdateGroupParticipants(groupJID, participants, whatsmeow.ParticipantChangeDemote)
	default:
		return nil, fmt.Errorf("ação inválida: %s", req.Action)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar participantes do grupo: %w", err)
	}

	logger.Info("Participantes do grupo atualizados - Session: %s, GroupID: %s, Action: %s, Count: %d",
		sessionID, req.GroupID, req.Action, len(participants))

	// Converter resultado para strings
	var resultParticipants []string
	for _, p := range result {
		resultParticipants = append(resultParticipants, p.JID.String())
	}

	return &responses.GroupResponse{
		GroupID:      req.GroupID,
		Participants: resultParticipants,
		Details:      fmt.Sprintf("Group participants %s successfully", req.Action),
	}, nil
}

// LeaveGroupUseCase representa o caso de uso para sair do grupo
type LeaveGroupUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewLeaveGroupUseCase cria uma nova instância do use case
func NewLeaveGroupUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *LeaveGroupUseCase {
	return &LeaveGroupUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a saída do grupo
func (uc *LeaveGroupUseCase) Execute(sessionID string, req *requests.LeaveGroupRequest) (*responses.GroupResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	err = client.GetClient().LeaveGroup(groupJID)
	if err != nil {
		return nil, fmt.Errorf("erro ao sair do grupo: %w", err)
	}

	logger.Info("Saiu do grupo - Session: %s, GroupID: %s", sessionID, req.GroupID)

	return &responses.GroupResponse{
		GroupID: req.GroupID,
		Details: "Left group successfully",
	}, nil
}

// JoinGroupUseCase representa o caso de uso para entrar no grupo via link
type JoinGroupUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewJoinGroupUseCase cria uma nova instância do use case
func NewJoinGroupUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *JoinGroupUseCase {
	return &JoinGroupUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a entrada no grupo via link
func (uc *JoinGroupUseCase) Execute(sessionID string, req *requests.JoinGroupRequest) (*responses.GroupResponse, error) {
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

	groupJID, err := client.GetClient().JoinGroupWithLink(req.Code)
	if err != nil {
		return nil, fmt.Errorf("erro ao entrar no grupo: %w", err)
	}

	logger.Info("Entrou no grupo - Session: %s, GroupID: %s", sessionID, groupJID.String())

	return &responses.GroupResponse{
		GroupID: groupJID.String(),
		Details: "Joined group successfully",
	}, nil
}

// GetGroupInfoUseCase representa o caso de uso para obter informações do grupo
type GetGroupInfoUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewGetGroupInfoUseCase cria uma nova instância do use case
func NewGetGroupInfoUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *GetGroupInfoUseCase {
	return &GetGroupInfoUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a obtenção de informações do grupo
func (uc *GetGroupInfoUseCase) Execute(sessionID string, req *requests.GetGroupInfoRequest) (*responses.GroupInfoResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	groupInfo, err := client.GetClient().GetGroupInfo(groupJID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter informações do grupo: %w", err)
	}

	// Converter participantes para strings
	var participants []responses.GroupParticipant
	for _, p := range groupInfo.Participants {
		participants = append(participants, responses.GroupParticipant{
			JID:          p.JID.String(),
			IsAdmin:      p.IsAdmin,
			IsSuperAdmin: p.IsSuperAdmin,
		})
	}

	logger.Info("Informações do grupo obtidas - Session: %s, GroupID: %s", sessionID, req.GroupID)

	return &responses.GroupInfoResponse{
		GroupID:      groupInfo.JID.String(),
		Name:         groupInfo.Name,
		Topic:        groupInfo.Topic,
		Owner:        groupInfo.OwnerJID.String(),
		CreatedAt:    groupInfo.GroupCreated.Unix(),
		Participants: participants,
		Size:         len(participants),
	}, nil
}

// ListGroupsUseCase representa o caso de uso para listar grupos
type ListGroupsUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewListGroupsUseCase cria uma nova instância do use case
func NewListGroupsUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *ListGroupsUseCase {
	return &ListGroupsUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a listagem de grupos
func (uc *ListGroupsUseCase) Execute(sessionID string) (*responses.GroupListResponse, error) {
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

	groups, err := client.GetClient().GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar grupos: %w", err)
	}

	var groupList []responses.GroupSummary
	for _, group := range groups {
		groupList = append(groupList, responses.GroupSummary{
			GroupID: group.JID.String(),
			Name:    group.Name,
		})
	}

	logger.Info("Grupos listados - Session: %s, Count: %d", sessionID, len(groupList))

	return &responses.GroupListResponse{
		Groups: groupList,
		Count:  len(groupList),
	}, nil
}

// GetGroupInviteLinkUseCase representa o caso de uso para obter link de convite do grupo
type GetGroupInviteLinkUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewGetGroupInviteLinkUseCase cria uma nova instância do use case
func NewGetGroupInviteLinkUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *GetGroupInviteLinkUseCase {
	return &GetGroupInviteLinkUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a obtenção do link de convite do grupo
func (uc *GetGroupInviteLinkUseCase) Execute(sessionID string, req *requests.GetGroupInviteLinkRequest) (*responses.GroupInviteLinkResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	link, err := client.GetClient().GetGroupInviteLink(groupJID, req.Reset)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter link de convite do grupo: %w", err)
	}

	logger.Info("Link de convite do grupo obtido - Session: %s, GroupID: %s", sessionID, req.GroupID)

	return &responses.GroupInviteLinkResponse{
		GroupID: req.GroupID,
		Link:    link,
	}, nil
}

// RevokeGroupInviteLinkUseCase representa o caso de uso para revogar link de convite do grupo
type RevokeGroupInviteLinkUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewRevokeGroupInviteLinkUseCase cria uma nova instância do use case
func NewRevokeGroupInviteLinkUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *RevokeGroupInviteLinkUseCase {
	return &RevokeGroupInviteLinkUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a revogação do link de convite do grupo
func (uc *RevokeGroupInviteLinkUseCase) Execute(sessionID string, req *requests.RevokeGroupInviteLinkRequest) (*responses.GroupInviteLinkResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	// Revogar link obtendo um novo (reset=true)
	_, err = client.GetClient().GetGroupInviteLink(groupJID, true)
	if err != nil {
		return nil, fmt.Errorf("erro ao revogar link de convite do grupo: %w", err)
	}

	logger.Info("Link de convite do grupo revogado - Session: %s, GroupID: %s", sessionID, req.GroupID)

	return &responses.GroupInviteLinkResponse{
		GroupID: req.GroupID,
		Link:    "revoked",
	}, nil
}

// SetGroupNameUseCase representa o caso de uso para definir nome do grupo
type SetGroupNameUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewSetGroupNameUseCase cria uma nova instância do use case
func NewSetGroupNameUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SetGroupNameUseCase {
	return &SetGroupNameUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a definição do nome do grupo
func (uc *SetGroupNameUseCase) Execute(sessionID string, req *requests.SetGroupNameRequest) (*responses.GroupResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	err = client.GetClient().SetGroupName(groupJID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("erro ao definir nome do grupo: %w", err)
	}

	logger.Info("Nome do grupo definido - Session: %s, GroupID: %s, Name: %s", sessionID, req.GroupID, req.Name)

	return &responses.GroupResponse{
		GroupID: req.GroupID,
		Name:    req.Name,
		Details: "Group name set successfully",
	}, nil
}

// SetGroupTopicUseCase representa o caso de uso para definir tópico do grupo
type SetGroupTopicUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionManager *whatsapp.SessionManager
}

// NewSetGroupTopicUseCase cria uma nova instância do use case
func NewSetGroupTopicUseCase(sessionRepo repository.SessionRepository, sessionManager *whatsapp.SessionManager) *SetGroupTopicUseCase {
	return &SetGroupTopicUseCase{
		sessionRepo:    sessionRepo,
		sessionManager: sessionManager,
	}
}

// Execute executa a definição do tópico do grupo
func (uc *SetGroupTopicUseCase) Execute(sessionID string, req *requests.SetGroupTopicRequest) (*responses.GroupResponse, error) {
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

	groupJID, err := parseJID(req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("ID do grupo inválido: %w", err)
	}

	err = client.GetClient().SetGroupTopic(groupJID, req.Topic, req.TopicID, "")
	if err != nil {
		return nil, fmt.Errorf("erro ao definir tópico do grupo: %w", err)
	}

	logger.Info("Tópico do grupo definido - Session: %s, GroupID: %s, Topic: %s", sessionID, req.GroupID, req.Topic)

	return &responses.GroupResponse{
		GroupID: req.GroupID,
		Details: "Group topic set successfully",
	}, nil
}
