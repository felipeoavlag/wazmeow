package base

import (
	"fmt"
	"reflect"
	"strings"
)

// Validator fornece funcionalidades de validação comuns
type Validator struct{}

// NewValidator cria uma nova instância do validador
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateRequired valida se os campos obrigatórios estão presentes
func (v *Validator) ValidateRequired(fields map[string]interface{}) error {
	var missingFields []string

	for fieldName, value := range fields {
		if v.isEmpty(value) {
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		return NewValidationError(
			strings.Join(missingFields, ", "),
			"é obrigatório",
		)
	}

	return nil
}

// ValidatePhone valida se um número de telefone é válido
func (v *Validator) ValidatePhone(phone string) error {
	if phone == "" {
		return NewValidationError("phone", "é obrigatório")
	}

	// Remove caracteres especiais comuns
	cleanPhone := strings.ReplaceAll(phone, "+", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")

	// Verifica se contém apenas números (exceto se for JID)
	if !strings.Contains(phone, "@") {
		for _, char := range cleanPhone {
			if char < '0' || char > '9' {
				return NewValidationError("phone", "deve conter apenas números")
			}
		}

		// Verifica comprimento mínimo
		if len(cleanPhone) < 10 {
			return NewValidationError("phone", "deve ter pelo menos 10 dígitos")
		}
	}

	return nil
}

// ValidateSessionID valida se um ID de sessão é válido
func (v *Validator) ValidateSessionID(sessionID string) error {
	if sessionID == "" {
		return ErrSessionNotFound
	}

	// Verifica se não contém caracteres especiais perigosos
	if strings.ContainsAny(sessionID, "/<>\"'&") {
		return NewValidationError("sessionID", "contém caracteres inválidos")
	}

	return nil
}

// ValidateEmail valida se um email é válido (básico)
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return NewValidationError("email", "é obrigatório")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return NewValidationError("email", "deve ter formato válido")
	}

	return nil
}

// ValidateURL valida se uma URL é válida (básico)
func (v *Validator) ValidateURL(url string) error {
	if url == "" {
		return nil // URL pode ser opcional
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return NewValidationError("url", "deve começar com http:// ou https://")
	}

	return nil
}

// ValidateLength valida se um campo tem o comprimento adequado
func (v *Validator) ValidateLength(fieldName, value string, min, max int) error {
	length := len(value)

	if min > 0 && length < min {
		return NewValidationError(fieldName, fmt.Sprintf("deve ter pelo menos %d caracteres", min))
	}

	if max > 0 && length > max {
		return NewValidationError(fieldName, fmt.Sprintf("deve ter no máximo %d caracteres", max))
	}

	return nil
}

// ValidateRange valida se um número está dentro de um intervalo
func (v *Validator) ValidateRange(fieldName string, value, min, max int) error {
	if min > 0 && value < min {
		return NewValidationError(fieldName, fmt.Sprintf("deve ser pelo menos %d", min))
	}

	if max > 0 && value > max {
		return NewValidationError(fieldName, fmt.Sprintf("deve ser no máximo %d", max))
	}

	return nil
}

// ValidateSliceLength valida se um slice tem o comprimento adequado
func (v *Validator) ValidateSliceLength(fieldName string, slice interface{}, min, max int) error {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return NewValidationError(fieldName, "deve ser uma lista")
	}

	length := val.Len()

	if min > 0 && length < min {
		return NewValidationError(fieldName, fmt.Sprintf("deve ter pelo menos %d itens", min))
	}

	if max > 0 && length > max {
		return NewValidationError(fieldName, fmt.Sprintf("deve ter no máximo %d itens", max))
	}

	return nil
}

// ValidateOneOf valida se um valor está dentro de uma lista de valores permitidos
func (v *Validator) ValidateOneOf(fieldName, value string, allowedValues []string) error {
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}

	return NewValidationError(
		fieldName,
		fmt.Sprintf("deve ser um dos valores: %s", strings.Join(allowedValues, ", ")),
	)
}

// isEmpty verifica se um valor está vazio
func (v *Validator) isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return strings.TrimSpace(val.String()) == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return val.Len() == 0
	case reflect.Ptr:
		return val.IsNil()
	case reflect.Interface:
		return val.IsNil()
	default:
		// Para tipos básicos como int, float, bool
		zero := reflect.Zero(val.Type())
		return reflect.DeepEqual(val.Interface(), zero.Interface())
	}
}

// Instância global do validador para uso direto
var GlobalValidator = NewValidator()

// Funções de conveniência para uso direto

// ValidateRequired valida campos obrigatórios usando o validador global
func ValidateRequired(fields map[string]interface{}) error {
	return GlobalValidator.ValidateRequired(fields)
}

// ValidatePhone valida telefone usando o validador global
func ValidatePhone(phone string) error {
	return GlobalValidator.ValidatePhone(phone)
}

// ValidateSessionID valida ID de sessão usando o validador global
func ValidateSessionID(sessionID string) error {
	return GlobalValidator.ValidateSessionID(sessionID)
}