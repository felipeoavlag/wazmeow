package middleware

import (
	"context"
	"net/http"

	"wazmeow/internal/http/handlers/base"

	"github.com/go-chi/chi/v5"
)

// ContextKey é o tipo para chaves do contexto
type ContextKey string

const (
	// SessionIDKey é a chave para o sessionID no contexto
	SessionIDKey ContextKey = "sessionID"
	// ValidatorKey é a chave para o validador no contexto
	ValidatorKey ContextKey = "validator"
)

// SessionValidator middleware que valida e injeta sessionID no contexto
func SessionValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := chi.URLParam(r, "sessionID")
		
		// Valida o sessionID
		if err := base.ValidateSessionID(sessionID); err != nil {
			base.SendError(w, err, base.GetHTTPStatus(err))
			return
		}

		// Adiciona sessionID ao contexto
		ctx := context.WithValue(r.Context(), SessionIDKey, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestValidator middleware que adiciona validador ao contexto
func RequestValidator(next http.Handler) http.Handler {
	validator := base.NewValidator()
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adiciona validador ao contexto
		ctx := context.WithValue(r.Context(), ValidatorKey, validator)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JSONContentType middleware que força Content-Type para JSON
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// ErrorRecovery middleware que captura panics e retorna erro JSON
func ErrorRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				base.SendInternalError(w, "processar requisição", 
					base.NewInternalError("panic recuperado", nil))
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// RequestLogging middleware que registra informações da requisição
func RequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseHandler := base.NewBaseHandler()
		baseHandler.LogRequest(r, "Incoming request")
		
		next.ServeHTTP(w, r)
	})
}

// GetSessionIDFromContext extrai sessionID do contexto
func GetSessionIDFromContext(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(SessionIDKey).(string)
	return sessionID, ok
}

// GetValidatorFromContext extrai validador do contexto
func GetValidatorFromContext(ctx context.Context) (*base.Validator, bool) {
	validator, ok := ctx.Value(ValidatorKey).(*base.Validator)
	return validator, ok
}

// RequireSessionID é um helper que extrai sessionID do contexto ou retorna erro
func RequireSessionID(w http.ResponseWriter, r *http.Request) (string, bool) {
	sessionID, ok := GetSessionIDFromContext(r.Context())
	if !ok {
		base.SendBadRequest(w, "Session ID não encontrado no contexto")
		return "", false
	}
	return sessionID, true
}

// PhoneValidation middleware específico para validação de telefone
func PhoneValidation(phoneField string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Este middleware pode ser usado para validar telefone em requests específicos
			// Por enquanto, apenas passa adiante - a validação será feita no handler
			next.ServeHTTP(w, r)
		})
	}
}

// CombinedValidation combina múltiplos middlewares de validação
func CombinedValidation(next http.Handler) http.Handler {
	return RequestValidator(
		JSONContentType(
			ErrorRecovery(
				RequestLogging(next),
			),
		),
	)
}

// SessionValidation combina validação de sessão com outras validações
func SessionValidation(next http.Handler) http.Handler {
	return SessionValidator(
		CombinedValidation(next),
	)
}