"use client";

import { useAuth } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LogOut, User, Sun, Moon, Check } from "lucide-react";
import { useTheme } from "next-themes";
import { usePresets } from "@/config/color-presets";
import { useColorPresetName, useColorPresets } from "@/hooks/use-theme-color";
import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";

export function UserSection() {
  const { user, logout } = useAuth();
  const { theme, setTheme } = useTheme();
  const COLOR_PRESETS = usePresets();
  const { colorPresets, setColorPresets } = useColorPresets();
  const { colorPresetName, setColorPresetName } = useColorPresetName();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const toggleTheme = () => {
    setTheme(theme === "dark" ? "light" : "dark");
  };

  const userName = user?.name || "Usuário";

  if (!mounted) return null;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          className="relative h-9 w-auto px-3 flex items-center gap-2 border border-border/40 hover:bg-accent rounded-full transition-all"
        >
          <div className="flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-primary">
            <User className="h-4 w-4" />
          </div>
          <span className="hidden sm:inline text-sm font-semibold text-foreground">
            {userName}
          </span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-64 p-2 bg-popover text-popover-foreground border border-border shadow-md">
        {/* Identificação do Usuário */}
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-semibold leading-none text-foreground">{userName}</p>
            {user?.email && (
              <p className="text-xs leading-none text-muted-foreground">
                {user.email}
              </p>
            )}
          </div>
        </DropdownMenuLabel>
        
        <DropdownMenuSeparator />

        {/* Alternância de Aparência (Tema) */}
        <DropdownMenuItem onClick={toggleTheme} className="cursor-pointer hover:bg-accent hover:text-accent-foreground py-2 rounded-md">
          {theme === "dark" ? (
            <Sun className="mr-2.5 h-4 w-4 text-yellow-500" />
          ) : (
            <Moon className="mr-2.5 h-4 w-4 text-indigo-400" />
          )}
          <span>Tema {theme === "dark" ? "Claro" : "Escuro"}</span>
        </DropdownMenuItem>

        <DropdownMenuSeparator />

        {/* Escolha das Cores de Destaque */}
        <div className="px-2 py-2.5 space-y-2">
          <span className="text-[11px] font-bold text-muted-foreground/85 uppercase tracking-wider block">
            Cor de Destaque
          </span>
          <div className="flex items-center justify-between gap-1.5 pt-1">
            {COLOR_PRESETS.map((preset) => {
              const isSelected =
                colorPresetName?.toLowerCase() === preset.name.toLowerCase();

              return (
                <button
                  key={preset.name}
                  title={preset.name}
                  onClick={() => {
                    setColorPresets(preset.colors);
                    setColorPresetName(preset.name.toLowerCase());
                  }}
                  className={cn(
                    "h-6 w-6 rounded-full border border-border/40 flex items-center justify-center transition-all duration-200 hover:scale-110 focus:outline-none shadow-sm relative",
                    isSelected ? "ring-2 ring-primary ring-offset-2 dark:ring-offset-gray-950" : ""
                  )}
                  style={{ backgroundColor: preset.colors.default }}
                >
                  {isSelected && (
                    <Check
                      className={cn(
                        "h-3.5 w-3.5 text-white",
                        preset.name === "Black" ? "text-black dark:text-white" : ""
                      )}
                    />
                  )}
                </button>
              );
            })}
          </div>
        </div>

        <DropdownMenuSeparator />

        {/* Opção de Logout (Sair) */}
        <DropdownMenuItem
          onClick={logout}
          className="text-destructive focus:bg-destructive focus:text-destructive-foreground cursor-pointer py-2 rounded-md"
        >
          <LogOut className="mr-2.5 h-4 w-4" />
          Sair
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
