# Arquitetura ERP (Iluminação e Construção) — Go & React

> [!TIP]
> Este documento traduz a arquitetura robusta de microserviços vista em ecossistemas .NET para o ecossistema e padrões da linguagem **Go (Golang)**. A base utilizada é a de um monorepo (Go Workspaces) consumido por um SPA em React.

## Visão Geral

O sistema ERP para o segmento de iluminação e construção será um sistema distribuído projetado com **Go** no backend e **React (Vite)** no frontend. 
Como Go não possui Solutions (`.sln`) ou o mesmo sistema de injeção de dependência "mágico" do .NET, a estrutura de pastas e organização de pacotes deve ser explícita, seguindo os princípios de *Clean Architecture* e *Domain-Driven Design (DDD)* adaptados para as convenções do Go.

---

## 1. Topologia do Repositório (Monorepo com Go Workspaces)

Para replicar a facilidade de ter uma Solution com projetos `.csproj` que conversam entre si, o Go introduziu o `go.work`. O workspace permite ter vários arquivos `go.mod` (módulos independentes) num mesmo repositório local sem precisar publicar os pacotes compartilhados na internet.

```text
meu-erp/
│
├── .cursor/rules/             # Mantém as regras da IA (como no C#)
│
├── frontend/                  # SPA React hospedado em CDN ou servido pelo Go
│   └── web-app/               # React + TypeScript + Vite + Zustand + Tailwind
│
└── backend/
    ├── go.work                # Arquivo de workspace do Go
    │
    ├── common/                # (Equivalente ao `Common/`) Módulos compartilhados
    │   ├── api/               # Middlewares (Recovery, CORS, Request ID)
    │   ├── auth/              # JWT, Validação de permissões/Claims
    │   ├── database/          # Conexão com banco (PostgreSQL/SQL Server via sqlx)
    │   ├── logger/            # RequestResponseLogging (via zap ou slog)
    │   ├── broker/            # Abstrações RabbitMQ
    │   └── asynq/             # Configurações do Asynq (Alternativa ao Hangfire)
    │
    └── microservices/         # (Equivalente ao `Microservices/`)
        ├── api.auth/          # Autenticação, Usuários, Permissões
        ├── api.estoque/       # Produtos, Lotes, Materiais de Construção/Iluminação
        ├── api.vendas/        # Orçamentos, Pedidos, PDV, Comissões
        ├── api.financeiro/    # Contas a Pagar/Receber, Fluxo de Caixa, NFe
        └── api.projetos/      # Gestão de obras, Cronogramas de instalação
```

---

## 2. Estrutura Interna de um Microserviço em Go

Diferente do .NET onde agrupamos por pastas genéricas (`Services`, `Controllers`, `Models`), em Go é fortemente recomendado agrupar por **Domínio/Feature** dentro de cada microserviço.

Exemplo de como ficaria o microserviço `api.estoque`:

```text
backend/microservices/api.estoque/
│
├── go.mod                     # Dependências locais deste serviço
├── main.go                    # Ponto de entrada (Inicia DB, Servidor, Injeta Deps)
│
├── internal/                  # Tudo aqui é privado ao microserviço (Padrão Go)
│   │
│   ├── config/                # Carregamento de variáveis de ambiente (.env)
│   │
│   ├── domain/                # Modelos, Interfaces e DTOs (Equivalente ao Common.Domain)
│   │   └── produto.go         # Structs: Produto, ProdutoDTO, ProdutoRepository, ProdutoService
│   │
│   ├── handlers/              # (Equivalente aos Controllers) Rotas HTTP
│   │   └── produto_handler.go # Gin/Echo Handlers parsing de JSON e Retornos
│   │
│   ├── services/              # Lógica de Negócio
│   │   └── produto_service.go # Regras (Ex: Validação se o item está em falta)
│   │
│   └── repository/            # Acesso a Dados (Equivalente ao DbContext/Dapper)
│       └── produto_repo.go    # Queries SQL puras ou usando sqlx/GORM
```

