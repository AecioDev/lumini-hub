// src/components/common/confirm-delete-dialog.tsx
"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Loader2 } from "lucide-react"; // Usaremos este ícone para o loader

interface ConfirmDeleteDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /**
   * O nome do item que está sendo excluído (ex: "o cliente João", "a permissão users.create").
   * Será inserido na descrição do diálogo.
   */
  itemName: string;
  /**
   * A função a ser chamada quando a exclusão é confirmada.
   * Deve ser uma Promise<void> para que o loader funcione corretamente.
   */
  onConfirm: () => Promise<void>;
  /**
   * Título customizado para o diálogo. Padrão: "Confirmar Exclusão".
   */
  title?: string;
  /**
   * Descrição customizada para o diálogo.
   * Por padrão, usará: "Tem certeza que deseja excluir [itemName]? Esta ação não pode ser desfeita."
   */
  description?: string;
  /**
   * Opcional: Texto customizado para o botão de confirmação. Padrão: "Excluir".
   */
  confirmButtonText?: string;
}

export function ConfirmDeleteDialog({
  open,
  onOpenChange,
  itemName,
  onConfirm,
  title = "Confirmar Exclusão",
  description, // Será construída se não for fornecida
  confirmButtonText = "Excluir",
}: ConfirmDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState(false);

  // Constrói a descrição padrão se uma customizada não for fornecida
  const defaultDescription = `Tem certeza que deseja excluir ${itemName}? Esta ação não pode ser desfeita.`;
  const dialogDescription = description || defaultDescription;

  const handleDelete = async () => {
    setIsDeleting(true);
    await onConfirm(); // Chama a função onConfirm passada como prop
    setIsDeleting(false);
    // onOpenChange(false); // Você pode fechar o diálogo aqui ou deixar a função pai fazer isso após o toast
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{dialogDescription}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isDeleting}
          >
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={isDeleting}
          >
            {isDeleting ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Excluindo...
              </>
            ) : (
              confirmButtonText // Usa o texto customizado ou o padrão "Excluir"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
