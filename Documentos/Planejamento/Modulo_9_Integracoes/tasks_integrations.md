# Tasks - api.integrations (Fase 1: Estrutura, Conexões e Gateway)

Acompanhamento da Fase 1 do `api.integrations` (porta 4007), conforme `Documentos/Diversos/Lumini Hub Project Initiation.md` e o plano registrado no `CLAUDE.md`.

## Backend

- [x] [CONCLUÍDO POR: Claude] Esqueleto do módulo (`go.mod`, registro em `go.work`)
- [x] [CONCLUÍDO POR: Claude] Domain Postgres (`IntegrationConfig`, `SyncLog`, `ProductMapping`, `WebhookEvent`) + `AutoMigrate` no `main.go`
- [x] [CONCLUÍDO POR: Claude] Repositório + `UnitOfWork` Postgres
- [x] [CONCLUÍDO POR: Claude] Camada de leitura do legado SQL Server (`LogAltera`, `CadSeq`/`GetSequencia`, `LocArm`, `Cencus`), sem UoW/AutoMigrate
- [x] [CONCLUÍDO POR: Claude] `ConfigService` com cache em memória (Reload/Get/SeedDefaults/UpdateSettings)
- [x] [CONCLUÍDO POR: Claude] Handlers + rotas (`/settings`, `/sync-logs/filter`, `/webhook-events/filter`, `/legacy/locations`, `/legacy/companies`, `/webhooks/loja-integrada`)
- [x] [CONCLUÍDO POR: Claude] Proxy no `api.gateway` + entrada no `run_services.bat`
- [x] [CONCLUÍDO POR: Claude] Build (`go build`) e testes (`go test`, incluindo `sequence_repository_test.go` skip-if-unconfigured) passando
- [x] [CONCLUÍDO POR: Claude] Swagger regenerado incluindo `api.integrations`
- [ ] **Pendente manual (não é código):** criar as permissões `integrations.view` e `integrations.edit` via `POST /api/permissions` (tela Configurações → Permissões ou Swagger UI) — `permissions` é tabela do `api.auth`, não deve ser escrita direto pelo `api.integrations`

## Frontend

- [x] [CONCLUÍDO POR: Claude] Services (`integration-config-service`, `sync-log-service`, `webhook-event-service`, `legacy-lookup-service`) + schemas Zod
- [x] [CONCLUÍDO POR: Claude] Tela Configurações → Integrações (`app/(system)/settings/integrations/page.tsx`): formulário de configurações + selects de local/empresa + tabelas de sync logs e webhook events
- [x] [CONCLUÍDO POR: Claude] Registro em `routes.ts` e `navigation.ts`
- [x] [CONCLUÍDO POR: Claude] `tsc --noEmit` sem novos erros introduzidos pelos arquivos de integrações

## Fase 2 — depende do backlog decomposto abaixo

A Fase 2 (`NOTAS`/`NOTAS1`/`NOTAS3`/`CADPRO`, cliente HTTP da Loja Integrada, scheduler de polling do `LogAltera`, processamento de `WebhookEvent`) foi escopada em `plano_integrations.md` § "🔌 Fase 2 — Integração real com a Loja Integrada (escopado 2026-09-27)" e decomposta em Épico → PBI → Tarefa na seção **"Backlog decomposto (2026-09-27)"** abaixo, que é agora a fonte de verdade pro que falta — não mantemos mais a lista solta que existia aqui antes. O bloqueio de chaves de API/App que travava esta seção foi resolvido (caminho Personal Token, sem fila de aprovação); os bloqueios que restam (política de estoque físico vs. virtual, regra de emissão de nota/despacho da pré-venda, confirmação prática de header vs. query string do webhook) estão marcados inline nas tarefas específicas que eles afetam, não como um bloqueio genérico de fase inteira.

---

## Backlog decomposto (2026-09-27) — esquema novo, Kanban orquestrado

Prefixo `INT` (Módulo 9 — Integrações). Convenção de IDs/tags conforme `.claude/skills/lumini_hub_dev_flow/SKILL.md` seção 2. Status de Épico/PBI é sempre a etapa mais atrasada entre as tarefas-filhas — hoje todos começam `[BACKLOG]`.

### EPIC INT-1: Fase 2 — Integração real com a Loja Integrada

