# Lumini Hub — Release 2026.2.2.0

**Versão do sistema:** 2026.2.2.0 · **Versão do banco de dados:** 2 (`dbversion.Required = 2`, tabela `schema_versions`)
**Tema da versão:** **banco que sobe sozinho** (tabelas, permissões, perfil ADMIN, menu e admin inicial criados pelo código), **versão do banco registrada** no próprio banco, **primeiro acesso guiado** (cadastro da primeira empresa) e **catálogo de permissões completo**.
*Documento gerado em 2026-10-22. Legenda:* ✅ testado · ⚠️ implementado, ainda não testado · ⏳ pendente / não incluído.

---

## 1. Banco de dados

Esta é a primeira versão que registra a versão do banco dentro dele (versão 2; o ambiente de testes da OCI estava na versão 1, criada por script SQL).

- **Nova tabela `schema_versions`** (`db_version`, `app_version`, `applied_at`): uma linha por publicação que chegou ao banco. O `api.auth` cria a tabela e grava a versão no fim de um boot bem-sucedido com `DB_AUTO_MIGRATE=true`; reiniciar não duplica a linha.
- Na subida, com `DB_AUTO_MIGRATE=true` (padrão), cada serviço cria as tabelas do seu schema (`api.auth`: 4; `api.core`: 10; `api.integrations`: 4) e o `api.auth` semeia o catálogo de permissões, o perfil `ADMIN`, o menu e o admin inicial.
- Com `DB_AUTO_MIGRATE=false` (produção com dados reais), nada é migrado nem semeado, e `api.auth`, `api.core` e `api.integrations` **conferem a versão do banco** e recusam subir se ele estiver atrás da versão 2 ou sem versão registrada ⚠️.

> **Atenção:** um banco da versão 1 (criado pelos scripts `deploy/db/*.sql`, com ids `integer`) **não deve ser atualizado por cima**: o seed só cria o perfil `ADMIN` quando não há nenhum perfil e pula itens de menu existentes. No ambiente da OCI o banco foi recriado do zero (ver seção 4).

---

## 2. Novidades

### 2.1 Banco que sobe sozinho ✅
- Banco vazio sobe apenas com o código: tabelas, **67 permissões**, perfil **ADMIN** com **58 vínculos** (todas, exceto as 9 do módulo `Develop`, que o ADMIN já alcança pelo bypass) e **29 itens de menu**.
- **Admin inicial** criado a partir de `BOOTSTRAP_ADMIN_USERNAME` (padrão `admin`) e `BOOTSTRAP_ADMIN_PASSWORD` (mínimo de 8 caracteres e máximo de 72 bytes). Só é criado com a tabela de usuários vazia, nasce sem empresa vinculada e a senha nunca vai para logs. Sem a senha, nenhum usuário é criado e o log avisa.
- Os scripts SQL provisórios (`deploy/db/01_schema.sql` e `02_catalog.sql`) foram removidos.

### 2.2 Versão do banco gravada no próprio banco ✅
- As constantes `Required` (versão do banco exigida) e `AppVersion` (versão do sistema) ficam em `Backend/common/dbversion/dbversion.go` e são atualizadas a cada release.
- Banco à frente do código não é erro (permite voltar uma versão do código).

### 2.3 Primeiro acesso: cadastro da primeira empresa ✅
- Num banco sem nenhuma empresa, o usuário autorizado (`companies.create`, inclusive ADMIN e DEVELOP) cadastra a empresa **direto na tela de seleção** (botão **Cadastrar Empresa**) e entra no sistema ao salvar. Sem permissão, a tela mostra uma mensagem e o botão sair.
- O roteiro completo (clonar, `.env`, subir, login do admin, cadastrar a empresa, remover a senha do `.env`) está no `deploy/README.md`.

### 2.4 Catálogo de permissões completo ✅
Entram as 8 permissões que as rotas já exigiam e que o catálogo não tinha:

- **Vendas:** `customers.create`, `customers.edit`, `customers.delete`.
- **Compras:** `suppliers.view`, `suppliers.create`, `suppliers.edit`, `suppliers.delete`.
- **Integrações:** `integrations.edit`.

Com isso toda permissão exigida por `RequirePermission` existe no catálogo (antes, essas rotas só funcionavam para ADMIN e DEVELOP).

### 2.5 Menu "Fornecedores" separado de Compras ⚠️
- O item de menu **Fornecedores** (cadastro) passa a usar `suppliers.view`, separado de `purchases.*` (a futura rotina de Compras, com entrada de notas). Vale para bancos novos.

---

## 3. Correções

- **Log do banco por ambiente** ✅: fora de `development`, o log do GORM deixa de mostrar o SQL com parâmetros (e com isso o hash de senha). Testado com `APP_ENV=staging`: só as requisições aparecem.
- **Menu "Perfis e Permissões" escondido do ADMIN** ⚠️: num banco semeado do zero o item nascia ligado a uma permissão que só o DEVELOP tem; agora usa `roles.view`.
- **`go run main.go` no `api.auth`** ✅: o comando voltou a compilar (a função de preparo do banco ficou de novo no `main.go`, e um teste protege isso).

---

