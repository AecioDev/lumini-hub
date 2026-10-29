# Tasks - Módulo 0 - Setup do Ant Design

Acompanhamento da configuração da biblioteca de UI.

> **2026-07-19**: o frontend deixou de ser Next.js e virou **Vite + React (SPA client-side)** — ver `Documentos/Planejamento/Historico/log_2026-07-19.md` e o `CLAUDE.md` ("Frontend Architecture"). As tasks abaixo foram cumpridas com esse stack novo; a task do `@ant-design/nextjs-registry` não se aplica mais (Vite não tem SSR/Server Components, então não precisa de registry de extração de CSS-in-JS).

- [x] [CONCLUÍDO POR: Claude] Instalar as dependências: `antd`, `@ant-design/icons` via pnpm no frontend (`@ant-design/nextjs-registry` não se aplica — stack virou Vite/SPA)
- [x] [CONCLUÍDO POR: Claude, adaptado] ~~Configurar o `AntdRegistry` no arquivo `app/layout.tsx`~~ — não aplicável a Vite (sem SSR); `ConfigProvider` é montado direto em `src/App.tsx`
- [x] [CONCLUÍDO POR: Claude] Implementar o `ConfigProvider` global com o locale de português do Brasil (`ptBR`) — `src/App.tsx`, `antd/locale/pt_BR`
- [x] [CONCLUÍDO POR: Claude] Customizar a paleta de cores (Theme Tokens) no `ConfigProvider` para o padrão de identidade visual do Lumini Hub — `src/theme/antd-theme.ts` (`theme.darkAlgorithm`/`defaultAlgorithm`, tokens de cor/raio replicados do design de referência)
- [x] [CONCLUÍDO POR: Claude] Testar a renderização de componentes básicos (Button, Select, Table) — usados extensivamente em Login, Dashboard e módulo Usuários; sem SSR não há risco de quebra de hidratação
