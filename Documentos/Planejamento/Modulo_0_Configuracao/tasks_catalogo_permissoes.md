# Tasks - Catálogo de Permissões (completar e limpar)

Criado em 2026-10-01 pelo orquestrador, por decisão do usuário, a partir de achados do EPIC CFG-9 (`tasks_schema_seed.md`). Prefixo `CFG`, esquema novo de Kanban. Status de Épico/PBI é sempre a etapa mais atrasada entre as tarefas-filhas. Tudo nasce `[BACKLOG]`. Fora do escopo do EPIC CFG-9 de propósito.

### EPIC CFG-10: Completar e limpar o catálogo de permissões

- [ ] [EM-ANDAMENTO] **EPIC CFG-10**: O catálogo de permissões em código (`api.auth/seed`, entregue na CFG-9.3.1 com 59 permissões herdadas do banco de dev) passa a cobrir tudo o que `RequirePermission` exige e fica com descrições/módulos consistentes. **Até lá, as rotas que exigem as permissões faltantes só funcionam para ADMIN/DEVELOP (bypass).** Escopo ainda não escopado em `plano_*.md`; as decisões de nome/módulo/descrição são do usuário.

#### PBI CFG-10.1: Cadastrar as 8 permissões exigidas pelo código e ausentes do catálogo

- [ ] [BACKLOG] **PBI CFG-10.1**: Toda permissão exigida por `RequirePermission` existe no catálogo.

  - [ ] [BACKLOG] CFG-10.1.1: Decidir com o usuário nome, módulo e descrição das 8 permissões ausentes: `customers.create`, `customers.edit`, `customers.delete`, `suppliers.view`, `suppliers.create`, `suppliers.edit`, `suppliers.delete`, `integrations.edit` (conferir por grep em `RequirePermission(` se há outras).
  - [ ] [BACKLOG] CFG-10.1.2: Adicionar as permissões decididas ao catálogo em `api.auth/seed`; atualizar as contagens dos testes (hoje 59 permissões / 50 vínculos ao ADMIN; mudam quando entrarem, a regra do ADMIN é `Module != "Develop"`). (depende de CFG-10.1.1, CFG-9.3.4)

#### PBI CFG-10.2: Limpar descrições e módulos herdados do dev

- [ ] [BACKLOG] **PBI CFG-10.2**: Catálogo com descrições e módulos consistentes.

  - [ ] [BACKLOG] CFG-10.2.1: Corrigir descrições/módulos inconsistentes herdados do dev (ex.: `roles.*` dizem "usuários"; módulo "dashboard" em minúsculo), sem renomear os códigos `permission`; avaliar impacto em bancos já semeados (o seed não altera permissões existentes). (depende de CFG-9.3.1)

#### PBI CFG-10.3: Higiene de segurança herdada (fora do escopo do catálogo, registrada aqui por decisão do usuário)

