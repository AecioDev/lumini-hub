// src/components/settings/permissions/edit-permission-modal.tsx
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
import { Loader2 } from "lucide-react"; // Importe o ícone de loading
import {
  Permission,
  UpdatePermissionFormData, // Assumindo que você terá um tipo para dados de atualização
} from "@/services/auth/permission-schema";
import permissionService from "@/services/auth/permission-service";
import { PermissionForm } from "../forms/permission-form"; // Reutilizando o formulário

interface EditPermissionModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  permissionId: string; // ID da permissão a ser editada
  onSuccess?: () => void; // Callback para quando a permissão for atualizada com sucesso
}

export function EditPermissionModal({
  open,
  onOpenChange,
  permissionId,
  onSuccess,
}: EditPermissionModalProps) {
  const { toast } = useToast();
  const [isSubmitting, setIsSubmitting] = useState(false); // Para o loading do formulário
  const [isLoadingData, setIsLoadingData] = useState(true); // Para o loading dos dados da permissão
  const [permissionData, setPermissionData] = useState<Permission | null>(null); // Dados da permissão carregada

  // Hook para carregar os dados da permissão quando o modal abre ou o ID muda
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
        console.error("Erro ao carregar dados da permissão:", error);
        toast({
          title: "Erro",
          description:
            error.message ||
            "Não foi possível carregar os dados da permissão para edição.",
          variant: "destructive",
        });
        onOpenChange(false); // Fecha o modal se houver erro ao carregar
      } finally {
        setIsLoadingData(false);
      }
    }

    loadPermissionData();
  }, [permissionId, open, onOpenChange, toast]); // Dependências do useEffect

  const handleSubmit = async (values: UpdatePermissionFormData) => {
    if (!permissionData) return; // Não envia se não há dados carregados

    setIsSubmitting(true);
    try {
      // Chama o serviço de atualização, passando o ID e os valores
      const updatedPermission = await permissionService.updatePermission(
        permissionData.id, // Usa o ID da permissão carregada
        values
      );

      toast({
        title: "Sucesso!",
        description: `Permissão "${updatedPermission.permission}" atualizada com sucesso.`,
        variant: "default",
      });
      onOpenChange(false); // Fecha o modal ao sucesso
      if (onSuccess) {
        onSuccess(); // Chama o callback para o componente pai recarregar a lista, etc.
      }
    } catch (error: any) {
      console.error("Erro ao atualizar permissão:", error);
      toast({
        title: "Erro",
        description:
          error.message ||
          "Ocorreu um erro inesperado ao atualizar a permissão.",
        variant: "destructive",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Editar Permissão</DialogTitle>
          <DialogDescription>
            Altere os dados da permissão existente.
          </DialogDescription>
        </DialogHeader>

        {isLoadingData ? (
          <div className="flex flex-col items-center justify-center h-48">
            <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
            <p className="text-muted-foreground">Carregando dados...</p>
          </div>
        ) : (
          <PermissionForm
            onSubmit={handleSubmit}
            isSubmitting={isSubmitting}
            initialValues={permissionData || undefined} // Passa os dados carregados como valores padrão
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
