import axios from "axios";
// Se você tiver acesso ao router do Next.js aqui (pode ser complexo,
// geralmente o router é um hook), seria ideal. Caso contrário, usaremos window.location.pathname.
// import router from 'next/router'; // Exemplo, não funcionará diretamente aqui

// Configuração base do axios
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api",
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

// Interceptor para tratar erros de resposta
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    const loginUrlPath = "/auth/login"; // Caminho da sua API de login
    // Ajuste se o AuthService.login usa um caminho diferente

    // Verifica se a URL original é a de login.
    // originalRequest.url pode ser o caminho relativo à baseURL.
    const isLoginAttempt = originalRequest.url?.split("?")[0] === loginUrlPath;

    // Se o erro for 401, não for uma tentativa de retry, E NÃO FOR UMA TENTATIVA DE LOGIN ORIGINAL
    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      !isLoginAttempt
    ) {
      originalRequest._retry = true;

      try {
        console.log(
          "[Interceptor] Tentando refresh token para:",
          originalRequest.url
        );
        await axios.post(
          `${api.defaults.baseURL}/auth/refresh-token`,
          {},
          {
            withCredentials: true,
          }
        );
        // Reenviar a requisição original
        return api(originalRequest);
      } catch (refreshError) {
        console.error(
          "[Interceptor] Falha ao tentar refresh token. Efetuando logout e redirecionando se necessário.",
          refreshError
        );
        // Limpar dados do usuário
        localStorage.removeItem("user_data");
        // Aqui, o ideal seria chamar uma função centralizada de logout que usa o router do Next.js.
        // Como alternativa, podemos usar window.location.href condicionalmente.
        // Verifique se a página atual JÁ NÃO É /login para evitar reload desnecessário.
        if (
          typeof window !== "undefined" &&
          window.location.pathname !== "/login"
        ) {
          window.location.href = "/login";
        }
        // Se já estiver em /login, não faz nada, o toast do AuthContext já foi exibido.
        return Promise.reject(refreshError); // Importante rejeitar para que a chamada original saiba da falha.
      }
    }

    // Se for uma tentativa de login que falhou, ou outro erro, apenas rejeita.
    return Promise.reject(error);
  }
);

export default api;
