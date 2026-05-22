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
import { LogOut, User, Sun, Moon, Settings } from "lucide-react"; // Adicionei Settings como sugestão
import { useTheme } from "next-themes";

export function UserSection() {
  const { user, logout } = useAuth();
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme(theme === "dark" ? "light" : "dark");
  };

  // Adicionar um placeholder se o nome do usuário não estiver carregado ainda ou não existir
  const userName = user?.name || "Usuário";

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          className="relative h-8 w-auto px-2 sm:px-3 flex items-center gap-2"
        >
          {" "}
          {/* Ajuste de padding e tamanho */}
          <User className="h-5 w-5 flex-shrink-0" />
          <span className="hidden sm:inline text-sm font-medium">
            {userName}
          </span>{" "}
          {/* Ocultar nome em telas pequenas se necessário */}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        {" "}
        {/* Definir uma largura pode ser bom */}
        <DropdownMenuLabel>
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">{userName}</p>
            {user?.email && ( // Mostra o email se disponível
              <p className="text-xs leading-none text-muted-foreground">
                {user.email}
              </p>
            )}
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={toggleTheme}>
          {theme === "dark" ? (
            <Sun className="mr-2 h-4 w-4" />
          ) : (
            <Moon className="mr-2 h-4 w-4" />
          )}
          <span>Tema {theme === "dark" ? "Claro" : "Escuro"}</span>
        </DropdownMenuItem>
        {/* Sugestão: Adicionar um link para Configurações de Perfil se aplicável */}
        {/* <DropdownMenuItem onClick={() => router.push('/profile/settings')}> // Exemplo de navegação
          <Settings className="mr-2 h-4 w-4" />
          <span>Configurações</span>
        </DropdownMenuItem> */}
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onClick={logout}
          className="text-destructive focus:bg-destructive focus:text-destructive-foreground" // Ajustado para Shadcn UI v0.8+
        >
          <LogOut className="mr-2 h-4 w-4" />
          Sair
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
