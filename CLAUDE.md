# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Lumini Hub is an ERP system with a Go microservices backend (`Backend/`) and a Next.js frontend (`Frontend/`). It is being migrated feature-by-feature from a legacy monolith, organized around ERP modules (CRM, Estoque/Produtos, Compras, Vendas/Caixas, Financeiro, Fiscal, Contabilidade) tracked in `Documentos/Planejamento/`.

## Commands

### Backend (Go, from `Backend/`)
```powershell
.\run_services.bat        # starts api.gateway (4000), api.auth (4001), api.core (4002) in separate windows
go build ./...             # build all workspace modules
go test ./...               # run tests (few exist today, e.g. common/database/db_sqlserver_test.go)
```
Requires a local PostgreSQL on port 5432 and a `Backend/.env` (see `Backend/.env.example`).

Regenerate Swagger docs after adding/changing handlers (run from `Backend/microservices/api.gateway/`):
```powershell
swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
```

### Frontend (Next.js, from `Frontend/`)
```powershell
pnpm install
pnpm dev      # dev server at http://localhost:3000
pnpm build
pnpm lint     # next lint / eslint
```
No frontend test runner is currently configured. Requires `Frontend/.env.local` with `NEXT_PUBLIC_API_URL` (defaults to `http://localhost:4000/api`, i.e. the gateway).

Local login: `admin` / `987321`.

## Backend Architecture

Go Workspace (`Backend/go.work`) with four modules:
- **`common/`** — shared package: `config` (env loading), `database` (GORM/Postgres pool), `middlewares` (`AuthMiddleware` JWT-from-cookie, `RequirePermission` RBAC), `utils` (JWT, bcrypt, `response.go`, pagination), `repository` (generic `Repository[T]` / `GormRepository[T]`).
- **`microservices/api.gateway`** (port 4000) — sole entrypoint. Reverse-proxies by path prefix to the other services and is the **only** place CORS is configured.
- **`microservices/api.auth`** (port 4001) — users, roles, permissions, menu items (dynamic sidebar tree), login/refresh (HTTP-only cookies).
- **`microservices/api.core`** (port 4002) — customers, suppliers, addresses, contacts, documents.

Both services share one physical Postgres database during this migration phase but never join across each other's tables — cross-service data (e.g. `User` inside a Customer DTO) is passed as plain IDs and re-hydrated into a local simplified struct (e.g. `ApiUser`), not a GORM relation.

### Mandatory backend patterns
- **Responses**: every handler returns via `utils.SuccessResponse` / `utils.ErrorResponse` / `utils.ValidationErrorResponse` ([response.go](Backend/common/utils/response.go)). Never hand-roll `gin.H` payloads.
- **Repository**: microservice-specific repositories extend the generic `commonrepo.Repository[T]` and override methods like `FindByID` when GORM `Preload` is needed.
- **Unit of Work**: services never touch a repository or `*gorm.DB` directly — they call `s.uow.Execute(func(uow repository.UnitOfWork) error { ... })` so writes are atomic. Each microservice defines its own `UnitOfWork` interface (see `api.core/internal/repository/repository.go`) exposing one method per aggregate (`Customers()`, `Suppliers()`, ...).
- **Search**: complex/paginated queries are POST endpoints ending in `/filter` (a `<Entity>FilterRequest` DTO with `page_no`, `page_size`, `order_by_column`, `is_asc`), not query-string GETs. Results go through `utils.Paginate`.
- **CORS**: only ever configured in `api.gateway/main.go`. Adding it to `api.auth` or `api.core` causes duplicate-header `AxiosError: Network Error` in the browser.
- **RBAC (hybrid)**: a `Role` is only a template/default — creating or editing a user picks a Role, which pre-fills a default set of permissions, but what `RequirePermission` middleware actually checks is the permission list embedded in the JWT, which comes from the user's own direct `user_permissions` link (`api.auth/internal/domain/user.go`), not from the Role. Operators can freely add/revoke individual permissions on a user beyond their Role's defaults (`PUT /users/:id/permissions`). `ADMIN` role still bypasses all permission checks. The Role/permission catalog itself (`/settings/roles`, "Perfis e Permissões") is a developer-only screen gated by the `admin.create_permissions` permission — creating a permission there has no effect unless a matching `RequirePermission("...")` already exists in code.
- **Routing**: new route prefixes must be registered both in the microservice's own routes file and proxied from `api.gateway/main.go` (via `router.NoRoute`, careful not to collide with `/swagger/*any`).

