# Exemplo de Uso da API WazMeow

Este documento demonstra como usar a API WazMeow para gerenciar sessões do WhatsApp.

## 1. Criar uma Sessão

```bash
curl -X POST http://localhost:8080/sessions/add \
  -H "Content-Type: application/json" \
  -d '{
    "name": "minha-sessao"
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Sessão criada com sucesso",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "minha-sessao",
    "status": "disconnected",
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

## 2. Conectar a Sessão

```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/connect
```

**Resposta:**
```json
{
  "success": true,
  "message": "Sessão conectada com sucesso"
}
```

## 3. Obter QR Code para Autenticação

```bash
curl -X GET http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/qr
```

**Resposta:**
```json
{
  "success": true,
  "data": {
    "qr_code": "2@BQcAEAYQAg==,f/9u+vz6zJTzOD0VGOEkjrU=,wU/DdpXJ0tPalzxUr6SQBlMAAAAAElFTkSuQmCC",
    "status": "qr_generated"
  }
}
```

## 4. Emparelhar por Telefone (Alternativa ao QR)

```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/pair \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+5511999999999"
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Código de emparelhamento gerado",
  "data": {
    "code": "ABCD-EFGH"
  }
}
```

## 5. Verificar Status da Sessão

```bash
curl -X GET http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000
```

**Resposta:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "minha-sessao",
    "status": "connected",
    "phone": "+5511999999999",
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:05:00Z",
    "is_connected": true,
    "is_logged_in": true
  }
}
```

## 6. Listar Todas as Sessões

```bash
curl -X GET http://localhost:8080/sessions
```

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "minha-sessao",
      "status": "connected",
      "phone": "+5511999999999",
      "created_at": "2025-01-09T10:00:00Z",
      "updated_at": "2025-01-09T10:05:00Z"
    }
  ]
}
```

## 7. Fazer Logout da Sessão

```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/logout
```

**Resposta:**
```json
{
  "success": true,
  "message": "Logout realizado com sucesso"
}
```

## 8. Deletar Sessão

```bash
curl -X DELETE http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000
```

**Resposta:**
```json
{
  "success": true,
  "message": "Sessão deletada com sucesso"
}
```

## 9. Health Check

```bash
curl -X GET http://localhost:8080/health
```

**Resposta:**
```json
{
  "status": "ok",
  "timestamp": 1704801600,
  "service": "wazmeow-api"
}
```

## Estados da Sessão

- **disconnected**: Sessão criada mas não conectada
- **connecting**: Tentando conectar ao WhatsApp
- **connected**: Conectada e autenticada no WhatsApp
- **logged_out**: Desconectada após logout

## Eventos WhatsApp

Quando uma sessão está conectada, a API processa automaticamente os seguintes eventos:

- **Mensagens**: Recebimento de mensagens de texto e mídia
- **Confirmações de leitura**: Status de entrega e leitura
- **Presença**: Status online/offline de contatos
- **Conexão**: Eventos de conexão e desconexão

## Configuração de Proxy (Opcional)

```bash
curl -X POST http://localhost:8080/sessions/550e8400-e29b-41d4-a716-446655440000/proxy \
  -H "Content-Type: application/json" \
  -d '{
    "proxy_type": "http",
    "proxy_host": "proxy.example.com",
    "proxy_port": 8080,
    "proxy_username": "user",
    "proxy_password": "pass"
  }'
```

## Notas Importantes

1. **QR Code**: O QR code deve ser escaneado no WhatsApp Web para autenticar a sessão
2. **Emparelhamento**: O código de emparelhamento deve ser inserido no WhatsApp do telefone
3. **Reconexão**: Sessões conectadas são automaticamente reconectadas quando o servidor reinicia
4. **Thread Safety**: Todas as operações são thread-safe e podem ser executadas concorrentemente
