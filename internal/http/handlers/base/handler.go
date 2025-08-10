package base

import (
	"encoding/json"
	"net/http"

	"wazmeow/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// BaseHandler fornece funcionalidades comuns para todos os handlers
type BaseHandler struct {
	validator *Validator
}

// NewBaseHandler cria uma nova instância do BaseHandler
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		validator: NewValidator(),
	}
}

// ExtractSessionID extrai e valida o sessionID da URL
func (h *BaseHandler) ExtractSessionID(r *http.Request) (string, error) {
	sessionID := chi.URLParam(r, "sessionID")
	if err := h.validator.ValidateSessionID(sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}

// ExtractSessionIDOrError extrai sessionID e envia erro se inválido
func (h *BaseHandler) ExtractSessionIDOrError(w http.ResponseWriter, r *http.Request) (string, bool) {
	sessionID, err := h.ExtractSessionID(r)
	if err != nil {
		SendError(w, err, GetHTTPStatus(err))
		return "", false
	}
	return sessionID, true
}

// DecodeJSON decodifica o corpo da requisição em uma estrutura
func (h *BaseHandler) DecodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		logger.Error("Erro ao decodificar JSON: %v", err)
		return NewHandlerError(http.StatusBadRequest, "Payload JSON inválido", err)
	}
	return nil
}

// DecodeJSONOrError decodifica JSON e envia erro se inválido
func (h *BaseHandler) DecodeJSONOrError(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := h.DecodeJSON(r, v); err != nil {
		SendError(w, err, GetHTTPStatus(err))
		return false
	}
	return true
}

// ValidateRequired valida campos obrigatórios
func (h *BaseHandler) ValidateRequired(fields map[string]interface{}) error {
	return h.validator.ValidateRequired(fields)
}

// ValidateRequiredOrError valida campos obrigatórios e envia erro se inválido
func (h *BaseHandler) ValidateRequiredOrError(w http.ResponseWriter, fields map[string]interface{}) bool {
	if err := h.ValidateRequired(fields); err != nil {
		SendError(w, err, GetHTTPStatus(err))
		return false
	}
	return true
}

// ValidatePhone valida um número de telefone
func (h *BaseHandler) ValidatePhone(phone string) error {
	return h.validator.ValidatePhone(phone)
}

// ValidatePhoneOrError valida telefone e envia erro se inválido
func (h *BaseHandler) ValidatePhoneOrError(w http.ResponseWriter, phone string) bool {
	if err := h.ValidatePhone(phone); err != nil {
		SendError(w, err, GetHTTPStatus(err))
		return false
	}
	return true
}

// ValidateRequestWithPhone valida campos obrigatórios incluindo telefone
func (h *BaseHandler) ValidateRequestWithPhone(w http.ResponseWriter, fields map[string]interface{}, phone string) bool {
	// Validar campos obrigatórios
	if !h.ValidateRequiredOrError(w, fields) {
		return false
	}
	
	// Validar telefone
	if !h.ValidatePhoneOrError(w, phone) {
		return false
	}
	
	return true
}

// SendSuccess envia uma resposta de sucesso
func (h *BaseHandler) SendSuccess(w http.ResponseWriter, data interface{}, message string) {
	SendSuccess(w, data, message)
}

// SendError envia uma resposta de erro
func (h *BaseHandler) SendError(w http.ResponseWriter, err error, statusCode int) {
	SendError(w, err, statusCode)
}

// SendValidationError envia uma resposta de erro de validação
func (h *BaseHandler) SendValidationError(w http.ResponseWriter, field, reason string) {
	SendValidationError(w, field, reason)
}

// SendInternalError envia uma resposta de erro interno
func (h *BaseHandler) SendInternalError(w http.ResponseWriter, operation string, cause error) {
	SendInternalError(w, operation, cause)
}

// SendNotFound envia uma resposta de não encontrado
func (h *BaseHandler) SendNotFound(w http.ResponseWriter, resource string) {
	SendNotFound(w, resource)
}

// SendBadRequest envia uma resposta de requisição inválida
func (h *BaseHandler) SendBadRequest(w http.ResponseWriter, message string) {
	SendBadRequest(w, message)
}

// HandleUseCaseExecution executa um use case e trata erros automaticamente
func (h *BaseHandler) HandleUseCaseExecution(
	w http.ResponseWriter,
	operation string,
	executor func() (interface{}, error),
	successMessage string,
) bool {
	result, err := executor()
	if err != nil {
		logger.Error("Erro ao %s: %v", operation, err)
		h.SendInternalError(w, operation, err)
		return false
	}

	h.SendSuccess(w, result, successMessage)
	return true
}

// LogRequest registra informações da requisição
func (h *BaseHandler) LogRequest(r *http.Request, operation string) {
	logger.Info("Requisição %s - Método: %s, URL: %s, IP: %s",
		operation, r.Method, r.URL.Path, r.RemoteAddr)
}

// LogSuccess registra sucesso da operação
func (h *BaseHandler) LogSuccess(operation, sessionID string, additionalInfo ...interface{}) {
	if len(additionalInfo) > 0 {
		logger.Info("%s realizado com sucesso - Session: %s, Info: %v",
			operation, sessionID, additionalInfo)
	} else {
		logger.Info("%s realizado com sucesso - Session: %s", operation, sessionID)
	}
}

// LogError registra erro da operação
func (h *BaseHandler) LogError(operation string, err error, additionalInfo ...interface{}) {
	if len(additionalInfo) > 0 {
		logger.Error("Erro ao %s: %v, Info: %v", operation, err, additionalInfo)
	} else {
		logger.Error("Erro ao %s: %v", operation, err)
	}
}

// GetValidator retorna o validador interno
func (h *BaseHandler) GetValidator() *Validator {
	return h.validator
}

// HandlerFunc é um tipo de função que representa um handler HTTP melhorado
type HandlerFunc func(w http.ResponseWriter, r *http.Request, h *BaseHandler)

// WrapHandler converte uma HandlerFunc em http.HandlerFunc
func WrapHandler(handlerFunc HandlerFunc) http.HandlerFunc {
	baseHandler := NewBaseHandler()
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r, baseHandler)
	}
}

// Middleware para adicionar BaseHandler ao contexto
func BaseHandlerMiddleware(next http.Handler) http.Handler {
	baseHandler := NewBaseHandler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adiciona o BaseHandler ao contexto da requisição se necessário
		// Por enquanto, apenas chama o próximo handler
		next.ServeHTTP(w, r)
		
		// Aqui podemos adicionar logging automático ou outras funcionalidades
		baseHandler.LogRequest(r, "Request processed")
	})
}