### Adding a new entity (Core or Auth microservice)
Follow this order (see `.agents/skills/lumini_hub_entity_creation/SKILL.md` for full detail):
1. `internal/models` (or `domain`) — GORM struct + `Create<Entity>Request`/`Update<Entity>Request` with `binding` validation tags.
2. `internal/models/<entity>_dto.go` — DTOs + `ToDTO()`/`ToDetailDTO()` mappers.
3. `internal/repository/<entity>_repository.go` — interface extending `commonrepo.Repository[domain.<Entity>]` + GORM impl, with `Preload` overrides as needed.
4. `internal/validator/<entity>_validator.go` — business-rule validation (uniqueness, nested items).
5. `internal/service/<entity>_service.go` — business logic, all writes wrapped in `s.uow.Execute(...)`.
6. `internal/handlers/<entity>.go` — CRUD + `/filter` search, Swagger annotations, standard response helpers.
7. `internal/routes/<entity>_routes.go` — register on the Gin router group with auth/role middleware.
8. Wire repository → validator → service → handler → routes in `internal/server/server.go` (or equivalent init).
9. Add the new struct to the `AutoMigrate` chain.
10. Seed default data if applicable.
11. Register `view_/create_/edit_/delete_<entity>` permissions and map them to the `ADMIN` role in seeders.
12. Regenerate Swagger (command above).

Filenames are snake_case (`customer_supplier_repository.go`), and controllers must never bind directly to repositories — always Service → UoW.

### Legacy SQL Server Integration (`api.integrations`, planned — port 4007)

Not yet implemented (design discussion in `Documentos/Planejamento/Modulo_1_CRM_Integracoes/`). It is a **bidirectional sync middleware** between the **Loja Integrada** e-commerce platform and the client's legacy SQL Server ERP (`FOCCO_ERP`), which this project is gradually replacing:
- **Loja Integrada → SQL Server**: a webhook receiver imports web orders as sales notes (`tipentsai=50`) into the legacy DB.
- **SQL Server → Loja Integrada**: polls the `LogAltera` change-log table and pushes product/price/stock updates to the LI REST API.
- Owns its own PostgreSQL persistence for sync state: `SyncLog`, `ProductMapping`, `WebhookEvent`, `IntegrationConfig` (API keys, `codtipnot`, `codlocarm_oficial`/`codlocarm_reserva`, `codemp` — cached in memory, not re-read per request).

Because it talks to the legacy DB, it breaks from the api.core/api.auth pattern: no `UnitOfWork`/`AutoMigrate` against SQL Server tables (schema is owned by the legacy system), and writes there (e.g. incrementing `CADSEQ` sequences) must run inside an explicit SQL Server transaction with `UPDLOCK` to stay atomic, mirroring the legacy `GetSequencia` logic (select → increment → update).

Legacy schema glossary (FOCCO_ERP / SQL Server) relevant to this integration:
- `LogAltera` — change log used for polling; `TipoAlt` (`P`=Produto, `V`=Preço Venda, `U`=Custo), `DatFimAlt` (polling cursor), `codRegAlt` (=`codpro`).
- `NOTAS` / `NOTAS1` / `NOTAS3` — order header / line items / payment plan; composite key `(tipentsai, nroentsai)`.
- `CADSEQ` — sequence generator table (`GetSequencia` pattern).
- `CADPRO.codint` — stores the product's Loja Integrada ID (the cross-system mapping field).
- `LOCARM` / `CENCUS` — stock locations and companies. Stock flow on a web sale: **Oficial → Reserva** (moved on web sale close) → **sent to customer** (on financeiro invoicing/faturamento).

