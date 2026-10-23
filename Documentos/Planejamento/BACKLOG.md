# Backlog - Lumini Hub

Painel do que já foi feito e do que falta. Quem conduz é o Aécio; o Claude só atualiza quando termina algo que ele pediu.

## Como usar

- `[ ]` = em aberto (ainda não feito, **ou** feito mas ainda não testado por você).
- `[x]` = feito **e testado por você**. Só você move para `[x]`; o Claude nunca marca sozinho.
- Quando algo está implementado mas ainda sem o seu teste, continua `[ ]` e leva a nota `(implementado, aguardando teste)`.
- **ID de tarefa**: todo item começa com um ID `PREFIXO-NNN` em negrito (ex.: **CFG-014**). O prefixo é o módulo (tabela abaixo); o número é sequencial dentro do prefixo, tem 3 dígitos e **nunca é reaproveitado**, nem quando o item é concluído ou cancelado. Para achar o próximo ID livre de um prefixo, procure o maior número dele neste arquivo e some 1. Para buscar uma tarefa, procure o ID (ex.: `CFG-014`) no repositório, no `PROXIMA_VERSAO.md` e nos releases.
- O que foi feito desde o último envio ao servidor fica listado em `PROXIMA_VERSAO.md`. Ao enviar, o Claude gera `releases/<versão>.md` a partir dele e zera o arquivo.
- Versão: `ano.versao_do_banco.sequencia_versao.sequencia_patch`. O segundo número é a versão do banco: sobe uma unidade por publicação que traz mudanças de banco (não por mudança individual), e cada release lista o que mudou nele. Próxima: `2026.2.2.0`.
- Os `tasks_*.md` e `plano_*.md` antigos continuam como histórico e consulta. Não são mais atualizados (o escopo e as regras de negócio de cada módulo continuam nos `plano_*.md`). Os códigos antigos desses arquivos (ex.: `CFG-10.2`) aparecem aqui como "antigo ...", só para você achar a história.

## Prefixos por módulo

| Prefixo | Módulo |
|---|---|
| `CFG` | 0 Configuração (empresa, usuários, perfis, banco, deploy, fluxo de trabalho) |
| `CRM` | 1 CRM |
| `EST` | 2 Produtos e Estoque |
| `COM` | 3 Compras (inclui cadastro de fornecedores) |
| `VND` | 4 Vendas e Caixas (inclui cadastro de clientes) |
| `FIN` | 5 Financeiro |
| `FIS` | 6 Fiscal |
| `CTB` | 7 Contabilidade |
| `MT` | 8 Multi-tenant |
| `INT` | 9 Integrações |

---

## 1. Aguardando o seu teste

Já implementado e na `develop`; falta você validar.

_Nada aguardando teste no momento._

## 2. Próximos (Configuração / base)

- [ ] **CFG-014** Descrições e módulos do catálogo de permissões herdados do dev: `roles.*` dizem "usuários", módulo "dashboard" em minúsculo, "Ver"/"Editar" em Integrações; sem renomear os códigos (antigo CFG-10.2)
- [ ] **CFG-015** Conferir os itens de menu do `menu_item_seeder.go` contra a tabela `menu_items` do dev (só `/settings/*` foi conferido) (antigo CFG-10.4)
- [ ] **CFG-016** Bancos já semeados (dev e VM da OCI): o item de menu "Fornecedores" continua em `purchases.view`; precisa de ajuste manual ou migração, porque o seed do menu pula itens existentes
- [ ] **CFG-017** `APP_ENV` não definido vira `development` (log verboso): exigir explícito fora de dev ou inverter o padrão para o modo seguro (antigo CFG-10.5.1)
- [ ] **CFG-018** Nunca logar `err` bruto nos handlers (mensagem do Postgres pode citar valores) (antigo CFG-10.5.2)
- [ ] **CFG-019** Dívida técnica: `CompanyService.GetCompanies` usa `GetDB()` cru, fora do `UoW.Execute()`
- [ ] **CFG-020** Gravar a versão do banco dentro do próprio banco: tabela `schema_versions` (`db_version`, `app_version`, `applied_at`, nota), uma linha por publicação; constante com a versão exigida pelo código num lugar só (`common`), subida à mão nas publicações com mudança de banco; só o `api.auth` grava (depois do `AutoMigrate` e do seed); no boot, com `DB_AUTO_MIGRATE=false` os serviços só conferem e recusam subir se o banco estiver atrás. A primeira versão gravada seria a 2, junto com a publicação `2026.2.2.0`. Limite: o `AutoMigrate` só adiciona; mudanças que mexem em dados existentes precisam de script próprio. Incluir "mudou o banco? subiu a versão?" no roteiro de release
- [ ] **CFG-021** Atualizar a VM da OCI para a versão nova (o deploy lá é antigo; para atualizar é preciso recriar o volume do Postgres). Conferir lá que o `docker-compose.prod.yml` novo sobe sem o `deploy/db`
- [ ] **CFG-022** Traduzir o `CLAUDE.md` e as skills para português
- [ ] **CFG-023** Arquivar o histórico dos `tasks_*.md` grandes (`tasks_empresa.md` tem ~330 KB)