- [ ] [BACKLOG] **EPIC INT-1**: Torna o `api.integrations` funcional de ponta a ponta — autenticação via Personal Token, recepção e processamento de webhooks de Pedido (criação de nota no legado, reserva/devolução de estoque, pré-venda pendente de faturamento), e o sincronismo de produtos/preço do legado pra Loja Integrada via polling do `LogAltera`. Ver `plano_integrations.md` § "🔌 Fase 2" pra todo o levantamento e trade-offs. Duas facetas ficam deliberadamente incompletas dentro deste Épico por dependerem de decisão do cliente/diretoria (não é esquecimento): o valor de estoque enviado à LI (INT-1.5.6) e a automação pós-pré-venda-paga (nota fiscal/despacho, mencionada como observação em INT-1.4.3) — ver `plano_integrations.md` § "Pendências em aberto".

#### PBI INT-1.1: Reformular `IntegrationConfig` e a tela de Configurações → Integrações pra Fase 2

- [ ] [BACKLOG] **PBI INT-1.1**: Sem isso nada mais do Épico roda — Personal Token, vendedor padrão e intervalo de polling são a base de configuração que as demais PBIs leem.

  - [ ] [BACKLOG] INT-1.1.1: Backend — ajustar `IntegrationConfig` (`Backend/microservices/api.integrations/internal/domain/integration_config.go`) e `ConfigService`: renomear `LiApiKey`→`LiPersonalToken`, remover `LiAppKey`, adicionar `LiPersonalTokenIssuedAt` (timestamp); adicionar `CodVendedorPadrao`; remover `CodTipNot` e `CodLocArmOficial` (mantêm `CodLocArmReserva`/`CodEmp` sem mudança); adicionar `PollingIntervalSeconds` (int, default 60). Cobre `UpdateIntegrationSettingsRequest`/`ApiIntegrationSettings`, `SeedDefaults`, `Reload`/`Get`/`UpdateSettings` do `ConfigService`, e os handlers/rotas de `/settings` já existentes (sem rota nova).
  - [ ] [BACKLOG] INT-1.1.2: Backend — cálculo de expiração do Personal Token: expor em `ApiIntegrationSettings` (via `GET /settings`) algo como `li_personal_token_expires_at` + `li_personal_token_expiring_soon` (true quando faltam ≤15 dias pros 90 desde `LiPersonalTokenIssuedAt`). (depende de INT-1.1.1)
  - [ ] [BACKLOG] INT-1.1.3: Backend — endpoint de lookup de vendedores do legado (padrão de `/legacy/locations`/`/legacy/companies` já existentes em `internal/handlers`/`internal/routes`), pra popular o select de `CodVendedorPadrao` na tela. Nome exato da tabela/view do legado (provavelmente algo como `VENDEDOR`/`CADVEN` — não confirmado no plano) é detalhe técnico a confirmar durante a implementação, não pendência de negócio. (depende de INT-1.1.1)
  - [ ] [BACKLOG] INT-1.1.4: Frontend — atualizar tela Configurações → Integrações (`integration-config-service`, schema Zod, formulário): campo Personal Token renomeado, campo App Key removido, novo select de "Vendedor Padrão" (consumindo INT-1.1.3) no lugar dos campos removidos (`codtipnot`/`codlocarm_oficial`), novo campo numérico de intervalo de polling, e alerta visual quando `li_personal_token_expiring_soon` vier `true`. (depende de INT-1.1.1, INT-1.1.2, INT-1.1.3)

#### PBI INT-1.2: Cliente HTTP da Loja Integrada

- [ ] [BACKLOG] **PBI INT-1.2**: Base de comunicação de saída com a LI, reaproveitada tanto pelo sincronismo de produtos (INT-1.5) quanto por qualquer necessidade futura de buscar detalhe de pedido a partir de um webhook incompleto.

  - [ ] [BACKLOG] INT-1.2.1: Cliente HTTP base com autenticação Personal Token (`Authorization: Basic <token>`, lido de `IntegrationConfig.LiPersonalToken` via `ConfigService`) — tratamento de erro HTTP padronizado, sem retry embutido aqui (retry é responsabilidade do scheduler, INT-1.5.5). Endpoints exatos levantados via doc de terceiro (`plano_integrations.md` § "Endpoints REST relevantes") — a confirmar/ajustar na prática. (depende de INT-1.1.1)
  - [ ] [BACKLOG] INT-1.2.2: Métodos de produto/preço: `GET/POST/PUT /produto`, `PUT /produto_preco`. Método de `PUT /produto_estoque` pode ser escrito (é só um método de cliente HTTP), mas **não deve ser chamado por nenhum código até INT-1.5.6 ser desbloqueada** — deixar claro no código/comentário que existe só como cliente, sem uso ainda. (depende de INT-1.2.1)
  - [ ] [BACKLOG] INT-1.2.3: Métodos de pedido: `GET/PUT /pedido` — reservado pra caso o payload do webhook não venha completo o suficiente pra montar a nota (a confirmar durante INT-1.4.2). (depende de INT-1.2.1)

