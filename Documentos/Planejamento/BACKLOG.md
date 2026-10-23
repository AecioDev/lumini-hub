# Backlog - Lumini Hub

Painel do que já foi feito e do que falta. Quem conduz é o Aécio; o Claude só atualiza quando termina algo que ele pediu.

## Como usar

- `[ ]` = em aberto (ainda não feito, **ou** feito mas ainda não testado por você).
- `[x]` = feito **e testado por você**. Só você move para `[x]`; o Claude nunca marca sozinho.
- Quando algo está implementado mas ainda sem o seu teste, continua `[ ]` e leva a nota `(implementado, aguardando teste)`.
- O que foi feito desde o último envio ao servidor fica listado em `PROXIMA_VERSAO.md`. Ao enviar, o Claude gera `releases/<versão>.md` a partir dele e zera o arquivo.
- Versão: `ano.versao_do_banco.sequencia_versao.sequencia_patch`. O segundo número é a versão do banco: sobe uma unidade por publicação que traz mudanças de banco (não por mudança individual), e cada release lista o que mudou nele. Próxima: `2026.2.2.0`.
- Os `tasks_*.md` e `plano_*.md` antigos continuam como histórico e consulta. Não são mais atualizados (o escopo e as regras de negócio de cada módulo continuam nos `plano_*.md`).

---

## 1. Aguardando o seu teste

Já implementado e na `develop`; falta você validar.

- [ ] Remoção de `deploy/db/01_schema.sql` e `02_catalog.sql` e do mount `./db` no `docker-compose.prod.yml`: os serviços passam a criar tudo sozinhos - CFG-9.5.3 *(implementado, aguardando confirmação sua)*

## 2. Próximos (Configuração / base)

- [ ] Descrições e módulos do catálogo de permissões herdados do dev: `roles.*` dizem "usuários", módulo "dashboard" em minúsculo; sem renomear os códigos - CFG-10.2
- [ ] Conferir os itens de menu do `menu_item_seeder.go` contra a tabela `menu_items` do dev (só `/settings/*` foi conferido) - CFG-10.4
- [ ] Bancos já semeados (dev e VM da OCI): o item "Fornecedores" continua em `purchases.view`; precisa de ajuste manual ou migração, porque o seed do menu pula itens existentes
- [ ] `APP_ENV` não definido vira `development` (log verboso): exigir explícito fora de dev ou inverter o padrão para o modo seguro - CFG-10.5.1
- [ ] Nunca logar `err` bruto nos handlers (mensagem do Postgres pode citar valores) - CFG-10.5.2
- [ ] Dívida técnica: `CompanyService.GetCompanies` usa `GetDB()` cru, fora do `UoW.Execute()`
- [ ] Gravar a versão do banco dentro do próprio banco: tabela `schema_versions` (`db_version`, `app_version`, `applied_at`, nota), uma linha por publicação; constante com a versão exigida pelo código num lugar só (`common`), subida à mão nas publicações com mudança de banco; só o `api.auth` grava (depois do `AutoMigrate` e do seed); no boot, com `DB_AUTO_MIGRATE=false` os serviços só conferem e recusam subir se o banco estiver atrás. A primeira versão gravada seria a 2, junto com a publicação `2026.2.2.0`. Limite: o `AutoMigrate` só adiciona; mudanças que mexem em dados existentes precisam de script próprio. Incluir "mudou o banco? subiu a versão?" no roteiro de release
- [ ] Atualizar a VM da OCI para a versão nova (o deploy lá é antigo; para atualizar é preciso recriar o volume do Postgres)
- [ ] Traduzir o `CLAUDE.md` e as skills para português
- [ ] Arquivar o histórico dos `tasks_*.md` grandes (`tasks_empresa.md` tem ~330 KB)

## 3. Módulos ainda não iniciados

Escopo de cada um em `Modulo_N_*/plano_*.md`. Nenhum tem código hoje; no menu caem em `PlaceholderPage`.

