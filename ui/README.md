# WazMeow Manager - Frontend

Interface web moderna para gerenciamento de sessões WhatsApp usando a API WazMeow.

## 🚀 Tecnologias

- **Next.js 14** - Framework React com App Router
- **TypeScript** - Tipagem estática
- **Tailwind CSS** - Framework CSS utilitário
- **Shadcn/ui** - Biblioteca de componentes
- **Zustand** - Gerenciamento de estado
- **React Query** - Cache e sincronização de dados
- **React Hook Form** - Gerenciamento de formulários
- **Zod** - Validação de esquemas
- **Sonner** - Sistema de notificações
- **Recharts** - Gráficos e visualizações

## 📁 Estrutura do Projeto

```
ui/
├── src/
│   ├── app/                    # App Router do Next.js
│   │   ├── (dashboard)/        # Layout do dashboard
│   │   │   ├── page.tsx        # Dashboard principal
│   │   │   ├── sessions/       # Gerenciamento de sessões
│   │   │   ├── webhooks/       # Configuração de webhooks
│   │   │   ├── monitoring/     # Monitoramento em tempo real
│   │   │   ├── logs/           # Logs e eventos
│   │   │   ├── proxy/          # Configuração de proxy
│   │   │   ├── notifications/  # Sistema de notificações
│   │   │   └── settings/       # Configurações gerais
│   │   ├── globals.css         # Estilos globais
│   │   └── layout.tsx          # Layout raiz
│   ├── components/             # Componentes reutilizáveis
│   │   ├── forms/              # Formulários
│   │   ├── layout/             # Componentes de layout
│   │   ├── providers/          # Providers de contexto
│   │   └── ui/                 # Componentes base (Shadcn/ui)
│   └── lib/                    # Utilitários e configurações
│       ├── api/                # Cliente HTTP e endpoints
│       ├── hooks/              # Hooks customizados
│       ├── stores/             # Gerenciamento de estado
│       └── types/              # Definições de tipos
├── public/                     # Arquivos estáticos
└── package.json               # Dependências e scripts
```

## 🎯 Funcionalidades

### 📱 Gerenciamento de Sessões
- **Listagem** - Visualização de todas as sessões com filtros e busca
- **Criação** - Formulário completo para criar novas sessões
- **Edição** - Modificação de configurações existentes
- **QR Code** - Geração e exibição de QR codes para autenticação
- **Emparelhamento** - Interface para emparelhamento por telefone
- **Status** - Monitoramento em tempo real do status das sessões

### 🔗 Webhooks
- **Configuração** - Interface para configurar URLs de webhook
- **Eventos** - Seleção de eventos para notificação
- **Teste** - Funcionalidade para testar webhooks
- **Logs** - Histórico de entregas e falhas

### 📊 Monitoramento
- **Métricas** - Gráficos de performance e uso
- **Status** - Indicadores de saúde do sistema
- **Alertas** - Notificações de problemas
- **Relatórios** - Análises detalhadas

### 🔧 Configurações
- **API** - Configuração de conexão com a API
- **Proxy** - Configuração de proxies para sessões
- **Temas** - Alternância entre modo claro/escuro
- **Notificações** - Preferências de notificação

## 🚀 Instalação e Execução

### Pré-requisitos
- Node.js 18+ 
- npm ou yarn

### Instalação
```bash
# Clone o repositório
git clone <repository-url>
cd wazmeow/ui

# Instale as dependências
npm install

# Configure as variáveis de ambiente
cp .env.example .env.local
```

### Configuração
Edite o arquivo `.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### Execução
```bash
# Desenvolvimento
npm run dev

# Build para produção
npm run build

# Executar produção
npm start

