import { Role } from "@/services/auth/role-service";
import { atom } from "jotai";

interface User {
  id: number;
  name: string;
  username: string;
  role: Role; // <--- Usar a interface Role importada
  email?: string;
}

// Função para obter o usuário do localStorage (só roda no cliente)
const getInitialUser = (): User | null => {
  if (typeof window !== "undefined") {
    // Verifica se está no ambiente do navegador
    try {
      const storedUser = localStorage.getItem("user_data"); // Supondo que você armazena o usuário aqui
      return storedUser ? JSON.parse(storedUser) : null;
    } catch (error) {
      console.error("Erro ao ler usuário do localStorage:", error);
      return null;
    }
  }
  return null; // Retorna null no SSR ou se houver erro
};

// O seu userAtom!
// Usamos a função getInitialUser para definir o valor inicial.
// No SSR, será null. No cliente, tentará ler do localStorage.
export const userAtom = atom<User | null>(getInitialUser());

// Opcional: Um átomo para gerenciar se o usuário está logado
export const isLoggedInAtom = atom(
  (get) => get(userAtom) !== null // Retorna true se userAtom não for null
);

// Opcional: Um átomo para setar o usuário e salvar no localStorage
// Se você for usar para login/logout, isso aqui é super útil
export const setUserAndLocalStorageAtom = atom(
  null, // O primeiro argumento é null porque este é um átomo de escrita
  (get, set, user: User | null) => {
    set(userAtom, user); // Atualiza o valor do userAtom

    if (typeof window !== "undefined") {
      if (user) {
        localStorage.setItem("user_data", JSON.stringify(user));
      } else {
        localStorage.removeItem("user_data");
      }
    }
  }
);
