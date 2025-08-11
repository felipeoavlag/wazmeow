package responses

// SimpleResponse representa uma resposta simples com apenas detalhes
type SimpleResponse struct {
	Details string `json:"details"`
}

// TimestampedResponse representa uma resposta com detalhes e timestamp
type TimestampedResponse struct {
	Details   string `json:"details"`
	Timestamp int64  `json:"timestamp"`
}
