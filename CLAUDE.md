# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Lumini Hub is an ERP system with a Go microservices backend (`Backend/`) and a Vite + React + Ant Design frontend (`Frontend/`). It is a **brand-new system built from scratch**, not a migration or rewrite of an existing one. It is organized around ERP modules (CRM, Estoque/Produtos, Compras, Vendas/Caixas, Financeiro, Fiscal, Contabilidade) tracked in `Documentos/Planejamento/`. The legacy SQL Server ERP of one client only appears in the Integrations module (`api.integrations`), which is paused and will move to a separate application; it is not part of how this system is designed.

## Commands

### Backend (Go, from `Backend/`)
```powershell
.\run_services.bat        # starts api.gateway (4000), api.auth (4001), api.core (4002) in separate windows
go build ./...             # build all workspace modules
go test ./...               # run tests; the ones that need Postgres skip unless TEST_DATABASE_URL is set (see "Database schema, seed and tests")
```
Requires a local PostgreSQL on port 5432 and a `Backend/.env` (see `Backend/.env.example`; it also documents `DB_AUTO_MIGRATE`, `BOOTSTRAP_ADMIN_*` and `TEST_DATABASE_URL`).

Each service is started with `go run main.go` (that is what `run_services.bat` does), which compiles **only `main.go`**: anything `main.go` uses must live in it or in an imported package, never in another file of `package main` (it passes `go build .` but breaks `go run main.go` with `undefined`). `api.auth` has a test (`TestMainCompilaSozinho`) guarding this.

Regenerate Swagger docs after adding/changing handlers (run from `Backend/microservices/api.gateway/`):
```powershell
swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
```

### Frontend (Vite + React, from `Frontend/`)
```powershell
pnpm install
pnpm dev      # Vite dev server at http://localhost:3000 (port pinned in vite.config.ts to match gateway CORS)
pnpm build    # tsc -b && vite build
pnpm lint     # eslint .
```
No frontend test runner is currently configured. Requires `Frontend/.env.local` with `VITE_API_URL` (defaults to `http://localhost:4000/api`, i.e. the gateway).

Local login: `admin` / `987321`, valid only for the existing development database. A brand-new empty database has no users until `BOOTSTRAP_ADMIN_*` is set (see "Database schema, seed and tests").

## Backend Architecture

