# Guia de Criação de Novas Entidades - Lumini Hub

Este manual descreve o passo a passo detalhado e padronizado para criar, implementar e registrar uma nova entidade de banco de dados (ex: Marcas, Produtos, etc.) no backend do Lumini Hub.

---

## 📌 Contexto Físico dos Microsserviços
Todos os caminhos descritos abaixo são relativos à pasta do microsserviço no qual a entidade está sendo implementada:
*   **Entidades Core (ERP/Cadastros)**: Ficam no diretório `Backend/microservices/api.core/`
*   **Entidades de Autenticação/Controle de Acesso**: Ficam no diretório `Backend/microservices/api.auth/`

---

## 🛠️ Passo a Passo para Desenvolvimento da Entidade

### 1. Criar o Modelo de Domínio Base (Model)
*   **Caminho:** `internal/models/<nome_da_entidade>.go` (ex: `brand.go`)
*   Definir a struct principal da tabela com tags GORM (ex: `gorm:"primaryKey"`, `gorm:"type:varchar(100)"`, `gorm:"uniqueIndex"`).
*   Definir os payloads de requisição de criação e atualização (ex: `CreateBrandRequest`, `UpdateBrandRequest`) usando tags de validação do binding (ex: `binding:"required,min=2"`).

### 2. Criar os DTOs do Modelo (Data Transfer Objects)
*   **Caminho:** `internal/models/<nome_da_entidade>_dto.go` (ex: `brand_dto.go`)
*   Criar DTOs otimizados para retornos em lista ou detalhes (ex: `BrandDTO`, `BrandDetailDTO`).
*   Implementar métodos/funções de conversão seguros para mapear a struct de banco para os DTOs correspondentes (ex: `ToDTO()` e `ToDetailDTO()`).

### 3. Criar o Repositório Específico
*   **Caminho:** `internal/repository/<nome_da_entidade>_repository.go` (ex: `brand_repository.go`)
*   Definir a interface específica herdando do repositório genérico: `commonrepo.Repository[domain.Brand]`.
*   Criar a implementação em GORM correspondente.
*   *Nota:* Caso a entidade dependa de preloads (ex: buscar marca trazendo outras tabelas associadas), sobrescreva o método `FindByID` ou crie métodos de busca específicos que apliquem a chamada `.Preload(...)` do GORM.

### 4. Criar o Validador de Negócio
*   **Caminho:** `internal/validator/<nome_da_entidade>_validator.go` (ex: `brand_validator.go`)
*   Implementar as regras de validação customizadas do domínio antes da persistência (ex: verificar duplicidade de nome, validar formato do código, validar regras cruzadas entre campos).

### 5. Criar a Camada de Serviço (Regras de Negócio)
*   **Caminho:** `internal/service/<nome_da_entidade>_service.go` (ex: `brand_service.go`)
*   Injetar a interface global do `repository.UnitOfWork` (UoW) e o validador específico.
*   **Regra Obrigatória:** Toda operação de escrita, atualização ou exclusão nos repositórios deve ser executada dentro do escopo transacional do Unit of Work:
    ```go
    err := s.uow.Execute(func(uow repository.UnitOfWork) error {
        return uow.Brands().Create(&brand)
    })
    ```

### 6. Criar a Camada de Controle (Handlers)
*   **Caminho:** `internal/api/handlers/<nome_da_entidade>.go` (ex: `brand.go`)
*   Injetar o Serviço específico da entidade.
*   Implementar os endpoints básicos: Criar, Atualizar, Deletar, Obter por ID e Busca com Filtro Avançado (`/filter` via método POST).
*   Utilizar os formatos de resposta padrão (`utils.SuccessResponse`, `utils.ValidationErrorResponse`, `utils.ErrorResponse`).
*   Adicionar as anotações declarativas do Swagger (`@Summary`, `@Tags`, `@Router`, etc.) acima de cada método.

### 7. Criar o Arquivo de Rotas
*   **Caminho:** `internal/api/routes/<nome_da_entidade>_routes.go` (ex: `brand_routes.go`)
*   Registrar os endpoints no grupo de rotas do Gin, configurando os middlewares de autenticação (JWT) e checagem de permissões necessárias.

### 8. Integrar no Inicializador do Servidor
*   **Caminho:** `internal/api/server/server.go`
*   No método `SetupServer()`, instancie o repositório, validador, serviço, handler e registre as rotas chamando a função criada no Passo 7.

### 9. Registrar nas Migrações Automáticas
*   **Caminho:** `internal/models/migrate.go` (ou arquivo de conexão correspondente no microsserviço)
*   Adicionar a nova struct na lista de structs enviadas para o `db.AutoMigrate(...)` para que a tabela seja criada automaticamente no banco de dados durante o boot da aplicação.

### 10. Configurar Cargas de Teste / Dados Iniciais (Seeder)
*   **Caminho:** `internal/models/seeder.go`
*   Caso a tabela precise de dados iniciais ou cadastros de sistema prévios, adicione a lógica de inserção inicial nela.

### 11. Registrar Novas Permissões (RBAC)
*   Criar ou registrar no script de migração/seeds as novas chaves de permissão do sistema para a nova entidade (ex: `view_brands`, `create_brands`, `edit_brands`, `delete_brands`) e associá-las aos perfis apropriados (ex: Administrador).

### 12. Escrever Testes Unitários
*   **Caminho:** `internal/service/<nome_da_entidade>_service_test.go`
*   Desenvolver testes mockados ou de integração leve para testar a lógica do serviço e regras do validador.

---

## 🚫 Erros Comuns a Evitar
1.  **Pular a Camada de Serviço**: Nunca chame repositórios diretamente nos Handlers. Toda lógica passa pela camada de Service.
2.  **Acessar GORM no Service**: Nunca injete `*gorm.DB` nos serviços. A comunicação persistente é gerenciada pelo `UnitOfWork`.
3.  **Não Atualizar o Swagger**: Sempre que criar ou editar rotas e handlers, execute o comando de geração do Swagger na pasta do gateway:
    ```powershell
    cd Backend/microservices/api.gateway
    swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
    ```
