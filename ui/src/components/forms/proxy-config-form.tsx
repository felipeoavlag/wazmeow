"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
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
  Globe, 
  Shield, 
  TestTube,
  CheckCircle,
  AlertCircle,
  Info,
  Eye,
  EyeOff
} from "lucide-react";

const proxyConfigSchema = z.object({
  enabled: z.boolean(),
  type: z.enum(["http", "https", "socks4", "socks5"]),
  host: z.string().min(1, "Host é obrigatório"),
  port: z.number().min(1, "Porta deve ser maior que 0").max(65535, "Porta deve ser menor que 65536"),
  username: z.string().optional(),
  password: z.string().optional(),
  timeout: z.number().min(1000, "Timeout mínimo é 1000ms").max(60000, "Timeout máximo é 60000ms"),
  retries: z.number().min(0, "Tentativas não pode ser negativo").max(10, "Máximo 10 tentativas"),
  testUrl: z.string().url("URL de teste deve ser válida").optional(),
});

type ProxyConfigData = z.infer<typeof proxyConfigSchema>;

interface ProxyConfigFormProps {
  initialData?: Partial<ProxyConfigData>;
  onSubmit: (data: ProxyConfigData) => void;
  onTest?: (data: ProxyConfigData) => Promise<boolean>;
  loading?: boolean;
}

