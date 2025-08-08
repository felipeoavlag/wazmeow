package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ConfigPrompt representa uma pergunta de configuração
type ConfigPrompt struct {
	Key          string
	Description  string
	DefaultValue string
	Required     bool
	Sensitive    bool
}

func main() {
	fmt.Println("🚀 WazMeow - Configurador de Servidor")
	fmt.Println("=====================================")
	fmt.Println()

	// Verificar se já existe um arquivo .env
	if _, err := os.Stat(".env"); err == nil {
		fmt.Print("Arquivo .env já existe. Deseja sobrescrever? (s/N): ")
		if !askConfirmation() {
			fmt.Println("Configuração cancelada.")
			return
		}
	}

	// Prompts de configuração
	prompts := []ConfigPrompt{
		// Banco de dados
		{Key: "DB_HOST", Description: "Host do PostgreSQL", DefaultValue: "localhost", Required: true},
		{Key: "DB_PORT", Description: "Porta do PostgreSQL", DefaultValue: "5432", Required: true},
		{Key: "DB_USER", Description: "Usuário do PostgreSQL", DefaultValue: "postgres", Required: true},
		{Key: "DB_PASSWORD", Description: "Senha do PostgreSQL", DefaultValue: "password", Required: true, Sensitive: true},
		{Key: "DB_NAME", Description: "Nome do banco de dados", DefaultValue: "wazmeow", Required: true},
		{Key: "DB_SSLMODE", Description: "Modo SSL do PostgreSQL", DefaultValue: "disable", Required: false},

		// Servidor
		{Key: "SERVER_HOST", Description: "Host do servidor HTTP", DefaultValue: "0.0.0.0", Required: false},
		{Key: "SERVER_PORT", Description: "Porta do servidor HTTP", DefaultValue: "8080", Required: true},
		{Key: "SERVER_READ_TIMEOUT", Description: "Timeout de leitura", DefaultValue: "30s", Required: false},
		{Key: "SERVER_WRITE_TIMEOUT", Description: "Timeout de escrita", DefaultValue: "30s", Required: false},

		// Logs
		{Key: "LOG_LEVEL", Description: "Nível de log (DEBUG, INFO, WARN, ERROR)", DefaultValue: "INFO", Required: false},
		{Key: "LOG_FORMAT", Description: "Formato do log (json, text)", DefaultValue: "json", Required: false},

		// CORS
		{Key: "CORS_ALLOWED_ORIGINS", Description: "Origens permitidas (separadas por vírgula)", DefaultValue: "*", Required: false},
		{Key: "CORS_ALLOWED_METHODS", Description: "Métodos HTTP permitidos", DefaultValue: "GET,POST,PUT,DELETE,OPTIONS", Required: false},
		{Key: "CORS_ALLOWED_HEADERS", Description: "Headers permitidos", DefaultValue: "Accept,Authorization,Content-Type,X-CSRF-Token", Required: false},

		// Sessões
		{Key: "MAX_SESSIONS", Description: "Máximo de sessões simultâneas", DefaultValue: "100", Required: false},
		{Key: "SESSION_TIMEOUT", Description: "Timeout das sessões", DefaultValue: "3600s", Required: false},

		// Segurança
		{Key: "API_KEY", Description: "Chave da API (deixe vazio para desabilitar)", DefaultValue: "", Required: false, Sensitive: true},
		{Key: "RATE_LIMIT_REQUESTS", Description: "Limite de requisições por janela", DefaultValue: "100", Required: false},
		{Key: "RATE_LIMIT_WINDOW", Description: "Janela de tempo para rate limit", DefaultValue: "1m", Required: false},

		// Aplicação
		{Key: "DEBUG", Description: "Modo debug (true/false)", DefaultValue: "false", Required: false},
		{Key: "ENVIRONMENT", Description: "Ambiente (development/production)", DefaultValue: "production", Required: false},
	}

	config := make(map[string]string)

	fmt.Println("📝 Configuração do Banco de Dados")
	fmt.Println("----------------------------------")
	for _, prompt := range prompts[:6] {
		config[prompt.Key] = askInput(prompt)
	}

	fmt.Println("\n🌐 Configuração do Servidor HTTP")
	fmt.Println("----------------------------------")
	for _, prompt := range prompts[6:10] {
		config[prompt.Key] = askInput(prompt)
	}

	fmt.Println("\n📋 Configuração de Logs")
	fmt.Println("------------------------")
	for _, prompt := range prompts[10:12] {
		config[prompt.Key] = askInput(prompt)
	}

	fmt.Println("\n🔗 Configuração de CORS")
	fmt.Println("------------------------")
	for _, prompt := range prompts[12:15] {
		config[prompt.Key] = askInput(prompt)
	}

	fmt.Println("\n👥 Configuração de Sessões")
	fmt.Println("--------------------------")
	for _, prompt := range prompts[15:17] {
		config[prompt.Key] = askInput(prompt)
	}

	fmt.Println("\n🔒 Configuração de Segurança")
	fmt.Println("-----------------------------")
	for _, prompt := range prompts[17:20] {
		config[prompt.Key] = askInput(prompt)
	}

	fmt.Println("\n⚙️  Configuração da Aplicação")
	fmt.Println("-----------------------------")
	for _, prompt := range prompts[20:] {
		config[prompt.Key] = askInput(prompt)
	}

	// Criar arquivo .env
	if err := createEnvFile(config); err != nil {
		fmt.Printf("❌ Erro ao criar arquivo .env: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Configuração concluída!")
	fmt.Println("📁 Arquivo .env criado com sucesso!")
	fmt.Println("\n🚀 Para iniciar o servidor, execute:")
	fmt.Println("   go run cmd/server/main.go")
	fmt.Println("\n   ou")
	fmt.Println("   go build -o bin/wazmeow cmd/server/main.go")
	fmt.Println("   ./bin/wazmeow")
}

func askInput(prompt ConfigPrompt) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s", prompt.Description)
		if prompt.DefaultValue != "" {
			fmt.Printf(" [%s]", prompt.DefaultValue)
		}
		if prompt.Required {
			fmt.Print(" (obrigatório)")
		}
		fmt.Print(": ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			if prompt.Required && prompt.DefaultValue == "" {
				fmt.Println("❌ Este campo é obrigatório!")
				continue
			}
			input = prompt.DefaultValue
		}

		// Validações específicas
		if err := validateInput(prompt.Key, input); err != nil {
			fmt.Printf("❌ %v\n", err)
			continue
		}

		return input
	}
}

