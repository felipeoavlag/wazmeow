import { apiClient } from './client';
import {
  Session,
  SessionInfo,
  CreateSessionRequest,
  QRResponse,
  PairPhoneRequest,
  PairCodeResponse,
  SetProxyRequest,
  SessionFilters,
  PaginatedResponse,
  PaginationParams,
} from '@/lib/types/api';

export class SessionsApi {
  // Listar sessões
  static async listSessions(
    filters?: SessionFilters,
    pagination?: PaginationParams
  ): Promise<Session[]> {
    const params = {
      ...filters,
      ...pagination,
    };
    
    const response = await apiClient.get<Session[]>('/sessions', params);
    return response.data || [];
  }

  // Criar sessão
  static async createSession(data: CreateSessionRequest): Promise<Session> {
    const response = await apiClient.post<Session>('/sessions/add', data);
    if (!response.success || !response.data) {
      throw new Error(response.error || 'Erro ao criar sessão');
    }
    return response.data;
  }

  // Obter informações da sessão
  static async getSessionInfo(sessionId: string): Promise<SessionInfo> {
    const response = await apiClient.get<SessionInfo>(`/sessions/${sessionId}`);
    if (!response.success || !response.data) {
      throw new Error(response.error || 'Erro ao obter informações da sessão');
    }
    return response.data;
  }

  // Deletar sessão
  static async deleteSession(sessionId: string): Promise<void> {
    const response = await apiClient.delete(`/sessions/${sessionId}`);
    if (!response.success) {
      throw new Error(response.error || 'Erro ao deletar sessão');
    }
  }

  // Conectar sessão
  static async connectSession(sessionId: string): Promise<void> {
    const response = await apiClient.post(`/sessions/${sessionId}/connect`);
    if (!response.success) {
      throw new Error(response.error || 'Erro ao conectar sessão');
    }
  }

  // Fazer logout da sessão
  static async logoutSession(sessionId: string): Promise<void> {
    const response = await apiClient.post(`/sessions/${sessionId}/logout`);
    if (!response.success) {
      throw new Error(response.error || 'Erro ao fazer logout da sessão');
    }
  }

  // Obter QR Code
  static async getQRCode(sessionId: string): Promise<QRResponse> {
    const response = await apiClient.get<QRResponse>(`/sessions/${sessionId}/qr`);
    if (!response.success || !response.data) {
      throw new Error(response.error || 'Erro ao obter QR Code');
    }
    return response.data;
  }

  // Emparelhar telefone
  static async pairPhone(
    sessionId: string,
    data: PairPhoneRequest
  ): Promise<PairCodeResponse> {
    const response = await apiClient.post<PairCodeResponse>(
      `/sessions/${sessionId}/pair`,
      data
    );
    if (!response.success || !response.data) {
      throw new Error(response.error || 'Erro ao emparelhar telefone');
    }
    return response.data;
  }

  // Configurar proxy
  static async setProxy(
    sessionId: string,
    data: SetProxyRequest
  ): Promise<void> {
    const response = await apiClient.post(`/sessions/${sessionId}/proxy`, data);
    if (!response.success) {
      throw new Error(response.error || 'Erro ao configurar proxy');
    }
  }

  // Verificar status de saúde da API
  static async healthCheck(): Promise<any> {
    const response = await apiClient.get('/health');
    return response.data;
  }
}