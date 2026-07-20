import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";

// Instância única do axios. baseURL aponta sempre para o api.gateway — nenhuma
// tela/serviço deve falar direto com api.auth/api.core (ver CLAUDE.md do projeto).
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

interface RetryableRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

// login/refresh-token nunca devem tentar se auto-renovar (evita loop).
function isLoginOrRefreshCall(url?: string): boolean {
  if (!url) return false;
  return url.includes("/auth/login") || url.includes("/auth/refresh-token");
}

// A checagem de sessão do AuthContext (bootstrap) já trata o 401 de forma
// gracil (usuário fica null, ProtectedRoute redireciona via SPA) — não faz
// sentido essa chamada específica forçar um reload duro da página.
function isMeCall(url?: string): boolean {
  return !!url && url.includes("/auth/me");
}

// Evita disparar N refreshes em paralelo quando várias requisições tomam 401 juntas.
let refreshInFlight: Promise<void> | null = null;

function redirectToLogin() {
  if (typeof window !== "undefined" && window.location.pathname !== "/login") {
    window.location.href = "/login";
  }
}

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as RetryableRequestConfig | undefined;

    const shouldTryRefresh =
      error.response?.status === 401 &&
      original &&
      !original._retry &&
      !isLoginOrRefreshCall(original.url);

    if (!shouldTryRefresh) {
      return Promise.reject(error);
    }

    original._retry = true;

    try {
      refreshInFlight ??= api.post("/auth/refresh-token").then(() => undefined);
      await refreshInFlight;
      refreshInFlight = null;
      return api(original);
    } catch (refreshError) {
      refreshInFlight = null;
      if (!isMeCall(original.url)) {
        redirectToLogin();
      }
      return Promise.reject(refreshError);
    }
  }
);
