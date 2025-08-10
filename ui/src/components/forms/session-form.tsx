"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Switch } from "@/components/ui/switch";
import { CreateSessionRequest } from "@/lib/types/api";

interface SessionFormProps {
  onSubmit: (data: CreateSessionRequest) => void;
  onCancel: () => void;
  loading?: boolean;
}

export function SessionForm({ 
  onSubmit, 
  onCancel, 
  loading = false
}: SessionFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    webhookUrl: "",
    useProxy: false,
    proxyType: "http" as "http" | "socks5",
    proxyHost: "",
    proxyPort: 8080,
    proxyUsername: "",
    proxyPassword: "",
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = "Nome é obrigatório";
    }

    if (formData.webhookUrl && !isValidUrl(formData.webhookUrl)) {
      newErrors.webhookUrl = "URL inválida";
    }

    if (formData.useProxy) {
      if (!formData.proxyHost.trim()) {
        newErrors.proxyHost = "Host do proxy é obrigatório";
      }
      if (formData.proxyPort < 1 || formData.proxyPort > 65535) {
        newErrors.proxyPort = "Porta deve estar entre 1 e 65535";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const isValidUrl = (url: string) => {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    const data: CreateSessionRequest = {
      name: formData.name,
      webhookUrl: formData.webhookUrl || undefined,
    };

    // Adicionar configuração de proxy se habilitado
    if (formData.useProxy && formData.proxyHost && formData.proxyPort) {
      data.proxy = {
        type: formData.proxyType,
        host: formData.proxyHost,
        port: formData.proxyPort,
        username: formData.proxyUsername || undefined,
        password: formData.proxyPassword || undefined,
      };
    }

    onSubmit(data);
  };

  const updateFormData = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Limpar erro do campo quando o usuário começar a digitar
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: "" }));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <Tabs defaultValue="basic" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="basic">Configurações Básicas</TabsTrigger>
          <TabsTrigger value="advanced">Configurações Avançadas</TabsTrigger>
        </TabsList>
        
        <TabsContent value="basic" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Informações da Sessão</CardTitle>
              <CardDescription>
                Configure as informações básicas da sessão
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Nome da Sessão</Label>
                <Input
                  id="name"
                  placeholder="Ex: Sessão Principal"
                  value={formData.name}
                  onChange={(e) => updateFormData("name", e.target.value)}
                  className={errors.name ? "border-red-500" : ""}
                />
                {errors.name && (
                  <p className="text-sm text-red-500">{errors.name}</p>
                )}
                <p className="text-sm text-muted-foreground">
                  Nome único para identificar esta sessão
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="webhookUrl">URL do Webhook (Opcional)</Label>
                <Input
                  id="webhookUrl"
                  placeholder="https://seu-site.com/webhook"
                  value={formData.webhookUrl}
                  onChange={(e) => updateFormData("webhookUrl", e.target.value)}
                  className={errors.webhookUrl ? "border-red-500" : ""}
                />
                {errors.webhookUrl && (
                  <p className="text-sm text-red-500">{errors.webhookUrl}</p>
                )}
                <p className="text-sm text-muted-foreground">
                  URL para receber eventos da sessão em tempo real
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="advanced" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Configuração de Proxy</CardTitle>
              <CardDescription>
                Configure um proxy para esta sessão (opcional)
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex flex-row items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
                  <Label className="text-base">Usar Proxy</Label>
                  <p className="text-sm text-muted-foreground">
                    Habilitar conexão através de proxy
                  </p>
                </div>
                <Switch
                  checked={formData.useProxy}
                  onCheckedChange={(checked) => updateFormData("useProxy", checked)}
                />
              </div>

              {formData.useProxy && (
                <div className="space-y-4 pl-4 border-l-2 border-muted">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="proxyType">Tipo do Proxy</Label>
                      <Select
                        value={formData.proxyType}
                        onValueChange={(value: "http" | "socks5") => updateFormData("proxyType", value)}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Selecione o tipo" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="http">HTTP</SelectItem>
                          <SelectItem value="socks5">SOCKS5</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="proxyPort">Porta</Label>
                      <Input
                        id="proxyPort"
                        type="number"
                        placeholder="8080"
                        value={formData.proxyPort}
                        onChange={(e) => updateFormData("proxyPort", parseInt(e.target.value) || 8080)}
                        className={errors.proxyPort ? "border-red-500" : ""}
                      />
                      {errors.proxyPort && (
                        <p className="text-sm text-red-500">{errors.proxyPort}</p>
                      )}
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="proxyHost">Host do Proxy</Label>
                    <Input
                      id="proxyHost"
                      placeholder="proxy.exemplo.com"
                      value={formData.proxyHost}
                      onChange={(e) => updateFormData("proxyHost", e.target.value)}
                      className={errors.proxyHost ? "border-red-500" : ""}
                    />
                    {errors.proxyHost && (
                      <p className="text-sm text-red-500">{errors.proxyHost}</p>
                    )}
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="proxyUsername">Usuário (Opcional)</Label>
                      <Input
                        id="proxyUsername"
                        placeholder="usuário"
                        value={formData.proxyUsername}
                        onChange={(e) => updateFormData("proxyUsername", e.target.value)}
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="proxyPassword">Senha (Opcional)</Label>
                      <Input
                        id="proxyPassword"
                        type="password"
                        placeholder="senha"
                        value={formData.proxyPassword}
                        onChange={(e) => updateFormData("proxyPassword", e.target.value)}
                      />
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <div className="flex justify-end space-x-2">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={loading}
        >
          Cancelar
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? "Criando..." : "Criar Sessão"}
        </Button>
      </div>
    </form>
  );
}