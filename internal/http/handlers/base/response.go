package base

import (
	"encoding/json"
	"net/http"
	"time"

	"wazmeow/pkg/logger"
)

// APIResponse representa uma resposta padrão da API
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// ResponseBuilder facilita a construção de respostas HTTP padronizadas
type ResponseBuilder struct {
	writer http.ResponseWriter
}

// NewResponseBuilder cria um novo builder de respostas
func NewResponseBuilder(w http.ResponseWriter) *ResponseBuilder {
	return &ResponseBuilder{writer: w}
}

// Success envia uma resposta de sucesso
func (rb *ResponseBuilder) Success(data interface{}, message string) {
	rb.sendJSON(http.StatusOK, &APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Error envia uma resposta de erro
func (rb *ResponseBuilder) Error(err error, statusCode int) {
	var errorMessage string
	
	if handlerErr, ok := err.(*HandlerError); ok {
		errorMessage = handlerErr.Message
		statusCode = handlerErr.Code
	} else {
		errorMessage = err.Error()
	}

	rb.sendJSON(statusCode, &APIResponse{
		Success:   false,
		Error:     errorMessage,
		Timestamp: time.Now().Unix(),
	})
}

// ValidationError envia uma resposta de erro de validação
func (rb *ResponseBuilder) ValidationError(field, reason string) {
	rb.Error(NewValidationError(field, reason), http.StatusBadRequest)
}

// InternalError envia uma resposta de erro interno
func (rb *ResponseBuilder) InternalError(operation string, cause error) {
	logger.Error("Erro interno em %s: %v", operation, cause)
	rb.Error(NewInternalError(operation, cause), http.StatusInternalServerError)
}

// NotFound envia uma resposta de recurso não encontrado
func (rb *ResponseBuilder) NotFound(resource string) {
	rb.Error(&HandlerError{
		Code:    http.StatusNotFound,
		Message: resource + " não encontrado",
	}, http.StatusNotFound)
}

// BadRequest envia uma resposta de requisição inválida
func (rb *ResponseBuilder) BadRequest(message string) {
	rb.Error(&HandlerError{
		Code:    http.StatusBadRequest,
		Message: message,
	}, http.StatusBadRequest)
}

// Unauthorized envia uma resposta de não autorizado
func (rb *ResponseBuilder) Unauthorized(message string) {
	if message == "" {
		message = "Operação não autorizada"
	}
	rb.Error(&HandlerError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}, http.StatusUnauthorized)
}

// sendJSON envia uma resposta JSON
func (rb *ResponseBuilder) sendJSON(statusCode int, response *APIResponse) {
	rb.writer.Header().Set("Content-Type", "application/json")
	rb.writer.WriteHeader(statusCode)

	if err := json.NewEncoder(rb.writer).Encode(response); err != nil {
		logger.Error("Erro ao codificar resposta JSON: %v", err)
		// Fallback para erro simples
		http.Error(rb.writer, "Erro interno do servidor", http.StatusInternalServerError)
	}
}

// Funções de conveniência para uso direto

// SendSuccess envia uma resposta de sucesso diretamente
func SendSuccess(w http.ResponseWriter, data interface{}, message string) {
	NewResponseBuilder(w).Success(data, message)
}

// SendError envia uma resposta de erro diretamente
func SendError(w http.ResponseWriter, err error, statusCode int) {
	NewResponseBuilder(w).Error(err, statusCode)
}

// SendValidationError envia uma resposta de erro de validação diretamente
func SendValidationError(w http.ResponseWriter, field, reason string) {
	NewResponseBuilder(w).ValidationError(field, reason)
}

// SendInternalError envia uma resposta de erro interno diretamente
func SendInternalError(w http.ResponseWriter, operation string, cause error) {
	NewResponseBuilder(w).InternalError(operation, cause)
}

// SendNotFound envia uma resposta de não encontrado diretamente
func SendNotFound(w http.ResponseWriter, resource string) {
	NewResponseBuilder(w).NotFound(resource)
}

// SendBadRequest envia uma resposta de requisição inválida diretamente
func SendBadRequest(w http.ResponseWriter, message string) {
	NewResponseBuilder(w).BadRequest(message)
}

// SendUnauthorized envia uma resposta de não autorizado diretamente
func SendUnauthorized(w http.ResponseWriter, message string) {
	NewResponseBuilder(w).Unauthorized(message)
}