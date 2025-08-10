"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { 
  Settings, 
  Globe, 
  Bell, 
  Shield, 
  Database, 
  Palette,
  Save,
  RefreshCw,
  AlertCircle,
  CheckCircle,
  Info
} from "lucide-react";

interface SettingsConfig {
  // API Settings
  apiUrl: string;
  apiTimeout: number;
  maxRetries: number;
  
  // Webhook Settings
  webhookTimeout: number;
  webhookMaxRetries: number;
  webhookRetryDelay: number;
  
  // Session Settings
  maxSessions: number;
  sessionTimeout: number;
  autoReconnect: boolean;
  
  // UI Settings
  theme: "light" | "dark" | "system";
  language: string;
  refreshInterval: number;
  
  // Notification Settings
  enableNotifications: boolean;
  notifyOnError: boolean;
  notifyOnSuccess: boolean;
  notifyOnDisconnect: boolean;
  
  // Security Settings
  enableApiKey: boolean;
  apiKey: string;
  enableRateLimit: boolean;
  rateLimitRequests: number;
  rateLimitWindow: number;
  
  // Logging Settings
  logLevel: "debug" | "info" | "warn" | "error";
  logRetention: number;
  enableFileLogging: boolean;
}

export default function SettingsPage() {
  const [settings, setSettings] = useState<SettingsConfig>({
    // API Settings
    apiUrl: "http://localhost:8080",
    apiTimeout: 30000,
    maxRetries: 3,
    
    // Webhook Settings
    webhookTimeout: 30000,
    webhookMaxRetries: 3,
    webhookRetryDelay: 5000,
    
    // Session Settings
    maxSessions: 100,
    sessionTimeout: 3600,
    autoReconnect: true,
    
    // UI Settings
    theme: "system",
    language: "pt-BR",
    refreshInterval: 30000,
    
    // Notification Settings
    enableNotifications: true,
    notifyOnError: true,
    notifyOnSuccess: false,
    notifyOnDisconnect: true,
    
    // Security Settings
    enableApiKey: false,
    apiKey: "",
    enableRateLimit: true,
    rateLimitRequests: 100,
    rateLimitWindow: 60,
    
    // Logging Settings
    logLevel: "info",
    logRetention: 30,
    enableFileLogging: true,
  });

  const [hasChanges, setHasChanges] = useState(false);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  const updateSetting = <K extends keyof SettingsConfig>(
    key: K,
    value: SettingsConfig[K]
  ) => {
    setSettings(prev => ({ ...prev, [key]: value }));
    setHasChanges(true);
    setSaved(false);
  };

  const handleSave = async () => {
    setSaving(true);
    
    // Simular salvamento
    await new Promise(resolve => setTimeout(resolve, 1500));
    
    setHasChanges(false);
    setSaved(true);
    setSaving(false);
    
    // Remover indicador de salvo após 3 segundos
    setTimeout(() => setSaved(false), 3000);
  };

  const handleReset = () => {
    // Reset para valores padrão
    setSettings({
      apiUrl: "http://localhost:8080",
      apiTimeout: 30000,
      maxRetries: 3,
      webhookTimeout: 30000,
      webhookMaxRetries: 3,
      webhookRetryDelay: 5000,
      maxSessions: 100,
      sessionTimeout: 3600,
      autoReconnect: true,
      theme: "system",
      language: "pt-BR",
      refreshInterval: 30000,
      enableNotifications: true,
      notifyOnError: true,
      notifyOnSuccess: false,
      notifyOnDisconnect: true,
      enableApiKey: false,
      apiKey: "",
      enableRateLimit: true,
      rateLimitRequests: 100,
      rateLimitWindow: 60,
      logLevel: "info",
      logRetention: 30,
      enableFileLogging: true,
    });
    setHasChanges(true);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Configurações</h1>
          <p className="text-muted-foreground">
            Configure as preferências e comportamentos do sistema
          </p>
        </div>
        
        <div className="flex items-center gap-2">
          {saved && (
            <Badge variant="default" className="flex items-center gap-1">
              <CheckCircle className="h-3 w-3" />
              Salvo
            </Badge>
          )}
          
          <Button variant="outline" onClick={handleReset}>
            <RefreshCw className="mr-2 h-4 w-4" />
            Restaurar Padrões
          </Button>
          
          <Button 
            onClick={handleSave} 
            disabled={!hasChanges || saving}
          >
            <Save className="mr-2 h-4 w-4" />
            {saving ? "Salvando..." : "Salvar Alterações"}
          </Button>
        </div>
      </div>

      {hasChanges && (
        <Alert>
          <Info className="h-4 w-4" />
          <AlertDescription>
            Você tem alterações não salvas. Clique em "Salvar Alterações" para aplicá-las.
          </AlertDescription>
        </Alert>
      )}

      <Tabs defaultValue="api" className="space-y-4">
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="api">API</TabsTrigger>
          <TabsTrigger value="sessions">Sessões</TabsTrigger>
          <TabsTrigger value="webhooks">Webhooks</TabsTrigger>
          <TabsTrigger value="ui">Interface</TabsTrigger>
          <TabsTrigger value="notifications">Notificações</TabsTrigger>
          <TabsTrigger value="security">Segurança</TabsTrigger>
        </TabsList>

        {/* API Settings */}
        <TabsContent value="api">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Globe className="h-5 w-5" />
                Configurações da API
              </CardTitle>
              <CardDescription>
                Configure a conexão e comportamento da API WazMeow
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="apiUrl">URL da API</Label>
                  <Input
                    id="apiUrl"
                    value={settings.apiUrl}
                    onChange={(e) => updateSetting("apiUrl", e.target.value)}
                    placeholder="http://localhost:8080"
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="apiTimeout">Timeout (ms)</Label>
                  <Input
                    id="apiTimeout"
                    type="number"
                    value={settings.apiTimeout}
                    onChange={(e) => updateSetting("apiTimeout", parseInt(e.target.value))}
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="maxRetries">Máximo de Tentativas</Label>
                  <Input
                    id="maxRetries"
                    type="number"
                    value={settings.maxRetries}
                    onChange={(e) => updateSetting("maxRetries", parseInt(e.target.value))}
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="refreshInterval">Intervalo de Atualização (ms)</Label>
                  <Input
                    id="refreshInterval"
                    type="number"
                    value={settings.refreshInterval}
                    onChange={(e) => updateSetting("refreshInterval", parseInt(e.target.value))}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Session Settings */}
        <TabsContent value="sessions">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5" />
                Configurações de Sessões
              </CardTitle>
              <CardDescription>
                Configure o comportamento das sessões WhatsApp
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="maxSessions">Máximo de Sessões</Label>
                  <Input
                    id="maxSessions"
                    type="number"
                    value={settings.maxSessions}
                    onChange={(e) => updateSetting("maxSessions", parseInt(e.target.value))}
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="sessionTimeout">Timeout da Sessão (segundos)</Label>
                  <Input
                    id="sessionTimeout"
                    type="number"
                    value={settings.sessionTimeout}
                    onChange={(e) => updateSetting("sessionTimeout", parseInt(e.target.value))}
                  />
                </div>
              </div>
              
              <div className="flex items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
                  <Label className="text-base">Reconexão Automática</Label>
                  <p className="text-sm text-muted-foreground">
                    Tentar reconectar automaticamente sessões desconectadas
                  </p>
                </div>
                <Switch
                  checked={settings.autoReconnect}
                  onCheckedChange={(checked) => updateSetting("autoReconnect", checked)}
                />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Webhook Settings */}
        <TabsContent value="webhooks">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Globe className="h-5 w-5" />
                Configurações de Webhooks
              </CardTitle>
              <CardDescription>
                Configure o comportamento dos webhooks
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="webhookTimeout">Timeout (ms)</Label>
                  <Input
                    id="webhookTimeout"
                    type="number"
                    value={settings.webhookTimeout}
                    onChange={(e) => updateSetting("webhookTimeout", parseInt(e.target.value))}
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="webhookMaxRetries">Máximo de Tentativas</Label>
                  <Input
                    id="webhookMaxRetries"
                    type="number"
                    value={settings.webhookMaxRetries}
                    onChange={(e) => updateSetting("webhookMaxRetries", parseInt(e.target.value))}
                  />
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="webhookRetryDelay">Delay entre Tentativas (ms)</Label>
                  <Input
                    id="webhookRetryDelay"
                    type="number"
                    value={settings.webhookRetryDelay}
                    onChange={(e) => updateSetting("webhookRetryDelay", parseInt(e.target.value))}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* UI Settings */}
        <TabsContent value="ui">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Palette className="h-5 w-5" />
                Configurações da Interface
              </CardTitle>
              <CardDescription>
                Personalize a aparência e comportamento da interface
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="theme">Tema</Label>
                  <Select 
                    value={settings.theme} 
                    onValueChange={(value: "light" | "dark" | "system") => updateSetting("theme", value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="light">Claro</SelectItem>
                      <SelectItem value="dark">Escuro</SelectItem>
                      <SelectItem value="system">Sistema</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="language">Idioma</Label>
                  <Select 
                    value={settings.language} 
                    onValueChange={(value) => updateSetting("language", value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="pt-BR">Português (Brasil)</SelectItem>
                      <SelectItem value="en-US">English (US)</SelectItem>
                      <SelectItem value="es-ES">Español</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Notification Settings */}
        <TabsContent value="notifications">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Bell className="h-5 w-5" />
                Configurações de Notificações
              </CardTitle>
              <CardDescription>
                Configure quando e como receber notificações
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-4">
                <div className="flex items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <Label className="text-base">Habilitar Notificações</Label>
                    <p className="text-sm text-muted-foreground">
                      Receber notificações do sistema
                    </p>
                  </div>
                  <Switch
                    checked={settings.enableNotifications}
                    onCheckedChange={(checked) => updateSetting("enableNotifications", checked)}
                  />
                </div>
                
                {settings.enableNotifications && (
                  <>
                    <div className="flex items-center justify-between rounded-lg border p-4">
                      <div className="space-y-0.5">
                        <Label className="text-base">Notificar Erros</Label>
                        <p className="text-sm text-muted-foreground">
                          Receber notificações quando ocorrerem erros
                        </p>
                      </div>
                      <Switch
                        checked={settings.notifyOnError}
                        onCheckedChange={(checked) => updateSetting("notifyOnError", checked)}
                      />
                    </div>
                    
                    <div className="flex items-center justify-between rounded-lg border p-4">
                      <div className="space-y-0.5">
                        <Label className="text-base">Notificar Sucessos</Label>
                        <p className="text-sm text-muted-foreground">
                          Receber notificações de operações bem-sucedidas
                        </p>
                      </div>
                      <Switch
                        checked={settings.notifyOnSuccess}
                        onCheckedChange={(checked) => updateSetting("notifyOnSuccess", checked)}
                      />
                    </div>
                    
                    <div className="flex items-center justify-between rounded-lg border p-4">
                      <div className="space-y-0.5">
                        <Label className="text-base">Notificar Desconexões</Label>
                        <p className="text-sm text-muted-foreground">
                          Receber notificações quando sessões se desconectarem
                        </p>
                      </div>
                      <Switch
                        checked={settings.notifyOnDisconnect}
                        onCheckedChange={(checked) => updateSetting("notifyOnDisconnect", checked)}
                      />
                    </div>
                  </>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Security Settings */}
        <TabsContent value="security">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="h-5 w-5" />
                Configurações de Segurança
              </CardTitle>
              <CardDescription>
                Configure as opções de segurança e autenticação
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
                  <Label className="text-base">Habilitar Chave API</Label>
                  <p className="text-sm text-muted-foreground">
                    Usar chave API para autenticação
                  </p>
                </div>
                <Switch
                  checked={settings.enableApiKey}
                  onCheckedChange={(checked) => updateSetting("enableApiKey", checked)}
                />
              </div>
              
              {settings.enableApiKey && (
                <div className="space-y-2">
                  <Label htmlFor="apiKey">Chave API</Label>
                  <Input
                    id="apiKey"
                    type="password"
                    value={settings.apiKey}
                    onChange={(e) => updateSetting("apiKey", e.target.value)}
                    placeholder="Digite sua chave API"
                  />
                </div>
              )}
              
              <div className="flex items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
                  <Label className="text-base">Habilitar Rate Limiting</Label>
                  <p className="text-sm text-muted-foreground">
                    Limitar número de requisições por período
                  </p>
                </div>
                <Switch
                  checked={settings.enableRateLimit}
                  onCheckedChange={(checked) => updateSetting("enableRateLimit", checked)}
                />
              </div>
              
              {settings.enableRateLimit && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="rateLimitRequests">Requisições por Janela</Label>
                    <Input
                      id="rateLimitRequests"
                      type="number"
                      value={settings.rateLimitRequests}
                      onChange={(e) => updateSetting("rateLimitRequests", parseInt(e.target.value))}
                    />
                  </div>
                  
                  <div className="space-y-2">
                    <Label htmlFor="rateLimitWindow">Janela de Tempo (segundos)</Label>
                    <Input
                      id="rateLimitWindow"
                      type="number"
                      value={settings.rateLimitWindow}
                      onChange={(e) => updateSetting("rateLimitWindow", parseInt(e.target.value))}
                    />
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}