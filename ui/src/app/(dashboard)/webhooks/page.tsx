"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import { 
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { 
  Webhook, 
  Plus, 
  Settings, 
  Trash2, 
  TestTube,
  CheckCircle,
  XCircle,
  AlertCircle,
  Globe,
  Activity,
  Zap
} from "lucide-react";
import { Session, WebhookConfig } from "@/lib/types/api";

// Mock data
const mockSessions: Session[] = [
  {
    id: "1",
    name: "Sessão Principal",
    status: "connected",
    webhook_url: "https://webhook.example.com/wazmeow",
    events: "message,receipt,connected",
    created_at: "2024-01-15T10:30:00Z",
    updated_at: "2024-01-15T14:20:00Z",
  },
  {
    id: "2",
    name: "Suporte",
    status: "connected",
    created_at: "2024-01-14T09:15:00Z",
    updated_at: "2024-01-15T13:45:00Z",
  },
  {
    id: "3",
    name: "Marketing",
    status: "disconnected",
    webhook_url: "https://api.marketing.com/webhook",
    events: "message,presence",
    created_at: "2024-01-13T16:20:00Z",
    updated_at: "2024-01-15T08:30:00Z",
  },
];

const supportedEvents = [
  { id: "message", name: "Mensagens", description: "Mensagens recebidas" },
  { id: "receipt", name: "Confirmações", description: "Confirmações de leitura" },
  { id: "connected", name: "Conexão", description: "Status de conexão" },
  { id: "disconnected", name: "Desconexão", description: "Perda de conexão" },
  { id: "qr", name: "QR Code", description: "Geração de QR code" },
  { id: "presence", name: "Presença", description: "Status online/offline" },
  { id: "chatpresence", name: "Presença no Chat", description: "Digitando, gravando áudio" },
  { id: "groupinfo", name: "Info do Grupo", description: "Mudanças em grupos" },
  { id: "picture", name: "Foto de Perfil", description: "Mudanças de foto" },
];

interface WebhookFormData {
  sessionId: string;
  webhookUrl: string;
  events: string[];
  enabled: boolean;
}

