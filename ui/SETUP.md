# Setup Local - WazMeow Manager Frontend

## 🚀 Guia de Instalação e Execução Local

### Pré-requisitos

- **Node.js** 18.0 ou superior
- **npm** ou **yarn**
- **Git**

### 1. Clone do Repositório

```bash
git clone <repository-url>
cd wazmeow/ui
```

### 2. Instalação das Dependências

```bash
# Usando npm
npm install

# Ou usando yarn
yarn install
```

### 3. Configuração do Ambiente

```bash
# Copie o arquivo de exemplo
cp .env.example .env.local

# Edite as configurações conforme necessário
nano .env.local
```

#### Variáveis de Ambiente

```env
# URL da API WazMeow (backend)
NEXT_PUBLIC_API_URL=http://localhost:8080

# Ambiente de desenvolvimento
NODE_ENV=development

# Porta do frontend (opcional)
PORT=3000

# Informações da aplicação
NEXT_PUBLIC_APP_NAME=WazMeow Manager
NEXT_PUBLIC_APP_VERSION=1.0.0
```

### 4. Execução em Desenvolvimento

```bash
# Iniciar servidor de desenvolvimento
npm run dev

# Ou com yarn
yarn dev
```

A aplicação estará disponível em: **http://localhost:3000**

### 5. Verificação da Instalação

Após iniciar o servidor, você deve ver:

1. ✅ Página de dashboard carregando
2. ✅ Sidebar com navegação
3. ✅ Tema claro/escuro funcionando
4. ✅ Responsividade mobile

## 🔧 Scripts Disponíveis

```bash
# Desenvolvimento
npm run dev          # Inicia servidor de desenvolvimento
npm run build        # Build para produção
npm run start        # Executa build de produção
npm run lint         # Executa linting
npm run lint:fix     # Corrige problemas de linting automaticamente

# Utilitários
npm run type-check   # Verifica tipos TypeScript
npm run clean        # Limpa cache e builds
```

## 🏗️ Estrutura de Desenvolvimento

### Fluxo de Trabalho

1. **Desenvolvimento**: `npm run dev`
2. **Linting**: `npm run lint`
3. **Build**: `npm run build`
4. **Teste**: `npm run start`

### Hot Reload

O Next.js oferece hot reload automático:
- ✅ Mudanças em componentes
- ✅ Mudanças em estilos
- ✅ Mudanças em páginas
- ✅ Mudanças em configurações

## 🔌 Integração com Backend

### Configuração da API

O frontend se conecta com a API WazMeow através da URL configurada em `NEXT_PUBLIC_API_URL`.

#### Endpoints Principais

```typescript
// Sessões
GET    /sessions           # Listar sessões
POST   /sessions           # Criar sessão
PUT    /sessions/:id       # Atualizar sessão
DELETE /sessions/:id       # Remover sessão
GET    /sessions/:id/qr    # Obter QR code
POST   /sessions/:id/pair  # Emparelhar por telefone

// Webhooks
GET    /webhooks           # Listar webhooks
POST   /webhooks           # Criar webhook
PUT    /webhooks/:id       # Atualizar webhook
DELETE /webhooks/:id       # Remover webhook
```

### Testando sem Backend

O frontend inclui dados mock para desenvolvimento sem backend:

```typescript
// Mock habilitado quando API não responde
const mockData = {
  sessions: [...],
  webhooks: [...],
  // ...
};
```

## 🎨 Desenvolvimento de Componentes

### Criando Novos Componentes

```bash
# Estrutura recomendada
src/components/
├── forms/              # Formulários específicos
├── layout/             # Componentes de layout
├── ui/                 # Componentes base (Shadcn/ui)
└── [feature]/          # Componentes por funcionalidade
```

### Exemplo de Componente

```typescript
// src/components/example/my-component.tsx
"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface MyComponentProps {
  title: string;
  children: React.ReactNode;
}

export function MyComponent({ title, children }: MyComponentProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        {children}
      </CardContent>
    </Card>
  );
}
```

