# Schema (AutoMigrate) e Dados Base (Seed) - Banco Novo Sobe Só com o Código

> Escopado em 2026-10-01 (cocriador, decisões fechadas com o usuário no chat). Mora em `Modulo_0_Configuracao/` por ser fundação. Ainda **não decomposto em backlog**, isso é do `orquestrador-lumini-hub`. Relacionado: `plano_onboarding_empresa.md` (depende deste: o admin de um banco novo hoje só existe por SQL manual) e `Modulo_8_MultiTenant/plano_multitenant.md` (consumidor futuro do pacote de bootstrap).

## Problema (descoberto no deploy na OCI, 2026-10-01)

Só `MenuItem` (api.auth), `Company`/`CompanyFiscalConfig` (api.core) e as 4 tabelas de api.integrations têm `AutoMigrate`. As tabelas centrais (`users`, `roles`, `permissions`, `role_permissions`, `user_permissions`, `customers`, `suppliers`, `country`/`state`/`city`/`address`, `contact`, `document`) vinham de um SQL externo (`Backend/migrations/create_database.sql`, do monólito, apagado do histórico). Num banco novo os serviços sobem mas o login falha (`relation "users" does not exist`). Solução provisória na branch `feature/deploy-oci`: `deploy/db/01_schema.sql` (pg_dump --schema-only do dev) e `02_catalog.sql`. Também faltam dados base: não existe seeder de permissions, roles nem admin (só `SeedMenuItems` e `ConfigService.SeedDefaults`).

## Objetivo

Um banco **vazio** sobe apenas com o código: os serviços criam as tabelas, o catálogo de permissions e o perfil ADMIN são semeados, e (opcionalmente, por env) o primeiro usuário admin é criado. Login funciona sem nenhum SQL manual.

## Decisões fechadas

1. **Estratégia híbrida.** AutoMigrate completo agora, com flag `DB_AUTO_MIGRATE` para desligar em produção. Migrations versionadas (golang-migrate/goose, SQL puro, tabela de controle por serviço, baseline = schema congelado) ficam como **evolução futura**, só quando houver dados reais em produção (fora de escopo, registrado aqui).
2. **Seed/bootstrap em pacote isolado e reutilizável**, não embutido no `main.go`: o app comercial do Módulo 8 vai provisionar tenants (criar banco, rodar migrations, criar o primeiro admin) e deve poder chamar o mesmo código. O bootstrap por variável de ambiente é **transitório**.
3. **Admin inicial** criado no primeiro boot **somente se não houver nenhum usuário**, com perfil **ADMIN**. O perfil **DEVELOP não é criado pelo seed** (equipe dev cria manualmente quando precisar, como já diz o `CLAUDE.md`).
4. **Roles semeadas: só ADMIN.** Vendas, Gerente, Estoque, Compras, Financeiro, "Teste RBAC" não entram (perfis por segmento virão depois).
5. **Permissions: todas as 59 do catálogo atual**, transcritas 1:1 do banco de dev, sem renomear (limpeza de nomes é tarefa separada), e vinculadas ao ADMIN como está hoje.
6. **Geografia sem seed.** Country/State/City serão criados sob demanda na consulta de CEP do cadastro de cliente/fornecedor (futuro, módulo de Clientes/Fornecedores, não implementar aqui).
7. `deploy/db/01_schema.sql` e `02_catalog.sql` (branch `feature/deploy-oci`) são **removidos ao final**, depois da validação (ver critérios de aceite). Não remover antes.

## Regras de negócio / técnicas

### AutoMigrate por serviço (ordem importa só dentro do serviço)

Não existem FKs GORM entre serviços (`User.CompanyID`/`ActiveCompanyID` são `*uint` puros apontando para `companies`), logo **não há dependência de ordem de subida entre serviços**.

