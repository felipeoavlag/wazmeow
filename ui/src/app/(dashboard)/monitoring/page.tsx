"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell
} from "recharts";
import { 
  Activity, 
  MessageSquare, 
  Webhook, 
  Users, 
  TrendingUp,
  TrendingDown,
  Wifi,
  WifiOff,
  AlertTriangle,
  CheckCircle,
  RefreshCw
} from "lucide-react";

// Mock data para gráficos
const messageData = [
  { time: "00:00", sent: 45, received: 32 },
  { time: "04:00", sent: 23, received: 18 },
  { time: "08:00", sent: 89, received: 67 },
  { time: "12:00", sent: 156, received: 134 },
  { time: "16:00", sent: 203, received: 189 },
  { time: "20:00", sent: 134, received: 98 },
];

const sessionStatusData = [
  { name: "Conectadas", value: 8, color: "#22c55e" },
  { name: "Desconectadas", value: 3, color: "#ef4444" },
  { name: "Conectando", value: 1, color: "#f59e0b" },
];

const webhookData = [
  { session: "Principal", success: 98.5, failed: 1.5 },
  { session: "Suporte", success: 99.2, failed: 0.8 },
  { session: "Marketing", success: 96.8, failed: 3.2 },
  { session: "Vendas", success: 97.9, failed: 2.1 },
];

interface MetricCardProps {
  title: string;
  value: string | number;
  change: number;
  icon: React.ReactNode;
  description: string;
}

function MetricCard({ title, value, change, icon, description }: MetricCardProps) {
  const isPositive = change > 0;
  
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        <div className="flex items-center text-xs text-muted-foreground">
          {isPositive ? (
            <TrendingUp className="mr-1 h-3 w-3 text-green-500" />
          ) : (
            <TrendingDown className="mr-1 h-3 w-3 text-red-500" />
          )}
          <span className={isPositive ? "text-green-500" : "text-red-500"}>
            {isPositive ? "+" : ""}{change}%
          </span>
          <span className="ml-1">{description}</span>
        </div>
      </CardContent>
    </Card>
  );
}

export default function MonitoringPage() {
  const [timeRange, setTimeRange] = useState("24h");
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [lastUpdate, setLastUpdate] = useState(new Date());

  // Simular atualização automática
  useEffect(() => {
    if (!autoRefresh) return;

    const interval = setInterval(() => {
      setLastUpdate(new Date());
    }, 30000); // Atualizar a cada 30 segundos

    return () => clearInterval(interval);
  }, [autoRefresh]);

  const handleRefresh = () => {
    setLastUpdate(new Date());
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Monitoramento</h1>
          <p className="text-muted-foreground">
            Acompanhe métricas e performance em tempo real
          </p>
        </div>
        
        <div className="flex items-center gap-2">
          <Select value={timeRange} onValueChange={setTimeRange}>
            <SelectTrigger className="w-32">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1h">1 hora</SelectItem>
              <SelectItem value="6h">6 horas</SelectItem>
              <SelectItem value="24h">24 horas</SelectItem>
              <SelectItem value="7d">7 dias</SelectItem>
            </SelectContent>
          </Select>
          
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            className="flex items-center gap-2"
          >
            <RefreshCw className="h-4 w-4" />
            Atualizar
          </Button>
        </div>
      </div>

      {/* Status Bar */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                <span className="text-sm">Sistema Online</span>
              </div>
              <div className="text-sm text-muted-foreground">
                Última atualização: {lastUpdate.toLocaleTimeString('pt-BR')}
              </div>
            </div>
            
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="flex items-center gap-1">
                <CheckCircle className="h-3 w-3 text-green-500" />
                Todos os serviços operacionais
              </Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Metrics Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          title="Mensagens/Hora"
          value="1,247"
          change={12.5}
          icon={<MessageSquare className="h-4 w-4 text-muted-foreground" />}
          description="vs hora anterior"
        />
        
        <MetricCard
          title="Sessões Ativas"
          value="8/12"
          change={-8.3}
          icon={<Wifi className="h-4 w-4 text-muted-foreground" />}
          description="vs ontem"
        />
        
        <MetricCard
          title="Webhooks Enviados"
          value="3,421"
          change={23.1}
          icon={<Webhook className="h-4 w-4 text-muted-foreground" />}
          description="últimas 24h"
        />
        
        <MetricCard
          title="Taxa de Sucesso"
          value="98.2%"
          change={1.2}
          icon={<Activity className="h-4 w-4 text-muted-foreground" />}
          description="webhooks"
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {/* Messages Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Fluxo de Mensagens</CardTitle>
            <CardDescription>
              Mensagens enviadas e recebidas nas últimas 24 horas
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={messageData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time" />
                <YAxis />
                <Tooltip />
                <Line 
                  type="monotone" 
                  dataKey="sent" 
                  stroke="#3b82f6" 
                  strokeWidth={2}
                  name="Enviadas"
                />
                <Line 
                  type="monotone" 
                  dataKey="received" 
                  stroke="#10b981" 
                  strokeWidth={2}
                  name="Recebidas"
                />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Session Status */}
        <Card>
          <CardHeader>
            <CardTitle>Status das Sessões</CardTitle>
            <CardDescription>
              Distribuição atual do status das sessões
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={sessionStatusData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {sessionStatusData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
            
            <div className="flex justify-center gap-4 mt-4">
              {sessionStatusData.map((entry, index) => (
                <div key={index} className="flex items-center gap-2">
                  <div 
                    className="w-3 h-3 rounded-full" 
                    style={{ backgroundColor: entry.color }}
                  />
                  <span className="text-sm">{entry.name}: {entry.value}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Webhook Performance */}
      <Card>
        <CardHeader>
          <CardTitle>Performance dos Webhooks</CardTitle>
          <CardDescription>
            Taxa de sucesso por sessão nas últimas 24 horas
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={webhookData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="session" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="success" fill="#22c55e" name="Sucesso %" />
              <Bar dataKey="failed" fill="#ef4444" name="Falha %" />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Real-time Activity */}
      <Card>
        <CardHeader>
          <CardTitle>Atividade em Tempo Real</CardTitle>
          <CardDescription>
            Últimos eventos do sistema
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {[
              {
                time: "16:32:15",
                type: "message",
                description: "Mensagem enviada via Sessão Principal",
                status: "success"
              },
              {
                time: "16:31:48",
                type: "webhook",
                description: "Webhook entregue para Marketing",
                status: "success"
              },
              {
                time: "16:31:22",
                type: "connection",
                description: "Sessão Suporte reconectada",
                status: "success"
              },
              {
                time: "16:30:55",
                type: "error",
                description: "Falha na entrega de webhook para Vendas",
                status: "error"
              },
              {
                time: "16:30:33",
                type: "message",
                description: "Mensagem recebida na Sessão Principal",
                status: "success"
              },
            ].map((activity, index) => (
              <div key={index} className="flex items-center gap-3 p-3 rounded-lg border">
                <div className="flex-shrink-0">
                  {activity.status === "success" ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <AlertTriangle className="h-4 w-4 text-red-500" />
                  )}
                </div>
                <div className="flex-1">
                  <p className="text-sm">{activity.description}</p>
                </div>
                <div className="flex-shrink-0">
                  <span className="text-xs text-muted-foreground">
                    {activity.time}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}