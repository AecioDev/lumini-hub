# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/<versão>.md` e zera esta lista.

- **Versão prevista:** `2026.2.2.0`
- **Última versão enviada:** a da VM da OCI hoje, com versão do banco 1 (`2026.1.x.x`)
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

- Os serviços agora criam as tabelas sozinhos (`DB_AUTO_MIGRATE`, padrão `true`). Um banco vazio sobe só com o código; o `api.auth` semeia permissões, role ADMIN e menu.
- Primeiro admin: definir `BOOTSTRAP_ADMIN_PASSWORD` (8 a 72 caracteres) no `deploy/.env` no primeiro boot e remover depois. Sem a senha nenhum usuário é criado.
- `deploy/.env` precisa de `APP_ENV` explícito: fora de `development` o log do GORM não mostra SQL com valores.
- `deploy/db/01_schema.sql` e `02_catalog.sql` foram removidos; o `docker-compose.prod.yml` não monta mais `./db`.
- Nova tabela `schema_versions` (`db_version`, `app_version`, `applied_at`): o banco passa a registrar a versão em que está. O `api-auth` cria a tabela e grava a versão 2 (`2026.2.2.0`) no fim de um boot com `DB_AUTO_MIGRATE=true`.
- Com `DB_AUTO_MIGRATE=false`, `api-auth`, `api-core` e `api-integrations` conferem a versão e **recusam subir** se o banco estiver atrás da 2 ou sem versão registrada. Para atualizar o banco: subir uma vez o `api-auth` com a flag ligada.
- Bancos já semeados (dev e OCI): o item de menu "Fornecedores" continua em `purchases.view`; ajustar a mão para `suppliers.view`.
- As 8 permissões novas entram no próximo boot do `api.auth` e vão automaticamente só para o ADMIN; usuários com permissões individuais precisam recebê-las pela tela de Usuários.

## Novidades e melhorias

- Cadastro da primeira empresa direto na tela de seleção, num banco sem empresas (CFG-8).
- Roteiro de primeiro acesso documentado em `deploy/README.md`.
- Catálogo de permissões cobre todas as rotas: entram `customers.create/edit/delete`, `suppliers.view/create/edit/delete` e `integrations.edit` (CFG-10.1). Catálogo: 67 permissões; o ADMIN recebe 58.
- Menu "Fornecedores" passa a usar `suppliers.view` (cadastro de fornecedores), separado de `purchases.*` (rotina de Compras).

## Correções

- Logger do GORM por ambiente: o SQL com parâmetros (inclusive hash de senha) deixa de ir para o log fora de `development` (CFG-10.3.1).
- Menu "Perfis e Permissões" nascia escondido para o ADMIN em bancos semeados do zero: agora usa `roles.view`.
- `go run main.go` voltou a compilar no `api.auth` (função movida de volta para o `main.go`).

## Interno (não afeta o servidor)

- Removidos subagentes, quadro Kanban e skills externas; o fluxo passa a ser conduzido por você, com `BACKLOG.md` e `PROXIMA_VERSAO.md`.
- `CLAUDE.md` e skill `lumini_hub_dev_flow` reescritos para esse fluxo.
