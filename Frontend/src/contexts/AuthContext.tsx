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
  // Identidade de papel "crua" — pra UI que precisa saber se é literalmente
  // o perfil DEVELOP (ex.: mostrar código técnico de permissão), não só "tem
  // essa permissão" (hasPermission pode dar falso positivo se o usuário
  // tiver um grant avulso de uma permissão normalmente exclusiva do DEVELOP).
  isDeveloper: boolean;
  login: (payload: LoginRequest) => Promise<ApiUserDetail>;
  logout: () => Promise<void>;
  hasPermission: (permissionCode: string) => boolean;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

// Espelha utils.IsDeveloperOnlyPermission (Backend/common/utils/rbac.go): o
// catálogo de Perfis (`roles.*` — criar/editar/excluir/atribuir permissões)
// é do ADMIN, sem restrição — quem manda é a permissão em si, checada em
// tempo real pelo backend. Ficam exclusivos do DEVELOP: criar/editar/excluir
// no catálogo de Permissões (`permissions.create/edit/delete`) e
// `admin.create_permissions` (CRUD de /menu-items). `permissions.view` fica
// de fora (mesmo motivo de `roles.view`): é só leitura, e o módulo Develop
// já não vem na resposta pra quem não é DEVELOP.
function isDeveloperOnlyPermission(permissionCode: string): boolean {
  if (permissionCode === "permissions.view") {
    return false;
  }
  return (
    permissionCode === "admin.create_permissions" ||
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
      isDeveloper,
      login,
      logout,
      hasPermission,
    }),
    [user, isBootstrapping, isDeveloper, login, logout, hasPermission]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth deve ser usado dentro de <AuthProvider>");
  return ctx;
}
