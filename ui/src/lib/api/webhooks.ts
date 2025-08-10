import { apiClient } from './client';
import {
  WebhookConfig,
  SetWebhookRequest,
  SupportedEventsResponse,
} from '@/lib/types/api';

export class WebhooksApi {
  // Configurar webhook para uma sessão
  static async setWebhook(
    sessionId: string,
    data: SetWebhookRequest
  ): Promise<void> {
    const response = await apiClient.post(
      `/sessions/${sessionId}/webhook/set`,
      data
    );
    if (!response.success) {
      throw new Error(response.error || 'Erro ao configurar webhook');
    }
  }

  // Obter configuração do webhook
  static async getWebhook(sessionId: string): Promise<WebhookConfig> {
    const response = await apiClient.get<WebhookConfig>(
      `/sessions/${sessionId}/webhook/find`
    );
    if (!response.success || !response.data) {
      throw new Error(response.error || 'Erro ao obter configuração do webhook');
    }
    return response.data;
  }

  // Obter eventos suportados
  static async getSupportedEvents(): Promise<SupportedEventsResponse> {
    const response = await apiClient.get<SupportedEventsResponse>('/webhook/events');
    if (!response.success || !response.data) {
      throw new Error(response.error || 'Erro ao obter eventos suportados');
    }
    return response.data;
  }

  // Testar conectividade do webhook
  static async testWebhook(
    sessionId: string,
    webhookUrl: string
  ): Promise<boolean> {
    try {
      // Implementar teste de conectividade
      // Por enquanto, apenas simular o teste
      const response = await fetch(webhookUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          test: true,
          timestamp: Date.now(),
          session_id: sessionId,
        }),
      });
      
      return response.ok;
    } catch (error) {
      return false;
    }
  }

  // Remover webhook
  static async removeWebhook(sessionId: string): Promise<void> {
    const response = await apiClient.post(
      `/sessions/${sessionId}/webhook/set`,
      {
        webhookurl: '',
        events: [],
        enabled: false,
      }
    );
    if (!response.success) {
      throw new Error(response.error || 'Erro ao remover webhook');
    }
  }
}