# Lumini Hub — Release 2026.4.5.0

**Versão do sistema:** 2026.4.5.0 · **Versão do banco de dados:** 4 (novas tabelas `system_logs` e `schema_migrations`, índices únicos parciais; `dbversion.Required = 4`)
**Tema da versão:** **auditoria** de tudo que é criado, alterado ou excluído (com o valor antes e depois), **migrações SQL versionadas** para o que o `AutoMigrate` não cobre, **integridade dos cadastros** (documento e CNPJ de registros excluídos voltam a ficar livres) e **testes** do `api.core`.
*Documento gerado em 2026-10-31. Legenda:* ✅ testado · ⚠️ implementado, ainda não testado · ⏳ pendente / não incluído.

---

## 1. Banco de dados

A versão do banco sobe de **3 para 4**. Há duas tabelas novas, índices novos e duas constraints antigas removidas; nenhum dado é reescrito.

| Tabela | Mudança | Detalhe |
|---|---|---|
| `system_logs` | **nova** ✅ | Auditoria (CFG-037): `id`, `created_at`, `user_id`, `action`, `entity_type`, `entity_id`, `details` (jsonb), `ip_address`. Índices em data, usuário e entidade. Sem exclusão lógica e sem chave estrangeira. |
| `schema_migrations` | **nova** ⚠️ | Migrações SQL já aplicadas (CFG-051): `name` e `applied_at`. |
| `customers` | índice ⚠️ | A constraint única comum de `document_number` (`uni_customers_document_number`) é trocada pelo índice único parcial `idx_customers_document_number` (`WHERE deleted_at IS NULL`). |
| `suppliers` | índice ⚠️ | Idem: `uni_suppliers_document_number` dá lugar a `idx_suppliers_document_number`. |
| `companies` | índice ⚠️ | Novo índice único parcial `idx_companies_cnpj` (antes só o validador conferia o CNPJ). |
| `permissions` | dados ✅ | Permissão nova `audit.view` (módulo Admin), vinculada ao ADMIN no boot. A migração `0001` corrige os textos de 12 permissões de bancos semeados antes do CFG-014 (só linhas com o texto antigo exato; sem efeito em banco que já recebeu a correção). |
| `menu_items` | dados ✅ | Item novo "Auditoria" (`/settings/audit`), semeado sozinho em bancos já existentes. |

- Tudo é feito pelo `api.auth` no primeiro boot com `DB_AUTO_MIGRATE=true`: `AutoMigrate`, seed, **migrações SQL** (`0001` e `0002`) e registro da versão 4.
- Com `DB_AUTO_MIGRATE=false`, os serviços **se recusam a subir** até o banco estar na versão 4 e sem migração pendente: rode o `api.auth` uma vez com a variável ligada (ver a seção 5).

---

## 2. Novidades

### 2.1 Auditoria (CFG-037) ✅
- Toda escrita bem-sucedida (criar, alterar, excluir) nos serviços `api.auth`, `api.core` e `api.integrations` é registrada com usuário, ação, entidade, ID e IP. Buscas de lista, login, sessão e webhooks ficam de fora, e o corpo da requisição nunca é gravado.
- Cada registro guarda **o que mudou**: numa alteração, só os campos que mudaram com o valor de antes e de depois (em listas como as permissões de um usuário, o que entrou e o que saiu); numa criação, o objeto criado; numa exclusão, o objeto excluído. Senhas, tokens, hashes e certificados nunca são gravados.
- Tela nova **Configurações → Auditoria**, só leitura, com busca por usuário e ação, filtros de entidade, ação e período, e a gaveta "Ver detalhes". Endpoints `POST /audit-logs/filter` e `GET /audit-logs/entity-types`, com a permissão `audit.view`.
- O que usuários do Perfil DEVELOP fazem só aparece para o próprio DEVELOP.

### 2.2 Migrações SQL versionadas (CFG-051) ⚠️
- Arquivos `NNNN_*.sql` em `Backend/common/dbversion/migrations/`, embutidos no binário, aplicados pelo `api.auth` no boot, **uma vez por banco**, em ordem e cada um numa transação (se falhar, desfaz e o `api.auth` não sobe, citando a migração).
- Substituem os scripts soltos em `deploy/sql/` para tudo o que o `AutoMigrate` não faz (remover constraint ou índice antigo, corrigir dados existentes). Nesta versão: `0001` (textos das permissões) e `0002` (índices únicos parciais).

---

## 3. Correções

