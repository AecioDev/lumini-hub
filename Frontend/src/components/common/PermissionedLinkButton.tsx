// components/common/PermissionedLinkButton.tsx
"use client";

import Link from "next/link";
import { Button, ButtonProps } from "@/components/ui/button";
import { useHasPermission } from "@/hooks/use-has-permission";
import { cn } from "@/lib/utils";
import { ReactNode } from "react";
import { Icon } from "@iconify/react";
import { SimpleTooltip } from "./SimpleTooltip";

interface PermissionedLinkButtonProps extends Omit<ButtonProps, "disabled"> {
  href: string;
  permission?: string;
  tooltipMessage?: string;
  children: ReactNode;
  className?: string; // Para o Link/Button principal
  iconName?: string;
  forceDisabled?: boolean;
  tooltipDelayDuration?: number;
}

export function PermissionedLinkButton({
  href,
  permission,
  tooltipMessage,
  children,
  className,
  iconName,
  forceDisabled = false,
  tooltipDelayDuration,
  ...buttonProps
}: PermissionedLinkButtonProps) {
  const hasPermission = permission ? useHasPermission(permission) : true;
  const isDisabled = !hasPermission || forceDisabled;
  const showTooltip = isDisabled && tooltipMessage;

  // Conteúdo do botão (ícone e texto)
  // As classes de espaçamento do ícone já vêm do buttonVariants do seu Button do shadcn/ui
  const buttonContent = (
    <>
      {iconName && <Icon icon={iconName} />}
      {children}
    </>
  );

  const renderableButton = (
    <Button
      disabled={isDisabled}
      {...buttonProps}
      className={cn("w-full @lg:w-auto", className)}
    >
      {buttonContent}
    </Button>
  );

  // Se o botão estiver desabilitado E houver uma mensagem de tooltip,
  // envolvemos o botão em um wrapper para que o tooltip funcione.
  // Caso contrário, o tooltip é aplicado diretamente ao Link/Button.
  const elementToWrapInTooltip = showTooltip ? (
    <span className={cn("inline-block", className)}>{renderableButton}</span>
  ) : (
    renderableButton
  );

  // Renderiza o Link apenas se o botão não estiver desabilitado.
  // A classe `inline-block` é importante para que o Link ocupe o espaço correto e o tooltip funcione.
  const finalElement = !isDisabled ? (
    <Link href={href} className={cn("inline-block", className)}>
      {elementToWrapInTooltip}
    </Link>
  ) : (
    elementToWrapInTooltip
  );

  // Se houver uma mensagem de tooltip e o botão estiver desabilitado,
  // ou se houver uma mensagem de tooltip e o botão estiver habilitado mas queremos ele,
  // envolvemos o elemento final no SimpleTooltip.
  if (showTooltip) {
    return (
      <SimpleTooltip
        label={tooltipMessage!}
        delayDuration={tooltipDelayDuration}
      >
        {finalElement}
      </SimpleTooltip>
    );
  }

  // Caso contrário, retorna o elemento final diretamente.
  return finalElement;
}
