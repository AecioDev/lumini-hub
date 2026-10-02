# Tasks - Schema (AutoMigrate) e Seed - Banco Novo Sobe Só com o Código

Backlog decomposto em 2026-10-01 pelo orquestrador a partir de `plano_schema_seed.md` (escopo fechado, não alterado aqui). Prefixo `CFG`, esquema novo de Kanban. Status de Épico/PBI é sempre a etapa mais atrasada entre as tarefas-filhas. Tudo nasce `[BACKLOG]`.

### EPIC CFG-9: Banco novo sobe só com o código (AutoMigrate completo + seed isolado + admin inicial)

- [ ] [BACKLOG] **EPIC CFG-9**: Banco vazio sobe apenas com o código: serviços criam as tabelas (sob `DB_AUTO_MIGRATE`), catálogo de 59 permissions e role ADMIN (50 vínculos) são semeados por pacote isolado e reutilizável, e o primeiro admin é criado por env se não houver usuários. Fecha removendo `deploy/db/01_schema.sql` e `02_catalog.sql`. Escopo: `plano_schema_seed.md`. Fora de escopo: migrations versionadas, seed de geografia, perfis por segmento/DEVELOP, limpeza de nomes do catálogo.

#### Decisões a tomar (pendências abertas do plano, NÃO decididas aqui)

Cada uma bloqueia ou orienta a tarefa indicada; o usuário decide, ninguém assume resposta.

- **P1 — Local exato do pacote de seed** (`common/bootstrap` vs. exportar o seeder de api.auth). `api.auth` usa `internal/` (não importável por outro módulo) e as structs `User`/`Role`/`Permission` estão em `api.auth/internal/domain`; `common` não deve importar microserviço. Decisão técnica, a validar com o usuário. **Bloqueia CFG-9.3.1.**
- **P2 — Regra de reaplicação do vínculo ADMIN em boots seguintes** quando o catálogo ganhar permissions novas. Proposta registrada no plano (ainda a confirmar): inserir a permission nova e vinculá-la ao ADMIN só se for permission inexistente no banco, sem reatribuir as que o operador removeu. **Bloqueia CFG-9.3.2.**
- **P3 — Valor padrão de `DB_AUTO_MIGRATE` quando ausente** (plano recomenda `true`, a confirmar). **Bloqueia CFG-9.1.1.**
- **P4 — Troca obrigatória de senha no primeiro login (`must_change_password`)**: não entra neste Épico; avaliar no futuro com o provisionamento do Módulo 8. Não bloqueia nada, só registrar que está fora.
- **P5 — Roteiro de primeiro acesso na doc de deploy da OCI** (admin por env, login, cadastro da empresa), a atualizar ao remover os arquivos provisórios. Coberta por CFG-9.5.3 (e CFG-8.2.1).

#### PBI CFG-9.1: Configuração e flag `DB_AUTO_MIGRATE`

- [ ] [BACKLOG] **PBI CFG-9.1**: Variáveis de ambiente novas disponíveis em `common/config` para as demais PBIs.

  - [ ] [BACKLOG] CFG-9.1.1: `Backend/common/config` + `Backend/.env.example`: `DB_AUTO_MIGRATE` (default conforme decisão P3), `BOOTSTRAP_ADMIN_USERNAME`, `BOOTSTRAP_ADMIN_PASSWORD`, `BOOTSTRAP_ADMIN_EMAIL`. (depende da decisão P3)

#### PBI CFG-9.2: AutoMigrate completo por serviço

- [ ] [BACKLOG] **PBI CFG-9.2**: Em banco vazio api.auth e api.core criam todas as suas tabelas, sem depender de ordem de subida entre serviços.

  - [ ] [BACKLOG] CFG-9.2.1: api.auth `main.go` — AutoMigrate `Permission`, `Role`, `User`, `MenuItem` (nessa ordem; join tables via `many2many`) sob `DB_AUTO_MIGRATE`; remover o código morto `RolePermissions` de `api.auth/internal/domain/role.go` (conferir por grep que ninguém referencia). (depende de CFG-9.1.1)
  - [ ] [BACKLOG] CFG-9.2.2: api.core `main.go` — AutoMigrate `Country`, `State`, `City`, depois `Customer`/`Supplier`, depois `Address`/`Document`/`Contact`, além de `Company`/`CompanyFiscalConfig` já migrados, sob `DB_AUTO_MIGRATE`. (depende de CFG-9.1.1)
  - [ ] [BACKLOG] CFG-9.2.3: api.integrations `main.go` — apenas respeitar `DB_AUTO_MIGRATE` (comportamento atual com a flag ligada); structs do legado SQL Server nunca entram no migrate. (depende de CFG-9.1.1)
  - [ ] [BACKLOG] CFG-9.2.4: Validar o schema gerado em banco vazio contra `git show feature/deploy-oci:deploy/db/01_schema.sql` (tipos, NOT NULL, defaults, únicos, FKs); divergências corrigidas nas tags das structs (`api.auth/internal/domain`, `api.core/internal/domain`) ou registradas como pendência, não ignoradas. (depende de CFG-9.2.1, CFG-9.2.2)