func validateInput(key, value string) error {
	switch key {
	case "DB_PORT", "SERVER_PORT":
		if _, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("deve ser um número válido")
		}
	case "LOG_LEVEL":
		validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
		upper := strings.ToUpper(value)
		for _, level := range validLevels {
			if level == upper {
				return nil
			}
		}
		return fmt.Errorf("deve ser um de: %s", strings.Join(validLevels, ", "))
	case "LOG_FORMAT":
		if value != "json" && value != "text" {
			return fmt.Errorf("deve ser 'json' ou 'text'")
		}
	case "DEBUG":
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("deve ser 'true' ou 'false'")
		}
	case "ENVIRONMENT":
		if value != "development" && value != "production" {
			return fmt.Errorf("deve ser 'development' ou 'production'")
		}
	}
	return nil
}

func askConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "s" || input == "sim" || input == "y" || input == "yes"
}

func createEnvFile(config map[string]string) error {
	file, err := os.Create(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	// Escrever cabeçalho
	file.WriteString("# ===========================================\n")
	file.WriteString("# CONFIGURAÇÃO DO WAZMEOW API\n")
	file.WriteString("# Gerado automaticamente pelo configurador\n")
	file.WriteString("# ===========================================\n\n")

	// Seções de configuração
	sections := []struct {
		title string
		keys  []string
	}{
		{"Configuração do Banco de Dados PostgreSQL", []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"}},
		{"Configuração do Servidor HTTP", []string{"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT"}},
		{"Configuração de Logs", []string{"LOG_LEVEL", "LOG_FORMAT"}},
		{"Configuração de CORS", []string{"CORS_ALLOWED_ORIGINS", "CORS_ALLOWED_METHODS", "CORS_ALLOWED_HEADERS"}},
		{"Configuração de Sessões WhatsApp", []string{"MAX_SESSIONS", "SESSION_TIMEOUT"}},
		{"Configuração de Segurança", []string{"API_KEY", "RATE_LIMIT_REQUESTS", "RATE_LIMIT_WINDOW"}},
		{"Configuração de Desenvolvimento", []string{"DEBUG", "ENVIRONMENT"}},
	}

	for _, section := range sections {
		file.WriteString(fmt.Sprintf("# %s\n", section.title))
		for _, key := range section.keys {
			if value, exists := config[key]; exists {
				file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
			}
		}
		file.WriteString("\n")
	}

	return nil
}