## 3. Módulos ainda não iniciados

O Lumini Hub é um sistema feito do zero. Escopo de cada módulo em `Modulo_N_*/plano_*.md`. Nenhum tem código hoje; no menu caem em `PlaceholderPage`.

- [ ] **VND-001** Tela de Clientes: lista, cadastro e edição em página cheia (o backend e o service do frontend já existem em `api.core`)
- [ ] **COM-001** Tela de Fornecedores: lista, cadastro e edição em página cheia (o backend e o service do frontend já existem em `api.core`)
- [ ] **CRM-001** CRM (`api.crm`, porta 4009): Leads, Pipeline Kanban, Oportunidades, Atividades, Pós-Venda, Churn. Hoje o Dashboard CRM é só mock
- [ ] **EST-001** Produtos e Estoque (`api.inventory`, 4003): Produto, Marca, Categoria, Depósito, saldo e movimentação, fotos, endereçamento físico
- [ ] **COM-002** Compras (`api.purchases`, 4004): pedidos de compra, entrada de notas, Manifesto SEFAZ, leitura de XML
- [ ] **VND-002** Vendas e Caixas (`api.sales` 4005 / `api.financial` 4006): venda de balcão por ambientes, projetos CAD/PDF, faturamento parcial, controle de caixas
- [ ] **FIN-001** Financeiro (`api.financial`, 4006): contas a pagar/receber, bancos, fluxo de caixa, conciliação, boletos (ACBrLib)
- [ ] **FIS-001** Fiscal (`api.fiscal`, 4008): tributação, certificados A1, NFe/NFSe/NFCe/MDFe via ACBrLib
- [ ] **CTB-001** Contabilidade: plano de contas, partida dobrada, Razão, Balancete, DRE. *Bloqueado por decisão: falta o plano de contas padrão (lista de contas)*
- [ ] **MT-001** Multi-tenant: banco físico separado por cliente. Só perguntas em aberto; revisitar quando houver um segundo cliente

## 4. Pausado

- [ ] **INT-001** Integrações com a Loja Integrada e com o legado: pausado em 01/10/2026. O acesso ao SQL Server do cliente vai virar uma aplicação separada, no servidor dele, fora deste sistema. Fase 1 está no código e não evolui; a Fase 2 (antigo EPIC INT-1, 28 tarefas) está em `Modulo_9_Integracoes/tasks_integrations.md`. Depende de decisões do cliente (prioridade de estoque físico x virtual, regra de nota/despacho da pré-venda)

---

## 5. Feito e testado

- [x] **CFG-001** Autenticação com cookies HTTP-only, refresh de token e `/auth/me`
- [x] **CFG-002** Usuários (lista, cadastro e edição em página cheia, com permissões individuais)
- [x] **CFG-003** Perfis e permissões (Roles), RBAC DEVELOP / ADMIN
- [x] **CFG-004** Menu lateral dinâmico vindo do banco (`menu_items`), filtrado por permissão
- [x] **CFG-005** Empresas: matriz e filiais, configuração fiscal, certificado digital A1, logo (antigos CFG-1/4/5/6/7)
- [x] **CFG-006** Banco vazio sobe só com o código: criação das tabelas, catálogo de permissões, role ADMIN e admin inicial via `BOOTSTRAP_ADMIN_*` (antigo CFG-9; cenários A a D testados em 03/10/2026)
- [x] **CFG-007** Deploy em Docker + Caddy na VM da OCI (versão antiga em produção)
- [x] **CFG-008** Gateway único com CORS centralizado; Swagger gerado
- [x] **CFG-009** Cadastro da primeira empresa num banco sem empresas, direto na tela de seleção (antigo CFG-8; testado)
- [x] **CFG-010** Roteiro de primeiro acesso no `deploy/README.md`: clonar, `.env`, subir, login do admin, cadastrar empresa (antigo CFG-8.2.1; testado)
- [x] **CFG-011** Remoção de `deploy/db/*.sql` e do mount `./db` do `docker-compose.prod.yml`: o banco vazio é criado só pelos serviços (antigo CFG-9.5.3; marcado como testado em 22/10/2026; a conferência na VM da OCI está em CFG-021)
- [x] **CFG-012** Logger do GORM por ambiente: fora de `development` o log não mostra SQL com valores nem hash de senha (antigo CFG-10.3.1; testado em 22/10/2026 com `APP_ENV=staging`)
- [x] **CFG-013** 8 permissões novas no catálogo (`customers.create/edit/delete`, `suppliers.view/create/edit/delete`, `integrations.edit`), 67 no total e 58 para o ADMIN (antigo CFG-10.1; testado em 22/10/2026 na tela de permissões do usuário). O menu "Fornecedores" ligado a `suppliers.view` só vale para bancos novos; o ajuste dos bancos já semeados está em CFG-016
