import { toast } from "sonner";
import { CheckCircle, AlertCircle, Info, AlertTriangle } from "lucide-react";

export interface NotificationOptions {
  title?: string;
  description?: string;
  duration?: number;
  action?: {
    label: string;
    onClick: () => void;
  };
}

export const useNotifications = () => {
  const success = (message: string, options?: NotificationOptions) => {
    toast.success(message, {
      description: options?.description,
      duration: options?.duration || 4000,
      action: options?.action ? {
        label: options.action.label,
        onClick: options.action.onClick,
      } : undefined,
    });
  };

  const error = (message: string, options?: NotificationOptions) => {
    toast.error(message, {
      description: options?.description,
      duration: options?.duration || 6000,
      action: options?.action ? {
        label: options.action.label,
        onClick: options.action.onClick,
      } : undefined,
    });
  };

  const warning = (message: string, options?: NotificationOptions) => {
    toast.warning(message, {
      description: options?.description,
      duration: options?.duration || 5000,
      action: options?.action ? {
        label: options.action.label,
        onClick: options.action.onClick,
      } : undefined,
    });
  };

  const info = (message: string, options?: NotificationOptions) => {
    toast.info(message, {
      description: options?.description,
      duration: options?.duration || 4000,
      action: options?.action ? {
        label: options.action.label,
        onClick: options.action.onClick,
      } : undefined,
    });
  };

  const loading = (message: string, options?: NotificationOptions) => {
    return toast.loading(message, {
      description: options?.description,
    });
  };

  const dismiss = (toastId?: string | number) => {
    if (toastId) {
      toast.dismiss(toastId);
    } else {
      toast.dismiss();
    }
  };

  const promise = <T>(
    promise: Promise<T>,
    {
      loading: loadingMessage,
      success: successMessage,
      error: errorMessage,
    }: {
      loading: string;
      success: string | ((data: T) => string);
      error: string | ((error: any) => string);
    }
  ) => {
    return toast.promise(promise, {
      loading: loadingMessage,
      success: successMessage,
      error: errorMessage,
    });
  };

  return {
    success,
    error,
    warning,
    info,
    loading,
    dismiss,
    promise,
  };
};

// Tipos específicos para notificações do sistema
export interface SystemNotification {
  id: string;
  type: "session_connected" | "session_disconnected" | "webhook_failed" | "api_error" | "system_update";
  title: string;
  message: string;
  timestamp: Date;
  read: boolean;
  sessionId?: string;
  webhookId?: string;
}

// Hook para notificações do sistema
export const useSystemNotifications = () => {
  const notifications = useNotifications();

  const notifySessionConnected = (sessionName: string, sessionId: string) => {
    notifications.success(`Sessão conectada`, {
      description: `A sessão "${sessionName}" foi conectada com sucesso`,
      action: {
        label: "Ver Sessão",
        onClick: () => {
          window.location.href = `/sessions/${sessionId}`;
        },
      },
    });
  };

  const notifySessionDisconnected = (sessionName: string, sessionId: string) => {
    notifications.warning(`Sessão desconectada`, {
      description: `A sessão "${sessionName}" foi desconectada`,
      action: {
        label: "Reconectar",
        onClick: () => {
          // Implementar lógica de reconexão
          console.log("Reconnecting session:", sessionId);
        },
      },
    });
  };

  const notifyWebhookFailed = (webhookUrl: string, error: string) => {
    notifications.error(`Webhook falhou`, {
      description: `Falha ao enviar para ${webhookUrl}: ${error}`,
      action: {
        label: "Ver Detalhes",
        onClick: () => {
          window.location.href = "/webhooks";
        },
      },
    });
  };

  const notifyApiError = (endpoint: string, error: string) => {
    notifications.error(`Erro na API`, {
      description: `Erro ao acessar ${endpoint}: ${error}`,
      duration: 8000,
    });
  };

  const notifySystemUpdate = (version: string) => {
    notifications.info(`Sistema atualizado`, {
      description: `WazMeow foi atualizado para a versão ${version}`,
      duration: 6000,
      action: {
        label: "Ver Changelog",
        onClick: () => {
          window.open("https://github.com/wazmeow/releases", "_blank");
        },
      },
    });
  };

  const notifyQRCodeReady = (sessionName: string) => {
    notifications.info(`QR Code disponível`, {
      description: `O QR Code para a sessão "${sessionName}" está pronto para escaneamento`,
      duration: 10000,
    });
  };

  const notifyMessageSent = (to: string, sessionName: string) => {
    notifications.success(`Mensagem enviada`, {
      description: `Mensagem enviada para ${to} via ${sessionName}`,
      duration: 3000,
    });
  };

  const notifyBulkOperation = (operation: string, count: number, success: number) => {
    if (success === count) {
      notifications.success(`${operation} concluída`, {
        description: `${success} de ${count} operações realizadas com sucesso`,
      });
    } else {
      notifications.warning(`${operation} parcialmente concluída`, {
        description: `${success} de ${count} operações realizadas com sucesso`,
      });
    }
  };

  return {
    notifySessionConnected,
    notifySessionDisconnected,
    notifyWebhookFailed,
    notifyApiError,
    notifySystemUpdate,
    notifyQRCodeReady,
    notifyMessageSent,
    notifyBulkOperation,
  };
};