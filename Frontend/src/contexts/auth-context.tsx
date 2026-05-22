"use client";

import type React from "react";

import { createContext, useState, useEffect, useCallback } from "react";
import AuthService from "@/services/auth/auth-service";
import { useRouter, usePathname } from "next/navigation";
import { useToast } from "@/components/ui/use-toast";
import { User } from "@/services/auth/user-schema";
import { Role } from "@/services/auth/role-schema";

// Tipos para o contexto de autenticação
export interface UserContext {
  id: number;
  name: string;
  username: string;
  role: Role; // <--- Usar a interface Role importada
  email?: string;
}

interface AuthContextType {
  user: UserContext | null;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
}

// Criação do contexto
export const AuthContext = createContext<AuthContextType>({
  user: null,
  isLoading: true,
  login: async () => false,
  logout: () => {},
});

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserContext | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();
  const pathname = usePathname(); // Hook para obter o pathname atual
  const { toast } = useToast();

  // Lista de rotas de autenticação (públicas)
  const AUTH_ROUTES = ["/login", "/signin"]; // Adicione outras rotas de auth se tiver

  // Função para mapear o usuário do backend para o contexto
  const mapUserToContext = useCallback((user: User): UserContext | null => {
    return {
      id: user.id,
      name: user.name,
      username: user.username,
      role: user.role,
      email: user.email,
    };
  }, []);

  useEffect(() => {
    //console.log("[EFFECT AuthProvider] Running. Pathname:", pathname,"Is Auth Route:", AUTH_ROUTES.includes(pathname));

    if (AUTH_ROUTES.includes(pathname)) {
      //console.log("[EFFECT AuthProvider] On auth route. Setting isLoading=false, user=null. Returning.");
      setIsLoading(false);
      setUser(null);
      return;
    }

    //console.log("[EFFECT AuthProvider] Not on auth route. Calling checkAuthStatus.");
    const checkAuthStatus = async () => {
      //console.log("[checkAuthStatus] Starting. Setting isLoading=true.");
      setIsLoading(true);
      try {
        const storedUserData = AuthService.getStoredUser();
        //console.log("Dados do usuário armazenados:", storedUserData);

        if (storedUserData) {
          setUser(mapUserToContext(storedUserData));
        }

        const isAuthenticated = await AuthService.checkSessionStatus();

        if (isAuthenticated) {
          const currentUserResponse = await AuthService.getCurrentUser();
          setUser(mapUserToContext(currentUserResponse.data.user));
        } else {
          //console.log("[checkAuthStatus] Not authenticated. Current pathname:",pathname,". Calling AuthService.logout() and router.push('/login').");
          AuthService.logout();
          setUser(null);
          if (pathname !== "/login") {
            router.push("/login");
          }
        }
      } catch (error) {
        console.error(
          "[checkAuthStatus] Error. Current pathname:",
          pathname,
          ". Calling AuthService.logout() and router.push('/login').",
          error
        );
        AuthService.logout();
        setUser(null);
        if (pathname !== "/login") {
          router.push("/login");
        }
      } finally {
        //console.log("[checkAuthStatus] Finished. Setting isLoading=false.");
        setIsLoading(false);
      }
    };

    checkAuthStatus();
  }, [mapUserToContext, pathname]);

  const login = async (
    username: string,
    password: string
  ): Promise<boolean> => {
    try {
      setIsLoading(true);
      const response = await AuthService.login({ username, password });
      //console.log("Resposta do login:", response);

      setUser(mapUserToContext(response.data.user));
      router.push("/dashboard");
      return true;
    } catch (error: any) {
      console.error("Erro ao fazer login:", error);
      toast({
        variant: "destructive",
        title: "Erro ao fazer login",
        description: error.response?.data?.message || "Credenciais inválidas",
      });

      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      await AuthService.logout();
      setUser(null);
      router.replace("/login");
    } catch (error) {
      console.error("Erro ao fazer logout:", error);
      toast({
        variant: "destructive",
        title: "Erro no Logout",
        description: "Não foi possível finalizar sua sessão.",
      });
    }
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
