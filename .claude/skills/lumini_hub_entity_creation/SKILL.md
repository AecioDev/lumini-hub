---
name: lumini-hub-entity-creation
description: Step-by-step guidelines for creating and registering new domain entities, repositories, services, and routes in the Lumini Hub backend.
---

# Lumini Hub Entity Creation Guidelines

This skill defines the step-by-step process for creating, writing, and registering a new database entity (such as Brands, Products, etc.) inside the microservices of the Lumini Hub Go backend.

Always perform these steps in order when the user asks to add or create a new entity.

---

## Overview of Microservice Context
All file paths below are relative to the specific microservice folder where the entity belongs:
*   **Core Entities (ERP/Admin/Registration)**: Locate in `Backend/microservices/api.core/`
*   **Auth/Role Entities**: Locate in `Backend/microservices/api.auth/`

---

## Step 0: Company Scoping — Ask Before Modeling

Lumini Hub is multi-company **within a single tenant database** (multi-tenant itself is one physical database per client, handled outside this skill). A holding company can have subsidiary/branch companies self-referencing under it (`Company.ParentID`, same shape as `MenuItem.ParentID`). Almost every entity has to take a position on how it relates to `Company`, and getting this wrong is expensive to retrofit once rows exist — **never assume, ask the user which of these three buckets the new entity falls into before writing Step 1's model**:

1. **Hard-scoped** — the entity has its own `CompanyID uint` FK and every row belongs to exactly one Company. Default assumption for anything operational/transactional that must never mix between companies: stock, sales, purchase orders, financial entries, invoices, cash registers, etc. Repositories/services for hard-scoped entities must filter by the requester's Company (see visibility rule below) in every query — this is as mandatory as the Unit of Work pattern, not optional hardening.
2. **Global with a visibility relationship** — no `CompanyID` on the entity itself (the catalog stays shared/global), but a weak many-to-many table (`<entity>_companies`) controls which Companies can see/use each row. Reference case: `Product` — one shared catalog, but a `product_companies` join table controls which products a given Company's sellers can see/sell.
3. **Global, unscoped** — no Company relationship at all. Reserved for genuinely cross-company reference data (e.g., the Permission catalog, generic lookup tables). Justify explicitly why an entity belongs here instead of bucket 1 or 2 — it should be the exception, not the default.

**Visibility rule to apply once `Company` exists** (confirmed 2026-07-21): a user has a default `CompanyID` on their own record. They only see data scoped to that Company, UNLESS they hold a permission that grants visibility into subsidiary Companies below their own in the self-referencing hierarchy (or, if the user has no Company linked, that same permission decides whether they see everything) — this is a permission-gated exception, not automatic just from being higher in the tree.

If `Company` itself doesn't exist yet in a given codebase snapshot, still ask the scoping question and record the answer in the feature's `plano_*.md` so the `CompanyID` FK/join table is added in the same pass once `Company` lands, instead of being silently forgotten.

## Step 0.5: 1:1 Config Entities Don't Get Delete