#### PBI INT-1.3: Domain legado `NOTAS`/`NOTAS1`/`NOTAS3` + criação de nota de venda

- [ ] [BACKLOG] **PBI INT-1.3**: A parte que efetivamente grava no SQL Server legado — sem UoW/AutoMigrate contra o schema legado (mesmo princípio já usado em `LogAltera`/`CadSeq`), escrita em transação explícita com `UPDLOCK`, espelhando `GetSequencia`.

  - [ ] [BACKLOG] INT-1.3.1: Domain de leitura/escrita do legado para `NOTAS`/`NOTAS1`/`NOTAS3` (structs Go, sem GORM/AutoMigrate — mesmo padrão dos structs de leitura do legado já existentes) + repositório com geração de `nroentsai` via `CADSEQ`/`GetSequencia` (select → incrementa → update, dentro de transação com `UPDLOCK`, igual ao legado faz).
  - [ ] [BACKLOG] INT-1.3.2: De-para de forma de pagamento — leitura (somente leitura, nunca escrita) do campo de vínculo entre a forma de pagamento da LI e a forma de pagamento do legado, dentro da própria tabela de formas de pagamento do legado. **Pré-condição operacional, não é código**: o campo já precisa existir no legado antes desta tarefa começar — é responsabilidade do usuário/cliente garantir isso (criar se não existir), conforme `plano_integrations.md` § "Configurações necessárias pra gerar a nota no legado". Confirmar com o usuário que o campo já existe antes de iniciar esta tarefa.
  - [ ] [BACKLOG] INT-1.3.3: Serviço de criação de `NOTAS`/`NOTAS1`/`NOTAS3` a partir de um pedido da LI: resolve o vendedor padrão (`CodVendedorPadrao`, cujo cadastro já traz tipo de nota/centro de custo/tabela de preços prontos), monta os itens (`NOTAS1`), monta pagamento/parcela (`NOTAS3`) — incluindo o caso de pagamento parcelado, que vira a forma de pagamento "Parcelado LI" com um único título a receber (a forma de pagamento "Parcelado LI" em si precisa existir cadastrada no legado, mesma pré-condição operacional da tarefa anterior). (depende de INT-1.1.1, INT-1.3.1, INT-1.3.2)

#### PBI INT-1.4: Processamento do `WebhookEvent` de Pedido por status

