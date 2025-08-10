# Arquitetura do Frontend - WazMeow Manager

## 📋 Visão Geral

O WazMeow Manager é uma aplicação web moderna construída com Next.js 14, seguindo os princípios de Clean Architecture e padrões de desenvolvimento React atuais.

## 🏗️ Arquitetura Geral

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
├─────────────────────────────────────────────────────────────┤
│  Pages & Components  │  Hooks  │  Providers  │  UI Library  │
├─────────────────────────────────────────────────────────────┤
│                    Application Layer                        │
├─────────────────────────────────────────────────────────────┤
│   State Management   │  API Client  │  Business Logic      │
├─────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                     │
├─────────────────────────────────────────────────────────────┤
│   HTTP Client   │   Local Storage   │   External APIs      │
└─────────────────────────────────────────────────────────────┘
```

## 🎯 Princípios Arquiteturais

### 1. Separation of Concerns
- **Presentation**: Componentes React focados apenas na UI
- **Business Logic**: Hooks customizados e stores
- **Data Access**: Cliente HTTP e cache

### 2. Dependency Inversion
- Componentes dependem de abstrações (hooks/stores)
- Implementações concretas são injetadas via providers

### 3. Single Responsibility
- Cada componente tem uma responsabilidade específica
- Hooks customizados encapsulam lógica reutilizável

### 4. Open/Closed Principle
- Componentes extensíveis via props e composition
- Sistema de temas e configurações flexível

## 📁 Estrutura Detalhada

### App Router (Next.js 14)
```
app/
├── layout.tsx              # Layout raiz com providers
├── page.tsx               # Página inicial (redirect)
├── globals.css            # Estilos globais
└── (dashboard)/           # Grupo de rotas do dashboard
    ├── layout.tsx         # Layout do dashboard
    ├── page.tsx          # Dashboard principal
    ├── sessions/         # Gerenciamento de sessões
    ├── webhooks/         # Configuração de webhooks
    ├── monitoring/       # Monitoramento
    ├── logs/            # Logs e eventos
    ├── proxy/           # Configuração de proxy
    ├── notifications/   # Sistema de notificações
    └── settings/        # Configurações
```

### Componentes
```
components/
├── forms/                 # Formulários específicos
│   ├── session-form.tsx   # Formulário de sessão
│   └── proxy-config-form.tsx # Configuração de proxy
├── layout/               # Componentes de layout
│   ├── header.tsx        # Cabeçalho
│   └── sidebar.tsx       # Barra lateral
├── providers/            # Context providers
│   ├── theme-provider.tsx # Tema
│   └── query-provider.tsx # React Query
└── ui/                   # Componentes base (Shadcn/ui)
    ├── button.tsx
    ├── card.tsx
    ├── input.tsx
    └── ...
```

### Biblioteca (lib)
```
lib/
├── api/                  # Cliente HTTP e endpoints
│   ├── client.ts         # Configuração do Axios
│   ├── sessions.ts       # Endpoints de sessões
│   └── webhooks.ts       # Endpoints de webhooks
├── hooks/                # Hooks customizados
│   ├── use-mobile.ts     # Detecção mobile
│   └── use-notifications.ts # Sistema de notificações
├── stores/               # Gerenciamento de estado
│   └── app-store.ts      # Store principal (Zustand)
├── types/                # Definições de tipos
│   └── api.ts            # Tipos da API
└── utils.ts              # Utilitários gerais
```

## 🔄 Fluxo de Dados

### 1. Requisições HTTP
```
Component → Hook → API Client → HTTP Request → Backend
                ↓
Component ← Hook ← React Query ← HTTP Response ← Backend
```

### 2. Gerenciamento de Estado
```
Component → Action → Zustand Store → State Update → Re-render
```

### 3. Notificações
```
Action → useNotifications → Sonner → Toast Display
```

## 🎨 Sistema de Design

### Design Tokens
```typescript
// Cores
const colors = {
  primary: 'hsl(221.2 83.2% 53.3%)',
  secondary: 'hsl(210 40% 98%)',
  muted: 'hsl(210 40% 96%)',
  // ...
}

// Espaçamento
const spacing = {
  xs: '0.25rem',
  sm: '0.5rem',
  md: '1rem',
  lg: '1.5rem',
  xl: '2rem',
  // ...
}
```

### Componentes Base
- **Atomic Design**: Atoms → Molecules → Organisms → Templates → Pages
- **Composition Pattern**: Componentes compostos via children
- **Render Props**: Flexibilidade para casos específicos

## 📱 Responsividade

### Breakpoints
```css
/* Mobile First */
.container {
  /* Mobile: < 768px */
  padding: 1rem;
}

