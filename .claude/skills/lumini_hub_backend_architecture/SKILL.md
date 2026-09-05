---
name: lumini-hub-backend-architecture
description: >-
  Architectural guidelines and code design standards for evolving and maintaining the Go backend of Lumini Hub.
---

# Lumini Hub Backend Architecture Guidelines

## Overview
This skill defines the mandatory architecture standards, design patterns, and coding guidelines for any AI agent working on the Lumini Hub Go backend. Follow these instructions strictly when creating new microservices, repositories, services, handlers, or routes.

---

## Quick Start
The Lumini Hub backend is structured as a **Go Workspace (`go.work`)** containing:
1.  **`common/`**: Shared module (database connections, JWT, pagination, standard responses).
2.  **`microservices/api.gateway/`** (Port `4000`): Reverse proxy entrypoint and CORS manager.
3.  **`microservices/api.auth/`** (Port `4001`): Users, Roles, Permissions (RBAC) microservice.
4.  **`microservices/api.core/`** (Port `4002`): Customers, Suppliers, Addresses, Contacts, Companies (multi-company foundation) microservice.
5.  **`microservices/api.integrations/`** (Port `4007`): Bidirectional sync middleware with the legacy SQL Server ERP and Loja Integrada. Phase 1 (config, sync-log/webhook-event tracking, legacy read-only lookups, gateway proxy) is implemented; Phase 2 (actual order/stock sync) is not — see `CLAUDE.md`'s dedicated section for the current split.

To run all backend microservices locally in parallel:
```powershell
cd Backend
.\run_services.bat
```

---

## Workflow & Design Patterns

### 1. Unified Response Format (ResponseDTO)
All handlers **MUST** use the standard `utils.Response` struct from [response.go](file:///c:/Projetos/lumini-hub/Backend/common/utils/response.go) to return JSON payloads. Never write ad-hoc map responses (`gin.H`).

#### JSON Payload Format:
- **Success Case (`SuccessResponse`)**:
  ```json
  {
    "success": true,
    "statusCode": 200,
    "message": "Resource found",
    "data": { ... }
  }
  ```
- **Error Case (`ErrorResponse`)**:
  ```json
  {
    "success": false,
    "statusCode": 500,
    "message": "Internal error occurred",
    "error": "exact db or system error details"
  }
  ```
- **Validation Errors (`ValidationErrorResponse`)**:
  Accepts either a single string or a slice of strings (`[]string`).
  ```json
  {
    "success": false,
    "statusCode": 400,
    "message": "Invalid model inputs",
    "validationErrors": [
      "first name is required",
      "invalid document number format"
    ]
  }
  ```

### 2. Base Repository Pattern (Go Generics)
All basic database CRUD logic is generic. 
- Extend `commonrepo.Repository[T]` for basic CRUD methods (`Create`, `Update`, `Delete`, `FindByID`, `FindAll`).
- Specific repository interfaces must compose the generic interface.

#### Example Repository Definition:
```go
type CustomerRepository interface {
	commonrepo.Repository[domain.Customer]
	FindByID(id uint) (*domain.Customer, error) // Override to apply specific GORM Preloads
	FindByDocument(document string) (*domain.Customer, error) // Domain-specific method
}
```

### 3. Unit of Work (UoW) Pattern
Services must never inject raw repositories or direct `*gorm.DB` connections. They must interact with `repository.UnitOfWork` to manage database operations and atomic transactions.

#### Atomic Transaction Example in Service:
```go
func (s *CustomerService) CreateCustomer(req domain.CreateCustomerRequest, userID uint) (*domain.ApiCustomer, error) {
	// ... validation logic ...
	customer := domain.Customer{ ... }

	// Execute inside an atomic transaction block managed by UoW
	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Customers().Create(&customer)
	})
	if err != nil {
		return nil, err
	}
	return &customerDTO, nil
}
```

### 4. Search via HTTP POST (PostFilter Pattern)
To handle complex queries with multiple parameters and pagination without cluttering URL query parameters, searches must be implemented as a **POST** endpoint ending with `/filter`.

- Create a filter request DTO (e.g., `CustomerFilterRequest` in domain).
- Mapped JSON body parameters include fields and paging variables (`page_no`, `page_size`, `order_by_column`, `is_asc`).
- In the repository, build GORM filters dynamically:
  ```go
  if filter.FirstName != "" {
      query = query.Where("first_name ILIKE ?", "%"+filter.FirstName+"%")
  }
  ```
- Page results using `utils.Paginate`.

### 5. Centralized CORS Management
- **Never** register CORS middlewares in the internal microservices (like `api.auth` or `api.core`). Doing so causes duplicate CORS headers at the browser level, triggering `AxiosError: Network Error`.
- CORS is central and managed exclusively in the **API Gateway** ([api.gateway/main.go](file:///c:/Projetos/lumini-hub/Backend/microservices/api.gateway/main.go)).

### 6. Swagger API Documentation & Routing in Gateway
- Document all HTTP methods in handlers using standard declarative Swagger annotations (`@Summary`, `@Tags`, `@Router`, etc.).
- When adding new routes, expose them through the gateway in `api.gateway/main.go` using the `router.NoRoute` proxy handler. Ensure the path does not conflict with `/swagger`.
- Regenerate swagger docs inside the Gateway directory:
  ```powershell
  swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
  ```

---

## Common Mistakes
- **CORS Duplication**: Adding CORS to any internal microservice. Only the Gateway handles CORS.
- **Gin Path Conflicts**: Registering `router.Any("/*path")` alongside `/swagger/*any` in the same Gin router group level. Always route unmapped proxy paths through `router.NoRoute(...)`.
- **Preload Omission**: Forgetting to override the generic `FindByID` or `FindByFilter` in Gorm repositories when entity preloads (e.g., Addresses, Contacts) are required.
