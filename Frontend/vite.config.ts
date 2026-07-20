import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Porta fixa em 3000: o CORS do api.gateway (Backend/microservices/api.gateway/main.go)
// só libera http://localhost:3000 e :3001.
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  server: {
    port: 3000,
    strictPort: true,
  },
});
