"use client";

import { NotificationDemo } from "@/components/notification-demo";

export default function NotificationsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Sistema de Notificações</h1>
        <p className="text-muted-foreground">
          Teste e configure o sistema de notificações do WazMeow
        </p>
      </div>
      
      <NotificationDemo />
    </div>
  );
}