#### PBI CFG-9.3: Seed isolado e idempotente (catálogo + role ADMIN)

- [ ] [BACKLOG] **PBI CFG-9.3**: Permissions e role ADMIN semeadas por pacote reutilizável pelo app comercial do Módulo 8, fora do `main.go`, recebendo `*gorm.DB` e config explícita.

  - [ ] [BACKLOG] CFG-9.3.1: Pacote de seed/bootstrap isolado (local conforme decisão P1) + catálogo em código (slice Go, fonte única) das 59 permissions (`permission`, `description`, `module`) transcritas 1:1 do dev, sem renomear; descrições transcritas do banco de dev (consulta somente leitura a `permissions`). (depende da decisão P1)
  - [ ] [BACKLOG] CFG-9.3.2: Seed idempotente por chave natural (`permissions.permission`, `roles.name`): upsert das permissions (nunca apaga as criadas à mão), cria role ADMIN e os 50 vínculos (as 59 menos as 9 do módulo "Develop") só quando a role foi criada neste seed; boots seguintes seguem a regra da decisão P2. Sem role DEVELOP/Vendas/Gerente/etc. (depende de CFG-9.3.1, decisão P2)
  - [ ] [BACKLOG] CFG-9.3.3: Ligar no boot do api.auth na ordem AutoMigrate, seed permissions + role ADMIN, `SeedMenuItems`, bootstrap do admin (permissions antes do menu, que as referencia por código); ajustar `menu_item_seeder.go` só se a dependência exigir. (depende de CFG-9.2.1, CFG-9.3.2)
  - [ ] [BACKLOG] CFG-9.3.4: Teste automatizado: banco com seed aplicado dá 59 permissions e role ADMIN com 50 vínculos; rodar duas vezes mantém o mesmo estado; permission criada manualmente antes da segunda execução permanece. (depende de CFG-9.3.2)

#### PBI CFG-9.4: Admin inicial por variável de ambiente

- [ ] [BACKLOG] **PBI CFG-9.4**: Login funciona num banco novo sem SQL manual.

  - [ ] [BACKLOG] CFG-9.4.1: Bootstrap do admin (pacote isolado): cria usuário com hash bcrypt (`common/utils`), perfil ADMIN, sem `company_id` (master), só se `users` estiver vazia; sem `BOOTSTRAP_ADMIN_PASSWORD` não cria e loga aviso; nunca há senha default, a senha não é logada; reinício não duplica nem altera a senha. (depende de CFG-9.1.1, CFG-9.3.2)

#### PBI CFG-9.5: Validação e fechamento

- [ ] [BACKLOG] **PBI CFG-9.5**: Prova de que o código basta, e só então remoção dos arquivos provisórios.

  - [ ] [BACKLOG] CFG-9.5.1: Validação ponta a ponta em banco vazio (ex.: Docker Postgres local) cobrindo os critérios 1 a 7 do plano: serviços sobem só com o código; admin criado e login funcional; sem `BOOTSTRAP_ADMIN_PASSWORD` nenhum usuário e aviso no log; 59 permissions + 1 role + 50 vínculos; seed idempotente; `DB_AUTO_MIGRATE=false` não migra; schema sem divergência injustificada. (depende de CFG-9.2.4, CFG-9.3.3, CFG-9.3.4, CFG-9.4.1)
  - [ ] [BACKLOG] CFG-9.5.2: Atualizar `CLAUDE.md` nas menções de seeder/AutoMigrate (hoje diz que Customer/Supplier não migram e que não há seeder de permissions/admin). Só ao entregar; edição de código/doc fora de `Documentos/Planejamento/`, feita pela sessão de execução. (depende de CFG-9.5.1)
  - [ ] [BACKLOG] CFG-9.5.3: **Tarefa de fechamento**: remover `deploy/db/01_schema.sql` e `deploy/db/02_catalog.sql` (branch `feature/deploy-oci` ou onde estiver então), ajustar `deploy/db/README.md`/docker compose e o roteiro de deploy incluindo o roteiro de primeiro acesso (pendência P5). Não remover antes da validação. (depende de CFG-9.5.1)
