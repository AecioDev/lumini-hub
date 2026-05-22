// components/ContainerLoader.tsx
import React from "react";
import { cn } from "@/lib/utils"; // Importe sua função cn do shadcn/ui
// import { seuSpinnerAqui } from "@/components/ui/spinner"; // Importe seu componente de spinner aqui

interface ContainerLoaderProps extends React.HTMLAttributes<HTMLDivElement> {
  children?: React.ReactNode; // Permite passar um texto ou outro elemento como children
  showText?: boolean; // Controla se o texto "carregando..." é exibido
}

export function ContainerLoader({
  children,
  showText = true, // Por padrão, mostra o texto
  className,
  ...props
}: ContainerLoaderProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center h-full w-full",
        className
      )}
      {...props}
    >
      {/* Coloque seu componente de spinner aqui */}
      {/* Exemplo: <seuSpinnerAqui className="h-8 w-8 text-primary" /> */}

      {/* Apenas um exemplo de spinner, substitua pelo seu! */}
      <svg
        className="animate-spin h-8 w-8 text-primary"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          className="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          strokeWidth="4"
        ></circle>
        <path
          className="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        ></path>
      </svg>

      {showText && (
        <div className="mt-4 text-center">
          {children ? children : <h1 className="text-xl">carregando...</h1>}
        </div>
      )}
    </div>
  );
}
