"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useNotifications, useSystemNotifications } from "@/lib/hooks/use-notifications";
import { 
  CheckCircle, 
  AlertCircle, 
  Info, 
  AlertTriangle,
  Loader2,
  Bell
} from "lucide-react";

export function NotificationDemo() {
  const notifications = useNotifications();
  const systemNotifications = useSystemNotifications();

  const handleSuccessNotification = () => {
    notifications.success("Operação realizada com sucesso!", {
      description: "Sua ação foi processada corretamente.",
    });
  };

  const handleErrorNotification = () => {
    notifications.error("Erro ao processar solicitação", {
      description: "Ocorreu um erro inesperado. Tente novamente.",
      action: {
        label: "Tentar Novamente",
        onClick: () => {
          notifications.info("Tentando novamente...");
        },
      },
    });
  };

  const handleWarningNotification = () => {
    notifications.warning("Atenção necessária", {
      description: "Esta ação pode ter consequências importantes.",
      action: {
        label: "Entendi",
        onClick: () => {
          notifications.success("Confirmado!");
        },
      },
    });
  };

  const handleInfoNotification = () => {
    notifications.info("Informação importante", {
      description: "Aqui está uma informação que você deve saber.",
    });
  };

  const handleLoadingNotification = () => {
    const loadingToast = notifications.loading("Processando...", {
      description: "Aguarde enquanto processamos sua solicitação.",
    });

    // Simular processo assíncrono
    setTimeout(() => {
      notifications.dismiss(loadingToast);
      notifications.success("Processo concluído!");
    }, 3000);
  };

  const handlePromiseNotification = () => {
    const mockPromise = new Promise((resolve, reject) => {
      setTimeout(() => {
        Math.random() > 0.5 ? resolve("Sucesso!") : reject("Erro!");
      }, 2000);
    });

    notifications.promise(mockPromise, {
      loading: "Executando operação...",
      success: "Operação concluída com sucesso!",
      error: "Falha na operação",
    });
  };

  const handleSessionConnected = () => {
    systemNotifications.notifySessionConnected("Minha Sessão", "session-123");
  };

  const handleSessionDisconnected = () => {
    systemNotifications.notifySessionDisconnected("Minha Sessão", "session-123");
  };

  const handleWebhookFailed = () => {
    systemNotifications.notifyWebhookFailed(
      "https://api.exemplo.com/webhook",
      "Timeout na conexão"
    );
  };

  const handleQRCodeReady = () => {
    systemNotifications.notifyQRCodeReady("Sessão Principal");
  };

  const handleMessageSent = () => {
    systemNotifications.notifyMessageSent("+55 11 99999-9999", "Sessão Principal");
  };

  const handleBulkOperation = () => {
    systemNotifications.notifyBulkOperation("Envio em massa", 10, 8);
  };

  return (
    <Card className="w-full max-w-4xl">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Bell className="h-5 w-5" />
          Demonstração do Sistema de Notificações
        </CardTitle>
        <CardDescription>
          Teste os diferentes tipos de notificações disponíveis no sistema
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Notificações Básicas */}
        <div className="space-y-3">
          <h3 className="text-lg font-semibold">Notificações Básicas</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
            <Button 
              onClick={handleSuccessNotification}
              variant="default"
              className="flex items-center gap-2"
            >
              <CheckCircle className="h-4 w-4" />
              Sucesso
            </Button>
            
            <Button 
              onClick={handleErrorNotification}
              variant="destructive"
              className="flex items-center gap-2"
            >
              <AlertCircle className="h-4 w-4" />
              Erro
            </Button>
            
            <Button 
              onClick={handleWarningNotification}
              variant="outline"
              className="flex items-center gap-2 text-yellow-600 border-yellow-300 hover:bg-yellow-50"
            >
              <AlertTriangle className="h-4 w-4" />
              Aviso
            </Button>
            
            <Button 
              onClick={handleInfoNotification}
              variant="secondary"
              className="flex items-center gap-2"
            >
              <Info className="h-4 w-4" />
              Info
            </Button>
          </div>
        </div>

        {/* Notificações Avançadas */}
        <div className="space-y-3">
          <h3 className="text-lg font-semibold">Notificações Avançadas</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
            <Button 
              onClick={handleLoadingNotification}
              variant="outline"
              className="flex items-center gap-2"
            >
              <Loader2 className="h-4 w-4" />
              Loading Toast
            </Button>
            
            <Button 
              onClick={handlePromiseNotification}
              variant="outline"
              className="flex items-center gap-2"
            >
              <Loader2 className="h-4 w-4" />
              Promise Toast
            </Button>
          </div>
        </div>

        {/* Notificações do Sistema */}
        <div className="space-y-3">
          <h3 className="text-lg font-semibold">Notificações do Sistema</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
            <Button 
              onClick={handleSessionConnected}
              variant="outline"
              size="sm"
            >
              Sessão Conectada
            </Button>
            
            <Button 
              onClick={handleSessionDisconnected}
              variant="outline"
              size="sm"
            >
              Sessão Desconectada
            </Button>
            
            <Button 
              onClick={handleWebhookFailed}
              variant="outline"
              size="sm"
            >
              Webhook Falhou
            </Button>
            
            <Button 
              onClick={handleQRCodeReady}
              variant="outline"
              size="sm"
            >
              QR Code Pronto
            </Button>
            
            <Button 
              onClick={handleMessageSent}
              variant="outline"
              size="sm"
            >
              Mensagem Enviada
            </Button>
            
            <Button 
              onClick={handleBulkOperation}
              variant="outline"
              size="sm"
            >
              Operação em Massa
            </Button>
          </div>
        </div>

        {/* Controles */}
        <div className="space-y-3">
          <h3 className="text-lg font-semibold">Controles</h3>
          <Button 
            onClick={() => notifications.dismiss()}
            variant="outline"
            className="w-full"
          >
            Limpar Todas as Notificações
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}