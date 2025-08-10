"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
  Smartphone, 
  Send, 
  CheckCircle, 
  AlertCircle,
  Loader2,
  Phone
} from "lucide-react";

interface PairPhoneModalProps {
  sessionId: string;
  sessionName: string;
  isOpen: boolean;
  onClose: () => void;
}

type PairStatus = "input" | "sending" | "code_sent" | "verifying" | "success" | "error";

export function PairPhoneModal({ 
  sessionId, 
  sessionName, 
  isOpen, 
  onClose 
}: PairPhoneModalProps) {
  const [phoneNumber, setPhoneNumber] = useState("");
  const [pairCode, setPairCode] = useState("");
  const [status, setStatus] = useState<PairStatus>("input");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const validatePhoneNumber = (phone: string) => {
    // Validação básica para número brasileiro
    const cleanPhone = phone.replace(/\D/g, '');
    return cleanPhone.length >= 10 && cleanPhone.length <= 13;
  };

  const formatPhoneNumber = (phone: string) => {
    const cleanPhone = phone.replace(/\D/g, '');
    
    if (cleanPhone.length <= 2) {
      return cleanPhone;
    } else if (cleanPhone.length <= 7) {
      return `(${cleanPhone.slice(0, 2)}) ${cleanPhone.slice(2)}`;
    } else if (cleanPhone.length <= 11) {
      return `(${cleanPhone.slice(0, 2)}) ${cleanPhone.slice(2, 7)}-${cleanPhone.slice(7)}`;
    } else {
      return `+${cleanPhone.slice(0, 2)} (${cleanPhone.slice(2, 4)}) ${cleanPhone.slice(4, 9)}-${cleanPhone.slice(9, 13)}`;
    }
  };

  const handleSendCode = async () => {
    if (!validatePhoneNumber(phoneNumber)) {
      setError("Número de telefone inválido");
      return;
    }

    setLoading(true);
    setError("");
    setStatus("sending");

    try {
      // Simular envio do código
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Simular código gerado
      const mockCode = Math.random().toString(36).substring(2, 8).toUpperCase();
      setPairCode(mockCode);
      setStatus("code_sent");
      
    } catch (error) {
      setStatus("error");
      setError("Erro ao enviar código de emparelhamento");
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyCode = async () => {
    setLoading(true);
    setStatus("verifying");

    try {
      // Simular verificação
      await new Promise(resolve => setTimeout(resolve, 1500));
      
      // Simular sucesso (80% de chance)
      if (Math.random() > 0.2) {
        setStatus("success");
        
        // Fechar modal após 2 segundos
        setTimeout(() => {
          onClose();
          resetModal();
        }, 2000);
      } else {
        setStatus("error");
        setError("Código inválido ou expirado");
      }
      
    } catch (error) {
      setStatus("error");
      setError("Erro ao verificar código");
    } finally {
      setLoading(false);
    }
  };

  const resetModal = () => {
    setPhoneNumber("");
    setPairCode("");
    setStatus("input");
    setError("");
    setLoading(false);
  };

  const handleClose = () => {
    onClose();
    resetModal();
  };

  const getStatusInfo = () => {
    switch (status) {
      case "input":
        return {
          icon: <Phone className="h-4 w-4" />,
          label: "Aguardando",
          variant: "secondary" as const,
          description: "Digite seu número de telefone para receber o código"
        };
      case "sending":
        return {
          icon: <Loader2 className="h-4 w-4 animate-spin" />,
          label: "Enviando",
          variant: "secondary" as const,
          description: "Enviando código de emparelhamento..."
        };
      case "code_sent":
        return {
          icon: <Send className="h-4 w-4" />,
          label: "Código Enviado",
          variant: "default" as const,
          description: "Código enviado! Verifique seu WhatsApp"
        };
      case "verifying":
        return {
          icon: <Loader2 className="h-4 w-4 animate-spin" />,
          label: "Verificando",
          variant: "secondary" as const,
          description: "Verificando código de emparelhamento..."
        };
      case "success":
        return {
          icon: <CheckCircle className="h-4 w-4" />,
          label: "Conectado",
          variant: "default" as const,
          description: "Emparelhamento realizado com sucesso!"
        };
      case "error":
        return {
          icon: <AlertCircle className="h-4 w-4" />,
          label: "Erro",
          variant: "destructive" as const,
          description: error || "Erro no emparelhamento"
        };
    }
  };

  const statusInfo = getStatusInfo();

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Smartphone className="h-5 w-5" />
            Emparelhar {sessionName}
          </DialogTitle>
          <DialogDescription>
            Emparelhe seu telefone com esta sessão usando um código
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

          {/* Status Description */}
          <Alert>
            <AlertDescription className="text-center">
              {statusInfo.description}
            </AlertDescription>
          </Alert>

          {/* Phone Input */}
          {(status === "input" || status === "sending") && (
            <div className="space-y-2">
              <Label htmlFor="phone">Número de Telefone</Label>
              <Input
                id="phone"
                type="tel"
                placeholder="+55 (11) 99999-9999"
                value={formatPhoneNumber(phoneNumber)}
                onChange={(e) => setPhoneNumber(e.target.value)}
                disabled={loading}
              />
              <p className="text-xs text-muted-foreground">
                Digite o número com código do país (ex: +55 11 99999-9999)
              </p>
            </div>
          )}

          {/* Code Display */}
          {status === "code_sent" && (
            <div className="space-y-2">
              <Label>Código de Emparelhamento</Label>
              <div className="p-4 bg-muted rounded-lg text-center">
                <p className="text-2xl font-mono font-bold tracking-wider">
                  {pairCode}
                </p>
              </div>
              <p className="text-xs text-muted-foreground text-center">
                Digite este código no seu WhatsApp para conectar
              </p>
            </div>
          )}

          {/* Instructions */}
          {status === "code_sent" && (
            <div className="text-sm text-muted-foreground space-y-2">
              <p className="font-medium">Como usar o código:</p>
              <ol className="list-decimal list-inside space-y-1 ml-4">
                <li>Abra o WhatsApp no seu celular</li>
                <li>Toque em "Mais opções" ou "Configurações"</li>
                <li>Toque em "Aparelhos conectados"</li>
                <li>Toque em "Conectar um aparelho"</li>
                <li>Escolha "Conectar com código do telefone"</li>
                <li>Digite o código mostrado acima</li>
              </ol>
            </div>
          )}

          {/* Success Message */}
          {status === "success" && (
            <div className="text-center space-y-2">
              <CheckCircle className="h-16 w-16 text-green-500 mx-auto" />
              <p className="text-lg font-medium">Emparelhamento Concluído!</p>
              <p className="text-sm text-muted-foreground">
                Sua sessão foi conectada com sucesso
              </p>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex justify-center gap-2">
            {status === "input" && (
              <>
                <Button variant="outline" onClick={handleClose}>
                  Cancelar
                </Button>
                <Button 
                  onClick={handleSendCode} 
                  disabled={loading || !validatePhoneNumber(phoneNumber)}
                  className="flex items-center gap-2"
                >
                  <Send className="h-4 w-4" />
                  Enviar Código
                </Button>
              </>
            )}

            {status === "code_sent" && (
              <>
                <Button variant="outline" onClick={() => setStatus("input")}>
                  Voltar
                </Button>
                <Button 
                  onClick={handleVerifyCode} 
                  disabled={loading}
                  className="flex items-center gap-2"
                >
                  <CheckCircle className="h-4 w-4" />
                  Verificar Conexão
                </Button>
              </>
            )}

            {status === "error" && (
              <>
                <Button variant="outline" onClick={handleClose}>
                  Cancelar
                </Button>
                <Button onClick={() => setStatus("input")}>
                  Tentar Novamente
                </Button>
              </>
            )}

            {(status === "sending" || status === "verifying") && (
              <Button variant="outline" onClick={handleClose}>
                Cancelar
              </Button>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}