Go Workspace (`Backend/go.work`) with five modules:
- **`common/`** — shared package: `config` (env loading), `database` (GORM/Postgres pool), `middlewares` (`AuthMiddleware` JWT-from-cookie, `RequirePermission` RBAC), `utils` (JWT, bcrypt, `response.go`, pagination), `repository` (generic `Repository[T]` / `GormRepository[T]`).
- **`microservices/api.gateway`** (port 4000) — sole entrypoint. Reverse-proxies by path prefix to the other services and is the **only** place CORS is configured.
- **`microservices/api.auth`** (port 4001) — users, roles, permissions, menu items (dynamic sidebar tree), login/refresh (HTTP-only cookies).
- **`microservices/api.core`** (port 4002) — customers, suppliers, addresses, contacts, documents, companies (multi-company, see `Company` below).
- **`microservices/api.integrations`** (port 4007) — bidirectional sync middleware with the legacy SQL Server ERP and Loja Integrada, see [dedicated section](#legacy-sql-server-integration-apiintegrations-port-4007) below.

```mermaid
graph TD
    Client[Frontend - Vite/React/antd<br>Porta 3000] -->|HTTP/API| Gateway[api.gateway<br>Porta 4000]
    Gateway -->|/api/auth, /api/users, /api/roles, /api/permissions, /api/menu-items| Auth[api.auth<br>Porta 4001]
    Gateway -->|/api/customers, /api/suppliers, /api/companies| Core[api.core<br>Porta 4002]
    Gateway -->|/api/integrations| Integrations[api.integrations<br>Porta 4007]
    Auth -.->|imports| Common[common]
    Core -.->|imports| Common
    Integrations -.->|imports| Common
    Auth --> PG[(PostgreSQL<br>Porta 5432)]
    Core --> PG
    Integrations --> PG
    Integrations -->|read/write, no UoW/AutoMigrate| Legacy[(SQL Server<br>FOCCO_ERP)]
    Integrations <-->|webhook + REST| LI[Loja Integrada]
```

`api.auth`, `api.core` and `api.integrations` share one physical Postgres database for now but never join across each other's tables — cross-service data (e.g. `User` inside a Customer DTO) is passed as plain IDs and re-hydrated into a local simplified struct (e.g. `ApiUser`), not a GORM relation.

### Mandatory backend patterns
- **Responses**: every handler returns via `utils.SuccessResponse` / `utils.ErrorResponse` / `utils.ValidationErrorResponse` ([response.go](Backend/common/utils/response.go)). Never hand-roll `gin.H` payloads.
- **Repository**: microservice-specific repositories extend the generic `commonrepo.Repository[T]` and override methods like `FindByID` when GORM `Preload` is needed.
- **Unit of Work**: services never touch a repository or `*gorm.DB` directly — they call `s.uow.Execute(func(uow repository.UnitOfWork) error { ... })` so writes are atomic. Each microservice defines its own `UnitOfWork` interface (see `api.core/internal/repository/repository.go`) exposing one method per aggregate (`Customers()`, `Suppliers()`, ...). Occasional, documented exceptions are acceptable (e.g. `CompanyService.GetCompanies` hands `GetDB()` to the shared read-only `utils.ResolveVisibleCompanyIDs`), but only for reads, with a comment saying why; they must never become the norm, and any write always goes through `s.uow.Execute(...)`.
- **Search**: complex/paginated queries are POST endpoints ending in `/filter` (a `<Entity>FilterRequest` DTO with `page_no`, `page_size`, `order_by_column`, `is_asc`), not query-string GETs. Results go through `utils.Paginate`.
- **CORS**: only ever configured in `api.gateway/main.go`. Adding it to `api.auth` or `api.core` causes duplicate-header `AxiosError: Network Error` in the browser.
- **RBAC (hybrid)**: a `Role` is only a template/default — creating or editing a user picks a Role, which pre-fills a default set of permissions, but what `RequirePermission` middleware actually checks is the permission list embedded in the JWT, which comes from the user's own direct `user_permissions` link (`api.auth/internal/domain/user.go`), not from the Role. Operators can freely add/revoke individual permissions on a user beyond their Role's defaults (`PUT /users/:id/permissions`). The Permissions catalog (`/settings/permissions`) is a developer-only screen (`permissions.*` + `admin.create_permissions`) — creating a permission there has no effect unless a matching `RequirePermission("...")` already exists in code. The Roles catalog (`/settings/roles`, shown as "Perfis e Permissões" in the sidebar, gated by `roles.view`) belongs to ADMIN.
- **RBAC bypass is two-tiered** (`common/utils/rbac.go`, 2026-07-19): `DEVELOP` bypasses every `RequirePermission` check and sees the full menu, no exceptions. `ADMIN` bypasses everything **except** *writing* to the Perfis/Permissões catalog itself (`utils.IsDeveloperOnlyPermission`: `admin.create_permissions` + any `permissions.*` code **except** `permissions.view`, which stays open to ADMIN because the Users create/edit screens need it to populate the permission picker) — creating/editing/deleting **Permissions** (`/settings/permissions`) is DEVELOP-only. The **Roles** catalog (`roles.*`) belongs to ADMIN: it is *not* developer-only in `rbac.go`, and its sidebar entry ("Perfis e Permissões", `/settings/roles`) is seeded with `roles.view` (corrected 2026-10-01: this text used to list `roles.*` as DEVELOP-only, which stopped being true on 2026-07-21, and the seeder still tied that menu entry to `admin.create_permissions`, so a freshly seeded database hid it from ADMIN). Reasoning: creating a permission with no matching `RequirePermission("...")` in code does nothing, so this is a developer tool, not a business-admin one. Both the backend (`middlewares.RequirePermission`, `domain.FilterMenuTreeForUser`) and the frontend (`AuthContext.tsx`'s `hasPermission`) implement this same tiering — keep them in sync if it changes. The seed (see "Database schema, seed and tests") creates only the `ADMIN` role and never the `DEVELOP` role or its grants: `DEVELOP` is created and populated manually in the DB (`user_permissions`) per-environment, by the dev team, not via code. The seed also does not link the 9 `Develop`-module permissions to `ADMIN` (it doesn't need them: ADMIN passes through the bypass).
- **The `DEVELOP` role and its users are invisible to everyone except DEVELOP itself** (`user_service.go`/`role_service.go`, 2026-07-20): `GET /users`/`GET /users/:id` never return a user whose Role is `DEVELOP` unless the requester's own role is `DEVELOP` (filtered at the SQL level for the list, via a `role_id NOT IN (SELECT id FROM roles WHERE name = 'DEVELOP')` subquery, so pagination totals stay correct; a direct `GET /users/:id` on a hidden user returns 404, not 403 — existence isn't revealed). Same treatment for `GET /roles`/`GET /roles/:id` and the Role's own CRUD (`UpdateRole`/`DeleteRole`/`UpdateRolePermissions` all 404 on the DEVELOP role for non-DEVELOP callers). `CreateUser`/`UpdateUser` also reject assigning the DEVELOP role's ID to a user unless the actor is DEVELOP (generic "perfil inválido" error — doesn't leak the role's name). Requester's role comes from `utils.RoleFromGinContext(c)`, threaded through the service layer as a plain string param — no repository/service method here takes `isAdmin bool` anymore, they take `requesterRole string`.
- **Routing**: new route prefixes must be registered both in the microservice's own routes file and proxied from `api.gateway/main.go` (via `router.NoRoute`, careful not to collide with `/swagger/*any`).

