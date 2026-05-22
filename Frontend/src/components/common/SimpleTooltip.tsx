// components/common/SimpleTooltip.tsx
"use client";

import * as React from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip"; // Certifique-se que o caminho está correto para o seu shadcn/ui tooltip

interface SimpleTooltipProps {
  label: React.ReactNode; // Pode ser string ou JSX
  children: React.ReactNode;
  delayDuration?: number; // Tempo de atraso para o tooltip aparecer (em ms)
  side?: "top" | "right" | "bottom" | "left"; // Posição do tooltip (opcional)
  align?: "start" | "center" | "end"; // Alinhamento do tooltip (opcional)
  className?: string; // Classes adicionais para o TooltipContent (opcional)
}

export function SimpleTooltip({
  label,
  children,
  delayDuration = 200, // Padrão Shadcn/UI é 700ms, mas 200ms é um bom padrão para tooltips simples
  side,
  align,
  className,
}: SimpleTooltipProps) {
  return (
    <TooltipProvider delayDuration={delayDuration}>
      <Tooltip>
        <TooltipTrigger asChild>{children}</TooltipTrigger>
        <TooltipContent side={side} align={align} className={className}>
          {label}
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
}