@media (min-width: 768px) {
  /* Tablet */
  .container {
    padding: 1.5rem;
  }
}

@media (min-width: 1024px) {
  /* Desktop */
  .container {
    padding: 2rem;
  }
}
```

### Layout Adaptativo
- **Mobile**: Sidebar overlay, header compacto
- **Tablet**: Sidebar colapsível, layout híbrido
- **Desktop**: Sidebar fixa, layout completo

## 🔧 Gerenciamento de Estado

### Zustand Store
```typescript
interface AppState {
  // Estado da aplicação
  sessions: Session[];
  selectedSession: Session | null;
  
  // Estado da UI
  theme: Theme;
  sidebarOpen: boolean;
  isMobile: boolean;
  
  // Ações
  setSessions: (sessions: Session[]) => void;
  toggleSidebar: () => void;
  // ...
}
```

### React Query
```typescript
// Cache e sincronização
const { data, isLoading, error } = useQuery({
  queryKey: ['sessions'],
  queryFn: fetchSessions,
  staleTime: 5 * 60 * 1000, // 5 minutos
  refetchInterval: 30 * 1000, // 30 segundos
});
```

## 🌐 Integração com API

### Cliente HTTP
```typescript
// Configuração base
const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 10000,
});

// Interceptors
apiClient.interceptors.request.use(addAuthToken);
apiClient.interceptors.response.use(handleResponse, handleError);
```

### Tratamento de Erros
```typescript
const handleApiError = (error: AxiosError) => {
  if (error.response?.status === 401) {
    // Redirect para login
  } else if (error.response?.status >= 500) {
    // Erro do servidor
    notifications.error('Erro interno do servidor');
  } else {
    // Outros erros
    notifications.error(error.message);
  }
};
```

## 🔔 Sistema de Notificações

### Arquitetura
```
Action → useNotifications → Sonner Provider → Toast Component
```

### Tipos de Notificação
- **Toast**: Notificações temporárias
- **Alert**: Avisos persistentes
- **Modal**: Confirmações importantes

## 🎯 Performance

### Otimizações
1. **Code Splitting**: Lazy loading de páginas
2. **Tree Shaking**: Eliminação de código não usado
3. **Image Optimization**: Next.js Image component
4. **Bundle Analysis**: Análise do tamanho do bundle

### Métricas
- **FCP**: First Contentful Paint < 1.5s
- **LCP**: Largest Contentful Paint < 2.5s
- **CLS**: Cumulative Layout Shift < 0.1
- **FID**: First Input Delay < 100ms

## 🔒 Segurança

### Medidas Implementadas
1. **CSP**: Content Security Policy
2. **HTTPS**: Comunicação segura
3. **Input Validation**: Validação com Zod
4. **XSS Protection**: Sanitização de dados

### Autenticação
```typescript
// Token JWT no header
const authToken = localStorage.getItem('auth_token');
apiClient.defaults.headers.Authorization = `Bearer ${authToken}`;
```

## 🧪 Testabilidade

### Estratégia de Testes
1. **Unit Tests**: Hooks e utilitários
2. **Component Tests**: Componentes isolados
3. **Integration Tests**: Fluxos completos
4. **E2E Tests**: Cenários de usuário

### Ferramentas
- **Jest**: Framework de testes
- **React Testing Library**: Testes de componentes
- **MSW**: Mock Service Worker
- **Cypress**: Testes E2E

## 📦 Build e Deploy

### Pipeline
```
Code → Lint → Test → Build → Deploy
```

### Ambientes
- **Development**: Local development
- **Staging**: Testes de integração
- **Production**: Ambiente de produção

### Configuração
```javascript
// next.config.js
module.exports = {
  output: 'standalone',
  images: {
    domains: ['api.wazmeow.com'],
  },
  env: {
    CUSTOM_KEY: process.env.CUSTOM_KEY,
  },
};
```

## 🔄 Versionamento

### Semantic Versioning
- **MAJOR**: Mudanças incompatíveis
- **MINOR**: Novas funcionalidades
- **PATCH**: Correções de bugs

### Git Flow
```
main ← develop ← feature/new-feature
     ← hotfix/critical-fix
```

## 📈 Monitoramento

### Métricas
- **Performance**: Web Vitals
- **Errors**: Error boundaries
- **Usage**: Analytics
- **API**: Response times

### Ferramentas
- **Vercel Analytics**: Performance
- **Sentry**: Error tracking
- **Google Analytics**: Usage
- **LogRocket**: Session replay

---

Esta arquitetura garante escalabilidade, manutenibilidade e performance para o WazMeow Manager.