### Adding a new entity (Core or Auth microservice)
Follow this order (see `.claude/skills/lumini_hub_entity_creation/SKILL.md` for full detail):
1. `internal/models` (or `domain`) — GORM struct + `Create<Entity>Request`/`Update<Entity>Request` with `binding` validation tags.
2. `internal/models/<entity>_dto.go` — DTOs + `ToDTO()`/`ToDetailDTO()` mappers.
3. `internal/repository/<entity>_repository.go` — interface extending `commonrepo.Repository[domain.<Entity>]` + GORM impl, with `Preload` overrides as needed.
4. `internal/validator/<entity>_validator.go` — business-rule validation (uniqueness, nested items).
5. `internal/service/<entity>_service.go` — business logic, all writes wrapped in `s.uow.Execute(...)`.
6. `internal/handlers/<entity>.go` — CRUD + `/filter` search, Swagger annotations, standard response helpers.
7. `internal/routes/<entity>_routes.go` — register on the Gin router group with auth/role middleware.
8. Wire repository → validator → service → handler → routes in `internal/server/server.go` (or equivalent init).
9. Add the new struct to the service's `AutoMigrate` chain (`api.auth/main.go`'s `prepareDatabase`, or `api.core/main.go`), in FK dependency order. It only runs when `DB_AUTO_MIGRATE` is on.
10. Seed default data if applicable: idempotent, and under the same `DB_AUTO_MIGRATE` flag (auth seeds live in `api.auth/seed`).
11. Register the entity's permissions (`<entity>.view/create/edit/delete`) in `PermissionCatalog` (`api.auth/seed/permission_catalog.go`). The next boot creates them and links them to `ADMIN` automatically (everything outside the `Develop` module); links that already exist are never touched. Every permission a route requires with `RequirePermission("...")` must be in the catalog, otherwise only ADMIN/DEVELOP can ever use that route (since 2026-10-20, CFG-10.1, the catalog covers every `RequirePermission` in the code — keep it that way when adding routes).
12. Regenerate Swagger (command above).

Filenames are snake_case (`customer_supplier_repository.go`), and controllers must never bind directly to repositories — always Service → UoW.

### Database schema, seed and tests (since 2026-10-01)

A new empty Postgres boots by itself: the services create their tables and `api.auth` seeds the base data. Planned in `Documentos/Planejamento/Modulo_0_Configuracao/plano_schema_seed.md` (EPIC CFG-9).