Open items still pending before full implementation: Loja Integrada API/App keys, the `codtipnot` value to use for LI-originated orders, and the stock-priority policy between the physical and web store when both sell the last unit concurrently (pending confirmation with the client's board).

## Frontend Architecture

Next.js 15 App Router, TypeScript, Tailwind, `rizzui`/Radix UI components, path alias `@/*` → `Frontend/src/*`.

> For the upcoming CRM module (`Documentos/Planejamento/Modulo_1_CRM_Integracoes/`), the agreed approach is **hybrid UI**: Ant Design for complex tables/forms (e.g. Kanban pipelines, filterable grids), Tailwind/rizzui for layout — this is not a full migration off the current stack, just an addition scoped to that module.

- **`app/`** — route groups `(auth)` (login/signin) and `(system)` (everything behind auth: dashboards, customers, suppliers, products, inventory, financial, sales, purchases, settings, ...). Each domain generally has list/`create`/`[id]`/`[id]/edit` routes.
- **`src/services/`** — one file per domain (`customer-service.ts`, `auth/user-service.ts`, ...), each wrapping the shared `axios` instance in `src/services/common/api.ts`. Zod schemas for form validation live alongside the service (e.g. `customerFormSchema`), and TS interfaces mirror the backend DTOs (`Customer`, `CustomerDetail`, `CustomerList` + `Pagination`).
- **`src/services/common/api.ts`** — the single axios client (`withCredentials: true`, base URL from `NEXT_PUBLIC_API_URL`). Its response interceptor auto-retries once via `/auth/refresh-token` on 401 (skipping the login call itself), then hard-redirects to `/login` on refresh failure.
- **`src/contexts/`** / **`src/atoms/`** — `auth-context.tsx` for session/user state, `jotai` atoms (`userAtom.ts`) for cross-component state.
- **`src/config/`** — `routes.ts` (route map), `site.config.tsx`, `color-presets.ts` (layout theming).
- **`src/components/`** — organized by domain (`customers/`, `settings/roles/`, `user/`, `inventory/`, `finance/`), plus `common/` (shared table, etc.), `ui/` (design-system primitives), `layout/`.
- All API calls go through the gateway (`http://localhost:4000/api/...`), never directly to `api.auth`/`api.core` ports.

### Sidebar menu is data-driven, not hardcoded

The sidebar (`src/components/layout/sidebar-menu.tsx`, recursive, arbitrary depth) renders from `user.menuItems` (`AuthContext`/`UserContext`), which comes from the login/refresh/me response body — **not** from a static config file. Source of truth is the `menu_items` table in `api.auth` (self-referencing `MenuItem` domain, `internal/domain/menu_item.go`), managed via `/menu-items` (dev-only, gated by `admin.create_permissions`, no admin UI yet — CRUD only). `AuthService.GetMenuItemsForUser` builds the tree and filters it per user (ADMIN sees everything; a node with an `href` needs its own linked permission or an ADMIN bypass; a pure group node with no `href` is visible only if it has a visible descendant — never grant it access on its own, since it has no permission of its own to check). Icons are Iconify strings (`ph:` prefix = Phosphor, matches the current visual set) resolved via `@iconify/react`, not JSX. To change the menu tree, edit rows in `menu_items` (or, for the seed defaults, `internal/seeder/menu_item_seeder.go`) — do not reintroduce a static `navigation.ts`-style config.

### Mandatory frontend CRUD pattern

Entity Cadastro (create) and Edição (edit) are **full pages**, not modals — `app/(system)/<domain>/create/page.tsx` and `app/(system)/<domain>/[id]/edit/page.tsx`, each rendering a form component from `src/components/<domain>/forms/`. Validation is `react-hook-form` + `zodResolver`, against a Zod schema colocated with the entity's types (`<entity>-schema.ts` or the service file — see `src/services/auth/role-schema.ts` / `src/services/auth/user-schema.ts`), reusing the shared `FormField` wrapper (`src/components/common/form-field.tsx`) for label+error display. On submit, `router.push` back to the list route; `useToast` for success/error feedback.

Modals are reserved for **confirming destructive/user actions** (`ConfirmDeleteDialog`, `src/components/common/confirm-delete-dialog.tsx`) — never for data entry forms.

Non-Zod-managed sub-sections (e.g. a permission picker tied to a record but not itself a form field) live as sibling state in the page/form component, not inside the RHF schema — see `PermissionsPicker` usage in `create-user-form.tsx`/`edit-user-form.tsx` for the pattern.

As of 2026-07-18 this pattern is fully implemented only for Usuários (`settings/users`) and Perfis (`settings/roles`) — most other domains still have `em-construcao` stub routes with no real CRUD yet (some without even a service layer). Apply this pattern when building each of those out; do not assume an existing modal needs converting without checking first.

## Development Workflow / Planning

`Documentos/Planejamento/` is the authoritative task tracker for this project and applies to **all agents and developers**, not just this session:
- `Documentos/Planejamento/README.md` explains the module breakdown (`Modulo_0_Configuracao` ... `Modulo_7_Contabilidade`).
- Before starting a feature, read the module's `plano_*.md` (scope, DB, business rules, affected files) and `tasks_*.md` (current status).
- Concurrent task locking in `tasks_*.md`: pick a free task (`- [ ]`), mark it `- [/] [EM EXECUÇÃO POR: NOME_DO_AGENTE]` and save *before* touching code; never start a task already tagged `[EM EXECUÇÃO POR: ...]` or `[BLOQUEADO POR: ...]`. On success mark `- [x] [CONCLUÍDO POR: NOME_DO_AGENTE]`; on interruption, revert to `- [ ]`.
- At the end of a working session, add/update `Documentos/Planejamento/Historico/log_YYYY-MM-DD.md` with what was done and what's pending, and update the relevant `tasks_*.md`.

Full skill definitions for the above (loaded automatically by Claude Code from `.agents/skills/`): `lumini_hub_backend_architecture`, `lumini_hub_dev_flow`, `lumini_hub_entity_creation`.
