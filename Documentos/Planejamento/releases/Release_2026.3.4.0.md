# Lumini Hub — Release 2026.3.4.0

**Versão do sistema:** 2026.3.4.0 · **Versão do banco de dados:** 3 (nova tabela `user_menu_favorites`; `dbversion.Required = 3`)
**Tema da versão:** **layout responsivo** (menu, cabeçalho e telas para computador, tablet e celular), **listas padronizadas** (cartões no celular, busca e filtros no servidor), **favoritos do menu** no cabeçalho e **limpeza** da leitura do banco legado no módulo de Integrações.
*Documento gerado em 2026-10-30. Legenda:* ✅ testado · ⚠️ implementado, ainda não testado · ⏳ pendente / não incluído.

---

## 1. Banco de dados

A versão do banco sobe de **2 para 3**. Há uma tabela nova; nenhuma tabela ou coluna existente muda e nenhum dado é reescrito.

| Tabela | Mudança | Detalhe |
|---|---|---|
| `user_menu_favorites` | **nova** ✅ | `id`, `user_id`, `menu_item_id`, `position` (ordem no cabeçalho), `created_at`. Índice único em `user_id` + `menu_item_id`. Sem chave estrangeira de propósito: um item que sumiu do menu, ou que o usuário deixou de enxergar, é ignorado na leitura. |

- Criada pelo `AutoMigrate` do `api.auth` no primeiro boot com `DB_AUTO_MIGRATE=true`, que também registra a versão 3 do banco.
- Com `DB_AUTO_MIGRATE=false`, os serviços **se recusam a subir** até o banco estar na versão 3: rode o `api.auth` uma vez com a variável ligada (ver a seção 5).

---

## 2. Novidades

### 2.1 Favoritos do menu no cabeçalho ✅
- Cada tela do menu lateral ganha uma **estrela**: aparece ao passar o mouse, fica dourada e fixa quando a tela é favorita, e no toque (celular) fica sempre visível.
- Até **5 favoritos por usuário**, gravados no banco (valem em qualquer máquina). Ao marcar o sexto, um aviso pede para remover um antes.
- Os favoritos aparecem **centralizados no cabeçalho**: ícone e nome em tela larga, só o ícone (com o nome numa dica) em tela média. Ao passar o mouse, o ícone dá um saltinho com zoom e uma linha se abre embaixo; a tela atual fica com a linha fixa.
- Abaixo do tablet (768 px) o cabeçalho não tem espaço: os favoritos aparecem no **topo da gaveta do menu**.
- Endpoint novo `PUT /auth/favorites` (troca a lista inteira; só aceita telas que o usuário enxerga, sem repetição, até 5). O login, o refresh e o `/auth/me` passam a devolver `favorite_menu_item_ids`.

### 2.2 Layout novo e responsivo (CFG-028) ✅
- O logo da marca sobe para o topo do menu lateral, na mesma altura do cabeçalho; a empresa ativa fica no pé do menu (ícone de troca, que abre a gaveta) e o nome dela também aparece na barra de status.
- O título sai do cabeçalho: o caminho da tela (breadcrumb) ganha uma faixa entre o cabeçalho e o corpo, e a busca do cabeçalho sai.
- Todas as páginas passam a usar o mesmo cabeçalho de página (título, descrição, botões de ação e Voltar), fixo no topo ao rolar.
- **Tablet e celular:** o menu lateral vira só ícones em tela média e uma gaveta em tela estreita; as tabelas rolam por dentro do cartão; a barra de status cabe numa linha; no celular o rótulo dos campos fica sempre em cima, os botões de ação das tabelas ficam maiores para o toque e os botões do cabeçalho ganham rótulos de acessibilidade.
- As bordas do menu seguem o tema claro e escuro, e o menu recolhido é lembrado ao recarregar.

### 2.3 Listas em cartões no celular (CFG-035) ✅
- Usuários, Empresas e Perfis e Permissões mostram, no celular, um cartão por registro (nome em destaque, Perfil, Status e dados principais) em vez de tabela que rola de lado. No computador e no tablet continua a tabela.
- As ações (ver, editar, excluir) ficam num menu de barras no canto do cartão, com confirmação para excluir.
- A coluna **ID** passa a aparecer em Usuários e em Perfis.
- Componente novo `ResponsiveTable`, que as próximas listas (Clientes, Fornecedores) já usam como padrão.

### 2.4 Busca e filtros no servidor (CFG-036) ✅
- A busca de **Usuários** e de **Permissões** passa a ser feita no servidor, em todas as páginas, e não só na página aberta (`POST /users/filter` e `POST /permissions/filter`).
- Nova **gaveta de filtros**, aberta pelo botão Filtros (com a quantidade de filtros ativos): Perfil e Status em Usuários; Módulo em Permissões.
- Todas as listas ganham a mesma barra (busca e Filtros à esquerda, botão de criar à direita) e o total de registros no pé, ao lado da paginação.

---

## 3. Correções

