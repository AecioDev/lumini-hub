// src/components/settings/permissions/create-permission-modal.tsx
"use client";

import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useToast } from "@/components/ui/use-toast";

import { CreatePermissionFormData } from "@/services/auth/permission-schema";
import permissionService from "@/services/auth/permission-service";
import { PermissionForm } from "../forms/permission-form";

interface CreatePermissionModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void; // Callback para quando uma permissão for criada com sucesso
}

export function CreatePermissionModal({
  open,
  onOpenChange,
  onSuccess,
}: CreatePermissionModalProps) {
  const { toast } = useToast();
  const [isSubmitting, setIsSubmitting] = useState(false); // Estado para controlar o loading do formulário

  const handleSubmit = async (values: CreatePermissionFormData) => {
    setIsSubmitting(true);
    try {
      const newPermission = await permissionService.createPermission(values);

      toast({
        title: "Sucesso!",
        description: `Permissão "${newPermission.permission}" criada com sucesso.`,
        variant: "default",
      });
      onOpenChange(false); // Fecha o modal ao sucesso
      if (onSuccess) {
        onSuccess(); // Chama o callback para o componente pai recarregar a lista, etc.
      }
    } catch (error: any) {
      console.error("Erro ao criar permissão:", error);
      toast({
        title: "Erro",
        description:
          error.message || "Ocorreu um erro inesperado ao criar a permissão.",
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
          <DialogTitle>Cadastrar Nova Permissão</DialogTitle>
          <DialogDescription>
            Preencha os dados para criar uma nova permissão no sistema.
          </DialogDescription>
        </DialogHeader>
        <PermissionForm onSubmit={handleSubmit} isSubmitting={isSubmitting} />
      </DialogContent>
    </Dialog>
  );
}
