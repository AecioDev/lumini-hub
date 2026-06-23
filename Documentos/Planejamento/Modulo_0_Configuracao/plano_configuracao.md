# Módulo 0 - Setup do Ant Design no Frontend

Este módulo é responsável por preparar o ambiente frontend do Next.js para utilizar os componentes e estilos do **Ant Design (antd)**, garantindo a renderização correta no lado do servidor (SSR/App Router), internacionalização para pt-BR e definição dos tokens globais de design.

## 🛠️ Especificações Técnicas

### 1. Dependências do Frontend
Adicionar no `package.json` do [Frontend](file:///c:/Projetos/lumini-hub/Frontend/package.json):
- `antd` (Biblioteca de UI principal)
- `@ant-design/icons` (Pacote oficial de ícones do AntD)
- `@ant-design/nextjs-registry` (Componente de compatibilidade para extração de CSS-in-JS no Next.js App Router/Server Components)

### 2. Integração SSR e Registro de Estilos
Criar/configurar o registro de estilos no Next.js no arquivo `app/layout.tsx` (ou arquivo de componente equivalente) utilizando o `@ant-design/nextjs-registry`:
```tsx
import { AntdRegistry } from '@ant-design/nextjs-registry';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="pt-BR">
      <body>
        <AntdRegistry>{children}</AntdRegistry>
      </body>
    </html>
  );
}
```

### 3. Configuração de Internacionalização (Locale) e Tema
Envolver a árvore de componentes com o `<ConfigProvider>` do Ant Design configurado com o locale `ptBR` importado de `'antd/locale/pt_BR'` e definir a paleta de cores corporativa personalizada nos design tokens da aplicação.