export function ProxyConfigForm({ 
  initialData, 
  onSubmit, 
  onTest,
  loading = false 
}: ProxyConfigFormProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<{
    success: boolean;
    message: string;
    latency?: number;
  } | null>(null);

  const form = useForm<ProxyConfigData>({
    resolver: zodResolver(proxyConfigSchema),
    defaultValues: {
      enabled: false,
      type: "http",
      host: "",
      port: 8080,
      username: "",
      password: "",
      timeout: 10000,
      retries: 3,
      testUrl: "https://httpbin.org/ip",
      ...initialData,
    },
  });

  const { register, handleSubmit, watch, setValue, formState: { errors } } = form;
  const watchEnabled = watch("enabled");
  const watchType = watch("type");

  const handleTestConnection = async () => {
    if (!onTest) return;
    
    setTesting(true);
    setTestResult(null);
    
    try {
      const formData = form.getValues();
      const startTime = Date.now();
      const success = await onTest(formData);
      const latency = Date.now() - startTime;
      
      setTestResult({
        success,
        message: success 
          ? `Conexão bem-sucedida! Latência: ${latency}ms`
          : "Falha na conexão com o proxy",
        latency: success ? latency : undefined,
      });
    } catch (error) {
      setTestResult({
        success: false,
        message: `Erro ao testar conexão: ${error instanceof Error ? error.message : 'Erro desconhecido'}`,
      });
    } finally {
      setTesting(false);
    }
  };

  const getProxyTypeDescription = (type: string) => {
    switch (type) {
      case "http":
        return "Proxy HTTP padrão para navegação web";
      case "https":
        return "Proxy HTTPS com criptografia SSL/TLS";
      case "socks4":
        return "SOCKS4 - Protocolo mais simples, sem autenticação";
      case "socks5":
        return "SOCKS5 - Protocolo avançado com autenticação";
      default:
        return "";
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Globe className="h-5 w-5" />
          Configuração de Proxy
        </CardTitle>
        <CardDescription>
          Configure um proxy para rotear o tráfego da sessão WhatsApp
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Enable Proxy */}
          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="space-y-0.5">
              <Label className="text-base">Habilitar Proxy</Label>
              <p className="text-sm text-muted-foreground">
                Usar proxy para conexões desta sessão
              </p>
            </div>
            <Switch
              checked={watchEnabled}
              onCheckedChange={(checked) => setValue("enabled", checked)}
            />
          </div>

          {watchEnabled && (
            <div className="space-y-4">
              {/* Proxy Type */}
              <div className="space-y-2">
                <Label htmlFor="type">Tipo de Proxy</Label>
                <Select 
                  value={watchType} 
                  onValueChange={(value: "http" | "https" | "socks4" | "socks5") => setValue("type", value)}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="http">HTTP</SelectItem>
                    <SelectItem value="https">HTTPS</SelectItem>
                    <SelectItem value="socks4">SOCKS4</SelectItem>
                    <SelectItem value="socks5">SOCKS5</SelectItem>
                  </SelectContent>
                </Select>
                <p className="text-xs text-muted-foreground">
                  {getProxyTypeDescription(watchType)}
                </p>
              </div>

              {/* Host and Port */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="md:col-span-2 space-y-2">
                  <Label htmlFor="host">Host/IP do Proxy</Label>
                  <Input
                    id="host"
                    {...register("host")}
                    placeholder="proxy.exemplo.com ou 192.168.1.100"
                    className={errors.host ? "border-red-500" : ""}
                  />
                  {errors.host && (
                    <p className="text-sm text-red-500">{errors.host.message}</p>
                  )}
                </div>
                
                <div className="space-y-2">
                  <Label htmlFor="port">Porta</Label>
                  <Input
                    id="port"
                    type="number"
                    {...register("port", { valueAsNumber: true })}
                    placeholder="8080"
                    className={errors.port ? "border-red-500" : ""}
                  />
                  {errors.port && (
                    <p className="text-sm text-red-500">{errors.port.message}</p>
                  )}
                </div>
              </div>

              {/* Authentication */}
              {(watchType === "http" || watchType === "https" || watchType === "socks5") && (
                <div className="space-y-4">
                  <div className="flex items-center gap-2">
                    <Shield className="h-4 w-4" />
                    <Label className="text-base">Autenticação (Opcional)</Label>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="username">Usuário</Label>
                      <Input
                        id="username"
                        {...register("username")}
                        placeholder="usuario"
                        className={errors.username ? "border-red-500" : ""}
                      />
                      {errors.username && (
                        <p className="text-sm text-red-500">{errors.username.message}</p>
                      )}
                    </div>
                    
                    <div className="space-y-2">
                      <Label htmlFor="password">Senha</Label>
                      <div className="relative">
                        <Input
                          id="password"
                          type={showPassword ? "text" : "password"}
                          {...register("password")}
                          placeholder="senha"
                          className={errors.password ? "border-red-500" : ""}
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                          onClick={() => setShowPassword(!showPassword)}
                        >
                          {showPassword ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </Button>
                      </div>
                      {errors.password && (
                        <p className="text-sm text-red-500">{errors.password.message}</p>
                      )}
                    </div>
                  </div>
                </div>
              )}

              {/* Advanced Settings */}
              <div className="space-y-4">
                <Label className="text-base">Configurações Avançadas</Label>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="timeout">Timeout (ms)</Label>
                    <Input
                      id="timeout"
                      type="number"
                      {...register("timeout", { valueAsNumber: true })}
                      className={errors.timeout ? "border-red-500" : ""}
                    />
                    {errors.timeout && (
                      <p className="text-sm text-red-500">{errors.timeout.message}</p>
                    )}
                  </div>
                  
                  <div className="space-y-2">
                    <Label htmlFor="retries">Tentativas</Label>
                    <Input
                      id="retries"
                      type="number"
                      {...register("retries", { valueAsNumber: true })}
                      className={errors.retries ? "border-red-500" : ""}
                    />
                    {errors.retries && (
                      <p className="text-sm text-red-500">{errors.retries.message}</p>
                    )}
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="testUrl">URL de Teste</Label>
                  <Input
                    id="testUrl"
                    {...register("testUrl")}
                    placeholder="https://httpbin.org/ip"
                    className={errors.testUrl ? "border-red-500" : ""}
                  />
                  {errors.testUrl && (
                    <p className="text-sm text-red-500">{errors.testUrl.message}</p>
                  )}
                  <p className="text-xs text-muted-foreground">
                    URL usada para testar a conectividade do proxy
                  </p>
                </div>
              </div>

              {/* Test Connection */}
              {onTest && (
                <div className="space-y-3">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={handleTestConnection}
                    disabled={testing}
                    className="w-full"
                  >
                    <TestTube className="mr-2 h-4 w-4" />
                    {testing ? "Testando Conexão..." : "Testar Conexão"}
                  </Button>

                  {testResult && (
                    <Alert variant={testResult.success ? "default" : "destructive"}>
                      {testResult.success ? (
                        <CheckCircle className="h-4 w-4" />
                      ) : (
                        <AlertCircle className="h-4 w-4" />
                      )}
                      <AlertDescription>
                        {testResult.message}
                      </AlertDescription>
                    </Alert>
                  )}
                </div>
              )}
            </div>
          )}

          {/* Submit Button */}
          <div className="flex justify-end gap-2">
            <Button type="submit" disabled={loading}>
              {loading ? "Salvando..." : "Salvar Configuração"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}