package requests

// CreateSessionRequest representa a requisição para criar uma sessão
type CreateSessionRequest struct {
	Name string `json:"name"`
}

// IsValidURLName verifica se o nome é válido para uso em URL
func (r *CreateSessionRequest) IsValidURLName() bool {
	if len(r.Name) < 3 || len(r.Name) > 50 {
		return false
	}
	
	// Permitir apenas letras, números, hífens e underscores
	for _, char := range r.Name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}
	
	// Não pode começar ou terminar com hífen ou underscore
	if r.Name[0] == '-' || r.Name[0] == '_' ||
		r.Name[len(r.Name)-1] == '-' || r.Name[len(r.Name)-1] == '_' {
		return false
	}
	
	return true
}

// PairPhoneRequest representa a requisição para emparelhar telefone
type PairPhoneRequest struct {
	Phone string `json:"phone"`
}

// SetProxyRequest representa a requisição para configurar proxy
type SetProxyRequest struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
