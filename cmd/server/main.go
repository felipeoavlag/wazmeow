package main

import (
	"fmt"

	"wazmeow/internal/http"
	"wazmeow/pkg/logger"

	"github.com/lib/pq"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {
	// Configurar PostgreSQL array wrapper
	sqlstore.PostgresArrayWrapper = pq.Array

	// Criar e iniciar servidor
	server, err := http.NewServer()
	if err != nil {
		fmt.Printf("❌ Erro ao criar servidor: %v\n", err)
		return
	}

	// Iniciar servidor com graceful shutdown
	if err := server.Start(); err != nil {
		logger.Fatal("Erro ao executar servidor: %v", err)
	}
}
