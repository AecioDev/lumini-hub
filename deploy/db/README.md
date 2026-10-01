# Banco de dados (deploy)

- `01_schema.sql`: schema atual (pg_dump --schema-only do Postgres de desenvolvimento).
- `02_catalog.sql`: catálogo de `roles`, `permissions` e `role_permissions` (sem dados de usuários/empresas).
- O Postgres executa esses arquivos só no **primeiro start com volume vazio**. Para refazer: `docker compose ... down -v`.
- Usuário inicial: criado manualmente (ver passo a passo no log da sessão), com senha nova, nunca a de desenvolvimento.
- Se o schema mudar no desenvolvimento, regerar `01_schema.sql`. Enquanto as tabelas centrais não tiverem AutoMigrate/migrations versionadas, este arquivo é a fonte do deploy.
