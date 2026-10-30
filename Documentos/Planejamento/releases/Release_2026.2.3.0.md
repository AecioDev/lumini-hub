# Lumini Hub — Release 2026.2.3.0

**Versão do sistema:** 2026.2.3.0 · **Versão do banco de dados:** 2 (sem mudança de esquema; `dbversion.Required = 2`)
**Tema da versão:** **segurança e validação da API** (erros do banco e detalhes internos deixam de chegar ao usuário, ordenação das listagens e tamanho dos campos validados antes do banco), **barra de status** com a versão do sistema e do banco, e acabamento do catálogo de permissões.
*Documento gerado em 2026-10-29. Legenda:* ✅ testado · ⚠️ implementado, ainda não testado · ⏳ pendente / não incluído.

---

## 1. Banco de dados

O esquema do banco **não muda** nesta versão: a versão do banco continua a 2, e não há tabela, coluna ou índice novo.

Há uma correção de **dados** em bancos que já existem:

- Os itens de menu "Fornecedores" (em Compras e em Financeiro > Contas a Pagar) passam a usar a permissão `suppliers.view`. O seed do menu pula itens que já existem, então um banco semeado antes precisa rodar o script idempotente `deploy/sql/CFG-016_menu_fornecedores_suppliers_view.sql` ✅. Bancos novos já nascem certos.
- Os textos e o módulo de algumas permissões do catálogo foram corrigidos no código, mas o seed nunca altera uma permissão que já existe: em bancos já semeados os textos antigos continuam (ver a seção 7).

---

## 2. Novidades

### 2.1 Barra de status no rodapé ✅
- Uma barra discreta, de ponta a ponta, na cor do menu lateral e do cabeçalho, mostra **"Versão X - BD vN"**.
- A versão é um **link** que abre as notas daquela versão no portal de releases, em nova guia.
- Se o banco estiver atrás da versão que o sistema exige, aparece um aviso ao lado.
- À direita, uma etiqueta do ambiente: Desenvolvimento, Teste ou Produção (sem etiqueta se o ambiente não estiver definido).
- Endpoint novo `GET /auth/version`, disponível só para quem está logado.

### 2.2 Painel de permissões do usuário ✅
- Na edição de usuário e no resumo, cada permissão aparece sempre com a **descrição completa** ("Excluir empresas", "Visualizar integrações"), na ordem ver, criar, editar, excluir. Antes, módulos com uma só permissão de uma ação mostravam apenas "Ver", "Criar", "Editar" ou "Excluir", misturado com os textos completos.

---

## 3. Correções

### Segurança
- **Ordenação das listagens** ✅: o campo de ordenação enviado pelo cliente (corpo do `/filter` e `sort` na query) não é mais usado diretamente na consulta. Só são aceitas colunas do próprio cadastro (nunca as escondidas da API, como o hash da senha), e qualquer outro valor usa a ordenação padrão. Vale para Clientes, Fornecedores, Usuários, Permissões, logs de sincronização e eventos de webhook.
- **Erros do banco de dados** ✅: o texto de um erro do banco (tipo e tamanho de coluna, código do Postgres, dados de conexão) deixa de aparecer na resposta, em qualquer status. O usuário vê "Não foi possível concluir a operação." e o erro completo fica no log do servidor. Um erro interno (500) também deixa de devolver o texto do erro.

### Validação
- **Tamanho dos campos** ✅: todo campo de texto com limite no banco passa a ser validado na API antes de chegar ao banco, para Clientes, Fornecedores, Empresas, Configuração fiscal, Itens de menu, Perfis e Usuários. As senhas ganham o limite de 72 caracteres. Os formulários das telas existentes (Usuários, Empresas, Configuração fiscal e Perfis) ganharam o mesmo limite, com a mensagem "Máximo N caracteres".
- **Mensagens de validação legíveis** ✅: os erros de validação passam a vir em frases com o nome do campo (`first_name: no máximo 100 caracteres`, `cnpj: campo obrigatório`), em vez do texto interno do validador do Go. JSON malformado, tipo errado e corpo vazio também têm frase própria. As rotas de permissões e de itens de menu, que respondiam num formato à parte nos erros de envio, passam a usar o formato padrão (`message` e `validationErrors`).

### Configuração e catálogo
- **`APP_ENV` não definido** ✅: o sistema sobe em modo seguro (log do banco sem valores), em vez de `development`. O `Backend/.env.example` explica cada opção.
- **Catálogo de permissões** ✅: `roles.*` deixam de se chamar "usuários" (agora "perfis"); `permissions.*` e `integrations.*` em minúsculas como as demais; `finance.receive_boleto` vira "Receber boletos"; `dashboard.view_default` passa para o módulo `Dashboard`.
- **Menu** ✅: o item "Financeiro > Contas a Pagar > Fornecedores" passa a usar `suppliers.view`, como o de Compras (bancos novos; ver a seção 1 para os já semeados).