- [ ] [BACKLOG] **PBI INT-1.4**: Reage ao pedido já recebido e persistido como `WebhookEvent` (mecanismo da Fase 1, `POST /integrations/webhooks/loja-integrada`) — cada status listado em `plano_integrations.md` § "Regras de negócio — reação por status do pedido" vira uma tarefa própria.
  > Nota (não é tarefa): a assinatura do webhook de Pedido (Criado+Editado) no painel da própria Loja Integrada é uma ação manual do usuário, não código — feita quando ele testar o cadastro do webhook na prática (ver INT-1.4.1).

  - [ ] [BACKLOG] INT-1.4.1: Ajuste técnico do segredo do webhook em `webhooks.go` (`Backend/microservices/api.integrations/internal/handlers/webhooks.go`) — hoje só lê `X-Webhook-Secret` do header. Ajustar pra aceitar **também** `?secret=` na query string, já que não há confirmação de que a LI permite header customizado ao cadastrar a URL de callback; o usuário decide qual usar de fato ao testar o cadastro na prática (`plano_integrations.md` § "Segredo do webhook").
  - [ ] [BACKLOG] INT-1.4.2: Processar evento "Pedido Criado": cria a nota no legado (via INT-1.3.3) e reserva estoque Oficial → `CodLocArmReserva`. Marca o `WebhookEvent` como `processed` (ou `error` com `ErrorMessage` preenchido em caso de falha). (depende de INT-1.3.3)
  - [ ] [BACKLOG] INT-1.4.3: Processar "Pedido Editado" com status **Pago**: marca a pré-venda como pendente de faturamento no Caixa, reaproveitando literalmente o mesmo mecanismo de uma venda física deixada sem faturar (não cria status novo). Não fatura automaticamente, não mexe em estoque além do já reservado em INT-1.4.2. **Nota**: o que acontece depois disso (emissão de nota fiscal/despacho automático a partir dessa pré-venda) é a pendência explícita em aberto no plano — esta tarefa cobre só marcar como pendente, não a automação seguinte, que fica fora de escopo até o usuário alinhar com o cliente. (depende de INT-1.4.2)
  - [ ] [BACKLOG] INT-1.4.4: Processar "Pedido Editado" com status **Cancelado**: devolve a quantidade reservada de volta pro estoque Oficial. (depende de INT-1.4.2)
  - [ ] [BACKLOG] INT-1.4.5: Processar "Pedido Editado" com status **Enviado/Entregue/Devolvido**: apenas atualiza o `WebhookEvent` (log/status), sem nenhuma ação automática no legado — a movimentação real de estoque pro cliente e a devolução física/fiscal continuam manuais (ver `plano_integrations.md`).

#### PBI INT-1.5: Sincronismo de Produtos (Legado → Loja Integrada)

- [ ] [BACKLOG] **PBI INT-1.5**: Scheduler de polling do `LogAltera` que empurra produto/preço pra LI, com criação automática de produto novo e retry+log de erro. A parte de estoque fica deliberadamente de fora (INT-1.5.6, bloqueada).

  - [ ] [BACKLOG] INT-1.5.1: Infraestrutura do scheduler — ciclo de polling configurável via `PollingIntervalSeconds` (start/stop, cursor de leitura do `LogAltera` por `DatFimAlt`, confirmar durante a implementação se reaproveita o mesmo cursor já usado na Fase 1 pra outra finalidade ou precisa de um cursor próprio pra esta direção), sem lógica de push ainda — só o esqueleto do loop + registro de execução via `SyncLog`. (depende de INT-1.1.1)
  - [ ] [BACKLOG] INT-1.5.2: Push de nome/descrição/categoria (`LogAltera.TipoAlt='P'`, confirmar durante a implementação se é de fato esse o tipo usado pra alteração cadastral) via `/produto` — inclui a detecção de produto novo (`CADPRO.codint` vazio) que cria o produto na LI e grava o `codint` retornado de volta em `CADPRO` (mecanismo de detecção exato — via `LogAltera` cruzado com `codint` vazio, ou varredura periódica separada — é detalhe de implementação a refinar, não bloqueia o escopo). (depende de INT-1.2.2, INT-1.5.1)
  - [ ] [BACKLOG] INT-1.5.3: Push de preço de venda (`LogAltera.TipoAlt='V'`) via `/produto_preco`. Custo (`TipoAlt='U'`) explicitamente **não** sincroniza (informação interna, não exposta na loja virtual). (depende de INT-1.2.2, INT-1.5.1)
  - [ ] [BACKLOG] INT-1.5.4: Push de imagens e variações/grade — endpoints exatos (`/grade` pra variação; imagem não confirmada na fonte usada) a confirmar durante a implementação, ajuste técnico, não pendência de negócio. (depende de INT-1.2.1, INT-1.5.1)
  - [ ] [BACKLOG] INT-1.5.5: Retry automático + log de erro persistente: qualquer falha de push (produto, preço, imagem/variação) fica pra retry na próxima rodada de polling, com erro visível na tela Configurações → Integrações via `SyncLog` (mecanismo já existente desde a Fase 1). (depende de INT-1.5.1)
  - [ ] [BACKLOG] INT-1.5.6: Push de estoque via `/produto_estoque` (`LogAltera.TipoAlt` a confirmar). **BLOQUEADA — não iniciar**: depende da decisão do cliente/diretoria sobre a política de prioridade de estoque físico vs. virtual, incluindo qual valor enviar pra LI (Oficial cru vs. colchão de segurança) — pendência explícita registrada em `plano_integrations.md` § "Pendências em aberto". Reavaliar assim que essa decisão existir.
