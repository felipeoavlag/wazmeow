package requests

import "go.mau.fi/whatsmeow/types"

// SendPresenceRequest representa a requisição para definir presença global
type SendPresenceRequest struct {
	Type string `json:"type" validate:"required,oneof=available unavailable"`
}

// ChatPresenceRequest representa a requisição para definir presença no chat
type ChatPresenceRequest struct {
	Phone string                  `json:"phone" validate:"required"`
	State string                  `json:"state" validate:"required,oneof=typing paused recording"`
	Media types.ChatPresenceMedia `json:"media,omitempty"`
}

// MarkReadRequest representa a requisição para marcar mensagens como lidas
type MarkReadRequest struct {
	ID     []string  `json:"id" validate:"required,min=1"`
	Chat   types.JID `json:"chat" validate:"required"`
	Sender types.JID `json:"sender,omitempty"`
}

// DownloadImageRequest representa a requisição para download de imagem
type DownloadImageRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}

// DownloadVideoRequest representa a requisição para download de vídeo
type DownloadVideoRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}

// DownloadAudioRequest representa a requisição para download de áudio
type DownloadAudioRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}

// DownloadDocumentRequest representa a requisição para download de documento
type DownloadDocumentRequest struct {
	URL           string `json:"url" validate:"required"`
	DirectPath    string `json:"directPath" validate:"required"`
	MediaKey      []byte `json:"mediaKey" validate:"required"`
	Mimetype      string `json:"mimetype" validate:"required"`
	FileEncSHA256 []byte `json:"fileEncSHA256" validate:"required"`
	FileSHA256    []byte `json:"fileSHA256" validate:"required"`
	FileLength    uint64 `json:"fileLength" validate:"required"`
}
