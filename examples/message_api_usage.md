# API de Mensagens - WazMeow

Este documento demonstra como usar a API de mensagens do WazMeow para enviar mensagens de texto e mídia.

## Pré-requisitos

1. Ter uma sessão criada e conectada
2. A sessão deve estar com status "connected"

## 1. Enviar Mensagem de Texto

### Endpoint
```
POST /message/{sessionID}/send/text
```

### Exemplo de Requisição
```bash
curl -X POST http://localhost:8080/message/14820c8d-50e2-42ed-a0bb-645d1b083bf7/send/text \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "body": "Olá! Esta é uma mensagem de teste do WazMeow.",
    "id": "msg-001"
  }'
```

### Parâmetros
- `phone` (obrigatório): Número do telefone de destino (com código do país)
- `body` (obrigatório): Texto da mensagem
- `id` (opcional): ID personalizado da mensagem. Se não fornecido, será gerado automaticamente
- `context_info` (opcional): Informações de contexto para reply ou menções

### Resposta de Sucesso
```json
{
  "success": true,
  "message": "Mensagem enviada com sucesso",
  "data": {
    "details": "Sent",
    "timestamp": 1704067200,
    "id": "msg-001"
  }
}
```

## 2. Enviar Mensagem com Reply

### Exemplo de Requisição
```bash
curl -X POST http://localhost:8080/message/14820c8d-50e2-42ed-a0bb-645d1b083bf7/send/text \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "body": "Esta é uma resposta à sua mensagem anterior.",
    "context_info": {
      "stanza_id": "3EB0C431C26A1916E07A",
      "participant": "5511999999999@s.whatsapp.net"
    }
  }'
```

## 3. Enviar Mídia (Imagem)

### Endpoint
```
POST /message/{sessionID}/send/media
```

### Exemplo de Requisição
```bash
curl -X POST http://localhost:8080/message/14820c8d-50e2-42ed-a0bb-645d1b083bf7/send/media \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "media_data": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
    "caption": "Esta é uma imagem de teste",
    "mime_type": "image/png"
  }'
```

### Parâmetros
- `phone` (obrigatório): Número do telefone de destino
- `media_data` (obrigatório): Dados da mídia em base64 (pode ser data URL ou base64 puro)
- `caption` (opcional): Legenda da mídia
- `mime_type` (opcional): Tipo MIME da mídia (será detectado automaticamente se não fornecido)
- `id` (opcional): ID personalizado da mensagem
- `context_info` (opcional): Informações de contexto

### Tipos de Mídia Suportados
- **Imagens**: PNG, JPEG, GIF, WebP
- **Áudios**: MP3, AAC, OGG, WAV
- **Vídeos**: MP4, AVI, MOV, WebM
- **Documentos**: PDF, DOC, DOCX, XLS, XLSX, etc.

### Resposta de Sucesso
```json
{
  "success": true,
  "message": "Mídia enviada com sucesso",
  "data": {
    "details": "Sent",
    "timestamp": 1704067200,
    "id": "3EB0C431C26A1916E07A"
  }
}
```

## 4. Enviar Áudio

### Exemplo de Requisição
```bash
curl -X POST http://localhost:8080/message/14820c8d-50e2-42ed-a0bb-645d1b083bf7/send/media \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "media_data": "data:audio/mp3;base64,SUQzBAAAAAAAI1RTU0UAAAAPAAADTGF2ZjU4Ljc2LjEwMAAAAAAAAAAAAAAA...",
    "mime_type": "audio/mp3"
  }'
```

## 5. Enviar Documento

### Exemplo de Requisição
```bash
curl -X POST http://localhost:8080/message/14820c8d-50e2-42ed-a0bb-645d1b083bf7/send/media \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "5511999999999",
    "media_data": "data:application/pdf;base64,JVBERi0xLjQKJcOkw7zDtsO8CjIgMCBvYmoKPDwKL0xlbmd0aCAzIDAgUgo...",
    "caption": "Documento importante",
    "mime_type": "application/pdf"
  }'
```

## Tratamento de Erros

### Sessão não encontrada
```json
{
  "error": "sessão não encontrada"
}
```

### Sessão não conectada
```json
{
  "error": "sessão 'minha-sessao' não está conectada"
}
```

### Número inválido
```json
{
  "error": "número de telefone inválido"
}
```

### Payload inválido
```json
{
  "error": "Payload inválido"
}
```

## Notas Importantes

1. **Formato do Telefone**: Use o formato internacional com código do país (ex: 5511999999999)
2. **Tamanho da Mídia**: Respeite os limites do WhatsApp para tamanho de arquivos
3. **Base64**: Certifique-se de que os dados em base64 estão corretos
4. **Sessão Ativa**: A sessão deve estar conectada e autenticada
5. **Rate Limiting**: Evite enviar muitas mensagens em sequência para não ser bloqueado

## Exemplo Completo em JavaScript

```javascript
// Função para enviar mensagem de texto
async function sendTextMessage(sessionId, phone, message) {
  const response = await fetch(`http://localhost:8080/message/${sessionId}/send/text`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      phone: phone,
      body: message
    })
  });
  
  return await response.json();
}

// Função para enviar imagem
async function sendImage(sessionId, phone, imageBase64, caption) {
  const response = await fetch(`http://localhost:8080/message/${sessionId}/send/media`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      phone: phone,
      media_data: imageBase64,
      caption: caption,
      mime_type: 'image/jpeg'
    })
  });
  
  return await response.json();
}

// Uso
sendTextMessage('14820c8d-50e2-42ed-a0bb-645d1b083bf7', '5511999999999', 'Olá!')
  .then(result => console.log('Mensagem enviada:', result))
  .catch(error => console.error('Erro:', error));
```
