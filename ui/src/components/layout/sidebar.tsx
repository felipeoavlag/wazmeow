"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { useUIStore } from "@/lib/stores/app-store";
import { Button } from "@/components/ui/button";
import {
  LayoutDashboard,
  Smartphone,
  Webhook,
  MessageSquare,
  BarChart3,
  Settings,
  FileText,
  Bell,
  Globe,
} from "lucide-react";

const navigation = [
  {
    name: "Dashboard",
    href: "/",
    icon: LayoutDashboard,
  },
  {
    name: "Sessões",
    href: "/sessions",
    icon: Smartphone,
  },
  {
    name: "Webhooks",
    href: "/webhooks",
    icon: Webhook,
  },
  {
    name: "Mensagens",
    href: "/messages",
    icon: MessageSquare,
  },
  {
    name: "Monitoramento",
    href: "/monitoring",
    icon: BarChart3,
  },
  {
    name: "Logs",
    href: "/logs",
    icon: FileText,
  },
  {
    name: "Proxy",
    href: "/proxy",
    icon: Globe,
  },
  {
    name: "Notificações",
    href: "/notifications",
    icon: Bell,
  },
  {
    name: "Configurações",
    href: "/settings",
    icon: Settings,
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const { sidebarOpen, isMobile, setSidebarOpen } = useUIStore();

  const handleLinkClick = () => {
    if (isMobile) {
      setSidebarOpen(false);
    }
  };

  return (
    <div
      className={cn(
        "fixed left-0 top-0 z-40 h-screen bg-background border-r transition-all duration-300",
        isMobile
          ? sidebarOpen
            ? "w-64 translate-x-0"
            : "w-64 -translate-x-full"
          : sidebarOpen
            ? "w-64"
            : "w-16"
      )}
    >
      {/* Logo */}
      <div className="flex h-16 items-center border-b px-4">
        <div className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
            <Smartphone className="h-4 w-4 text-primary-foreground" />
          </div>
          {sidebarOpen && (
            <span className="font-semibold text-lg">WazMeow</span>
          )}
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 p-2">
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          
          return (
            <Link key={item.name} href={item.href} onClick={handleLinkClick}>
              <Button
                variant={isActive ? "secondary" : "ghost"}
                className={cn(
                  "w-full justify-start gap-3 h-11",
                  !sidebarOpen && !isMobile && "px-2"
                )}
              >
                <item.icon className="h-5 w-5 flex-shrink-0" />
                {(sidebarOpen || isMobile) && (
                  <span className="truncate">{item.name}</span>
                )}
              </Button>
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="border-t p-2">
        <div className={cn(
          "text-xs text-muted-foreground",
          (sidebarOpen || isMobile) ? "px-3 py-2" : "text-center py-2"
        )}>
          {(sidebarOpen || isMobile) ? "v1.0.0" : "v1"}
        </div>
      </div>
    </div>
  );
}