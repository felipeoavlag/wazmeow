package requests

// SessionRequest representa a requisição para operações de sessão
type SessionRequest struct {
	// Nome da sessão
	Name string `json:"name" example:"minha-sessao"`
}

// CreateSessionRequest representa a requisição para criar uma sessão
type CreateSessionRequest struct {
	// Nome único da sessão
	Name string `json:"name" example:"minha-sessao" validate:"required"`
}

// PairPhoneRequest representa a requisição para emparelhar telefone
type PairPhoneRequest struct {
	// Número de telefone com código do país (formato: +5511999999999)
	Phone string `json:"phone" example:"+5511999999999" validate:"required"`
}

// SetProxyRequest representa a requisição para configurar proxy
type SetProxyRequest struct {
	// Tipo do proxy (http, socks5)
	Type string `json:"type" example:"http" validate:"required"`
	// Host do servidor proxy
	Host string `json:"host" example:"proxy.example.com" validate:"required"`
	// Porta do servidor proxy
	Port int `json:"port" example:"8080" validate:"required,min=1,max=65535"`
	// Nome de usuário para autenticação (opcional)
	Username string `json:"username,omitempty" example:"usuario"`
	// Senha para autenticação (opcional)
	Password string `json:"password,omitempty" example:"senha"`
}