- **api.auth**: `Permission`, `Role`, `User`, `MenuItem` (nessa ordem). As join tables `role_permissions` e `user_permissions` são criadas pelo GORM via `many2many`. `MenuItem` tem FK para `Permission`, por isso vem depois.
- **api.core**: `Country`, `State`, `City`, depois `Customer`/`Supplier`, depois `Address`/`Document`/`Contact`, além dos já migrados `Company`/`CompanyFiscalConfig`.
- **api.integrations**: já completo, não muda. As structs de `internal/legacy/domain` são do SQL Server legado e **nunca** entram no migrate.
- **Flag `DB_AUTO_MIGRATE`** (env, via `common/config`): `true` executa o migrate no boot; `false` pula. Padrão a definir na implementação (recomendado `true` para não quebrar dev; o `.env` de produção seta `false` quando migrations versionadas existirem).
- **Código morto `RolePermissions`** (`api.auth/internal/domain/role.go:21`): campos não exportados (`role_id`, `permission_id`), não funciona como modelo GORM nem é usado na relação (a join table vem do `many2many`). Remover (confirmar com grep que ninguém referencia) ou, se algo precisar dela, corrigir com campos exportados e tags. Não migrar.
- **Validação contra o schema atual**: comparar o que o AutoMigrate gera em banco vazio com `deploy/db/01_schema.sql` (`git show feature/deploy-oci:deploy/db/01_schema.sql`) para tipos, NOT NULL, defaults, índices únicos e FKs não divergirem do banco de dev. Divergências são corrigidas nas tags das structs ou registradas como pendência, não ignoradas.

### Seed (pacote isolado, idempotente)

- Local sugerido: pacote reutilizável (ex.: `Backend/common/bootstrap` ou `api.auth/internal/seeder` exportado de forma importável pelo app comercial; escolha final na implementação, desde que **não fique dentro do `main.go`** e **não dependa de estado global do serviço**: recebe `*gorm.DB` e config explícita).
- **Idempotência**: upsert por chave natural (`permissions.permission`, `roles.name`, `users.username`). Rodar N vezes dá o mesmo estado.
- **Nunca apaga** permissions criadas à mão pela tela de Perfis/Permissões. O seed só insere/atualiza as do catálogo (descrição/módulo) e garante o vínculo.
- **Catálogo em código** (slice Go, fonte única): `permission` (código), `description`, `module`, transcritos 1:1 do dev. Resolve a lacuna do `CLAUDE.md` ("registrar permissões nos seeders").
- **Vínculo ADMIN**: o vínculo role-permission é feito para as permissões do catálogo que o ADMIN tem hoje (50). Para não desfazer ajustes feitos pelo operador depois, o vínculo só é criado quando a role ADMIN foi **criada neste seed** (primeiro boot); em boots seguintes, só se adicionam permissions novas do catálogo ainda inexistentes no banco (pendência P2 abaixo detalha a regra).
- **Ordem no boot do api.auth**: AutoMigrate, depois seed de permissions + role ADMIN, depois `SeedMenuItems` (o menu referencia permission por código, então o seed de permissions **tem** que rodar antes), depois bootstrap do admin.
- **Admin inicial** (env, transitório): `BOOTSTRAP_ADMIN_USERNAME`, `BOOTSTRAP_ADMIN_PASSWORD`, `BOOTSTRAP_ADMIN_EMAIL`. Cria o usuário (hash bcrypt via `common/utils`, perfil ADMIN, sem `company_id`, ou seja, master) **só se `users` estiver vazia**. Sem `BOOTSTRAP_ADMIN_PASSWORD`: não cria e loga aviso. **Nunca** há senha default, **nunca** `987321`. A senha não é logada. Usuário criado fica "master", compatível com `plano_onboarding_empresa.md` (primeiro login com zero empresas cai no cadastro de empresa).
- Fonte do catálogo: ver seção a seguir.

### Fonte e conteúdo do catálogo (consulta somente leitura ao Postgres de dev `erp_system`, 2026-10-01)

Consultado apenas `permissions`, `roles` e `role_permissions`, **sem copiar nenhum dado de usuários ou empresas**. Total: **59** permissions ativas (`deleted_at is null`). Roles no dev e quantidade de vínculos: ADMIN 50, Compras 5, Estoque 9, Financeiro 8, Gerente 9, Vendas 9, DEVELOP 0 (só ADMIN entra no seed).

Permissions por módulo (nome do `module` como está hoje, inclusive inconsistências de caixa):

