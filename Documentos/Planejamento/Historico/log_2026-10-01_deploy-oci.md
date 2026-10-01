# Log 2026-10-01: ambiente de dev na Oracle Cloud (OCI)

Branch: `feature/deploy-oci` (a partir da `develop`).

## Infra criada (console OCI, feita pelo usuário)
- Tenancy `espirandams`, região Brazil East (São Paulo), VM `lumini-hub-dev`.
- Shape `VM.Standard.A1.Flex` (Always Free, arm64), 2 OCPU / 12 GB, Ubuntu 22.04, boot 50 GB.
- IP público **reservado** `64.181.183.243`. Security List com 22, 80 e 443 de entrada.
- Docker instalado na VM; iptables liberando 80/443.

## Código (pasta `deploy/`)
- `Dockerfile.backend` (um só, `--build-arg SERVICE=api.x`), `Dockerfile.frontend` (pnpm fixado em 10.5.2: o pnpm 12 ignora `pnpm.onlyBuiltDependencies` e barra o build do esbuild), `Caddyfile`, `docker-compose.prod.yml`, `.env.example`.
- Só o Caddy expõe portas. `api.integrations` fica no profile `integrations` (desligado: depende do SQL Server legado).
- Gateway: `AUTH_SERVICE_URL`, `CORE_SERVICE_URL`, `INTEGRATIONS_SERVICE_URL` e `CORS_ALLOWED_ORIGINS` agora vêm do ambiente (defaults = desenvolvimento).
- `deploy/db/01_schema.sql` + `02_catalog.sql`: schema (pg_dump --schema-only do Postgres local) e catálogo de roles/permissions, executados no primeiro start do Postgres.

## Como subir
```
docker compose -f deploy/docker-compose.prod.yml --env-file deploy/.env up -d --build
```
`deploy/.env` é criado na VM a partir do `.env.example`, com segredos gerados por `openssl rand` (nunca commitado).

## Pós-deploy manual (feito uma vez)
- Usuário `admin` criado por SQL com senha nova (hash bcrypt via `htpasswd`), não a de desenvolvimento.
- Uma empresa inserida por SQL (`companies`), para destravar o seletor de empresa do primeiro login.

## Pendências / decisões
- **Dívida:** as tabelas centrais (users, roles, permissions, customers...) não têm AutoMigrate nem migrations versionadas. `01_schema.sql` é a fonte do deploy e precisa ser regerado quando o schema mudar. Vale decidir um mecanismo de migrations.
- **Backlog:** "primeiro acesso / bootstrap de empresa". Com zero empresas o seletor trava o login e o cadastro de empresa fica atrás do login. Escopar via `cocriador-lumini-hub` (toca o Módulo 8 multi-tenant).
- **Cookies:** ambiente roda com `APP_ENV=staging` e HTTP, porque cookies `Secure` (`production`) não funcionam sem HTTPS.
- **Futuro:** domínio + HTTPS (`SITE_ADDRESS` com o domínio, `APP_ENV=production`, ajustar `CORS_ALLOWED_ORIGINS`).
- `api.integrations` fora do ar neste ambiente.
- A senha do `admin` e o `deploy/.env` existem só na VM.
- Chave SSH movida do OneDrive para `C:\Users\espir\.ssh\` (estava sincronizando na nuvem).
