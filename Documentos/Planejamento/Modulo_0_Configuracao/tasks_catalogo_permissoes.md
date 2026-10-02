# Tasks - Catálogo de Permissões (completar e limpar)

Criado em 2026-10-01 pelo orquestrador, por decisão do usuário, a partir de achados do EPIC CFG-9 (`tasks_schema_seed.md`). Prefixo `CFG`, esquema novo de Kanban. Status de Épico/PBI é sempre a etapa mais atrasada entre as tarefas-filhas. Tudo nasce `[BACKLOG]`. Fora do escopo do EPIC CFG-9 de propósito.

### EPIC CFG-10: Completar e limpar o catálogo de permissões

- [ ] [BACKLOG] **EPIC CFG-10**: O catálogo de permissões em código (`api.auth/seed`, entregue na CFG-9.3.1 com 59 permissões herdadas do banco de dev) passa a cobrir tudo o que `RequirePermission` exige e fica com descrições/módulos consistentes. **Até lá, as rotas que exigem as permissões faltantes só funcionam para ADMIN/DEVELOP (bypass).** Escopo ainda não escopado em `plano_*.md`; as decisões de nome/módulo/descrição são do usuário.

#### PBI CFG-10.1: Cadastrar as 8 permissões exigidas pelo código e ausentes do catálogo

- [ ] [BACKLOG] **PBI CFG-10.1**: Toda permissão exigida por `RequirePermission` existe no catálogo.

  - [ ] [BACKLOG] CFG-10.1.1: Decidir com o usuário nome, módulo e descrição das 8 permissões ausentes: `customers.create`, `customers.edit`, `customers.delete`, `suppliers.view`, `suppliers.create`, `suppliers.edit`, `suppliers.delete`, `integrations.edit` (conferir por grep em `RequirePermission(` se há outras).
  - [ ] [BACKLOG] CFG-10.1.2: Adicionar as permissões decididas ao catálogo em `api.auth/seed`; atualizar as contagens dos testes (hoje 59 permissões / 50 vínculos ao ADMIN; mudam quando entrarem, a regra do ADMIN é `Module != "Develop"`). (depende de CFG-10.1.1, CFG-9.3.4)

#### PBI CFG-10.2: Limpar descrições e módulos herdados do dev

- [ ] [BACKLOG] **PBI CFG-10.2**: Catálogo com descrições e módulos consistentes.

  - [ ] [BACKLOG] CFG-10.2.1: Corrigir descrições/módulos inconsistentes herdados do dev (ex.: `roles.*` dizem "usuários"; módulo "dashboard" em minúsculo), sem renomear os códigos `permission`; avaliar impacto em bancos já semeados (o seed não altera permissões existentes). (depende de CFG-9.3.1)

#### PBI CFG-10.3: Higiene de segurança herdada (fora do escopo do catálogo, registrada aqui por decisão do usuário)

- [ ] [BACKLOG] **PBI CFG-10.3**: Logs de produção não vazam dados sensíveis. Follow-up pré-existente, achado do revisor na CFG-9.4.1, não introduzido pela feature de seed.

  - [ ] [BACKLOG] CFG-10.3.1: `Backend/common/database/db.go` (linhas ~15 e ~27) fixa o logger do GORM em `LogMode(logger.Info)` em qualquer ambiente, então todo SQL com parâmetros vai para o log, inclusive `password_hash` e dados pessoais. Fazer o nível do log depender de `APP_ENV` (production = Warn/Error). **Prioridade ALTA (reforçada em 2026-10-01 com evidência concreta):** nos logs do container api-auth da validação do deploy apareceu o `UPDATE users SET ... "password_hash"='$2a$10$...'` completo a cada login (era um admin descartável), ou seja, em produção o hash bcrypt de TODO usuário que loga vai para o log. Recomendado fazer ANTES de qualquer deploy com dados reais.

#### PBI CFG-10.4: Conferir o seed do menu contra o banco de dev

- [ ] [BACKLOG] **PBI CFG-10.4**: `SeedMenuItems` do código reflete o menu real do dev.

  - [ ] [BACKLOG] CFG-10.4.1: Conferir os itens de menu do `menu_item_seeder.go` contra a tabela `menu_items` do banco de dev (permissão vinculada por item). A conferência feita na CFG-9.5.2 cobriu só `/settings/*` (`companies.view`, `integrations.view`, `roles.view`, `users.view`) e achou uma divergência (`/settings/roles` ligado a `admin.create_permissions` no código vs. `roles.view` no dev, já corrigida); os demais itens não foram conferidos. Divergências corrigidas no seeder ou registradas.
