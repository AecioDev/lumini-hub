# Lumini Hub - Backend (Go Workspace)

Este é o repositório de backend do **Lumini Hub**, um sistema ERP moderno reestruturado a partir de um monólito em Go para uma **arquitetura de microsserviços modular** baseada em **Go Workspaces (`go.work`)**.

---

## 🚀 Arquitetura do Projeto

O backend está subdividido em módulos independentes sob o mesmo workspace:

*   **`common/`**: Pacote de código compartilhado que contém configurações, conexões de banco de dados, middlewares utilitários de autorização e paginação comuns.
*   **`microservices/api.gateway/`**: Proxy reverso transparente centralizado rodando na porta **`4000`**. Roteia as requisições para os microsserviços correspondentes e gerencia a segurança CORS de forma única.
*   **`microservices/api.auth/`**: Microsserviço rodando na porta **`4001`**, responsável por autenticação de usuários, renovação de tokens (JWT), gerenciamento de usuários e controle de acesso RBAC.
*   **`microservices/api.core/`**: Microsserviço rodando na porta **`4002`**, contendo os cadastros essenciais do ERP: Clientes, Fornecedores, Endereços, Contatos, Documentos e Empresas (fundação multi-empresa).
*   **`microservices/api.integrations/`**: Microsserviço rodando na porta **`4007`**, middleware de sincronização com a Loja Integrada e o ERP legado em SQL Server.

Os padrões obrigatórios de código (Response DTO, Repository genérico, Unit of Work, RBAC, etc.) e o detalhamento de cada microsserviço estão documentados em [`CLAUDE.md`](../CLAUDE.md), na raiz do repositório — é a referência técnica canônica do projeto.

---

## 🛠️ Requisitos

- **Go** v1.21 ou superior
- **PostgreSQL** v12 ou superior (rodando localmente ou via container na porta `5432`)

---

## ⚙️ Configuração

1. Certifique-se de ter o PostgreSQL rodando localmente na porta `5432`.
2. Configure o arquivo `.env` na pasta `Backend` com suas credenciais do banco. Exemplo:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=sua_senha
   DB_NAME=lumini-hub
   DB_SSLMODE=disable
   
   JWT_SECRET=seu_segredo_jwt_super_seguro
   JWT_EXPIRES_IN=24h
   ```

---

## 🚦 Execução Local

Para facilitar a inicialização simultânea dos serviços do backend, disponibilizamos o script `run_services.bat` na pasta raiz.

1. Abra o terminal na pasta `Backend`.
2. Execute o comando:
   ```powershell
   .\run_services.bat
   ```
   *(Isso abrirá automaticamente novas janelas do Prompt de Comando para executar individualmente a `api.auth`, `api.core`, `api.integrations` e o `api.gateway`)*.

Agora, todas as requisições do frontend devem ser feitas diretamente à porta **`4000`** (ex: `http://localhost:4000/api/...`).

---

## 🧪 Estrutura de Diretórios

```bash
Backend/
├── common/                # Módulo compartilhado importável
│   ├── config/            # Leitura do .env
│   ├── database/          # Conexão GORM
│   ├── middlewares/       # Auth e RBAC middlewares
│   └── utils/             # JWT, Criptografia, Paginação
├── microservices/
│   ├── api.gateway/       # Gateway de Entrada (Porta 4000)
│   ├── api.auth/          # Auth, Users, Roles & Perms (Porta 4001)
│   ├── api.core/          # Clientes, Fornecedores, Empresas & Cadastros (Porta 4002)
│   └── api.integrations/  # Sync Loja Integrada <-> SQL Server legado (Porta 4007)
├── go.work                # Configuração do Workspace do Go
├── run_services.bat       # Script utilitário de execução paralela
└── README.md              # Este arquivo
```
