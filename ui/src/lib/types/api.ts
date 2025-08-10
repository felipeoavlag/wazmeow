// Tipos baseados na API WazMeow

// Status da sessão
export type SessionStatus = 'disconnected' | 'connecting' | 'connected';

// Configuração de Proxy
export interface ProxyConfig {
  type: 'http' | 'socks5';
  host: string;
  port: number;
  username?: string;
  password?: string;
}

// Sessão
export interface Session {
  id: string;
  name: string;
  status: SessionStatus;
  phone?: string;
  device_jid?: string;
  proxy_config?: ProxyConfig;
  webhook_url?: string;
  events?: string;
  created_at: string;
  updated_at: string;
}

// Informações da sessão com status de conexão
export interface SessionInfo extends Session {
  is_connected: boolean;
  is_logged_in: boolean;
}

// Resposta de QR Code
export interface QRResponse {
  qr_code?: string;
  status: string;
}

// Resposta de código de emparelhamento
export interface PairCodeResponse {
  code: string;
  status: string;
}

// Configuração de Webhook
export interface WebhookConfig {
  webhookurl: string;
  events: string[];
  enabled?: boolean;
}

// Resposta padrão da API
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

// Resposta simples
export interface SimpleResponse {
  details: string;
}

// Resposta com timestamp
export interface TimestampedResponse {
  details: string;
  timestamp: number;
}

// Resposta de envio de mensagem
export interface SendMessageResponse {
  details: string;
  timestamp: number;
  id: string;
}

// Requests para criação de sessão
export interface CreateSessionRequest {
  name: string;
  webhookUrl?: string;
  proxy?: ProxyConfig;
}

// Request para emparelhar telefone
export interface PairPhoneRequest {
  phone: string;
}

// Request para configurar proxy
export interface SetProxyRequest {
  type: 'http' | 'socks5';
  host: string;
  port: number;
  username?: string;
  password?: string;
}

// Request para configurar webhook
export interface SetWebhookRequest {
  webhookurl: string;
  events: string[];
  enabled?: boolean;
}

// Contexto de mensagem (para reply e menções)
export interface ContextInfo {
  stanza_id?: string;
  participant?: string;
  mentioned_jid?: string[];
}

// Request para envio de mensagem de texto
export interface SendTextMessageRequest {
  phone: string;
  body: string;
  id?: string;
  context_info?: ContextInfo;
}

// Request para envio de mídia
export interface SendMediaMessageRequest {
  phone: string;
  media_data: string;
  caption?: string;
  mime_type?: string;
  id?: string;
  context_info?: ContextInfo;
}

// Estrutura de botão
export interface ButtonStruct {
  button_id: string;
  button_text: string;
}

// Request para envio de botões
export interface SendButtonsMessageRequest {
  phone: string;
  title: string;
  buttons: ButtonStruct[];
  id?: string;
  context_info?: ContextInfo;
}

// Item de lista
export interface ListItem {
  title: string;
  desc?: string;
  row_id: string;
}

// Seção de lista
export interface ListSection {
  title: string;
  rows: ListItem[];
}

// Request para envio de lista
export interface SendListMessageRequest {
  phone: string;
  button_text: string;
  desc: string;
  top_text: string;
  sections: ListSection[];
  footer_text?: string;
  id?: string;
}

// Request para envio de enquete
export interface SendPollMessageRequest {
  phone: string;
  header: string;
  options: string[];
  id?: string;
}

// Request para deletar mensagem
export interface DeleteMessageRequest {
  phone: string;
  id: string;
}

// Request para reagir a mensagem
export interface ReactMessageRequest {
  phone: string;
  body: string; // Emoji ou "remove"
  id: string;
}

// Eventos suportados pelo webhook
export interface SupportedEventsResponse {
  events: string[];
  groups: Record<string, string[]>;
  wildcards: string[];
  examples: Record<string, any>;
}

// Métricas da sessão
export interface SessionMetrics {
  messages_sent: number;
  messages_received: number;
  uptime: number;
  last_activity: string;
  webhook_success_rate: number;
}

// Dashboard data
export interface DashboardData {
  total_sessions: number;
  active_sessions: number;
  total_messages: number;
  webhook_events: number;
  recent_activities: Activity[];
}

// Atividade recente
export interface Activity {
  id: string;
  type: 'session_created' | 'session_connected' | 'message_sent' | 'webhook_configured';
  description: string;
  timestamp: string;
  session_id?: string;
}

// Filtros para listagem
export interface SessionFilters {
  status?: SessionStatus;
  search?: string;
  has_webhook?: boolean;
  has_proxy?: boolean;
}

// Paginação
export interface PaginationParams {
  page: number;
  limit: number;
}

// Resposta paginada
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}