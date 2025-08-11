# Plano de Refatoração para Multi-Sessão WhatsApp

## Problemas Identificados

### 1. Duplicações no Código
- **Gerenciamento de clientes duplicado:**
  - `service.go:29`: `clients map[string]*whatsmeow.Client`
  - `client_manager.go:13`: `whatsmeowClients map[string]*whatsmeow.Client`

- **Kill channels duplicados:**
  - `service.go:31`: `killChannels map[string]chan bool`
  - `myclient.go:18`: `killChannel chan bool`

- **Dependência inexistente:**
  - `client_manager.go:6`: Import `github.com/go-resty/resty/v2` não existe no go.mod

- **Funcionalidade QR duplicada:**
  - `service.go:30`: `qrChannels map[string]<-chan whatsmeow.QRChannelItem`
  - `myclient.go:110-121`: Método `GenerateQR()` similar

### 2. Arquitetura Atual vs Desejada

**Atual:**
```
Service
├── clients: map[string]*whatsmeow.Client
├── qrChannels: map[string]<-chan whatsmeow.QRChannelItem  
├── killChannels: map[string]chan bool
└── métodos de gerenciamento
```

**Desejada:**
```
Service
└── clientManager: *ClientManager
    ├── myClients: map[string]*MyClient
    │   └── MyClient (encapsula WAClient + metadata)
    └── métodos de gerenciamento
```

## Solução Proposta

### Etapa 1: Corrigir client_manager.go
- Remover import `github.com/go-resty/resty/v2`
- Simplificar para focar apenas em `MyClient` 
- Remover `httpClients` (não usado no wuzapi original)

### Etapa 2: Consolidar Service
- Substituir `clients`, `qrChannels`, `killChannels` por um único `clientManager`
- Usar `MyClient` que encapsula toda funcionalidade de sessão
- Manter interface pública inalterada

### Etapa 3: Ajustar MyClient
- Integrar funcionalidade de QR generation
- Manter kill channel interno
- Adicionar webhook e subscriptions

## Estrutura Final

```go
// Service simplificado
type Service struct {
    sessionRepo   repositories.SessionRepository
    container     *sqlstore.Container
    clientManager *ClientManager
    mu            sync.RWMutex
}

// ClientManager focado
type ClientManager struct {
    sync.RWMutex
    myClients map[string]*MyClient
}

// MyClient completo  
type MyClient struct {
    WAClient      *whatsmeow.Client
    UserID        string
    Token         string
    webhook       string
    subscriptions []string
    killChannel   chan bool
    qrChannel     <-chan whatsmeow.QRChannelItem
    db            *sql.DB
    mutex         sync.RWMutex
}
```

## Benefícios da Refatoração

1. **Eliminação de duplicações** - Um único local para cada responsabilidade
2. **Melhor encapsulamento** - MyClient contém todo contexto de sessão
3. **Thread-safety** - Mutex apropriado em cada nível
4. **Extensibilidade** - Fácil adicionar novos tipos de cliente
5. **Compatibilidade** - Interface Service mantém mesmos métodos

## Próximos Passos

1. ✅ Identificar todas as duplicações
2. 🔄 Refatorar ClientManager (remover resty)
3. 🔄 Refatorar Service para usar ClientManager  
4. 🔄 Consolidar MyClient com funcionalidade QR
5. ⏳ Implementar cache de usuários
6. ⏳ Adicionar event handling por sessão
7. ⏳ Implementar subscription system
8. ⏳ Testar multi-sessão