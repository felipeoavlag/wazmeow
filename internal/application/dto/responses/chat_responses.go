package responses

// ChatPresenceResponse representa a resposta de presença no chat
type ChatPresenceResponse = SimpleResponse

// MarkReadResponse representa a resposta de marcar como lida
type MarkReadResponse = SimpleResponse

// DownloadResponse representa a resposta de download de mídia
type DownloadResponse struct {
	Mimetype string `json:"mimetype"`
	Data     string `json:"data"` // Base64 encoded data
}

// HistorySyncResponse representa a resposta de sincronização de histórico
type HistorySyncResponse = TimestampedResponse