- [ ] **Clientes e Fornecedores - telas** (o backend e os services do frontend existem em `api.core`; faltam as telas de cadastro/edição no padrão de página cheia)
- [ ] **M1 CRM** (`api.crm`, porta 4009): Leads, Pipeline Kanban, Oportunidades, Atividades, Pós-Venda, Churn. Hoje o Dashboard CRM é só mock
- [ ] **M2 Produtos e Estoque** (`api.inventory`, 4003): Produto, Marca, Categoria, Depósito, saldo e movimentação, fotos, endereçamento físico
- [ ] **M3 Compras** (`api.purchases`, 4004): pedidos de compra, entrada de notas, Manifesto SEFAZ, leitura de XML
- [ ] **M4 Vendas e Caixas** (`api.sales` 4005 / `api.financial` 4006): venda de balcão por ambientes, projetos CAD/PDF, faturamento parcial, controle de caixas
- [ ] **M5 Financeiro** (`api.financial`, 4006): contas a pagar/receber, bancos, fluxo de caixa, conciliação, boletos (ACBrLib)
- [ ] **M6 Fiscal** (`api.fiscal`, 4008): tributação, certificados A1, NFe/NFSe/NFCe/MDFe via ACBrLib
- [ ] **M7 Contabilidade**: plano de contas, partida dobrada, Razão, Balancete, DRE. *Bloqueado por decisão: falta o plano de contas padrão (lista de contas)*
- [ ] **M8 Multi-tenant**: banco físico separado por cliente. Só perguntas em aberto; revisitar quando houver um segundo cliente

## 4. Pausado

- [ ] **M9 Integrações (Loja Integrada ↔ SQL Server legado)** - pausado em 01/10/2026: o acesso ao legado vira uma aplicação separada no servidor do cliente. Fase 1 está no código; a Fase 2 (EPIC INT-1, 28 tarefas) está em `Modulo_9_Integracoes/tasks_integrations.md`. Depende de decisões do cliente (prioridade de estoque físico x virtual, regra de nota/despacho da pré-venda)

---

## 5. Feito e testado

- [x] Autenticação com cookies HTTP-only, refresh de token e `/auth/me`
- [x] Usuários (lista, cadastro e edição em página cheia, com permissões individuais)
- [x] Perfis e permissões (Roles), RBAC DEVELOP / ADMIN
- [x] Menu lateral dinâmico vindo do banco (`menu_items`), filtrado por permissão
- [x] Empresas: matriz e filiais, configuração fiscal, certificado digital A1, logo (CFG-1/4/5/6/7)
- [x] Banco vazio sobe só com o código: criação das tabelas, catálogo de permissões, role ADMIN e admin inicial via `BOOTSTRAP_ADMIN_*` (CFG-9; cenários A a D testados em 03/10/2026)
- [x] Deploy em Docker + Caddy na VM da OCI (versão antiga em produção)
- [x] Gateway único com CORS centralizado; Swagger gerado
- [x] Cadastro da primeira empresa num banco sem empresas, direto na tela de seleção (CFG-8; testado)
- [x] Roteiro de primeiro acesso no `deploy/README.md`: clonar, `.env`, subir, login do admin, cadastrar empresa (CFG-8.2.1; testado)
- [x] Logger do GORM por ambiente: fora de `development` o log não mostra SQL com valores nem hash de senha (CFG-10.3.1; testado em 22/10/2026 com `APP_ENV=staging`)
- [x] 8 permissões novas no catálogo (`customers.create/edit/delete`, `suppliers.view/create/edit/delete`, `integrations.edit`), 67 no total e 58 para o ADMIN (CFG-10.1; testado em 22/10/2026 na tela de permissões do usuário). O menu "Fornecedores" ligado a `suppliers.view` só vale para bancos novos; o ajuste dos bancos já semeados está em "Próximos"