export default function WebhooksPage() {
  const [sessions, setSessions] = useState<Session[]>(mockSessions);
  const [showConfigDialog, setShowConfigDialog] = useState(false);
  const [selectedSession, setSelectedSession] = useState<Session | null>(null);
  const [formData, setFormData] = useState<WebhookFormData>({
    sessionId: "",
    webhookUrl: "",
    events: [],
    enabled: true,
  });
  const [testResults, setTestResults] = useState<Record<string, boolean>>({});

  const handleConfigureWebhook = (session: Session) => {
    setSelectedSession(session);
    setFormData({
      sessionId: session.id,
      webhookUrl: session.webhook_url || "",
      events: session.events ? session.events.split(",") : [],
      enabled: !!session.webhook_url,
    });
    setShowConfigDialog(true);
  };

  const handleSaveWebhook = () => {
    if (!selectedSession) return;

    const updatedSessions = sessions.map(session => 
      session.id === selectedSession.id 
        ? {
            ...session,
            webhook_url: formData.enabled ? formData.webhookUrl : undefined,
            events: formData.enabled ? formData.events.join(",") : undefined,
            updated_at: new Date().toISOString(),
          }
        : session
    );

    setSessions(updatedSessions);
    setShowConfigDialog(false);
    setSelectedSession(null);
  };

  const handleTestWebhook = async (sessionId: string, webhookUrl: string) => {
    setTestResults(prev => ({ ...prev, [sessionId]: false }));
    
    // Simular teste de webhook
    setTimeout(() => {
      const success = Math.random() > 0.3; // 70% de chance de sucesso
      setTestResults(prev => ({ ...prev, [sessionId]: success }));
    }, 2000);
  };

  const handleRemoveWebhook = (sessionId: string) => {
    const updatedSessions = sessions.map(session => 
      session.id === sessionId 
        ? {
            ...session,
            webhook_url: undefined,
            events: undefined,
            updated_at: new Date().toISOString(),
          }
        : session
    );

    setSessions(updatedSessions);
  };

  const toggleEvent = (eventId: string) => {
    setFormData(prev => ({
      ...prev,
      events: prev.events.includes(eventId)
        ? prev.events.filter(e => e !== eventId)
        : [...prev.events, eventId]
    }));
  };

  const getWebhookStatus = (session: Session) => {
    if (!session.webhook_url) {
      return { status: "none", label: "Não configurado", variant: "secondary" as const };
    }
    
    if (session.status === "connected") {
      return { status: "active", label: "Ativo", variant: "default" as const };
    }
    
    return { status: "inactive", label: "Inativo", variant: "destructive" as const };
  };

  const sessionsWithWebhooks = sessions.filter(s => s.webhook_url);
  const totalEvents = sessionsWithWebhooks.reduce((acc, s) => 
    acc + (s.events ? s.events.split(",").length : 0), 0
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Webhooks</h1>
          <p className="text-muted-foreground">
            Configure webhooks para receber eventos em tempo real
          </p>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Webhooks Ativos
            </CardTitle>
            <Webhook className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{sessionsWithWebhooks.length}</div>
            <p className="text-xs text-muted-foreground">
              de {sessions.length} sessões
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Eventos Configurados
            </CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{totalEvents}</div>
            <p className="text-xs text-muted-foreground">
              tipos de eventos
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Taxa de Sucesso
            </CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">98.5%</div>
            <p className="text-xs text-muted-foreground">
              últimas 24 horas
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Webhooks Table */}
      <Card>
        <CardHeader>
          <CardTitle>Configurações de Webhook</CardTitle>
          <CardDescription>
            Gerencie os webhooks de suas sessões
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Sessão</TableHead>
                <TableHead>URL do Webhook</TableHead>
                <TableHead>Eventos</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Ações</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sessions.map((session) => {
                const webhookStatus = getWebhookStatus(session);
                const testResult = testResults[session.id];
                
                return (
                  <TableRow key={session.id}>
                    <TableCell className="font-medium">
                      {session.name}
                    </TableCell>
                    <TableCell>
                      {session.webhook_url ? (
                        <div className="flex items-center gap-2">
                          <Globe className="h-4 w-4 text-muted-foreground" />
                          <span className="truncate max-w-xs">
                            {session.webhook_url}
                          </span>
                        </div>
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {session.events ? (
                        <div className="flex flex-wrap gap-1">
                          {session.events.split(",").slice(0, 3).map(event => (
                            <Badge key={event} variant="outline" className="text-xs">
                              {event}
                            </Badge>
                          ))}
                          {session.events.split(",").length > 3 && (
                            <Badge variant="outline" className="text-xs">
                              +{session.events.split(",").length - 3}
                            </Badge>
                          )}
                        </div>
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Badge variant={webhookStatus.variant}>
                        {webhookStatus.label}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleConfigureWebhook(session)}
                        >
                          <Settings className="h-4 w-4" />
                        </Button>
                        
                        {session.webhook_url && (
                          <>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleTestWebhook(session.id, session.webhook_url!)}
                              disabled={testResult !== undefined}
                            >
                              {testResult === undefined ? (
                                <TestTube className="h-4 w-4" />
                              ) : testResult ? (
                                <CheckCircle className="h-4 w-4 text-green-500" />
                              ) : (
                                <XCircle className="h-4 w-4 text-red-500" />
                              )}
                            </Button>
                            
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleRemoveWebhook(session.id)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Configuration Dialog */}
      <Dialog open={showConfigDialog} onOpenChange={setShowConfigDialog}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              Configurar Webhook - {selectedSession?.name}
            </DialogTitle>
            <DialogDescription>
              Configure a URL e eventos para receber notificações
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-6">
            {/* Enable/Disable */}
            <div className="flex items-center justify-between rounded-lg border p-4">
              <div className="space-y-0.5">
                <Label className="text-base">Habilitar Webhook</Label>
                <p className="text-sm text-muted-foreground">
                  Ativar ou desativar o webhook para esta sessão
                </p>
              </div>
              <Switch
                checked={formData.enabled}
                onCheckedChange={(checked) => 
                  setFormData(prev => ({ ...prev, enabled: checked }))
                }
              />
            </div>

            {formData.enabled && (
              <>
                {/* Webhook URL */}
                <div className="space-y-2">
                  <Label htmlFor="webhookUrl">URL do Webhook</Label>
                  <Input
                    id="webhookUrl"
                    placeholder="https://seu-site.com/webhook"
                    value={formData.webhookUrl}
                    onChange={(e) => 
                      setFormData(prev => ({ ...prev, webhookUrl: e.target.value }))
                    }
                  />
                  <p className="text-sm text-muted-foreground">
                    URL que receberá os eventos via POST
                  </p>
                </div>

                {/* Events Selection */}
                <div className="space-y-4">
                  <Label>Eventos para Receber</Label>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    {supportedEvents.map((event) => (
                      <div
                        key={event.id}
                        className={`flex items-center space-x-3 rounded-lg border p-3 cursor-pointer transition-colors ${
                          formData.events.includes(event.id)
                            ? "border-primary bg-primary/5"
                            : "border-muted hover:border-primary/50"
                        }`}
                        onClick={() => toggleEvent(event.id)}
                      >
                        <div className="flex-shrink-0">
                          <div className={`w-4 h-4 rounded border-2 flex items-center justify-center ${
                            formData.events.includes(event.id)
                              ? "border-primary bg-primary"
                              : "border-muted-foreground"
                          }`}>
                            {formData.events.includes(event.id) && (
                              <CheckCircle className="w-3 h-3 text-primary-foreground" />
                            )}
                          </div>
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium">{event.name}</p>
                          <p className="text-xs text-muted-foreground">
                            {event.description}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                  
                  <Alert>
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>
                      Selecione apenas os eventos que você realmente precisa para otimizar o desempenho.
                    </AlertDescription>
                  </Alert>
                </div>
              </>
            )}

            {/* Actions */}
            <div className="flex justify-end space-x-2">
              <Button
                variant="outline"
                onClick={() => setShowConfigDialog(false)}
              >
                Cancelar
              </Button>
              <Button onClick={handleSaveWebhook}>
                Salvar Configuração
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}