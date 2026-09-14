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

## Fora do escopo desta fase (Fase 2 — depende de pendências)

- Domain/repositório de `NOTAS`/`NOTAS1`/`NOTAS3`/`CADPRO` e lógica real de criação de nota de venda a partir de pedido da LI
- Cliente HTTP da Loja Integrada (produtos/pedidos)
- Scheduler de polling que efetivamente empurra mudanças do `LogAltera` para a LI
- Processamento do `WebhookEvent` recebido (transformar em pedido de venda)

**Pendências que bloqueiam a Fase 2:** chaves de API/App da Loja Integrada, valor de `codtipnot`, política de prioridade de estoque loja física vs. virtual (a confirmar com a diretoria do cliente).