- **Busca de Permissões** ✅: a busca por "código/descrição" só procurava no código; agora procura nos dois.

---

## 4. Infraestrutura

- **Módulo de Integrações (INT-002)** ✅: o `api.integrations` deixou de ler o SQL Server legado. Removidos a camada `internal/legacy`, as rotas `/integrations/legacy/*`, o driver do SQL Server e as variáveis `MSSQL_*` (configuração, `.env.example` e `docker-compose.prod.yml`). A tela Configurações → Integrações perdeu o bloco "ERP Legado" e o aviso de indisponibilidade. As chaves `codlocarm_oficial`, `codlocarm_reserva` e `codemp` continuam em `integration_configs`, mas a tela não as edita mais.
- Swagger do gateway regerado (novos endpoints `/auth/favorites`, `/users/filter` e `/permissions/filter`).
- Nenhuma mudança nos Dockerfiles nem no Caddy.

---

## 5. Configuração necessária no servidor

- **Atualizar o código, reconstruir as imagens e reiniciar** `api-auth`, `api-core`, `api-integrations` e o frontend. Sem recriar o banco: os dados são mantidos.
- **Banco:** com `DB_AUTO_MIGRATE=true`, o `api-auth` cria a tabela `user_menu_favorites` e registra a versão 3 sozinho no boot. Com `DB_AUTO_MIGRATE=false`, rode o `api-auth` uma vez com a variável ligada antes de subir os demais serviços.
- **`.env` do servidor:** as variáveis `MSSQL_HOST`, `MSSQL_PORT`, `MSSQL_USER`, `MSSQL_PASSWORD` e `MSSQL_NAME` não são mais lidas e podem ser removidas.
- Depois de subir, a barra de status deve mostrar "Versão 2026.3.4.0 - BD v3".

---

## 6. Arquivos principais alterados / criados

**Biblioteca comum (`Backend/common`):** `utils/pagination.go` (`PaginationFromFilter`) e `pagination_filter_test.go` (novo), `config/config.go`, `database/db.go`, `dbversion/dbversion.go`, `go.mod`, `go.sum`; removido `database/db_sqlserver_test.go`.
**api.auth:** `internal/domain/` (`menu_favorite.go` novo, `user.go`, `permission.go`), `internal/service/` (`menu_favorites.go` novo e `menu_favorites_test.go`, `auth_service.go`, `user_service.go`, `permission_service.go`), `internal/repository/` (`user_repository.go`, `permission_repository.go`, `filter_repository_test.go` novo), `internal/handlers/` (`auth.go`, `users.go`, `permissions.go`), `internal/routes/routes.go`, `main.go`.
**api.integrations:** `main.go`, `internal/routes/routes.go`; removidos `internal/legacy/**` e `internal/handlers/legacy_lookups.go`.
**api.core:** `go.mod`, `go.sum`.
**api.gateway:** `docs/` (Swagger regerado).
**Frontend:** `components/common/` (`ResponsiveTable.tsx`, `ListToolbar.tsx`, `FilterDrawer.tsx` novos; `PageHeader.tsx`, `FormField.tsx`), `components/layout/` (`FavoriteLink.tsx`, `HeaderFavorites.tsx`, `FavoritesShortcuts.tsx`, `MobileMenu.tsx`, `SiderContent.tsx`, `AppBreadcrumb.tsx` novos; `AppHeader.tsx`, `AppSider.tsx`, `SidebarMenu.tsx`, `StatusBar.tsx`, `ActiveCompanySwitcher.tsx`), `components/users|roles|companies|permissions|integrations/`, `hooks/` (`useMenuFavorites.ts`, `useDebouncedValue.ts`, `useMediaQuery.ts`, `useActionButtonSize.ts`), `contexts/AuthContext.tsx`, `routes/AppLayout.tsx`, `pages/users/UsersListPage.tsx`, `pages/settings/`, `services/` (`auth`, `users`, `permissions`), `types/`, `utils/breadcrumb.ts`, `index.css`.
**Deploy e configuração:** `deploy/docker-compose.prod.yml`, `deploy/.env.example`, `Backend/.env.example`.

---

## 7. Pendências conhecidas (próximas versões)

| Item | Situação |
|---|---|
| **VND-001 / COM-001** Telas de Clientes e de Fornecedores (a API já valida tamanho e formato; a lista já nasce com `ResponsiveTable`, busca e filtros) | ⏳ |
| **Textos das permissões em bancos já semeados**: os textos e o módulo corrigidos no catálogo só valem em banco novo; nos já semeados é preciso um `UPDATE` à parte | ⏳ |
| **`DB_AUTO_MIGRATE=false`** recusando subir quando o banco está atrás da versão exigida | ⚠️ passou nos testes automáticos, ainda não testado num servidor real |
| **Integrações, Fase 2** (cliente da Loja Integrada e processamento dos webhooks) e destino das chaves `codlocarm_*` e `codemp` | ⏳ pausado |
| **HTTPS e domínio** no ambiente de testes | ⏳ Hoje o acesso é por IP, em HTTP |