---

## 3. Mapeamento de Tecnologias (.NET → Go)

> [!IMPORTANT]
> Em Go, a filosofia é utilizar ferramentas menores e combináveis em vez de frameworks gigantes. O código é ligeiramente mais verboso (explícito), mas incrivelmente performático e fácil de depurar.

| Componente | Ferramenta em .NET (Anterior) | Substituição em Go recomendada | Motivo da escolha |
| :--- | :--- | :--- | :--- |
| **API/Roteamento** | ASP.NET Core MVC (Controllers) | **Gin** ou **Echo** | São os frameworks web mais rápidos, possuem excelente ecossistema de middlewares. |
| **Acesso a Dados** | EF Core + Dapper | **sqlx** ou **GORM** | O `sqlx` mapeia as consultas SQL puras diretamente para as *Structs* em Go (Idêntico ao Dapper). Se preferir um ORM completo, o `GORM` é o padrão de mercado. |
| **Agendamento** | Hangfire | **Asynq** | Backed by Redis, possui interface gráfica Web idêntica ao Hangfire e suporta retries, cron jobs e workers paralelos. |
| **Autenticação** | JwtBearer Middleware nativo | **golang-jwt/jwt** | Biblioteca padrão para gerar e assinar os tokens. Middlewares customizados do Gin interceptam a rota e validam os *Claims*. |
| **Tempo Real** | SignalR | **Centrifugo** ou **Gorilla WebSockets** | O Centrifugo é excelente para gerenciar milhões de conexões de forma agnóstica via RPC. Gorilla é mais baixo nível e padrão. |
| **Injeção de Dependência** | Nativas no `Program.cs` | **Injeção Manual** ou **google/wire** | O mais comum em Go é inicializar o Repo, injetar no Service, e injetar no Handler diretamente no `main.go`. |
| **Armazenamento/Arquivos**| MinIO | **minio-go** | O SDK oficial da MinIO em Go é completo e idêntico. |

---

## 4. O Fluxo de uma Requisição (Inversão de Controle em Go)

Para não acoplar código, o Go depende pesadamente de **Interfaces**.

1. O **Handler (Controller)** recebe a requisição (`POST /produtos`).
2. Ele chama a interface do **Service** (`CriarProduto(dto)`).
3. O **Service** executa regras (ex: validar se a Lâmpada LED já existe) e chama a interface do **Repository** (`Salvar(produto)`).
4. O **Repository** conecta ao banco usando `sqlx` ou `database/sql` e executa o `INSERT`.

> [!TIP]  
> Ao usar Interfaces para os *Services* e *Repositories*, fica extremamente fácil gerar *Mocks* para testes unitários.

---

## 5. FrontEnd (React)

A estrutura do React via Vite continua a mesma. Como será um ERP (muitos dados em tela), recomenda-se:
* **UI/UX:** Como é iluminação/construção, focar em interfaces escuras limpas (Dark Mode) ou tons industriais. Ferramentas como TailwindCSS + Shadcn/UI criam interfaces premium extremamente velozes.
* **Axios/Fetch:** Apontando para o gateway do Go ou para os microserviços diretamente usando CORS.
* **State Management:** Zustand para dados de usuário, e React Query (@tanstack/react-query) para cache de requisições do ERP (Orçamentos, Estoque), aliviando o banco de dados e os microserviços Go.

---

## Próximos Passos Sugeridos

Se o projeto já está iniciado, os passos ideais de adequação arquitetural seriam:

1. Executar `go work init` na pasta do Backend para criar o monorepo.
2. Criar a pasta `common/` e mover as funções básicas (conexão com BD, validação JWT) para lá.
3. Criar o primeiro microserviço (`microservices/api.usuarios` ou `api.estoque`), e adicionar esse caminho no `go.work`.
4. Definir o padrão de *Handlers -> Services -> Repositories* dentro do `internal/` desse primeiro módulo.
