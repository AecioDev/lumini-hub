# Lumini Hub - Sistema ERP

O **Lumini Hub** é um sistema ERP moderno, desenvolvido com uma arquitetura distribuída no backend em **Go (Golang)** estruturado em microsserviços e um frontend interativo e responsivo construído em **React/Next.js**.

---

## 📂 Estrutura do Monorepo

O projeto está organizado em formato monorepo para facilitar o gerenciamento de código:

```bash
lumini-hub/
├── Backend/               # API do Sistema baseada em Go Workspaces
│   ├── common/            # Pacote compartilhado (banco de dados, middlewares, utilitários)
│   └── microservices/     # Microsserviços e API Gateway (Portas 4000, 4001, 4002)
├── Frontend/              # Interface Web em Next.js / Tailwind CSS (Porta 3000)
├── Documentos/            # Manuais e documentações arquiteturais do projeto
│   └── exemplos/
└── README.md              # Este arquivo
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
Certifique-se de que o seu serviço PostgreSQL local está ativo na porta `5432` e possui a base de dados configurada conforme os arquivos de ambiente do backend.

---

### 2. Backend (Microsserviços)
Abra um terminal na pasta [Backend](file:///c:/Projetos/lumini-hub/Backend) e inicialize todos os serviços rodando o script utilitário:

```powershell
cd Backend
.\run_services.bat
```
*(Isso iniciará automaticamente o API Gateway na porta **4000**, o microsserviço de Autenticação na porta **4001** e o Core na porta **4002** em janelas de terminal dedicadas).*

---

### 3. Frontend (Interface Web)
Abra um **segundo terminal** na pasta [Frontend](file:///c:/Projetos/lumini-hub/Frontend) para instalar as dependências e iniciar o servidor Next.js:

```powershell
cd Frontend
pnpm install
pnpm dev
```
*(O frontend será iniciado e ficará disponível em `http://localhost:3000`).*

---

## 🔑 Credenciais para Acesso Local

Após iniciar os servidores, abra o seu navegador e acesse: **[http://localhost:3000](http://localhost:3000)**

- **Usuário:** `admin`
- **Senha:** `987321`

---

## 📑 Documentação Recomendada
Para entender mais a fundo o funcionamento dos microsserviços, acoplamento de banco de dados e guias de expansão de código, consulte a [Documentação da Arquitetura de Microsserviços](file:///c:/Projetos/lumini-hub/Documentos/arquitetura_microsservicos.md).
