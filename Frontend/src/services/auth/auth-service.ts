import api from "../common/api";
import { User } from "./user-schema";

export interface LoginCredentials {
  username: string;
  password: string;
}

// A AuthResponse será simplificada, pois os tokens não virão mais no body
export interface AuthResponse {
  success: boolean;
  message: string;
  data: {
    user: User; // Agora o backend só envia os dados do usuário no body
  };
}

const AuthService = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>("/auth/login", credentials);

    // Apenas armazenar as informações do usuário (não sensíveis) no localStorage para a UI
    // O backend já terá definido os cookies HttpOnly aqui.
    localStorage.setItem("user_data", JSON.stringify(response.data.data.user));

    return response.data;
  },

  async logout(): Promise<void> {
    try {
      // Faz uma requisição ao backend para limpar os cookies de autenticação.
      // O backend em Go agora tem a lógica para limpar os cookies HttpOnly.
      await api.post("/auth/logout");
    } catch (error) {
      console.error("Erro ao fazer logout no backend:", error);
    } finally {
      // Remover apenas os dados do usuário do localStorage no frontend
      localStorage.removeItem("user_data");
    }
  },

  // Esta função fará uma chamada para uma rota protegida. Se a chamada for bem-sucedida,
  // significa que o cookie de acesso foi enviado e é válido.
  async checkSessionStatus(): Promise<boolean> {
    try {
      // Tenta obter os dados do usuário atual de uma rota protegida (ex: /auth/me)
      // O cookie HttpOnly será enviado automaticamente pelo Axios.
      await api.get("/auth/me"); // Se o backend responder 200 OK, a sessão é válida.
      return true;
    } catch (error) {
      // Qualquer erro (401 Unauthorized, 403 Forbidden, erro de rede, etc.)
      // indica que a sessão não é válida ou o token expirou/está ausente.
      console.error("Erro ao verificar status da sessão:", error);
      return false;
    }
  },

  async getCurrentUser(): Promise<AuthResponse> {
    // Esta rota retorna os dados do usuário autenticado.
    // O cookie HttpOnly será enviado automaticamente.
    const response = await api.get<AuthResponse>("/auth/me");
    // Opcional: Atualizar os dados do usuário no localStorage com os mais recentes do backend
    localStorage.setItem("user_data", JSON.stringify(response.data.data.user));
    return response.data;
  },

  // getStoredUser agora lida apenas com os dados do usuário para a UI
  getStoredUser(): User | null {
    // Alterado para retornar apenas o tipo User
    const userData = localStorage.getItem("user_data"); // Pega de 'user_data'
    return userData ? JSON.parse(userData) : null;
  },
};

export default AuthService;
