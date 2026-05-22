// src/components/settings/permissions/view-permission-modal.tsx
"use client";

import { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useToast } from "@/components/ui/use-toast";
import { Loader2 } from "lucide-react";

import { Permission } from "@/services/auth/permission-schema";
import permissionService from "@/services/auth/permission-service";
import { PermissionForm } from "../forms/permission-form";

interface ViewPermissionModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  permissionId: string; // ID da permissão a ser visualizada
}

export function ViewPermissionModal({
  open,
  onOpenChange,
  permissionId,
}: ViewPermissionModalProps) {
  const { toast } = useToast();
  const [isLoadingData, setIsLoadingData] = useState(true); // Para o loading dos dados da permissão
  const [permissionData, setPermissionData] = useState<Permission | null>(null); // Dados da permissão carregada

  useEffect(() => {
    async function loadPermissionData() {
      if (!permissionId || !open) {
        setPermissionData(null);
        setIsLoadingData(false);
        return;
      }

      setIsLoadingData(true);
      try {
        const data = await permissionService.getPermissionById(permissionId);
        setPermissionData(data);
      } catch (error: any) {
        console.error(
          "Erro ao carregar dados da permissão para visualização:",
          error
        );
        toast({
          title: "Erro",
          description:
            error.message ||
            "Não foi possível carregar os dados da permissão para visualização.",
          variant: "destructive",
        });
        onOpenChange(false); // Fecha o modal se houver erro ao carregar
      } finally {
        setIsLoadingData(false);
      }
    }

    loadPermissionData();
  }, [permissionId, open, onOpenChange, toast]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Detalhes da Permissão</DialogTitle>
          <DialogDescription>
            Informações detalhadas da permissão selecionada.
          </DialogDescription>
        </DialogHeader>

        {isLoadingData ? (
          <div className="flex flex-col items-center justify-center h-48">
            <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
            <p className="text-muted-foreground">Carregando dados...</p>
          </div>
        ) : (
          <PermissionForm
            initialValues={permissionData || undefined} // Passa os dados carregados como valores padrão
            readOnly={true} // <-- Define o formulário como somente leitura
            // Não passa onSubmit nem isSubmitting, pois não há botão de submit
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
