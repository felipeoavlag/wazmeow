package requests

// SessionRequest representa a requisição para operações de sessão
type SessionRequest struct {
	Name string `json:"name"`
}

// CreateSessionRequest representa a requisição para criar uma sessão
type CreateSessionRequest struct {
	Name string `json:"name"`
}

// Validação movida para SessionDomainService.ValidateSessionName()

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
