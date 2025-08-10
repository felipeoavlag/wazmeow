"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { 
  Smartphone, 
  Wifi, 
  WifiOff, 
  MessageSquare, 
  Webhook,
  Plus,
  Activity,
  Users,
  TrendingUp
} from "lucide-react";

// Mock data - será substituído por dados reais da API
const mockStats = {
  totalSessions: 12,
  activeSessions: 8,
  totalMessages: 1247,
  webhookEvents: 342,
};

const mockSessions = [
  {
    id: "1",
    name: "Sessão Principal",
    status: "connected" as const,
    phone: "+55 11 99999-9999",
    lastActivity: "2 min atrás",
  },
  {
    id: "2",
    name: "Suporte",
    status: "connected" as const,
    phone: "+55 11 88888-8888",
    lastActivity: "5 min atrás",
  },
  {
    id: "3",
    name: "Marketing",
    status: "disconnected" as const,
    phone: "+55 11 77777-7777",
    lastActivity: "1 hora atrás",
  },
];

const mockActivities = [
  {
    id: "1",
    type: "session_connected",
    description: "Sessão 'Principal' conectada com sucesso",
    timestamp: "há 5 minutos",
  },
  {
    id: "2",
    type: "message_sent",
    description: "Mensagem enviada para +55 11 99999-9999",
    timestamp: "há 10 minutos",
  },
  {
    id: "3",
    type: "webhook_configured",
    description: "Webhook configurado para sessão 'Suporte'",
    timestamp: "há 15 minutos",
  },
];

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">
            Visão geral das suas sessões WhatsApp
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Nova Sessão
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total de Sessões
            </CardTitle>
            <Smartphone className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{mockStats.totalSessions}</div>
            <p className="text-xs text-muted-foreground">
              +2 desde o último mês
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Sessões Ativas
            </CardTitle>
            <Wifi className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{mockStats.activeSessions}</div>
            <p className="text-xs text-muted-foreground">
              {Math.round((mockStats.activeSessions / mockStats.totalSessions) * 100)}% do total
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Mensagens Enviadas
            </CardTitle>
            <MessageSquare className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{mockStats.totalMessages}</div>
            <p className="text-xs text-muted-foreground">
              +180 desde ontem
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Eventos Webhook
            </CardTitle>
            <Webhook className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{mockStats.webhookEvents}</div>
            <p className="text-xs text-muted-foreground">
              +12% desde a última semana
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        {/* Recent Sessions */}
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Sessões Recentes</CardTitle>
            <CardDescription>
              Status das suas sessões WhatsApp
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {mockSessions.map((session) => (
                <div
                  key={session.id}
                  className="flex items-center justify-between p-4 border rounded-lg"
                >
                  <div className="flex items-center space-x-4">
                    <div className="flex items-center space-x-2">
                      {session.status === "connected" ? (
                        <Wifi className="h-4 w-4 text-green-500" />
                      ) : (
                        <WifiOff className="h-4 w-4 text-red-500" />
                      )}
                      <div>
                        <p className="text-sm font-medium">{session.name}</p>
                        <p className="text-xs text-muted-foreground">
                          {session.phone}
                        </p>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Badge
                      variant={session.status === "connected" ? "default" : "secondary"}
                    >
                      {session.status === "connected" ? "Conectado" : "Desconectado"}
                    </Badge>
                    <p className="text-xs text-muted-foreground">
                      {session.lastActivity}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Recent Activity */}
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Atividade Recente</CardTitle>
            <CardDescription>
              Últimas ações realizadas
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {mockActivities.map((activity) => (
                <div key={activity.id} className="flex items-start space-x-3">
                  <div className="flex-shrink-0">
                    <Activity className="h-4 w-4 text-muted-foreground mt-0.5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-foreground">
                      {activity.description}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {activity.timestamp}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}