### Adicionando Páginas

```typescript
// src/app/(dashboard)/nova-pagina/page.tsx
export default function NovaPaginaPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Nova Página</h1>
        <p className="text-muted-foreground">
          Descrição da nova página
        </p>
      </div>
      
      {/* Conteúdo da página */}
    </div>
  );
}
```

## 🎯 Hooks Customizados

### Criando Hooks

```typescript
// src/lib/hooks/use-example.ts
"use client";

import { useState, useEffect } from "react";

export function useExample() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Lógica do hook
    setLoading(false);
  }, []);

  return { data, loading };
}
```

## 🔄 Gerenciamento de Estado

### Adicionando ao Store

```typescript
// src/lib/stores/app-store.ts
interface AppState {
  // Novo estado
  newFeature: boolean;
  
  // Nova ação
  setNewFeature: (enabled: boolean) => void;
}

// Implementação
export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      newFeature: false,
      setNewFeature: (enabled) => set({ newFeature: enabled }),
    }),
    { name: 'app-store' }
  )
);
```

## 🎨 Estilos e Temas

### Customizando Cores

```css
/* src/app/globals.css */
:root {
  --custom-color: 220 14.3% 95.9%;
}

.dark {
  --custom-color: 220 14.3% 4.1%;
}
```

### Usando no Tailwind

```typescript
// tailwind.config.ts
module.exports = {
  theme: {
    extend: {
      colors: {
        custom: "hsl(var(--custom-color))",
      }
    }
  }
}
```

## 🔍 Debugging

### Ferramentas de Debug

1. **React DevTools** - Extensão do navegador
2. **Next.js DevTools** - Debug de performance
3. **Zustand DevTools** - Estado da aplicação
4. **Network Tab** - Requisições HTTP

### Logs de Debug

```typescript
// Habilitar logs detalhados
if (process.env.NODE_ENV === 'development') {
  console.log('Debug info:', data);
}
```

## 📱 Teste Mobile

### Testando Responsividade

1. **Chrome DevTools** - Device simulation
2. **Dispositivos reais** - Teste em smartphones/tablets
3. **Browserstack** - Teste em múltiplos dispositivos

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

## 🚨 Troubleshooting

### Problemas Comuns

#### 1. Erro de Porta em Uso
```bash
Error: Port 3000 is already in use
```
**Solução:**
```bash
# Matar processo na porta 3000
npx kill-port 3000

# Ou usar porta diferente
PORT=3001 npm run dev
```

#### 2. Erro de Dependências
```bash
Module not found: Can't resolve '@/components/...'
```
**Solução:**
```bash
# Reinstalar dependências
rm -rf node_modules package-lock.json
npm install
```

#### 3. Erro de TypeScript
```bash
Type error: Property 'x' does not exist
```
**Solução:**
```bash
# Verificar tipos
npm run type-check

# Reiniciar TypeScript server no VS Code
Ctrl+Shift+P > "TypeScript: Restart TS Server"
```

#### 4. Erro de Build
```bash
Build failed with errors
```
**Solução:**
```bash
# Limpar cache
npm run clean
rm -rf .next

# Build novamente
npm run build
```

### Logs Úteis

```bash
# Logs detalhados do Next.js
DEBUG=* npm run dev

# Logs apenas do Next.js
DEBUG=next:* npm run dev
```

## 📞 Suporte

### Recursos de Ajuda

1. **Documentação Next.js**: https://nextjs.org/docs
2. **Documentação Shadcn/ui**: https://ui.shadcn.com
3. **Documentação Tailwind**: https://tailwindcss.com/docs
4. **Issues do GitHub**: Para reportar bugs

### Contato

- **Email**: dev@wazmeow.com
- **Discord**: WazMeow Community
- **GitHub Issues**: Para bugs e features

---

**Dica**: Mantenha sempre as dependências atualizadas e siga as boas práticas de desenvolvimento React/Next.js!