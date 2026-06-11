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
Create the handler in `internal/api/handlers/<entity>.go`.
*   Inject the Service.
*   Implement endpoints for Create, Update, Delete, GetByID, and Search (`/filter` POST method).
*   Always use `utils.SuccessResponse`, `utils.ErrorResponse`, or `utils.ValidationErrorResponse` for returning payloads.
*   Document every endpoint using standard declarative Swagger annotations (`@Summary`, `@Tags`, `@Router`, etc.).

### Step 7: Routes Definition
Create the routes file in `internal/api/routes/<entity>_routes.go`.
*   Register endpoints to the Gin router group, applying necessary Auth or Role-based middleware.

### Step 8: Server Integration
Update the router setup in `internal/api/server/server.go`.
*   Instantiate the repository, validator, service, handler, and call the route registration function inside the server initialization.

### Step 9: Automatic DB Migrations
Update the database migration process in `internal/models/migrate.go` (or `main.go` database initialization) to include the new domain struct in the `AutoMigrate` chain.

### Step 10: Database Seeder
Add default seed data or system configurations (like system status) in `internal/models/seeder.go` (or equivalent seeder script).

### Step 11: Security & RBAC Permissions
Register the new entity permissions (e.g., `view_<entities>`, `create_<entities>`, `edit_<entities>`, `delete_<entities>`) in the permissions table via seeders or initial SQL scripts. Map them to default user roles (Admin, Manager).

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
