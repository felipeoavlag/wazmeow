package requests

// GetUserInfoRequest representa a requisição para obter informações do usuário
type GetUserInfoRequest struct {
	Phone []string `json:"phone" validate:"required,min=1"`
}

// CheckUserRequest representa a requisição para verificar se usuário está no WhatsApp
type CheckUserRequest struct {
	Phone []string `json:"phone" validate:"required,min=1"`
}

// GetAvatarRequest representa a requisição para obter avatar do usuário
type GetAvatarRequest struct {
	Phone   string `json:"phone" validate:"required"`
	Preview bool   `json:"preview"`
}
