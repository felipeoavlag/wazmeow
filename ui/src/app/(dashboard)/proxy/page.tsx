"use client";

import { useState } from "react";
import { ProxyConfigForm } from "@/components/forms/proxy-config-form";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useNotifications } from "@/lib/hooks/use-notifications";
import { 
  Globe, 
  Shield, 
  Activity,
  CheckCircle,
  AlertCircle,
  Info,
  Trash2,
  Plus
} from "lucide-react";

interface ProxyConfig {
  id: string;
  name: string;
  enabled: boolean;
  type: "http" | "https" | "socks4" | "socks5";
  host: string;
  port: number;
  username?: string;
  password?: string;
  timeout: number;
  retries: number;
  testUrl?: string;
  status: "active" | "inactive" | "error";
  lastTested?: Date;
  latency?: number;
}

export default function ProxyPage() {
  const notifications = useNotifications();
  const [configs, setConfigs] = useState<ProxyConfig[]>([
    {
      id: "1",
      name: "Proxy Principal",
      enabled: true,
      type: "http",
      host: "proxy.exemplo.com",
      port: 8080,
      username: "usuario",
      timeout: 10000,
      retries: 3,
      testUrl: "https://httpbin.org/ip",
      status: "active",
      lastTested: new Date(Date.now() - 300000), // 5 minutos atrás
      latency: 150,
    },
    {
      id: "2",
      name: "Proxy Backup",
      enabled: false,
      type: "socks5",
      host: "backup.exemplo.com",
      port: 1080,
      timeout: 15000,
      retries: 2,
      status: "inactive",
    },
  ]);

  const [editingConfig, setEditingConfig] = useState<ProxyConfig | null>(null);
  const [showForm, setShowForm] = useState(false);

  const handleSaveConfig = async (data: any) => {
    try {
      if (editingConfig) {
        // Atualizar configuração existente
        setConfigs(prev => prev.map(config => 
          config.id === editingConfig.id 
            ? { ...config, ...data, id: editingConfig.id }
            : config
        ));
        notifications.success("Configuração atualizada", {
          description: "As configurações do proxy foram atualizadas com sucesso",
        });
      } else {
        // Criar nova configuração
        const newConfig: ProxyConfig = {
          ...data,
          id: Date.now().toString(),
          name: `Proxy ${configs.length + 1}`,
          status: data.enabled ? "active" : "inactive",
        };
        setConfigs(prev => [...prev, newConfig]);
        notifications.success("Configuração criada", {
          description: "Nova configuração de proxy foi criada com sucesso",
        });
      }
      
      setEditingConfig(null);
      setShowForm(false);
    } catch (error) {
      notifications.error("Erro ao salvar", {
        description: "Não foi possível salvar a configuração do proxy",
      });
    }
  };

  const handleTestProxy = async (data: any): Promise<boolean> => {
    // Simular teste de proxy
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    // Simular resultado aleatório
    const success = Math.random() > 0.3;
    
    if (success) {
      notifications.success("Teste bem-sucedido", {
        description: "O proxy está funcionando corretamente",
      });
    } else {
      notifications.error("Teste falhou", {
        description: "Não foi possível conectar através do proxy",
      });
    }
    
    return success;
  };

  const handleDeleteConfig = (id: string) => {
    setConfigs(prev => prev.filter(config => config.id !== id));
    notifications.success("Configuração removida", {
      description: "A configuração do proxy foi removida com sucesso",
    });
  };

  const handleToggleConfig = (id: string) => {
    setConfigs(prev => prev.map(config => 
      config.id === id 
        ? { 
            ...config, 
            enabled: !config.enabled,
            status: !config.enabled ? "active" : "inactive"
          }
        : config
    ));
  };

  const getStatusBadge = (status: ProxyConfig["status"]) => {
    switch (status) {
      case "active":
        return <Badge variant="default" className="flex items-center gap-1">
          <CheckCircle className="h-3 w-3" />
          Ativo
        </Badge>;
      case "inactive":
        return <Badge variant="secondary" className="flex items-center gap-1">
          <Activity className="h-3 w-3" />
          Inativo
        </Badge>;
      case "error":
        return <Badge variant="destructive" className="flex items-center gap-1">
          <AlertCircle className="h-3 w-3" />
          Erro
        </Badge>;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Configuração de Proxy</h1>
          <p className="text-muted-foreground">
            Configure proxies para rotear o tráfego das sessões WhatsApp
          </p>
        </div>
        
        <Button onClick={() => setShowForm(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Nova Configuração
        </Button>
      </div>

      {/* Info Alert */}
      <Alert>
        <Info className="h-4 w-4" />
        <AlertDescription>
          Os proxies configurados aqui podem ser utilizados nas sessões WhatsApp para rotear o tráfego. 
          Certifique-se de testar a conectividade antes de usar em produção.
        </AlertDescription>
      </Alert>

      {/* Existing Configurations */}
      {configs.length > 0 && (
        <div className="space-y-4">
          <h2 className="text-xl font-semibold">Configurações Existentes</h2>
          
          <div className="grid gap-4">
            {configs.map((config) => (
              <Card key={config.id}>
                <CardHeader className="pb-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Globe className="h-5 w-5" />
                      <div>
                        <CardTitle className="text-lg">{config.name}</CardTitle>
                        <CardDescription>
                          {config.type.toUpperCase()} - {config.host}:{config.port}
                        </CardDescription>
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2">
                      {getStatusBadge(config.status)}
                      
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => {
                          setEditingConfig(config);
                          setShowForm(true);
                        }}
                      >
                        Editar
                      </Button>
                      
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleDeleteConfig(config.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </CardHeader>
                
                <CardContent className="pt-0">
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <span className="text-muted-foreground">Tipo:</span>
                      <p className="font-medium">{config.type.toUpperCase()}</p>
                    </div>
                    
                    <div>
                      <span className="text-muted-foreground">Timeout:</span>
                      <p className="font-medium">{config.timeout}ms</p>
                    </div>
                    
                    <div>
                      <span className="text-muted-foreground">Tentativas:</span>
                      <p className="font-medium">{config.retries}</p>
                    </div>
                    
                    {config.latency && (
                      <div>
                        <span className="text-muted-foreground">Latência:</span>
                        <p className="font-medium">{config.latency}ms</p>
                      </div>
                    )}
                  </div>
                  
                  {config.lastTested && (
                    <div className="mt-3 pt-3 border-t">
                      <span className="text-xs text-muted-foreground">
                        Último teste: {config.lastTested.toLocaleString()}
                      </span>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Configuration Form */}
      {showForm && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold">
              {editingConfig ? "Editar Configuração" : "Nova Configuração"}
            </h2>
            
            <Button 
              variant="outline" 
              onClick={() => {
                setShowForm(false);
                setEditingConfig(null);
              }}
            >
              Cancelar
            </Button>
          </div>
          
          <ProxyConfigForm
            initialData={editingConfig || undefined}
            onSubmit={handleSaveConfig}
            onTest={handleTestProxy}
          />
        </div>
      )}

      {/* Empty State */}
      {configs.length === 0 && !showForm && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Globe className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">Nenhum proxy configurado</h3>
            <p className="text-muted-foreground text-center mb-4">
              Configure um proxy para rotear o tráfego das suas sessões WhatsApp
            </p>
            <Button onClick={() => setShowForm(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Criar Primeira Configuração
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  );
}