- **Recadastro após exclusão (CFG-050)** ⚠️: recadastrar o documento de um cliente ou fornecedor excluído, ou o CNPJ de uma empresa excluída, agora funciona (antes passava na validação e falhava no banco com "Não foi possível concluir a operação."). Duas criações simultâneas com o mesmo CNPJ de empresa também deixam de duplicar.
- **Autoria do fornecedor (CFG-049, segurança)** ✅: o cadastro de fornecedor deixou de aceitar `created_by` no corpo da requisição; o autor é sempre o usuário logado (antes, quem chamasse a API podia atribuir a criação a outro usuário).
- **`LoggerMiddleware`** ✅: existia mas não estava ligado em nenhum serviço, e a tabela `system_logs` nunca era criada; agora é (ver Auditoria).

---

## 4. Infraestrutura

- **Testes do `api.core` (CFG-044)** ✅: 22 testes de service (empresa, configuração fiscal com certificado A1, cliente e fornecedor) e testes das migrações; precisam de `TEST_DATABASE_URL`. Certificado autoassinado só de teste em `testdata/`.
- Swagger do gateway regerado (endpoints da auditoria; `created_by` sai da criação de fornecedor).
- Gateway: novo prefixo `/api/audit-logs`, repassado ao `api.auth`.
- Nenhuma mudança nos Dockerfiles nem no Caddy.

---

## 5. Configuração necessária no servidor

- **Atualizar o código, reconstruir as imagens e reiniciar** `gateway`, `api-auth`, `api-core`, `api-integrations` e `web`. Sem recriar o banco: os dados são mantidos.
- **Banco:** com `DB_AUTO_MIGRATE=true`, o `api-auth` cria as tabelas, aplica as migrações `0001` e `0002` e registra a versão 4 sozinho no boot (o log mostra "Migração SQL aplicada" para cada uma e "Banco registrado na versão 4"). Reiniciar não repete nada. Com `DB_AUTO_MIGRATE=false`, rode o `api-auth` uma vez com a variável ligada antes de subir os demais serviços.
- **Antes de subir:** se houver CNPJ de empresa ativa repetido no banco, a migração `0002` falha e o `api-auth` não sobe; corrija os duplicados e suba de novo.
- **Sem passo manual:** o script de textos das permissões (CFG-041) e o dos índices (CFG-050) rodam como migrações; não é preciso usar o `psql`.
- Depois de subir, a barra de status deve mostrar "Versão 2026.4.5.0 - BD v4".

---

## 6. Arquivos principais alterados / criados

**Biblioteca comum (`Backend/common`):** `database/models.go` (`SystemLog`), `middlewares/` (`logger.go`, `audit_changes.go`, `logger_test.go`), `dbversion/` (`migrations.go`, `migrations_test.go`, `migrations/0001_*.sql`, `migrations/0002_*.sql`, `dbversion.go`).
**api.auth:** `main.go`, `internal/domain/audit_log.go`, `internal/repository/audit_log_repository.go`, `internal/service/audit_log_service.go`, `internal/handlers/audit_logs.go`, `internal/routes/routes.go`, `seed/permission_catalog.go`, `internal/seeder/menu_item_seeder.go`.
**api.core:** `main.go`, `internal/domain/` (`customer.go`, `supplier.go`, `company.go`), `internal/service/supplier_service.go` e os testes de service (novos).
**api.integrations:** `main.go`.
**api.gateway:** `main.go`, `docs/` (Swagger regerado).
**Frontend:** `pages/settings/audit/AuditLogsPage.tsx`, `components/audit/AuditDetailDrawer.tsx`, `services/audit/audit-service.ts`, `types/audit.ts`, `routes/router.tsx`.

---

## 7. Pendências conhecidas (próximas versões)

| Item | Situação |
|---|---|
| **VND-001 / COM-001** Telas de Clientes e de Fornecedores | ⏳ |
| **Recuperar senha** de ponta a ponta (a tela existe, falta o endpoint e o canal de e-mail) | ⏳ |
| **Retenção da auditoria**: limpeza automática de registros antigos de `system_logs`, e auditar login e logout | ⏳ |
| **`DB_AUTO_MIGRATE=false`** recusando subir com banco atrás da versão ou com migração pendente | ⚠️ passou nos testes automáticos, ainda não testado num servidor real |
| **Integrações, Fase 2** (cliente da Loja Integrada e processamento dos webhooks) | ⏳ pausado |
| **HTTPS e domínio** no ambiente de testes | ⏳ Hoje o acesso é por IP, em HTTP |