- [ ] [FINALIZADO] **PBI CFG-10.3**: Logs de produção não vazam dados sensíveis. Follow-up pré-existente, achado do revisor na CFG-9.4.1, não introduzido pela feature de seed.

  - [ ] [FINALIZADO] CFG-10.3.1: `Backend/common/database/db.go` (linhas ~15 e ~27) fixa o logger do GORM em `LogMode(logger.Info)` em qualquer ambiente, então todo SQL com parâmetros vai para o log, inclusive `password_hash` e dados pessoais. Fazer o nível do log depender de `APP_ENV` (production = Warn/Error). **Prioridade ALTA (reforçada em 2026-10-01 com evidência concreta):** nos logs do container api-auth da validação do deploy apareceu o `UPDATE users SET ... "password_hash"='$2a$10$...'` completo a cada login (era um admin descartável), ou seja, em produção o hash bcrypt de TODO usuário que loga vai para o log. Recomendado fazer ANTES de qualquer deploy com dados reais.
    - Iniciada em 2026-10-01, direto na branch `develop` (worktree `C:\Projetos\lumini-hub-develop`), por autorização explícita do usuário ("corrige o logger do GORM antes do push, e faz o push"), em desvio consciente da prática de branch de feature.
    - Critério de aceite:
      - [x] O nível do log do GORM depende de `APP_ENV`: `Info` só em `development`; qualquer outro ambiente (incluindo staging e production) em `Warn`.
      - [x] Fora de `development`, as consultas logadas por erro/lentidão não carregam valores de parâmetros (`ParameterizedQueries`).
      - [x] `record not found` não vira linha de erro no log (`IgnoreRecordNotFoundError`).
      - [x] Teste automatizado do nível por ambiente.
      - [x] Verificação no deploy real (pilha Docker com `APP_ENV=staging`): nenhum hash bcrypt (`$2a$`) nos logs depois de um login.
      - [x] `go build ./...` e `go vet` limpos.
    - Finalizada em 2026-10-01, direto na `develop` (commit 303c442 + ajuste de higiene de teste). Revisor aprovou, sem bloqueantes. Nível por `APP_ENV`: `Info` só em `development` (case-insensitive); staging/production/vazio/desconhecido = `Warn`, `ParameterizedQueries`, `IgnoreRecordNotFoundError`; aplicado ao Postgres e ao SQL Server.
    - ACHADO IMPORTANTE DURANTE A EXECUÇÃO: a 1ª versão (só `ParameterizedQueries`) NÃO bastava. Um teste com valor secreto numa consulta lenta via `Raw().Scan()` mostrou o hash vazando no log de SLOW SQL, porque o `Scan()` do GORM v1.30.0 grava o SQL por um `traceRecorder` interno que ignora o filtro do logger. A correção sobrescreve a variável GLOBAL `logger.RecorderParamsFilter` fora de `development` (e restaura em `development`).
    - Provas: 5 testes de logger com banco real (erro, consulta lenta, `Scan` com erro, `record not found` sem log, `development` mostrando valores inclusive no `Scan`); mutações (remover `ParameterizedQueries`; remover a sobrescrita global) fazem testes falharem; suíte completa passa; pilha Docker com `APP_ENV=staging` a partir de volume vazio: 13 PASS, incluindo "nenhum hash bcrypt nos logs depois dos logins", "nenhuma linha com `password_hash`" e "sem `SELECT * FROM` de rotina" (rodada COM a correção final). Registro honesto: uma rodada intermediária deu 1 FAIL numa checagem combinada (hash OU nome da coluna) cujo conteúdo não foi capturado (containers já removidos), provavelmente o nome da coluna num log de consulta lenta; a checagem foi separada em duas e a rodada seguinte passou com 0 linhas. Aguardando teste do usuário (EM-TESTE/TESTADO-USUARIO ainda não percorridos).

#### PBI CFG-10.4: Conferir o seed do menu contra o banco de dev

- [ ] [BACKLOG] **PBI CFG-10.4**: `SeedMenuItems` do código reflete o menu real do dev.

  - [ ] [BACKLOG] CFG-10.4.1: Conferir os itens de menu do `menu_item_seeder.go` contra a tabela `menu_items` do banco de dev (permissão vinculada por item). A conferência feita na CFG-9.5.2 cobriu só `/settings/*` (`companies.view`, `integrations.view`, `roles.view`, `users.view`) e achou uma divergência (`/settings/roles` ligado a `admin.create_permissions` no código vs. `roles.view` no dev, já corrigida); os demais itens não foram conferidos. Divergências corrigidas no seeder ou registradas.

#### PBI CFG-10.5: Notas e follow-ups do revisor sobre o log do GORM (da CFG-10.3.1; não bloqueantes)

- [ ] [BACKLOG] **PBI CFG-10.5**: Endurecer o restante da higiene de log.

  - [ ] [BACKLOG] CFG-10.5.1: `APP_ENV` não definido vira `development` em `config.Load` (`Backend/common/config/config.go` ~148), então um deploy que esqueça `APP_ENV` volta ao log verboso. Decidir entre exigir `APP_ENV` explícito fora de dev ou inverter o padrão para o modo seguro.
  - [ ] [BACKLOG] CFG-10.5.2: A mensagem de erro do Postgres (pgconn) pode citar valores em alguns casos (ex.: `invalid input syntax for type integer: "abc"`; o `Detail` da violação de unique NÃO entra na string do erro). Residual baixo; mitigação: nunca logar `err` bruto nos handlers.
  - [ ] [BACKLOG] CFG-10.5.3: `logger.RecorderParamsFilter` é variável GLOBAL escrita sem sincronização: ok hoje (um `InitDB` por serviço na subida, testes sem `t.Parallel`); se algum dia houver `InitDB` fora do boot, proteger com `sync.Once`/mutex.
  - [ ] [BACKLOG] CFG-10.5.4: SQL com literais embutidos via `Raw` concatenado não é filtrado por nenhum mecanismo de parâmetros (o código do projeto usa `?`); manter a regra de sempre usar placeholders.
