import { useEffect, useState } from "react";

// Devolve `value` só depois de ele parar de mudar por `delayMs`: usado na busca por texto das listas,
// para não chamar o servidor a cada tecla.
export function useDebouncedValue<T>(value: T, delayMs = 400): T {
  const [debounced, setDebounced] = useState(value);

  useEffect(() => {
    const timer = window.setTimeout(() => setDebounced(value), delayMs);
    return () => window.clearTimeout(timer);
  }, [value, delayMs]);

  return debounced;
}
