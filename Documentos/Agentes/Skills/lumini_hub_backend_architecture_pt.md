# Diretrizes de Arquitetura do Backend - Lumini Hub

Esta documentação descreve as diretrizes de design, padrões de projeto e regras de estilo adotadas no backend em Go do **Lumini Hub**. Ela serve como especificação oficial para agentes de inteligência artificial (IAs) e desenvolvedores humanos manterem e evoluírem o monorepo local.

---

## 🚀 Visão Geral e Estrutura do Monorepo

O backend é organizado no modelo **Go Workspaces (`go.work`)**, contendo os seguintes módulos principais:
1.  **`common/`**: Módulo compartilhado importado por todos os microsserviços. Contém configurações de banco, utilitários JWT, helpers de paginação e o formato de resposta padrão.
2.  **`microservices/api.gateway/`** (Porta `4000`): Ponto de entrada unificado da aplicação. Atua como proxy reverso e gerenciador centralizado de CORS.
3.  **`microservices/api.auth/`** (Porta `4001`): Microsserviço de gestão de usuários, papéis (Roles) e controle de acesso RBAC.
4.  **`microservices/api.core/`** (Porta `4002`): Microsserviço contendo as regras de negócio de cadastros do ERP (Clientes, Fornecedores, Endereços, Contatos).

Para executar o ecossistema completo de backend em paralelo:
```powershell
cd Backend
.\run_services.bat
```

---

## 🎨 Padrões de Projeto Obrigatórios

### 1. Padrão de Resposta Unificada (ResponseDTO)
Todos os controladores (Handlers) devem utilizar exclusivamente o utilitário `utils.Response` do arquivo [response.go](file:///c:/Projetos/lumini-hub/Backend/common/utils/response.go). Nunca devem ser emitidas respostas em formatos livres como `gin.H`.

#### Formatos de Payload JSON:
-   **Sucesso (`SuccessResponse`):**
    ```json
    {
      "success": true,
      "statusCode": 200,
      "message": "Operação realizada com sucesso",
      "data": { ... }
    }
    ```
-   **Erro Comum (`ErrorResponse`):**
    ```json
    {
      "success": false,
      "statusCode": 500,
      "message": "Falha na operação",
      "error": "detalhes técnicos da falha"
    }
    ```
-   **Erro de Validação (`ValidationErrorResponse`):**
    Aceita tanto uma string de erro única quanto uma lista de strings (`[]string`).
    ```json
    {
      "success": false,
      "statusCode": 400,
      "message": "Dados inválidos fornecidos",
      "validationErrors": [
        "o nome é obrigatório",
        "formato de documento incorreto"
      ]
    }
    ```

### 2. Base Repository Genérico (Go Generics)
Operações básicas de banco de dados (CRUD) são genéricas para evitar código repetitivo:
- Herdar a interface `commonrepo.Repository[T]` no repositório do microsserviço para ganhar métodos como `Create`, `Update`, `Delete`, `FindByID`, `FindAll` automaticamente.
- Se a entidade exigir relacionamentos pré-carregados (Preloads), sobrescreva localmente o método `FindByID(id uint)` aplicando o preload específico.

### 3. Unit of Work (Unidade de Trabalho)
Os serviços (Service) nunca devem interagir diretamente com `*gorm.DB` ou instanciar conexões de transação de banco. Eles devem injetar a interface `repository.UnitOfWork` no construtor para coordenar as escritas e atomicidade das operações.

#### Exemplo de Transação com UoW:
```go
err := s.uow.Execute(func(uow repository.UnitOfWork) error {
    // Ambas as escritas abaixo rodam sob a mesma transação física no banco
    if err := uow.Customers().Create(&customer); err != nil {
        return err // Rollback automático
    }
    if err := uow.Suppliers().Create(&supplier); err != nil {
        return err // Rollback automático
    }
    return nil // Commit automático
})
```

### 4. Consultas e Filtros Complexos via POST (PostFilter)
Para evitar query strings poluídas e lentas em requisições `GET`, todas as listagens com filtragem dinâmica complexa devem ser expostas através de um endpoint do tipo **POST** terminando com `/filter`.

- Definir um DTO de request de filtro (ex: `CustomerFilterRequest`) que encapsula todos os parâmetros de busca mais dados de paginação (`page_no`, `page_size`, `order_by_column`, `is_asc`).
- Mapear a query dinamicamente no repositório usando cláusulas condicionais do GORM.
- Utilizar `utils.Paginate` para realizar a paginação automática e o cálculo de contagem no banco.

### 5. CORS Centralizado no Gateway
- **Proibido** adicionar middlewares de CORS nos microsserviços internos (como `api.auth` e `api.core`). Cabeçalhos de CORS duplicados causam erros de segurança e impedem chamadas Ajax do navegador.
- O CORS é tratado centralmente no **API Gateway** na porta `4000` ([api.gateway/main.go](file:///c:/Projetos/lumini-hub/Backend/microservices/api.gateway/main.go)).

### 6. Swagger API UI no Gateway
- Todos os endpoints devem ser documentados com comentários estruturados no padrão do Swagger (`// @Summary`, `// @Router`, etc.) diretamente nos Handlers de Go.
- O Gateway hospeda a interface do Swagger UI em: **`http://localhost:4000/swagger/index.html`**
- Para registrar novos caminhos, atualize a rota de proxy `NoRoute` no Gateway.
- Sempre regere o Swagger na pasta do Gateway rodando:
  ```powershell
  swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
  ```

---

## 🚫 Erros Comuns a Evitar
- **Duplicação de CORS**: Registrar CORS em microsserviços internos. Mantenha isso exclusivo no Gateway.
- **Pânico no Roteador do Gin**: Adicionar rotas curinga (`router.Any("/*path")`) no mesmo nível que `/swagger/*any`. Sempre utilize `router.NoRoute(...)` para capturar e repassar chamadas de proxy.
- **Esquecimento de Preloads**: Usar o `FindByID` genérico padrão em entidades do Core que exigem relacionamentos vinculados (como Contatos e Endereços). Certifique-se de sobrescrever o método para incluir os `Preload` necessários.
