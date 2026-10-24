# Backlog - Lumini Hub

Painel do que já foi feito e do que falta. Quem conduz é o Aécio; o Claude só atualiza quando termina algo que ele pediu.

## Como usar

- `[ ]` = em aberto (ainda não feito, **ou** feito mas ainda não testado por você).
- `[x]` = feito **e testado por você**. Só você move para `[x]`; o Claude nunca marca sozinho.
- Quando algo está implementado mas ainda sem o seu teste, continua `[ ]` e leva a nota `(implementado, aguardando teste)`.
- **ID de tarefa**: todo item começa com um ID `PREFIXO-NNN` em negrito (ex.: **CFG-014**). O prefixo é o módulo (tabela abaixo); o número é sequencial dentro do prefixo, tem 3 dígitos e **nunca é reaproveitado**, nem quando o item é concluído ou cancelado (ex.: o **CFG-019** foi descartado e o número ficou sem uso). Para achar o próximo ID livre de um prefixo, procure o maior número dele neste arquivo e some 1. Para buscar uma tarefa, procure o ID (ex.: `CFG-014`) no repositório, no `PROXIMA_VERSAO.md` e nos releases.
- O que foi feito desde o último envio ao servidor fica listado em `PROXIMA_VERSAO.md`. Ao enviar, o Claude gera `releases/Release_<versão>.md` a partir dele e zera o arquivo.
- Versão: `ano.versao_do_banco.sequencia_versao.sequencia_patch`. O segundo número é a versão do banco: sobe uma unidade por publicação que traz mudanças de banco (não por mudança individual), e cada release lista o que mudou nele. A versão atual não fica escrita aqui: veja `Backend/common/dbversion/dbversion.go`, `PROXIMA_VERSAO.md` e o último arquivo em `releases/`.
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

- [ ] **CFG-016** Bancos já semeados (dev e VM da OCI): os itens de menu "Fornecedores" (em Compras e em Financeiro > Contas a Pagar) continuam em `purchases.view`; precisa de ajuste manual ou migração, porque o seed do menu pula itens existentes
- [ ] **CFG-027** Revisar os cerca de 25 casos de resposta 400 genérica nos handlers (ex.: "Erro ao criar usuário" depois de uma falha que não é de validação), que podem devolver ao usuário texto de erro do banco (nome de tabela ou constraint). O 500 já está coberto pelo CFG-018. Primeiro testar à mão como o sistema se comporta (provocar uma constraint violada, como e-mail duplicado, e ver o que aparece na tela), depois decidir o tratamento por entidade
- [ ] **CFG-023** Arquivar o histórico dos `tasks_*.md` grandes (`tasks_empresa.md` tem ~330 KB)
- [ ] **CFG-024** Barra de status no layout do sistema (frontend), numa área fixa da página, mostrando a versão do sistema e a versão do banco (precisa de um endpoint que leia `dbversion.AppVersion` e a tabela `schema_versions`, ex.: `/api/auth/version`) e outras informações que forem surgindo com o tempo (ambiente, empresa ativa, etc.). Ligada ao CFG-020
- [ ] **CFG-025** Scripts de migração por versão do banco, para o que o `AutoMigrate` não resolve (remover ou renomear coluna, ajustar dados existentes, triggers/funções, `NOT NULL` novo em tabela com linhas). Cada script leva o banco de uma versão à seguinte (`002_...`, `003_...`) e roda só os acima do número gravado em `schema_versions.db_version`, para atualizar um banco muito antigo em ordem. Só fazer quando aparecer a primeira mudança desse tipo. Ligada ao CFG-020

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
- [x] **CFG-020** Versão do banco gravada no próprio banco: tabela `schema_versions` (uma linha por publicação), constantes `Required` e `AppVersion` em `Backend/common/dbversion/dbversion.go` (subidas à mão a cada release), o `api.auth` registra no fim do boot com `DB_AUTO_MIGRATE=true`, e com `DB_AUTO_MIGRATE=false` os três serviços conferem e recusam subir se o banco estiver atrás ou sem versão (testado em 22/10/2026: registrou a versão 2 no primeiro boot e não duplicou ao reiniciar)
- [x] **CFG-021** Ambiente da OCI atualizado para a `2026.2.2.0`: código na `develop`, banco recriado do zero, admin inicial, primeira empresa cadastrada e permissões conferidas (testado em 22/10/2026). Release em `releases/Release_2026.2.2.0.md`
- [x] **CFG-014** Textos e módulos do catálogo de permissões corrigidos no código (`roles.*` = perfis, `permissions.*`/`integrations.*` em minúsculas, "Receber boletos", módulo `Dashboard`); `roles.*` fica no módulo Develop de propósito (antigo CFG-10.2; testado em 23/10/2026 num banco de dev novo)
- [x] **CFG-026** Painel de permissões do usuário (editar usuário e o resumo): todo item mostra a descrição completa da permissão, na ordem ver, criar, editar, excluir e demais, no lugar dos rótulos curtos "Ver"/"Criar"/"Editar"/"Excluir" que misturavam com os completos (testado em 23/10/2026)
- [x] **CFG-015** Menu do seeder conferido contra o menu do banco de dev antigo (`erp_system`): hrefs, ícones e permissões iguais, exceto Fornecedores em Compras (já tratado) e uma posição repetida em Configurações. Achado e corrigido: "Financeiro > Contas a Pagar > Fornecedores" ainda usava `purchases.view`; agora usa `suppliers.view`. Novos testes garantem que toda permissão do menu existe no catálogo, que todo item com tela tem permissão e que Fornecedores usa `suppliers.view` (antigo CFG-10.4) (testado em 23/10/2026)
- [x] **CFG-017** `APP_ENV` não definido deixa de virar `development`: sem a variável o sistema sobe em modo seguro (log do banco sem valores) e para ver o SQL no log é preciso `APP_ENV=development` no `.env`. As opções (development, staging, production, em branco) e o que cada uma faz estão explicadas no `Backend/.env.example`, no lugar de um aviso no console (antigo CFG-10.5.1) (testado em 23/10/2026)
- [x] **CFG-018** Erro interno (500) não devolve mais o texto do erro ao cliente: `utils.ErrorResponse` grava o erro completo no log do servidor (método, caminho e erro, sem a query string) e a resposta leva só a `message`. Cobre os 40 pontos dos handlers que passavam o erro cru (o diagnóstico era de resposta, não de log). Ficam de fora os 25 casos de 400 genérico, que podem carregar texto de banco e precisam de revisão por entidade (antigo CFG-10.5.2) (testado em 23/10/2026: aceito pelos testes automáticos)
- [x] **CFG-022** `CLAUDE.md` e as skills `lumini_hub_backend_architecture` e `lumini_hub_entity_creation` traduzidos para português (a `lumini_hub_dev_flow` já estava; a `agent-browser` é documentação de terceiros e ficou como veio). Aproveitei para corrigir três pontos desatualizados das skills: o nome das permissões e o lugar delas (agora o `PermissionCatalog`, e não seed solto), o exemplo de `ErrorResponse` (500 não devolve o erro) e o status pausado do `api.integrations` Na mesma passada: o `CLAUDE.md` deixou de descrever o banco legado (o módulo de Integrações é do sistema; o acesso ao legado é uma aplicação separada) e o login local passou a depender do `BOOTSTRAP_ADMIN_PASSWORD`. (testado em 23/10/2026)
