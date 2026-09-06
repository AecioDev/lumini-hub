import { isAxiosError } from "axios";

interface ApiErrorBody {
  message?: string;
  error?: string;
  validationErrors?: string[];
}

// Extrai a mensagem mais específica que o backend mandou (envelope padrão em
// Backend/common/utils/response.go) — útil pra erros de conflito/validação
// que já vêm com um texto acionável (ex.: "não é possível excluir uma
// permissão que está sendo usada em um papel"), em vez de sempre mostrar um
// texto genérico.
export function getApiErrorMessage(error: unknown, fallback: string): string {
  if (isAxiosError(error)) {
    const body = error.response?.data as ApiErrorBody | undefined;
    if (body?.validationErrors?.length) return body.validationErrors[0];
    if (body?.error) return body.error;
    if (body?.message) return body.message;
  }
  return fallback;
}

// 403 = a sessão é válida mas falta permissão (RequirePermission do backend
// recusou). Diferente de 401 (sessão expirada, já tratado pelo interceptor
// do axios) — aqui a tela deve parar de tentar renderizar e mostrar
// explicitamente "sem acesso" em vez de um formulário pela metade.
export function isForbiddenError(error: unknown): boolean {
  return isAxiosError(error) && error.response?.status === 403;
}

// 404 = o recurso simplesmente não existe ainda (ex.: empresa sem
// Configuração Fiscal cadastrada) — não é erro de verdade, é um estado
// "vazio" que a tela deve tratar como formulário de criação, não mostrar
// mensagem de erro.
export function isNotFoundError(error: unknown): boolean {
  return isAxiosError(error) && error.response?.status === 404;
}
