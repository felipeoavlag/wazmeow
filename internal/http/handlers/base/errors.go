package base

import (
	"fmt"
	"net/http"
)

// HandlerError representa um erro específico de handler com código HTTP
type HandlerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

// Error implementa a interface error
func (e *HandlerError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap retorna o erro original para compatibilidade com errors.Unwrap
func (e *HandlerError) Unwrap() error {
	return e.Cause
}

// NewHandlerError cria um novo HandlerError
func NewHandlerError(code int, message string, cause error) *HandlerError {
	return &HandlerError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Erros pré-definidos comuns
var (
	// ErrSessionNotFound indica que a sessão não foi encontrada
	ErrSessionNotFound = &HandlerError{
		Code:    http.StatusBadRequest,
		Message: "Session ID é obrigatório",
	}

	// ErrInvalidJSON indica que o JSON da requisição é inválido
	ErrInvalidJSON = &HandlerError{
		Code:    http.StatusBadRequest,
		Message: "Payload JSON inválido",
	}

	// ErrValidationFailed indica que a validação dos dados falhou
	ErrValidationFailed = &HandlerError{
		Code:    http.StatusBadRequest,
		Message: "Dados de entrada inválidos",
	}

	// ErrInternalServer indica um erro interno do servidor
	ErrInternalServer = &HandlerError{
		Code:    http.StatusInternalServerError,
		Message: "Erro interno do servidor",
	}

	// ErrSessionNotConnected indica que a sessão não está conectada
	ErrSessionNotConnected = &HandlerError{
		Code:    http.StatusBadRequest,
		Message: "Sessão não está conectada",
	}

	// ErrUnauthorized indica que a operação não é autorizada
	ErrUnauthorized = &HandlerError{
		Code:    http.StatusUnauthorized,
		Message: "Operação não autorizada",
	}

	// ErrNotFound indica que o recurso não foi encontrado
	ErrNotFound = &HandlerError{
		Code:    http.StatusNotFound,
		Message: "Recurso não encontrado",
	}
)

// NewSessionNotFoundError cria um erro de sessão não encontrada com contexto
func NewSessionNotFoundError(sessionID string) *HandlerError {
	return &HandlerError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Sessão '%s' não encontrada", sessionID),
	}
}

// NewValidationError cria um erro de validação com detalhes específicos
func NewValidationError(field string, reason string) *HandlerError {
	return &HandlerError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Campo '%s' %s", field, reason),
	}
}

// NewInternalError cria um erro interno com causa específica
func NewInternalError(operation string, cause error) *HandlerError {
	return &HandlerError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("Erro ao %s", operation),
		Cause:   cause,
	}
}

// IsHandlerError verifica se um erro é do tipo HandlerError
func IsHandlerError(err error) bool {
	_, ok := err.(*HandlerError)
	return ok
}

// GetHTTPStatus extrai o código de status HTTP de um erro
func GetHTTPStatus(err error) int {
	if handlerErr, ok := err.(*HandlerError); ok {
		return handlerErr.Code
	}
	return http.StatusInternalServerError
}