- **`DB_AUTO_MIGRATE`** (`common/config`, default **true**; an invalid value fails `config.Load()` instead of silently falling back). `true`: each service runs `AutoMigrate` for its own tables (`api.auth`: Permission, Role, User, MenuItem; `api.core`: Company, CompanyFiscalConfig, Country, State, City, Customer, Supplier, Address, Document, Contact; `api.integrations`: its 4 Postgres tables — never the legacy SQL Server structs). `false` (production with real data): no service runs `AutoMigrate`, and `api.auth` also skips all of its seeding (`prepareDatabase` returns immediately), so the operator owns schema and data. Then the permission catalog and the menu are **not** applied by themselves: apply them by other means (e.g. run once with the flag on, or the future `api.auth` binary command). Not covered by the flag: `api.integrations` still runs `ConfigService.SeedDefaults()` (default keys in `integration_configs`), which needs that table to exist.
- **Boot of `api.auth`** (`prepareDatabase`, `api.auth/main.go`), all under that flag and in this order: `AutoMigrate` → `seed.SyncCatalog` → `seed.BootstrapAdmin` → `seeder.SeedMenuItems` (the menu references permissions by code, so the catalog comes first) → `dbversion.Record` (see below). With the flag off, `prepareDatabase` only runs `dbversion.Check`.
- **Database version** (`common/dbversion`, table `schema_versions`, CFG-020): the database records which version it is at. `dbversion.Required` (the database version this code needs, the 2nd number of the system version) and `dbversion.AppVersion` (the full system version) are two constants edited **by hand at every release**: `AppVersion` always, `Required` only when the release has database changes (+1; patches and releases without database changes keep it). `api.auth` calls `dbversion.Record` at the end of a successful boot with `DB_AUTO_MIGRATE=true` (one row per release, restarting never duplicates). With `DB_AUTO_MIGRATE=false` every service (`api.auth`, `api.core`, `api.integrations`) calls `dbversion.Check` and **refuses to start** if the database is behind `Required` or has no version recorded; a database ahead of the code is not an error. The operator brings the database up to date by running `api.auth` once with the flag on (or applying the update another way). `AutoMigrate` only adds, so changes that rewrite existing data need their own script.
- **`seed.SyncCatalog`** (public package `api.auth/seed`, importable as `lumini-hub/api.auth/seed`): `PermissionCatalog` is the single in-code source of the permissions (transcribed from the dev database on 2026-10-01, plus the 8 permissions the code required and dev never had, added 2026-10-20; naming cleanup of the inherited entries is EPIC CFG-10, CFG-10.2). Idempotent. First boot (no role at all): creates the `ADMIN` role (the only role the seed creates) and links it to every permission whose `Module != "Develop"` (58 of 67). Later boots: creates only the *new* catalog permissions and links the new non-Develop ones to ADMIN; it never re-adds a link the operator removed, never touches or deletes existing permissions (including hand-made ones and soft-deleted ones), and never creates another role.
- **`seed.BootstrapAdmin`** (transitional — in the multi-tenant phase a commercial provisioning app will create the first admin): from `BOOTSTRAP_ADMIN_USERNAME` (default `admin`), `BOOTSTRAP_ADMIN_PASSWORD` (**no default**), `BOOTSTRAP_ADMIN_EMAIL` (empty becomes `<username>@lumini.local`, because `users.email` is unique). Only runs if the `users` table is empty (soft-deleted rows count), so a restart never duplicates the user or changes its password. Without a password it creates nothing and just warns. Password: at least 8 characters and at most 72 bytes (bcrypt limit); anything else is a configuration error that stops the boot. The user is `ADMIN`, active, with no `company_id` (master: sees every company, which is why a database with no company lands on the company-registration step of the frontend). The password never reaches logs or error messages, and `config.Config`/`AdminBootstrap` print it as `<omitida>`. Remove the variable from the environment after the first boot.
- **No geography seed**: countries, states and cities are meant to be created on demand when a CEP is looked up in the Customer/Supplier screens (future feature). Note `Document.StateID` is a plain `uint` with an FK to `state` (its zero value violates the constraint, and the `state` table starts empty; planned: make it an optional `*uint` when that module is built), and `person_type` (today `text`, unvalidated) should become a validated `F`/`J`.
- **Tests that need Postgres** use `testutil.OpenTestDB` (`common/testutil`): set `TEST_DATABASE_URL=postgres://user:pass@host:5432/lumini_test?sslmode=disable` (see `Backend/.env.example`). Without it they **skip**. The helper refuses any database whose real name (`current_database()`) does not end in `_test`, requires the URL format, and gives each test a throwaway schema (`t_<hex>`), so packages run in parallel without clashing and nothing else in the database is touched. Never point it at the development database. `PermissionCatalog` is a global, so the `api.auth/seed` tests cannot use `t.Parallel`.

