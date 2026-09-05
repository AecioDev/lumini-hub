# Lumini Hub - Sistema ERP

O **Lumini Hub** é um sistema ERP moderno, desenvolvido com uma arquitetura distribuída no backend em **Go (Golang)** estruturado em microsserviços e um frontend em **Vite + React + Ant Design**.

---

## 📂 Estrutura do Monorepo

O projeto está organizado em formato monorepo para facilitar o gerenciamento de código:

```bash
lumini-hub/
├── Backend/               # API do Sistema baseada em Go Workspaces
│   ├── common/            # Pacote compartilhado (banco de dados, middlewares, utilitários)
│   └── microservices/     # api.gateway (4000), api.auth (4001), api.core (4002), api.integrations (4007)
├── Frontend/              # Interface Web em Vite + React + Ant Design (Porta 3000)
├── Documentos/            # Planejamento, histórico de desenvolvimento e materiais de referência
│   ├── Planejamento/      # Rastreador de tarefas por módulo (fonte da verdade do que falta fazer)
│   ├── exemplos/          # Trechos de código do ERP legado usados como referência de integração
│   └── Imagens/           # Identidade visual
├── .claude/skills/        # Skills carregadas automaticamente pelo Claude Code (regras de arquitetura e fluxo)
└── README.md               # Este arquivo
```

---

## 🛠️ Pré-requisitos Gerais

Para rodar o projeto localmente, você precisará ter instalado em sua máquina:

- **Go (Golang)** v1.21 ou superior
- **Node.js** v18 ou superior + **pnpm** (gerenciador de pacotes)
- **PostgreSQL** v12 ou superior (rodando localmente na porta `5432`)

---

## ⚡ Inicialização Rápida do Sistema

Siga os passos abaixo para colocar o sistema completo para rodar na sua máquina local:

### 1. Banco de Dados (PostgreSQL)
Certifique-se de que o seu serviço PostgreSQL local está ativo na porta `5432` e possui a base de dados configurada conforme os arquivos de ambiente do backend (`Backend/.env`, veja `Backend/.env.example`).

---

### 2. Backend (Microsserviços)
Abra um terminal na pasta [Backend](Backend) e inicialize todos os serviços rodando o script utilitário:

```powershell
cd Backend
.\run_services.bat
```
*(Isso iniciará automaticamente o API Gateway na porta **4000**, os microsserviços de Autenticação (**4001**), Core (**4002**) e Integrações (**4007**) em janelas de terminal dedicadas).*

---

### 3. Frontend (Interface Web)
Abra um **segundo terminal** na pasta [Frontend](Frontend) para instalar as dependências e iniciar o servidor Vite:

```powershell
cd Frontend
pnpm install
pnpm dev
```
*(O frontend será iniciado e ficará disponível em `http://localhost:3000`, porta fixada em `vite.config.ts` para casar com o CORS do gateway).*

---

## 🔑 Credenciais para Acesso Local

Após iniciar os servidores, abra o seu navegador e acesse: **[http://localhost:3000](http://localhost:3000)**

- **Usuário:** `admin`
- **Senha:** `987321`

---

## 📑 Documentação Recomendada

- **[CLAUDE.md](CLAUDE.md)** — referência técnica canônica do projeto (arquitetura de backend e frontend, padrões obrigatórios, comandos). É o primeiro documento que qualquer agente de IA lê antes de mexer no código, e deve ser mantido atualizado junto com o código.
- **[Documentos/Planejamento/README.md](Documentos/Planejamento/README.md)** — rastreador de tarefas por módulo do ERP (o que já foi feito, o que está pendente).
- **`.claude/skills/lumini_hub_*`** — regras detalhadas de arquitetura de backend, fluxo de desenvolvimento e criação de entidades, carregadas automaticamente pelo Claude Code.