Found 2026-09-05 building `CompanyFiscalConfig`: a soft-deleted row plus a plain unique index (e.g. `CompanyID uniqueIndex`) lets the deleted row keep occupying the unique value, blocking recreation later (`ExistsByX` on the app side correctly ignores soft-deleted rows via GORM's default scope, but the raw DB unique *index* doesn't unless it's a partial index with `where:deleted_at IS NULL`). The real fix isn't a smarter index (though that's cheap insurance, add it anyway) — it's recognizing that **a 1:1 configuration record intrinsic to a parent (e.g. `CompanyFiscalConfig`, `CompanyVisualConfig`) has no business meaning for "delete"**: as long as the parent (`Company`) exists, it always has (or should have) that configuration, even if empty. Deleting the row doesn't correspond to anything a user would actually want — "clear the certificate" or "reset the colors" are `Update`s with empty/default values, or a dedicated action, never a row deletion.

**Before Step 6 (Handler), decide**: is this entity a real list (rows genuinely come and go — `Company`, `ChartOfAccounts`, `CompanyDocumentIssuanceConfig`) or a 1:1 config glued to a parent? For the latter, skip Delete entirely: no `Delete<Entity>` in the service, no handler method, no route, and no `<entity>.delete` permission in Step 11 — only `view`/`create`/`edit`. Don't build it "for completeness" and remove it later; ask this before Step 1, same spirit as Step 0.

---

## Step-by-Step Entity Creation Process

### Step 1: Base Domain Model
Create the base model structure in `internal/models/<entity>.go`.
*   Define the struct representing the database table with GORM tags (e.g., `gorm:"primaryKey"`, `gorm:"type:varchar(100)"`, `gorm:"uniqueIndex"`).
*   Create request payloads for creation and updating (e.g., `Create<Entity>Request`, `Update<Entity>Request`) with validation tags (e.g., `binding:"required,min=2"`).

### Step 2: Model DTOs & Mappers
Create the DTO structure in `internal/models/<entity>_dto.go`.
*   Include DTOs optimized for lists and details (e.g., `<Entity>DTO`, `<Entity>DetailDTO`).
*   Implement converter methods/functions to safely map entities to DTOs (e.g., `ToDTO()` and `ToDetailDTO()`).

### Step 3: Specific Repository Interface and GORM Implementation
Create the repository in `internal/repository/<entity>_repository.go`.
*   Define a specific interface that extends `commonrepo.Repository[domain.<Entity>]`.
*   Provide a GORM struct implementation.
*   If the entity requires preloaded relationships (e.g., `Addresses` or `Contacts`), override `FindByID` or write a custom search method to apply `preload`.

### Step 4: Input and Business Validator
Create the validator in `internal/validator/<entity>_validator.go`.
*   Implement validations to enforce business rules before persistency (e.g., uniqueness of a code or document, validating nested items).

### Step 5: Service Layer (Business Logic)
Create the service in `internal/service/<entity>_service.go`.
*   Inject the unit of work interface `repository.UnitOfWork` and the validator.
*   **Strict rule**: Use the Unit of Work (`s.uow.Execute(func(uow repository.UnitOfWork) error { ... })`) to wrap all repository insert, update, or delete operations. This guarantees transactions are committed atomically.

### Step 6: Handler Layer (Controller)
Create the handler in `internal/handlers/<entity>.go`.
*   Inject the Service.
*   Implement endpoints for Create, Update, Delete, GetByID, and Search (`/filter` POST method).
*   Always use `utils.SuccessResponse`, `utils.ErrorResponse`, or `utils.ValidationErrorResponse` for returning payloads.
*   Document every endpoint using standard declarative Swagger annotations (`@Summary`, `@Tags`, `@Router`, etc.).

### Step 7: Routes Definition
Create the routes file in `internal/routes/<entity>_routes.go`.
*   Register endpoints to the Gin router group, applying necessary Auth or Role-based middleware.

### Step 8: Server Integration
Wire the repository, validator, service, handler and route registration together in the microservice's own startup code (`main.go` today — neither `api.core` nor `api.auth` has a separate `internal/server/` package yet, despite `CLAUDE.md` naming one as the convention; check the current microservice before assuming the file exists).

### Step 9: Automatic DB Migrations
Add the new struct to the `AutoMigrate(...)` call in the microservice's `main.go` (e.g. `api.core/main.go` for the `Company` entity). There is no dedicated `migrate.go` file in either `api.core` or `api.auth` today.

### Step 10: Database Seeder
Add default seed data (permissions, menu items, system config) via the microservice's `internal/seeder/` package if it has one (`api.auth` does, e.g. `menu_item_seeder.go`; `api.core` doesn't yet — follow the same package-per-microservice pattern if you're the first entity there needing seed data).

### Step 11: Security & RBAC Permissions
Register the new entity permissions (e.g., `view_<entities>`, `create_<entities>`, `edit_<entities>`, `delete_<entities>`) in the permissions table via seeders or initial SQL scripts. Map them to default user roles (Admin, Manager).
Whenever you add standard permissions, make sure to also write seed logic to insert these permissions and link them directly to the `ADMIN` role.

### Step 12: Unit Tests (Recommended)
Create unit/integration tests in `internal/service/<entity>_service_test.go` or `internal/validator/<entity>_validator_test.go` to test logic isolation.

---

## Important Rules for AI Code Generation
1.  **Do not skip layers**: Never bind controllers directly to repositories. Always implement the Service layer and route operations via the Unit of Work.
2.  **No direct GORM DB in Services**: Services must only interact with `repository.UnitOfWork`.
3.  **Strict Path Naming**: Align snake_case filenames (e.g., `customer_supplier_repository.go`) to match existing file name conventions.
4.  **Swagger Regeneration**: Whenever you modify or add handlers, regenerate swagger documentation in the gateway directory:
    ```powershell
    swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
    ```
5.  **Automatic Permission Seeding**: When generating files for a new entity, you **MUST** automatically write/propose the GORM seeder logic (e.g. in `seeder.go`) or the equivalent database migration scripts to seed the new permissions (`view_<entity>`, `create_<entity>`, `edit_<entity>`, `delete_<entity>`) and map them to the `ADMIN` role template immediately.
6.  **Ask the Company scoping question before Step 1, every time**: hard-scoped (own `CompanyID`), global-with-visibility-relationship (join table), or global-unscoped — see Step 0 above. Do not default to "unscoped" just because `Company` doesn't exist in the codebase yet; record the intended answer in `plano_*.md` regardless.
7.  **1:1 config entities skip Delete entirely** (service, handler, route, and `.delete` permission) — see Step 0.5. Only real list-shaped entities keep the full CRUD.
