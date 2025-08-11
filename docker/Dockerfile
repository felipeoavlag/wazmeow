# Dockerfile para WazMeow API
# Imagem multi-stage para produzir uma imagem pequena e eficiente

# ========================================
# STAGE 1: Build
# ========================================
FROM golang:1.23-alpine AS builder

# Instalar dependências necessárias para build
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências primeiro (para cache do Docker)
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Gerar documentação Swagger se necessário
RUN if [ -f "cmd/server/main.go" ]; then \
        go install github.com/swaggo/swag/cmd/swag@latest && \
        mkdir -p docs && \
        swag init -g cmd/server/main.go -o docs/ --parseDependency --parseInternal; \
    fi



# Compilar aplicação
# Flags de otimização:
# -ldflags="-w -s" remove informações de debug e reduz tamanho
# CGO_ENABLED=0 para binary estático
# -a força rebuild de todos os pacotes
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o wazmeow \
    ./cmd/server

# ========================================
# STAGE 2: Runtime
# ========================================
FROM alpine:3.19

# Instalar certificados CA e timezone data
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root para segurança
RUN addgroup -g 1001 -S wazmeow && \
    adduser -u 1001 -S wazmeow -G wazmeow

# Definir diretório de trabalho
WORKDIR /app

# Copiar binary da stage de build
COPY --from=builder /app/wazmeow .

# Copiar documentação Swagger
COPY --from=builder /app/docs ./docs

# Definir permissões corretas
RUN chown -R wazmeow:wazmeow /app && \
    chmod +x wazmeow

# Mudar para usuário não-root
USER wazmeow

# Expor porta padrão da aplicação
EXPOSE 8080

# Configurações de ambiente padrão para container
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV LOG_LEVEL=INFO
ENV LOG_FORMAT=json

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Labels para metadados
LABEL maintainer="WazMeow Team"
LABEL description="WazMeow API - WhatsApp API usando whatsmeow"
LABEL version="1.0"

# Comando padrão
CMD ["./wazmeow"]
