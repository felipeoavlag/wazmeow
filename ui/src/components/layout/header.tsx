"use client";

import { Button } from "@/components/ui/button";
import { useUIStore } from "@/lib/stores/app-store";
import { ThemeToggle } from "@/components/theme-toggle";
import { Menu, Settings } from "lucide-react";

export function Header() {
  const { toggleSidebar } = useUIStore();

  return (
    <header className="h-16 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-full items-center justify-between px-3 md:px-6">
        {/* Left side */}
        <div className="flex items-center gap-2 md:gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleSidebar}
            className="h-9 w-9"
          >
            <Menu className="h-4 w-4" />
          </Button>
          
          <div className="hidden sm:block">
            <h1 className="text-lg font-semibold">WazMeow Manager</h1>
            <p className="text-sm text-muted-foreground">
              Gerenciador de Sessões WhatsApp
            </p>
          </div>
          
          <div className="sm:hidden">
            <h1 className="text-base font-semibold">WazMeow</h1>
          </div>
        </div>

        {/* Right side */}
        <div className="flex items-center gap-1 md:gap-2">
          <ThemeToggle />
          
          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-9"
          >
            <Settings className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </header>
  );
}