# Linting
npm run lint
```

## 🎨 Sistema de Design

### Cores
- **Primary** - Azul (#0070f3)
- **Secondary** - Cinza (#6b7280)
- **Success** - Verde (#10b981)
- **Warning** - Amarelo (#f59e0b)
- **Error** - Vermelho (#ef4444)

### Tipografia
- **Font Family** - Inter (Google Fonts)
- **Sizes** - Sistema baseado em rem (0.875rem - 2.25rem)

### Componentes
Todos os componentes seguem o padrão do Shadcn/ui com customizações:
- **Button** - Variantes: default, destructive, outline, secondary, ghost
- **Card** - Container principal para conteúdo
- **Input** - Campos de entrada com validação
- **Select** - Seleção com busca
- **Dialog** - Modais e overlays

## 📱 Responsividade

O sistema é totalmente responsivo com breakpoints:
- **Mobile** - < 768px
- **Tablet** - 768px - 1024px  
- **Desktop** - > 1024px

### Comportamento Mobile
- Sidebar colapsível com overlay
- Header compacto
- Formulários adaptados
- Tabelas com scroll horizontal

## 🔄 Gerenciamento de Estado

### Zustand Store
```typescript
interface AppState {
  // Sessões
  sessions: Session[];
  selectedSession: Session | null;
  
  // UI State
  theme: 'light' | 'dark' | 'system';
  sidebarOpen: boolean;
  isMobile: boolean;
  
  // Configurações
  apiUrl: string;
  refreshInterval: number;
}
```

### React Query
- Cache automático de dados da API
- Sincronização em background
- Retry automático em falhas
- Invalidação inteligente

## 🌐 Integração com API

### Cliente HTTP
```typescript
// Configuração base
const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 10000,
});

// Interceptors para tratamento de erros
apiClient.interceptors.response.use(
  (response) => response,
  (error) => handleApiError(error)
);
```

### Endpoints Principais
- `GET /sessions` - Listar sessões
- `POST /sessions` - Criar sessão
- `PUT /sessions/:id` - Atualizar sessão
- `DELETE /sessions/:id` - Remover sessão
- `GET /sessions/:id/qr` - Obter QR code
- `POST /sessions/:id/pair` - Emparelhar por telefone

## 🔔 Sistema de Notificações

### Tipos de Notificação
- **Success** - Operações bem-sucedidas
- **Error** - Erros e falhas
- **Warning** - Avisos importantes
- **Info** - Informações gerais

### Uso
```typescript
import { useNotifications } from '@/lib/hooks/use-notifications';

const notifications = useNotifications();

// Notificação simples
notifications.success('Sessão criada com sucesso!');

// Notificação com ação
notifications.error('Erro ao conectar', {
  action: {
    label: 'Tentar Novamente',
    onClick: () => retry()
  }
});
```

## 🎯 Hooks Customizados

### useMobile
Detecta dispositivos móveis e atualiza o estado global:
```typescript
const isMobile = useMobile(); // boolean
```

### useNotifications
Sistema completo de notificações:
```typescript
const { success, error, warning, info } = useNotifications();
```

### useSystemNotifications
Notificações específicas do sistema:
```typescript
const { 
  notifySessionConnected,
  notifySessionDisconnected,
  notifyWebhookFailed 
} = useSystemNotifications();
```

## 🔧 Configuração de Desenvolvimento

### ESLint
```json
{
  "extends": ["next/core-web-vitals"],
  "rules": {
    "@typescript-eslint/no-unused-vars": "error",
    "prefer-const": "error"
  }
}
```

### Tailwind CSS
```javascript
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        border: "hsl(var(--border))",
        background: "hsl(var(--background))",
        // ... outras cores
      }
    }
  }
}
```

## 📦 Build e Deploy

### Build
```bash
npm run build
```

### Deploy
O projeto pode ser deployado em:
- **Vercel** (recomendado)
- **Netlify**
- **Docker**
- **Servidor próprio**

### Docker
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🆘 Suporte

Para suporte e dúvidas:
- Abra uma issue no GitHub
- Consulte a documentação da API
- Entre em contato com a equipe de desenvolvimento

---

Desenvolvido com ❤️ pela equipe WazMeow
