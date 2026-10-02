# Deploy (Docker + Caddy)

Sobe o Lumini numa VM só: Postgres, `api.auth`, `api.core`, gateway e o frontend servido pelo Caddy
(que também faz proxy de `/api` para o gateway). Só o Caddy expõe portas (80/443). O
`api.integrations` fica fora por padrão (profile `integrations`): o acesso ao SQL Server legado vai
virar uma aplicação separada, no servidor do cliente.

## Primeiro acesso (banco vazio)

Não há script de banco: os serviços criam as tabelas sozinhos (`DB_AUTO_MIGRATE=true`, o padrão) e o
`api.auth` semeia o catálogo de permissões, o perfil `ADMIN` (com as permissões que ele recebe) e o menu.

1. Na VM, com Docker instalado, clone o repositório e entre nele.
2. Crie o `deploy/.env` a partir do `deploy/.env.example` e preencha:
   `DB_PASSWORD`, `JWT_SECRET`, `CERTIFICATE_ENCRYPTION_KEY` (veja os comandos `openssl` no próprio
   arquivo), `CORS_ALLOWED_ORIGINS` (o endereço público) e `BOOTSTRAP_ADMIN_PASSWORD` (8 a 72 caracteres,
   uma senha forte). Proteja o arquivo: `chmod 600 deploy/.env`. Nunca o commite.
3. Suba tudo:
   ```bash
   docker compose -f deploy/docker-compose.prod.yml --env-file deploy/.env up -d --build
   ```
4. Abra o endereço da VM e entre com o usuário `admin` (ou o `BOOTSTRAP_ADMIN_USERNAME` que você definiu)
   e a senha do `BOOTSTRAP_ADMIN_PASSWORD`.
5. Num banco novo ainda não existe nenhuma empresa. O sistema mostra o seletor de empresa com o botão
   **Cadastrar Empresa**: cadastre a primeira e o acesso segue direto para o sistema.
6. **Remova `BOOTSTRAP_ADMIN_PASSWORD` do `deploy/.env`** e recrie o container do `api-auth`
   (`docker compose ... up -d api-auth`). A senha só é lida no primeiro boot, e deixá-la no arquivo é
   risco desnecessário.

O admin só é criado se a tabela de usuários estiver vazia e a senha estiver definida. Reiniciar não o
duplica nem muda a senha. Sem a senha, nenhum usuário é criado (há um aviso no log do `api-auth`).

## Atualizar a versão

```bash
git pull
docker compose -f deploy/docker-compose.prod.yml --env-file deploy/.env up -d --build
```

Permissões novas criadas no código entram no banco (e no perfil `ADMIN`, exceto as do módulo `Develop`) no
boot seguinte, sem tocar nas que já existem nem nos ajustes feitos pelo operador.

## `DB_AUTO_MIGRATE=false` (produção com dados reais)

Com a flag desligada nenhum serviço migra e o `api.auth` não semeia nada. O catálogo de permissões e o
menu **não** são aplicados sozinhos: aplique-os por outro meio (por exemplo, subindo uma vez com a flag
ligada). O `api.integrations`, se estiver no ar, ainda precisa que a tabela `integration_configs` exista.

## HTTP, domínio e HTTPS

Sem domínio, o acesso é por IP em HTTP (`SITE_ADDRESS=:80`, `APP_ENV=staging`). Os cookies de login só
ganham o atributo `Secure` com `APP_ENV=production`, o que exige HTTPS. Com um domínio apontando para a
VM, troque `SITE_ADDRESS` pelo domínio (o Caddy emite o certificado sozinho), ajuste
`CORS_ALLOWED_ORIGINS` para `https://seu-dominio` e use `APP_ENV=production`.

## Notas de infraestrutura (OCI)

- VM `VM.Standard.A1.Flex` (Ampere, arm64): as imagens são compiladas na própria VM.
- Portas 80 e 443 liberadas na Security List da VCN e no iptables do Ubuntu.
- A chave SSH da VM fica fora do repositório.
