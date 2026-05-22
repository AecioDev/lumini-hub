// components/layout/sidebarItem.tsx
"use client";

import Link from "next/link";
import { useState } from "react";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
// MUDANÇA 1: Remover imports de lucide-react para ícones de seta
// import { ChevronDown, ChevronRight } from "lucide-react";

// MUDANÇA 2: Importar o componente Icon do Iconify
import { Icon } from "@iconify/react";

// MUDANÇA 3: Importar a interface NavItem do navigation.ts (já está correto)
import { NavItem } from "@/config/navigation";

interface SidebarItemProps {
  item: NavItem;
  collapsed?: boolean;
  // level?: number; // Opcional: para controlar indentação ou estilos de nível
}

export function SidebarItem({
  item,
  collapsed /*, level = 0*/,
}: SidebarItemProps) {
  const pathname = usePathname();
  const [open, setOpen] = useState(false); // Estado para submenu expandido/colapsado
  const [popoverOpen, setPopoverOpen] = useState(false); // Estado para Popover (quando sidebar está colapsado)

  // Verifica se o item de navegação está ativo
  const isActive = (href?: string) =>
    href && (pathname === href || pathname.startsWith(`${href}/`));

  // Lógica para alternar a abertura do submenu
  const toggleOpen = () => setOpen((prev) => !prev);

  // Se o item tem filhos, ele é um item de menu pai (com submenu)
  if (item.children && item.children.length > 0) {
    // Opção 1: Sidebar colapsado - Usar Popover para mostrar submenus
    if (collapsed) {
      return (
        <li key={item.title}>
          <Popover open={popoverOpen} onOpenChange={setPopoverOpen}>
            <PopoverTrigger asChild>
              <Button
                variant={
                  isActive(item.href) || popoverOpen ? "secondary" : "ghost"
                }
                className="w-full justify-center px-0"
              >
                {/* MUDANÇA 4: Usar componente Icon do Iconify para o ícone principal */}
                {item.icon && (
                  <Icon icon={item.icon} className="h-5 w-5 mx-auto" />
                )}
                <span className="sr-only">{item.title}</span>
              </Button>
            </PopoverTrigger>
            <PopoverContent
              side="right"
              align="start"
              className="w-48 p-2"
              onMouseLeave={() => setPopoverOpen(false)}
            >
              {item.children.map((childItem) => (
                <SidebarItem
                  key={childItem.title}
                  item={childItem}
                  collapsed={false}
                  // level={level + 1}
                />
              ))}
            </PopoverContent>
          </Popover>
        </li>
      );
    }
    // Opção 2: Sidebar expandido - Usar expansão normal de submenu
    else {
      return (
        <li key={item.title}>
          <Button
            variant={isActive(item.href) || open ? "secondary" : "ghost"}
            onClick={toggleOpen}
            className="w-full justify-start pr-2"
          >
            {/* MUDANÇA 5: Usar componente Icon do Iconify para o ícone principal */}
            {item.icon && <Icon icon={item.icon} className="h-5 w-5 mr-2" />}
            <span className="flex-1 text-left">{item.title}</span>
            {/* MUDANÇA 6: Usar ícones MDI para as setas de expansão */}
            {open ? (
              <Icon
                icon="mdi:chevron-up"
                className="ml-auto h-4 w-4 transition-transform rotate-180"
              /> // Seta para cima
            ) : (
              <Icon
                icon="mdi:chevron-down"
                className="ml-auto h-4 w-4 transition-transform rotate-0"
              /> // Seta para baixo
            )}
          </Button>
          {open && (
            <ul
              className={cn("mt-1 space-y-1", item.children ? "ml-4" : "ml-0")}
            >
              {item.children.map((childItem) => (
                <SidebarItem
                  key={childItem.title}
                  item={childItem}
                  collapsed={false}
                  // level={level + 1}
                />
              ))}
            </ul>
          )}
        </li>
      );
    }
  }

  // Se o item NÃO tem filhos, ele é um item de menu folha (final)
  return (
    <li key={item.title}>
      <Link href={item.href || "#"}>
        <Button
          variant={isActive(item.href) ? "secondary" : "ghost"}
          className={cn(
            "w-full justify-start",
            collapsed ? "px-0 justify-center" : "pl-2"
            // level > 0 && !collapsed && `pl-${level * 4}`
          )}
        >
          {/* MUDANÇA 7: Usar componente Icon do Iconify para o ícone principal */}
          {item.icon && (
            <Icon
              icon={item.icon}
              className={cn("h-5 w-5", collapsed ? "mx-auto" : "mr-2")}
            />
          )}
          {!collapsed && <span>{item.title}</span>}
        </Button>
      </Link>
    </li>
  );
}
