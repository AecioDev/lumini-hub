---
name: lumini-hub-backend-architecture
description: >-
  Diretrizes de arquitetura e padrões de código para evoluir e manter o backend Go do Lumini Hub.
---

# Diretrizes de arquitetura do backend do Lumini Hub

## Visão geral
Esta skill define os padrões obrigatórios de arquitetura, os padrões de projeto e as diretrizes de código para qualquer agente de IA que trabalhe no backend Go do Lumini Hub. Siga estas instruções à risca ao criar microsserviços, repositories, services, handlers ou rotas novos.

---

## Início rápido
O backend do Lumini Hub é estruturado como um **Go Workspace (`go.work`)** que contém:
1.  **`common/`**: módulo compartilhado (conexões com o banco, JWT, paginação, respostas padrão).
2.  **`microservices/api.gateway/`** (porta `4000`): ponto de entrada com proxy reverso e gerenciador de CORS.
3.  **`microservices/api.auth/`** (porta `4001`): microsserviço de usuários, perfis e permissões (RBAC).
4.  **`microservices/api.core/`** (porta `4002`): microsserviço de clientes, fornecedores, endereços, contatos e empresas (base multiempresa).
5.  **`microservices/api.integrations/`** (porta `4007`): integrações do sistema com serviços externos (hoje, a Loja Integrada). **Pausado desde 2026-10-01**: a Fase 1 (configuração, acompanhamento de sync-log/webhook-event, proxy no gateway) está no código e não evolui; a Fase 2 (cliente da Loja Integrada e processamento dos webhooks) não foi feita. Veja a seção dedicada do `CLAUDE.md`.

Para subir todos os microsserviços do backend em paralelo, localmente:
```powershell
cd Backend
.\run_services.bat
```

---

## Fluxo de trabalho e padrões de projeto

### 1. Formato único de resposta (ResponseDTO)
Todos os handlers **DEVEM** usar a struct padrão `utils.Response` de [response.go](file:///c:/Projetos/lumini-hub/Backend/common/utils/response.go) para devolver JSON. Nunca escreva respostas improvisadas com map (`gin.H`).

#### Formato do JSON:
- **Sucesso (`SuccessResponse`)**:
  ```json
  {
    "success": true,
    "statusCode": 200,
    "message": "Recurso encontrado",
    "data": { ... }
  }
  ```
- **Erro (`ErrorResponse`)**:
  ```json
  {
    "success": false,
    "statusCode": 404,
    "message": "Usuário não encontrado",
    "error": "texto da regra de negócio, mostrado ao usuário"
  }
  ```
  Num **500**, o `ErrorResponse` nunca devolve o texto do erro ao cliente: grava o erro completo no log do servidor (método, caminho e erro) e a resposta leva só a `message`, sem o campo `error`. Pode passar o `err.Error()` cru num 500, mas nunca ponha erro bruto do banco numa resposta 4xx, onde o texto é mostrado ao usuário.
- **Erros de validação (`ValidationErrorResponse`)**:
  Aceita uma única string ou uma lista de strings (`[]string`).
  ```json
  {
    "success": false,
    "statusCode": 400,
    "message": "Dados inválidos",
    "validationErrors": [
      "o nome é obrigatório",
      "formato de documento inválido"
    ]
  }
  ```

### 2. Padrão Repository base (Go Generics)
Toda a lógica básica de CRUD no banco é genérica.
- Estenda `commonrepo.Repository[T]` para os métodos básicos de CRUD (`Create`, `Update`, `Delete`, `FindByID`, `FindAll`).
- As interfaces de repository específicas devem compor a interface genérica.

#### Exemplo de definição de repository:
```go
type CustomerRepository interface {
	commonrepo.Repository[domain.Customer]
	FindByID(id uint) (*domain.Customer, error) // Sobrescreve para aplicar Preloads específicos do GORM
	FindByDocument(document string) (*domain.Customer, error) // Método específico do domínio
}
```

### 3. Padrão Unit of Work (UoW)
Os services nunca devem receber repositories soltos nem conexões `*gorm.DB` diretas. Devem usar o `repository.UnitOfWork` para gerenciar as operações no banco e as transações atômicas.

#### Exemplo de transação atômica num service:
```go
func (s *CustomerService) CreateCustomer(req domain.CreateCustomerRequest, userID uint) (*domain.ApiCustomer, error) {
	// ... lógica de validação ...
	customer := domain.Customer{ ... }

	// Executa dentro de um bloco de transação atômica gerenciado pelo UoW
	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Customers().Create(&customer)
	})
	if err != nil {
		return nil, err
	}
	return &customerDTO, nil
}
```

### 4. Busca por HTTP POST (padrão PostFilter)
Para tratar consultas complexas, com vários parâmetros e paginação, sem encher a URL de query parameters, as buscas devem ser implementadas como um endpoint **POST** terminado em `/filter`.

- Crie um DTO de filtro (por exemplo, `CustomerFilterRequest` no domain).
- Os parâmetros do corpo JSON incluem os campos de filtro e as variáveis de paginação (`page_no`, `page_size`, `order_by_column`, `is_asc`).
- No repository, monte os filtros do GORM dinamicamente:
  ```go
  if filter.FirstName != "" {
      query = query.Where("first_name ILIKE ?", "%"+filter.FirstName+"%")
  }
  ```
- Pagine os resultados com `utils.Paginate`.

### 5. Gerenciamento centralizado de CORS
- **Nunca** registre middlewares de CORS nos microsserviços internos (como `api.auth` ou `api.core`). Isso gera cabeçalhos de CORS duplicados no navegador e dispara `AxiosError: Network Error`.
- O CORS é central e gerenciado exclusivamente no **API Gateway** ([api.gateway/main.go](file:///c:/Projetos/lumini-hub/Backend/microservices/api.gateway/main.go)).

### 6. Documentação Swagger e roteamento no Gateway
- Documente todos os métodos HTTP nos handlers com as anotações declarativas padrão do Swagger (`@Summary`, `@Tags`, `@Router` etc.).
- Ao adicionar rotas novas, exponha-as pelo gateway em `api.gateway/main.go`, usando o handler de proxy do `router.NoRoute`. Garanta que o caminho não colida com `/swagger`.
- Regere a documentação do Swagger dentro da pasta do Gateway:
  ```powershell
  swag init -g main.go -d ./,../../common,../api.auth,../api.core --parseDependency
  ```

---

## Erros comuns
- **CORS duplicado**: adicionar CORS a qualquer microsserviço interno. Só o Gateway cuida do CORS.
- **Conflito de caminhos no Gin**: registrar `router.Any("/*path")` junto com `/swagger/*any` no mesmo nível de grupo do roteador Gin. Direcione sempre os caminhos de proxy não mapeados por `router.NoRoute(...)`.
- **SQL montado com texto**: concatenar ou usar `fmt.Sprintf` para pôr dado do usuário dentro de um `Where`/`Raw`/`Exec` abre a porta para SQL Injection. Valores sempre como parâmetro (`?`); a ordenação das listagens passa pelo `utils.Paginate`, que só aceita colunas reais do modelo; nunca monte um `Order(...)` com texto da requisição fora dele.
- **Esquecer o Preload**: não sobrescrever o `FindByID` ou o `FindByFilter` genéricos nos repositories do GORM quando a entidade precisa de preloads (por exemplo, endereços e contatos).