- Admin: `dashboard.admin.view`
- Compras: `purchases.create`, `purchases.delete`, `purchases.edit`, `purchases.reports`, `purchases.view`
- dashboard: `dashboard.view_default`
- Develop (9): `admin.create_permissions`, `permissions.create`, `permissions.delete`, `permissions.edit`, `permissions.view`, `roles.create`, `roles.delete`, `roles.edit`, `roles.view`
- Empresas: `companies.create`, `companies.delete`, `companies.edit`, `companies.fiscal_config.create`, `companies.fiscal_config.edit`, `companies.fiscal_config.view`, `companies.hierarchy.view`, `companies.view`
- Estoque: `dashboard.inventory.view`, `inventory.create`, `inventory.delete`, `inventory.edit`, `inventory.reports`, `inventory.view`, `prices_promotions.view`, `product_location.view`, `products.view`, `stock_locations.view`, `supplier_codes.view`, `taxation.view`
- Financeiro: `dashboard.finance.view`, `finance.create`, `finance.delete`, `finance.edit`, `finance.receive_boleto`, `finance.reports`, `finance.view`, `finance.view_pendencies`
- Gerencial: `dashboard.manager.view`
- Integrações: `integrations.view`
- Usuários: `users.create`, `users.delete`, `users.edit`, `users.view`
- Vendas: `customers.view`, `dashboard.sales.view`, `orders.view`, `payment_plans.view`, `sales.create`, `sales.delete`, `sales.edit`, `sales.reports`, `sales.view`

**Vínculo do ADMIN hoje = 50 permissions = as 59 menos as 9 do módulo "Develop"** (as 9 listadas acima, incluindo `roles.view`/`permissions.view`). Isto é consistente com o `CLAUDE.md`: ADMIN passa por bypass em tudo exceto escrita no catálogo de Perfis/Permissões, e `roles.view`/`permissions.view` ficam abertos a ADMIN pelo bypass em `common/utils/rbac.go`, não pelo vínculo. A implementação deve **reproduzir exatamente esses 50 vínculos** (a transcrição em código é o que vale; conferir a contagem em teste). O nome do campo na tabela é `permission` (coluna do código), com `description` e `module`; as descrições devem ser transcritas do dev na implementação (não listadas aqui).

## Dado / DB

Nenhuma tabela nova. Pode exigir ajuste de tags GORM nas structs existentes para igualar o schema de dev (ver validação). Nenhum campo novo em `users` (ex.: "forçar troca de senha no primeiro login" **não** existe e não entra aqui, ver pendências).

## Arquivos afetados (previstos)

- `Backend/microservices/api.auth/main.go`: AutoMigrate completo na ordem acima, sob `DB_AUTO_MIGRATE`; chamada ao pacote de seed/bootstrap antes de `SeedMenuItems`.
- `Backend/microservices/api.core/main.go`: AutoMigrate completo na ordem acima, sob `DB_AUTO_MIGRATE`.
- `Backend/microservices/api.integrations/main.go`: apenas respeitar `DB_AUTO_MIGRATE` (mesmo comportamento hoje com a flag ligada).
- `Backend/common/config` (+ `Backend/.env.example`): `DB_AUTO_MIGRATE`, `BOOTSTRAP_ADMIN_USERNAME`, `BOOTSTRAP_ADMIN_PASSWORD`, `BOOTSTRAP_ADMIN_EMAIL`.
- Novo pacote de seed/bootstrap isolado (catálogo de permissions em código, seed da role ADMIN, bootstrap do admin).
- `Backend/microservices/api.auth/internal/seeder/menu_item_seeder.go`: só se a ordem/dependência exigir ajuste.
- `Backend/microservices/api.auth/internal/domain/role.go`: remoção/correção de `RolePermissions`.
- Structs de `api.auth/internal/domain` e `api.core/internal/domain` (tags GORM) conforme validação contra `01_schema.sql`.
- `CLAUDE.md`: atualizar as menções de "seeder"/AutoMigrate (hoje diz que Customer/Supplier não migram e que não há seeder de permissions/admin) quando entregar.
- Remoção final (na branch `feature/deploy-oci` ou onde estiver então): `deploy/db/01_schema.sql`, `deploy/db/02_catalog.sql`, ajuste de `deploy/db/README.md`/docker compose e do roteiro de deploy.

## Critérios de aceite

1. **Banco vazio (ex.: Docker Postgres local)**: subir api.auth/api.core/api.integrations só com o código, sem `01_schema.sql`/`02_catalog.sql`: todas as tabelas criadas, sem erro.
2. Com `BOOTSTRAP_ADMIN_*` definidos: admin criado no primeiro boot e **login funciona**; reiniciar não duplica nada nem altera a senha.
3. Sem `BOOTSTRAP_ADMIN_PASSWORD`: nenhum usuário criado, aviso no log, serviços sobem normalmente.
4. 59 permissions e 1 role (ADMIN) com 50 vínculos no banco novo; nenhum perfil DEVELOP/Vendas/Gerente/etc.
5. Rodar o seed duas vezes dá o mesmo estado; uma permission criada manualmente antes da segunda execução permanece.
6. Com `DB_AUTO_MIGRATE=false`, nenhum serviço executa migrate.
7. Comparação do schema gerado com `01_schema.sql` sem divergência não justificada nas tabelas cobertas.
8. **Tarefa de fechamento**: só após 1 a 7, remover os arquivos provisórios `deploy/db/01_schema.sql` e `02_catalog.sql` e atualizar o roteiro/README de deploy.