---

## 4. Infraestrutura

- Nenhuma mudança nos arquivos de deploy (compose, Dockerfiles, Caddy).
- Novo script de dados em `deploy/sql/` (ver seção 1).
- **Ambiente de testes** atualizado para esta versão em 2026-10-29 ✅: código atualizado, imagens reconstruídas e serviços reiniciados, **sem recriar o banco** (os dados foram mantidos). Conferido no servidor: os cinco serviços no ar, a versão do banco 2 registrada com a nova versão do sistema (o histórico guarda as duas linhas), o frontend respondendo e o `GET /auth/version` recusando acesso sem login. A barra de status mostra "Versão 2026.2.3.0 - BD v2" no navegador.
- O `go.mod` do módulo comum passa a declarar o `go-playground/validator` e suas dependências indiretas (nenhuma versão mudou).

---

## 5. Configuração necessária no servidor

- **Reiniciar `api-auth`, `api-core` e `api-integrations`** (e reconstruir o frontend): a mudança está no código dos três serviços e das telas.
- **`APP_ENV`** no `.env`: defina explicitamente (`staging` enquanto o acesso for por HTTP, `production` com HTTPS). Em branco, o sistema sobe em modo seguro, sem o SQL no log.
- **Bancos já semeados:** rodar uma vez o script `deploy/sql/CFG-016_menu_fornecedores_suppliers_view.sql` (idempotente).
- **Frontend (opcional):** `VITE_RELEASES_URL` muda o endereço do portal de releases usado pelo link da versão; sem ela vale o endereço padrão do portal.

---

## 6. Arquivos principais alterados / criados

**Biblioteca comum (`Backend/common`):** `config/config.go`, `utils/response.go`, `utils/response_dberror.go` (novo), `utils/response_validation.go` (novo), `utils/pagination.go`, `utils/pagination_order.go` (novo), `testutil/request_size.go` (novo), `dbversion/dbversion.go`, `go.mod`.
**api.auth:** `internal/domain/` (`user.go`, `role.go`, `menu_item.go`, `version.go` novo), `internal/handlers/` (`auth.go`, `users.go`, `roles.go`, `permissions.go`, `menu_items.go`), `internal/routes/routes.go`, `internal/service/auth_service.go`, `internal/seeder/menu_item_seeder.go`, `seed/permission_catalog.go`.
**api.core:** `internal/domain/` (`customer.go`, `supplier.go`, `company.go`, `company_fiscal_config.go`), `internal/handlers/` (`customers.go`, `suppliers.go`, `companies.go`, `company_fiscal_configs.go`).
**api.integrations:** `internal/handlers/` (`config.go`, `sync_logs.go`, `webhook_events.go`).
**Frontend:** `components/layout/StatusBar.tsx` (novo), `routes/AppLayout.tsx`, `components/users/forms/ModulePermissionsPanel.tsx`, `services/auth/auth-service.ts`, `types/version.ts` (novo), `schemas/` (`create-user-schema.ts`, `update-user-schema.ts`, `company-schema.ts`, `company-fiscal-config-schema.ts`, `role-schema.ts`), `vite-env.d.ts`, `.env.example`.
**Banco, deploy e configuração:** `deploy/sql/CFG-016_menu_fornecedores_suppliers_view.sql` (novo), `Backend/.env.example`.
**Documentação:** o Swagger do gateway foi regerado (novo endpoint e os limites de tamanho dos campos).

---

## 7. Pendências conhecidas (próximas versões)

| Item | Situação |
|---|---|
| **CFG-028** Refatorar o layout principal (cabeçalho, rodapé/barra de status e menu lateral): cores, bordas e alturas ainda não seguem uma regra única | ⏳ |
| **Textos das permissões em bancos já semeados**: os textos e o módulo corrigidos no catálogo só valem em banco novo; nos já semeados é preciso um `UPDATE` à parte | ⏳ |
| **`DB_AUTO_MIGRATE=false`** recusando subir quando o banco está atrás da versão exigida | ⚠️ passou nos testes automáticos, ainda não testado num servidor real |
| **VND-001 / COM-001** Telas de Clientes e de Fornecedores (a API já valida tamanho e formato) | ⏳ |
| **INT-002** Remover do módulo de Integrações a camada de leitura do banco legado herdada da Fase 1 | ⏳ |
| **HTTPS e domínio** no ambiente de testes | ⏳ Hoje o acesso é por IP, em HTTP |
