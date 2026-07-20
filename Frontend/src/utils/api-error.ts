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