### Legacy SQL Server Integration (`api.integrations`, port 4007)

> **Status: PAUSED since 2026-10-01.** Decided by the product owner: access to the legacy SQL Server will be a **separate application** running on the client's legacy server (the cloud cannot reach the client's network), so the Phase 2 backlog (EPIC INT-1) is frozen here and may be migrated back later. Phase 1 stays in the code but is not evolved, and the service is not part of the cloud deployment. Details and the preserved design notes are in `Documentos/Planejamento/Modulo_9_Integracoes/`.

**Phase 1 is implemented** (see `Documentos/Planejamento/Modulo_9_Integracoes/tasks_integrations.md` for the exact checklist — split out from the former `Modulo_1_CRM_Integracoes` on 2026-09-13, since CRM and Integrations are unrelated microservices that only ended up sharing a module folder by early-planning coincidence): the module skeleton, Postgres domain/repositories/UoW, read-only legacy SQL Server layer, `ConfigService`, handlers/routes (`/settings`, `/sync-logs/filter`, `/webhook-events/filter`, `/legacy/locations`, `/legacy/companies`, `/webhooks/loja-integrada`), the gateway proxy, `run_services.bat` entry, and a frontend Configurações → Integrações screen are all done. It is a **bidirectional sync middleware** between the **Loja Integrada** e-commerce platform and the client's legacy SQL Server ERP (`FOCCO_ERP`):
- **Loja Integrada → SQL Server**: a webhook receiver imports web orders as sales notes (`tipentsai=50`) into the legacy DB. *(Phase 2, not yet implemented — see below.)*
- **SQL Server → Loja Integrada**: polls the `LogAltera` change-log table and pushes product/price/stock updates to the LI REST API. *(Phase 2, not yet implemented.)*
- Owns its own PostgreSQL persistence for sync state: `SyncLog`, `ProductMapping`, `WebhookEvent`, `IntegrationConfig` (API keys, `codtipnot`, `codlocarm_oficial`/`codlocarm_reserva`, `codemp` — cached in memory, not re-read per request).

Because it talks to the legacy DB, it breaks from the api.core/api.auth pattern: no `UnitOfWork`/`AutoMigrate` against SQL Server tables (schema is owned by the legacy system), and writes there (e.g. incrementing `CADSEQ` sequences) must run inside an explicit SQL Server transaction with `UPDLOCK` to stay atomic, mirroring the legacy `GetSequencia` logic (select → increment → update).

Legacy schema glossary (FOCCO_ERP / SQL Server) relevant to this integration:
- `LogAltera` — change log used for polling; `TipoAlt` (`P`=Produto, `V`=Preço Venda, `U`=Custo), `DatFimAlt` (polling cursor), `codRegAlt` (=`codpro`).
- `NOTAS` / `NOTAS1` / `NOTAS3` — order header / line items / payment plan; composite key `(tipentsai, nroentsai)`.
- `CADSEQ` — sequence generator table (`GetSequencia` pattern).
- `CADPRO.codint` — stores the product's Loja Integrada ID (the cross-system mapping field).
- `LOCARM` / `CENCUS` — stock locations and companies. Stock flow on a web sale: **Oficial → Reserva** (moved on web sale close) → **sent to customer** (on financeiro invoicing/faturamento).

