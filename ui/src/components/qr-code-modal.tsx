"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { 
  QrCode, 
  RefreshCw, 
  CheckCircle, 
  AlertCircle,
  Smartphone,
  Loader2
} from "lucide-react";
import QRCodeLib from "qrcode";

interface QRCodeModalProps {
  sessionId: string;
  sessionName: string;
  isOpen: boolean;
  onClose: () => void;
}

type QRStatus = "generating" | "ready" | "expired" | "connected" | "error";

export function QRCodeModal({ 
  sessionId, 
  sessionName, 
  isOpen, 
  onClose 
}: QRCodeModalProps) {
  const [qrCode, setQrCode] = useState<string>("");
  const [status, setStatus] = useState<QRStatus>("generating");
  const [loading, setLoading] = useState(false);
  const [countdown, setCountdown] = useState(60);

  // Simular geração de QR code
  const generateQRCode = async () => {
    setLoading(true);
    setStatus("generating");
    
    try {
      // Simular delay da API
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Gerar QR code mock
      const mockQRData = `2@BQcAEAYQAg==,f/9u+vz6zJTzOD0VGOEkjrU=,wU/DdpXJ0tPalzxUr6SQBlMAAAAAElFTkSuQmCC,${sessionId},${Date.now()}`;
      const qrCodeDataURL = await QRCodeLib.toDataURL(mockQRData, {
        width: 256,
        margin: 2,
        color: {
          dark: '#000000',
          light: '#FFFFFF'
        }
      });
      
      setQrCode(qrCodeDataURL);
      setStatus("ready");
      setCountdown(60);
      
      // Simular expiração do QR code após 60 segundos
      const timer = setInterval(() => {
        setCountdown(prev => {
          if (prev <= 1) {
            clearInterval(timer);
            setStatus("expired");
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
      
      // Simular conexão bem-sucedida (50% de chance após 10-30 segundos)
      if (Math.random() > 0.5) {
        setTimeout(() => {
          clearInterval(timer);
          setStatus("connected");
        }, Math.random() * 20000 + 10000);
      }
      
    } catch (error) {
      setStatus("error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen && !qrCode) {
      generateQRCode();
    }
  }, [isOpen]);

  const handleRefresh = () => {
    setQrCode("");
    generateQRCode();
  };

  const getStatusInfo = () => {
    switch (status) {
      case "generating":
        return {
          icon: <Loader2 className="h-4 w-4 animate-spin" />,
          label: "Gerando",
          variant: "secondary" as const,
          description: "Gerando QR code para autenticação..."
        };
      case "ready":
        return {
          icon: <QrCode className="h-4 w-4" />,
          label: "Pronto",
          variant: "default" as const,
          description: `Escaneie o QR code com seu WhatsApp. Expira em ${countdown}s`
        };
      case "expired":
        return {
          icon: <AlertCircle className="h-4 w-4" />,
          label: "Expirado",
          variant: "destructive" as const,
          description: "QR code expirado. Clique em 'Gerar Novo' para tentar novamente."
        };
      case "connected":
        return {
          icon: <CheckCircle className="h-4 w-4" />,
          label: "Conectado",
          variant: "default" as const,
          description: "WhatsApp conectado com sucesso!"
        };
      case "error":
        return {
          icon: <AlertCircle className="h-4 w-4" />,
          label: "Erro",
          variant: "destructive" as const,
          description: "Erro ao gerar QR code. Tente novamente."
        };
    }
  };

  const statusInfo = getStatusInfo();

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Smartphone className="h-5 w-5" />
            Conectar {sessionName}
          </DialogTitle>
          <DialogDescription>
            Escaneie o QR code com seu WhatsApp para conectar a sessão
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Status Badge */}
          <div className="flex items-center justify-center">
            <Badge variant={statusInfo.variant} className="flex items-center gap-2">
              {statusInfo.icon}
              {statusInfo.label}
            </Badge>
          </div>

          {/* QR Code Display */}
          <div className="flex justify-center">
            <div className="relative">
              {status === "generating" ? (
                <div className="w-64 h-64 border-2 border-dashed border-muted-foreground/25 rounded-lg flex items-center justify-center">
                  <div className="text-center space-y-2">
                    <Loader2 className="h-8 w-8 animate-spin mx-auto text-muted-foreground" />
                    <p className="text-sm text-muted-foreground">
                      Gerando QR code...
                    </p>
                  </div>
                </div>
              ) : qrCode ? (
                <div className="relative">
                  <img 
                    src={qrCode} 
                    alt="QR Code" 
                    className={`w-64 h-64 rounded-lg ${
                      status === "expired" ? "opacity-50 grayscale" : ""
                    }`}
                  />
                  {status === "connected" && (
                    <div className="absolute inset-0 bg-green-500/20 rounded-lg flex items-center justify-center">
                      <CheckCircle className="h-16 w-16 text-green-500" />
                    </div>
                  )}
                </div>
              ) : null}
            </div>
          </div>

          {/* Status Description */}
          <Alert>
            <AlertDescription className="text-center">
              {statusInfo.description}
            </AlertDescription>
          </Alert>

          {/* Instructions */}
          {status === "ready" && (
            <div className="text-sm text-muted-foreground space-y-2">
              <p className="font-medium">Como conectar:</p>
              <ol className="list-decimal list-inside space-y-1 ml-4">
                <li>Abra o WhatsApp no seu celular</li>
                <li>Toque em "Mais opções" ou "Configurações"</li>
                <li>Toque em "Aparelhos conectados"</li>
                <li>Toque em "Conectar um aparelho"</li>
                <li>Aponte a câmera para este QR code</li>
              </ol>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex justify-center gap-2">
            {(status === "expired" || status === "error") && (
              <Button 
                onClick={handleRefresh} 
                disabled={loading}
                className="flex items-center gap-2"
              >
                <RefreshCw className="h-4 w-4" />
                Gerar Novo
              </Button>
            )}
            
            {status === "connected" && (
              <Button onClick={onClose} className="flex items-center gap-2">
                <CheckCircle className="h-4 w-4" />
                Concluir
              </Button>
            )}
            
            {status !== "connected" && (
              <Button variant="outline" onClick={onClose}>
                Cancelar
              </Button>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}