## Fora de escopo

- Migrations versionadas (golang-migrate/goose): evolução futura, quando houver dados reais em produção.
- Seed de países/estados/cidades: virá sob demanda pela consulta de CEP no cadastro de Clientes/Fornecedores (módulo futuro).
- Perfis por segmento (Vendas, Gerente, Estoque, Compras, Financeiro) e role/usuário DEVELOP.
- Limpeza/renomeação do catálogo de permissions e do campo `module` (inconsistências de caixa, ex.: "dashboard" vs "Admin"), tarefa separada.
- Provisionamento de tenant (criar banco, rodar migrations, admin por plano), Módulo 8; este plano só garante que o pacote de seed seja reutilizável.
- Tabelas existentes no dump de dev sem struct no código (`products`, `sales`, `purchases`, `financial_transactions`, `inventory_movements`, `payment_methods`, `product_categories`, `measurement_units`, `system_logs`, `sale_items`, `purchase_items`): legado do monólito, **não** serão criadas por este escopo; cada módulo cria as suas quando forem implementados.
- Cadastro da primeira empresa pela UI: ver `plano_onboarding_empresa.md`.

## Riscos

- **AutoMigrate em produção com dados**: só adiciona/altera, mas pode tentar alterar tipos, constraints e índices de forma surpreendente em tabela povoada; por isso a flag `DB_AUTO_MIGRATE=false` em produção e a migração futura para versionamento. Testar sempre primeiro num dump de dev.
- **`unique` globais em `Document.Number` e `Contact.Contact`** (api.core): impedem o mesmo número/contato em entidades diferentes (ex.: um mesmo e-mail para cliente e fornecedor). É regra pré-existente; este escopo **não muda**, mas o AutoMigrate passa a criar essas constraints em banco novo, então a divergência de comportamento vira definitiva. Registrado para decisão futura no módulo de Clientes/Fornecedores.
- **Divergência entre struct e schema de dev**: constraints/índices/nomes gerados pelo GORM podem diferir do dump (ex.: nomes de FK/índices `uni_*`, `idx_*`); a validação do critério 7 mitiga.
- **`MenuItem` depende de `Permission`** (FK): erro de ordem quebra o boot em banco vazio; coberto por critério 1 e pela ordem explícita.
- **Seed vs. edição manual**: reaplicar vínculos a cada boot desfaria ajustes do operador no ADMIN; mitigado pela regra "vínculo só na criação da role" (pendência P2).
- **Admin por env**: senha em variável de ambiente fica exposta a quem acessa o ambiente/compose; é transitório, aceito, e a senha nunca é logada nem tem default.
- **Banco com tabelas legadas órfãs** (`empresas`, tabelas listadas em "fora de escopo"): inofensivas, mas podem confundir comparação de schema.

## Pendências em aberto

- [ ] **P1 - Local exato do pacote de seed** (`common/bootstrap` vs. exportar o seeder de api.auth): como `api.auth` usa `internal/`, não é importável por outro módulo; a escolha final (provavelmente fora de `internal/`, em `common/`) é decisão técnica da implementação, validar contra a restrição de que `common` não importe nada de microserviço (as structs `User`/`Role`/`Permission` estão em `api.auth/internal/domain`).
- [ ] **P2 - Regra de reaplicação do vínculo ADMIN em boots seguintes** quando o catálogo ganhar permissions novas (ex.: nova feature adiciona `fiscal.view`): proposta é "inserir a permission nova e vinculá-la ao ADMIN só se for permission ainda inexistente no banco", sem reatribuir as que o operador removeu do ADMIN. Confirmar com o usuário.
- [ ] **P3 - Valor padrão de `DB_AUTO_MIGRATE`** (recomendado `true` quando ausente, para não quebrar dev).
- [ ] **P4 - Troca obrigatória de senha no primeiro login** (`must_change_password`): não existe hoje, não entra aqui; avaliar no futuro junto do provisionamento do Módulo 8.
- [ ] **P5 - Roteiro de primeiro acesso na doc de deploy da OCI** (admin por env, login, cadastro da empresa), a atualizar ao remover os arquivos provisórios.