**Phase 2 (not yet implemented)**: `NOTAS`/`NOTAS1`/`NOTAS3`/`CADPRO` domain + real sales-note creation from an LI order, the Loja Integrada HTTP client (products/orders), the polling scheduler that actually pushes `LogAltera` changes to LI, and processing a received `WebhookEvent` into a sales order. Blocked on: Loja Integrada API/App keys, the `codtipnot` value to use for LI-originated orders, and the stock-priority policy between the physical and web store when both sell the last unit concurrently (pending confirmation with the client's board).

## Frontend Architecture

Vite + React 18 + TypeScript + React Router v6 + Ant Design v5 (`antd`). The frontend was rebuilt from scratch on 2026-07-19 — the previous Next.js App Router + Tailwind + `rizzui`/Radix stack was fully deleted and replaced, not extended. Path alias `@/*` → `Frontend/src/*`. The dev server is pinned to port 3000 in `vite.config.ts` — required because the gateway's CORS only allows `http://localhost:3000`/`:3001`.

There is no more "hybrid UI" split for the CRM module — `antd` is the sole UI library for the whole app now (no Tailwind, no rizzui/Radix).

- **`src/types/`** — one file per entity, TS interfaces mirroring the backend Go DTOs (`auth.ts`, `user.ts`, `role.ts`, `permission.ts`, `menu.ts`, `customer.ts`, `supplier.ts`), plus `common.ts` for the shared `ApiResponse<T>`/`ApiPagination` envelope. `crm.ts` is explicitly mock-only (see below).
- **`src/services/`** — one file per entity (`auth/auth-service.ts`, `users/user-service.ts`, `roles/role-service.ts`, `permissions/permission-service.ts`, `customers/customer-service.ts`, `suppliers/supplier-service.ts`), each wrapping the shared axios instance. No page/component calls axios directly.
- **`src/services/common/api.ts`** — the single axios client (`withCredentials: true`, base URL from `VITE_API_URL`). Response interceptor auto-retries once via `/auth/refresh-token` on 401 (skipping the login/refresh calls themselves); on refresh failure it hard-redirects to `/login`, **except** for the initial `/auth/me` session-bootstrap call, which is left to fail gracefully through `AuthContext` + `ProtectedRoute`'s own SPA redirect instead of forcing a full page reload.
- **`src/schemas/`** — one Zod schema per form (`login-schema.ts`, `forgot-password-schema.ts`, `create-user-schema.ts`, `update-user-schema.ts`), validation aligned to the backend's `binding` tags.
- **`src/contexts/`** — `AuthContext.tsx` (bootstraps session via `GET /auth/me` on load, exposes `user`, `menuItems`, `login()`, `logout()`, `hasPermission()`) and `ThemeContext.tsx` (dark/light, persisted in `localStorage`, feeds `ConfigProvider`'s `theme.darkAlgorithm`/`defaultAlgorithm` via `src/theme/antd-theme.ts`).
- **`src/routes/`** — `router.tsx` (route tree), `ProtectedRoute.tsx` (session gate), `RequirePermission.tsx` (per-route permission gate — UI convenience only, the backend `RequirePermission` middleware is the real enforcement), `AuthLayout.tsx` (centered card shell for Login/Recuperar senha), `AppLayout.tsx` (`Layout.Sider` 240/76px collapsible + `Layout.Header` + `Layout.Content`).
- **`src/pages/`** / **`src/components/`** — organized by domain (`pages/auth/`, `pages/dashboard/`, `pages/users/`; `components/layout/`, `components/dashboard/`, `components/users/forms/`, `components/common/`). `pages/PlaceholderPage.tsx` is the catch-all fallback for any DB-seeded menu `href` that doesn't have a real screen yet.
- All API calls go through the gateway (`http://localhost:4000/api/...`), never directly to `api.auth`/`api.core` ports.

### Sidebar menu is data-driven, not hardcoded

The sidebar (`src/components/layout/SidebarMenu.tsx`, recursive antd `Menu`, arbitrary depth) renders from `user.menuItems`, which comes from the login/refresh/me response body — **not** from a static config file. Source of truth is the `menu_items` table in `api.auth` (self-referencing `MenuItem` domain, `internal/domain/menu_item.go`), managed via `/menu-items` (dev-only, gated by `admin.create_permissions`, no admin UI yet — CRUD only). `AuthService.GetMenuItemsForUser` builds the tree and filters it per user (ADMIN sees everything; a node with an `href` needs its own linked permission or an ADMIN bypass; a pure group node with no `href` is visible only if it has a visible descendant — never grant it access on its own, since it has no permission of its own to check). Icons are Iconify strings (`ph:` prefix = Phosphor) resolved via `@iconify/react`'s `<Icon>` component, not JSX. **antd `Menu` item `key`s must be built from the item's `id`, not its `href`** — the real seeded tree has duplicate hrefs at different tree positions (e.g. "Clientes" appears both under Vendas > Cadastros and under Financeiro > Contas a Receber, both pointing at `/customers`), which collides if `href` is used as the key. Seeded hrefs are nested paths like `/settings/users`, not flat ones — always check the actual login response's `menu_items` before wiring a new frontend route, don't assume the path. To change the menu tree, edit rows in `menu_items` (or, for the seed defaults, `internal/seeder/menu_item_seeder.go`).

### Mandatory frontend CRUD pattern

Entity Cadastro (create) and Edição (edit) are **full pages**, not modals — e.g. `src/pages/users/CreateUserPage.tsx` / `EditUserPage.tsx`, registered in `src/routes/router.tsx`, each rendering a form component from `src/components/<domain>/forms/` (e.g. `CreateUserForm.tsx`/`EditUserForm.tsx`). Validation is `react-hook-form` + `zodResolver` against a schema in `src/schemas/`, reusing the shared `FormField` wrapper (`src/components/common/FormField.tsx`, thin layer over antd `Form.Item`) for label+error display. Complex forms are split into `antd Tabs` (e.g. "Dados do Usuário" vs "Permissões") rather than one long scroll. On submit, `navigate()` back to the list route; `App.useApp()`'s `message`/`notification` for feedback.

Modals are reserved for **confirming destructive/user actions** — never for data entry forms.

Non-Zod-managed sub-sections (e.g. the Perfil/permissions picker tied to a user but not itself a form field) live as sibling `useState` in the form component, not inside the RHF schema — see `CreateUserForm.tsx`/`EditUserForm.tsx`: the Perfil preset buttons (`RolePresetPicker.tsx`, roles fetched live from `GET /roles`) drive `ModulePermissionsPanel.tsx`, which shows only the permissions of the module whose name matches the selected Perfil, falling back to the first module with a selected permission when the names don't match (the Role/module naming in the seeded catalog isn't fully consistent yet — e.g. "Gerente"/"ADMIN" roles have no same-named module).

This pattern is fully implemented for **Usuários** (`/settings/users`) and **Empresas** (`/settings/companies`, added 2026-07-27 as part of the multi-company foundation — `src/pages/settings/companies/`, `company-service.ts`, `company-schema.ts`). The frontend was rebuilt from scratch on 2026-07-19; Login, Recuperar senha (UI-only, no backend endpoint yet), Dashboard CRM (mocked, no `api.crm` backend yet) and Configurações → Integrações (real, see `api.integrations` above) round out what's real today — every other sidebar destination falls through to `PlaceholderPage`. Apply this pattern when building each of those out.

### Masked fields (CPF/CNPJ/telefone/CEP): mask on screen, digits in the database

Decided 2026-09-05, applies to every field of this kind, past and future (`Company.cnpj`, any future phone/CEP field, etc.): the frontend shows and types with the mask (better UX, and the mask shape itself is free client-side format validation — e.g. rejecting a CNPJ with the wrong digit count as you type), but **only the raw digits are ever sent to the backend and stored in the database** — never persist the punctuation. No masking library is installed yet (checked 2026-09-05); when a field like this is built or touched, add one (or a small local mask helper) rather than storing the masked string as-is. This is a frontend-only concern — backend request/DTO fields for these values are always plain digit strings, no format validation needed there beyond length/checksum if applicable.

### Brazilian document/address fields keep their Brazilian name, not a literal English translation

Decided 2026-09-05: the project's naming convention is English identifiers, but a handful of Brazilian legal/postal document types are so specific and universally referred to by their Portuguese abbreviation (in conversation, in the UI, in every Brazilian dev's head) that translating them literally actively hurts clarity instead of helping — `CNPJ`/`CPF`/`CEP`/`RG` are not "TaxID"/"PersonalID"/"PostalCode"/anything else, they're domain vocabulary. Case in point: `Company.TaxID` was renamed to `Company.CNPJ` (Go field, JSON key `cnpj`, DB column `cnpj`, TS field, Zod schema key) — same rule applies to CPF/CEP/RG wherever they show up next (e.g. a future Address entity's CEP field). Keep the abbreviation itself in Portuguese/uppercase as the identifier at every layer (Go struct field, JSON tag, DB column, TS type) — don't invent an English gloss for it. This does **not** extend to ordinary business nouns that happen to be Portuguese in conversation (`empresa` still becomes `Company`, `usuário` still becomes `User`) — it's specific to these fixed legal/postal document abbreviations.

