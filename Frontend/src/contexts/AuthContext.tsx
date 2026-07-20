import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { authService } from "@/services/auth/auth-service";
import type { ApiUserDetail, LoginRequest } from "@/types/auth";
import type { ApiUserMenuItem } from "@/types/menu";

interface AuthContextValue {
  user: ApiUserDetail | null;
  menuItems: ApiUserMenuItem[];
  isBootstrapping: boolean;
  isAuthenticated: boolean;
  login: (payload: LoginRequest) => Promise<ApiUserDetail>;
  logout: () => Promise<void>;
  hasPermission: (permissionCode: string) => boolean;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

// Espelha utils.IsDeveloperOnlyPermission (Backend/common/utils/rbac.go): o
// catálogo de Perfis e Permissões só faz bypass pro perfil "Desenvolvedor" —
// o ADMIN passa em tudo, EXCETO nessas.
function isDeveloperOnlyPermission(permissionCode: string): boolean {
  return (
    permissionCode === "admin.create_permissions" ||
    permissionCode.startsWith("roles.") ||
    permissionCode.startsWith("permissions.")
  );
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<ApiUserDetail | null>(null);
  const [isBootstrapping, setIsBootstrapping] = useState(true);

  useEffect(() => {
    let cancelled = false;

    authService
      .me()
      .then((loggedUser) => {
        if (!cancelled) setUser(loggedUser);
      })
      .catch(() => {
        if (!cancelled) setUser(null);
      })
      .finally(() => {
        if (!cancelled) setIsBootstrapping(false);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  const login = useCallback(async (payload: LoginRequest) => {
    const loggedUser = await authService.login(payload);
    setUser(loggedUser);
    return loggedUser;
  }, []);

  const logout = useCallback(async () => {
    try {
      await authService.logout();
    } finally {
      setUser(null);
    }
  }, []);

  const isAdmin = user?.role?.name?.toUpperCase() === "ADMIN";
  const isDeveloper = user?.role?.name === "DEVELOP";

  const hasPermission = useCallback(
    (permissionCode: string) => {
      if (isDeveloper) return true;
      if (isAdmin && !isDeveloperOnlyPermission(permissionCode)) return true;
      return (user?.permissions ?? []).some((p) => p.permission === permissionCode);
    },
    [isAdmin, isDeveloper, user]
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      menuItems: user?.menu_items ?? [],
      isBootstrapping,
      isAuthenticated: user !== null,
      login,
      logout,
      hasPermission,
    }),
    [user, isBootstrapping, login, logout, hasPermission]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth deve ser usado dentro de <AuthProvider>");
  return ctx;
}