## 4. Infraestrutura

- **`deploy/docker-compose.prod.yml`**: deixou de montar `./db`; repassa `DB_AUTO_MIGRATE` a todos os serviços e `BOOTSTRAP_ADMIN_*` ao `api-auth`. `deploy/.env.example` e `deploy/README.md` atualizados.
- **Ambiente da OCI** (VM `lumini-hub-dev`, `http://64.181.183.243`) atualizado para esta versão em 2026-10-22 ✅: código na `develop`, **banco recriado do zero** (os dados de teste anteriores foram apagados: usuário `admin` antigo, `TESTE` e a empresa "Minha Empresa"), imagens reconstruídas na própria VM (arm64) e primeiro acesso refeito com o roteiro do `deploy/README.md`. Conferido: 67 permissões, 58 vínculos do ADMIN, versão do banco 2, um usuário (`admin`), sem senha nos logs.
- O `api.integrations` continua **fora** do compose (profile `integrations`); o acesso ao legado vai virar uma aplicação separada.
- Acesso por IP, em HTTP (`APP_ENV=staging`): sem domínio e sem HTTPS ainda.

---

## 5. Configuração necessária no servidor

- **`deploy/.env`** (em um servidor novo ou recriado):
  - `APP_ENV`: `staging` enquanto o acesso for por HTTP. Com `production` os cookies de login exigem HTTPS.
  - `DB_AUTO_MIGRATE`: `true` (padrão) para um banco novo ou para atualizar; `false` em produção com dados reais.
  - `BOOTSTRAP_ADMIN_USERNAME` e `BOOTSTRAP_ADMIN_PASSWORD`: só para o primeiro boot. **Remover a senha do `.env`** depois e recriar só o `api-auth`.
- **Bancos já semeados (ambiente de desenvolvimento e qualquer um que não seja novo):**
  - O item de menu "Fornecedores" continua em `purchases.view`; ajustar a mão para `suppliers.view`, porque o seed do menu pula itens existentes (ver 7).
  - As 8 permissões novas entram no próximo boot do `api.auth` e vão automaticamente só para o ADMIN. Usuários com permissões individuais precisam recebê-las pela tela de Usuários.
- **Servidor com `DB_AUTO_MIGRATE=false`:** para subir esta versão, rodar o `api-auth` uma vez com a flag ligada (ele migra, semeia e registra a versão 2) e depois voltar a flag para `false`.

---

## 6. Arquivos principais alterados / criados

**Biblioteca comum (`Backend/common`):** `dbversion/dbversion.go` (novo), `config/config.go`, `database/db.go`, `database/log.go`, `testutil/postgres.go`.
**api.auth:** `main.go`, `seed/` (`admin.go`, `catalog.go`, `permission_catalog.go`, `permissions.go`), `internal/seeder/menu_item_seeder.go`, `internal/domain/role.go`, `internal/repository/role_repository.go`.
**api.core e api.integrations:** `main.go` (conferência da versão do banco).
**Frontend:** `components/companies/CompanySelectionGate.tsx`, `contexts/AuthContext.tsx`, `hooks/useCreateCompany.ts`, `pages/settings/companies/CreateCompanyPage.tsx`.
**Deploy:** `deploy/docker-compose.prod.yml`, `deploy/.env.example`, `deploy/README.md`, `Backend/.env.example`; removidos `deploy/db/01_schema.sql`, `02_catalog.sql` e `README.md`.

---

## 7. Pendências conhecidas (próximas versões)

| Item | Situação |
|---|---|
| **CFG-014** Descrições e módulos do catálogo de permissões herdados do dev (`roles.*` dizem "usuários", módulo `dashboard` em minúsculo, "Ver"/"Editar" em Integrações) | ⏳ Só texto; no banco já semeado precisa de um `UPDATE` à parte, porque o seed não altera permissões existentes |
| **CFG-015** Conferir os itens de menu do `menu_item_seeder.go` contra o `menu_items` do dev | ⏳ Só `/settings/*` foi conferido |
| **CFG-016** Item "Fornecedores" em bancos já semeados continua em `purchases.view` | ⏳ Ajuste manual ou migração |
| **CFG-017** `APP_ENV` não definido vira `development` (log verboso) | ⏳ Exigir explícito fora de dev ou inverter o padrão |
| **CFG-018** Nunca logar `err` bruto nos handlers (a mensagem do Postgres pode citar valores) | ⏳ |
| **CFG-024** Barra de status no layout com a versão do sistema, a versão do banco e outras informações | ⏳ Precisa de um endpoint de versão |
| **CFG-025** Scripts de migração por versão do banco (para o que o `AutoMigrate` não faz: remover/renomear, ajustar dados, triggers) | ⏳ Só quando aparecer a primeira mudança desse tipo |
| **CFG-022 / CFG-023** Tradução do `CLAUDE.md` e das skills; arquivar o histórico dos `tasks_*.md` | ⏳ |
| **VND-001 / COM-001** Telas de Clientes e de Fornecedores | ⏳ O backend e os services do frontend já existem |
| **HTTPS e domínio** no ambiente da OCI | ⏳ Hoje por IP, em HTTP |