### Every list/grid shows the record's ID

Decided 2026-09-05: every table/grid listing a saved entity (`CompaniesTable`, `UsersListPage`, etc.) shows the record's `id` as a column — first column, unless there's a good reason not to (the ID doesn't exist as a stable concept for that view, or the user explicitly asks for it to be left out for a specific screen). It's cheap, and it's the fastest way to cross-reference a row with logs, the DB, or a support conversation.

## Development Workflow / Planning

The user drives the work: every task starts from a chat with the user, and nothing runs on its own. There are no subagents, no Kanban board and no automatic task tracking in this project (removed 2026-10-20). Do not spawn agents unless the user asks.

- `Documentos/Planejamento/` holds the scope and history of each module (`Modulo_0_Configuracao` ... `Modulo_9_Integracoes`): `plano_*.md` (scope, DB, business rules, affected files) and `tasks_*.md` (task list, statuses). `Documentos/Planejamento/README.md` explains the layout.
- Read those files only when the user points at a module/feature, or when you need the business rules for the task at hand. Some `tasks_*.md` are very large (`tasks_empresa.md` is ~330 KB): read the relevant section, never the whole file.
- **Tracking lives in two small files in `Documentos/Planejamento/`**: `BACKLOG.md` (what is done and what is left, `[ ]`/`[x]`; every item has an ID `PREFIX-NNN` such as `CFG-014`, the prefix being the module from the table at the top of that file, the number sequential per prefix and never reused; give every new item the next free ID and quote it in `PROXIMA_VERSAO.md`; only the user marks `[x]`, after testing it himself) and `PROXIMA_VERSAO.md` (changes made since the last deploy to the server). When you finish a task the user asked for, add one line to `PROXIMA_VERSAO.md` and show it; if the task is in `BACKLOG.md`, keep it `[ ]` with the note `(implementado, aguardando teste)`. When the user asks for a release, read `PROXIMA_VERSAO.md`, update the constants in `Backend/common/dbversion/dbversion.go` (`AppVersion` always, `Required` only if the release has database changes), write `releases/<version>.md`, then clear `PROXIMA_VERSAO.md`. Version format: `year.database_version.version_sequence.patch_sequence`: the 2nd number is the database version and goes up by one per published release that contains database changes (not per individual change: ten new columns and a table in the same release count as one), and each release must list what changed in the database. Next version: `2026.2.2.0`.
- Do not edit the old `tasks_*.md` or write `Historico/log_*.md` unless the user asks. Their statuses (`[EM-TESTE]`, `[ENTREGUE]`, ...) are historical; treat them as information, not as a workflow to keep feeding.
- **Browser automation (`agent-browser`, the built-in browser, screenshots, exploratory testing) only when the user asks for it**, or when it is clearly the only way to answer the question. Prefer reading code, logs and tests first; it is a large token cost.

Skills (loaded from `.claude/skills/`): `lumini_hub_backend_architecture`, `lumini_hub_entity_creation`, `lumini_hub_dev_flow` (short notes on how to work in this repo) and `agent-browser` (see the rule above).
