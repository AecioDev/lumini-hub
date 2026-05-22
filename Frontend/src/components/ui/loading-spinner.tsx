// src/components/ui/loading-spinner.tsx
import { Loader2 } from "lucide-react";
import React from "react";
import { cn } from "@/lib/utils"; // Assumindo que você tem um utilitário de classes

interface LoadingSpinnerProps {
  /**
   * Mensagem opcional a ser exibida abaixo do spinner.
   */
  message?: string;
  /**
   * Classes adicionais para o container principal.
   */
  containerClassName?: string;
  /**
   * Classes adicionais para o ícone do spinner.
   */
  spinnerClassName?: string;
  /**
   * Classes adicionais para o texto da mensagem.
   */
  messageClassName?: string;
}

export function LoadingSpinner({
  message = "Carregando...", // Mensagem padrão
  containerClassName,
  spinnerClassName,
  messageClassName,
}: LoadingSpinnerProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center space-y-3 p-6", // Flexbox, centralizado, espaçamento e padding
        containerClassName
      )}
    >
      <Loader2
        className={cn("h-10 w-10 animate-spin text-primary", spinnerClassName)}
      />
      {message && (
        <p
          className={cn(
            "text-lg font-medium text-muted-foreground",
            messageClassName
          )}
        >
          {message}
        </p>
      )}
    </div>
  );
}
