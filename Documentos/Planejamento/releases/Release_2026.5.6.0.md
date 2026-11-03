# Lumini Hub — Release 2026.5.6.0

**Versão do sistema:** 2026.5.6.0 · **Versão do banco de dados:** 5 (tabelas novas `password_reset_tokens`, `system_settings` e `company_email_configs`; colunas novas em `users`; `dbversion.Required = 5`)
**Tema da versão:** **e-mail do sistema** (configuração de servidor SMTP por empresa, recuperação de senha e convite de novos usuários), **login por usuário ou e-mail** e **segurança do acesso** (login que não revela quais usuários existem, limite de tentativas e bloqueio imediato de quem é desativado).
*Documento gerado em 2026-11-02. Legenda:* ✅ testado · ⚠️ implementado, ainda não testado · ⏳ pendente / não incluído.

---

## 1. Banco de dados

A versão do banco sobe de **4 para 5**. Há três tabelas novas e colunas novas; nenhuma tabela ou coluna existente é removida e nenhum dado é reescrito.

| Tabela | Mudança | Detalhe |
|---|---|---|
| `system_settings` | **nova** ✅ | Parâmetros do sistema inteiro, **uma única linha** (garantida por constraint): servidor SMTP padrão (`smtp_...`, senha criptografada) e `app_base_url` (endereço público do sistema, usado nos links dos e-mails). Criada pelo `AutoMigrate` do `api.core`. |
| `company_email_configs` | **nova** ✅ | Servidor SMTP próprio de uma empresa (1:1, opcional), com índice único parcial em `company_id`. Sem linha, a empresa herda a da Matriz e, por fim, a do sistema. Criada pelo `api.core`. |
| `password_reset_tokens` | **nova** ✅ | Pedidos de redefinição de senha e convites: só o **hash** do token é gravado, com validade, uso único e a finalidade (`purpose`: `reset` ou `invite`). Criada pelo `api.auth`. |
| `users` | coluna ✅ | `must_change_password` (obriga a troca da senha no primeiro acesso). |
| `permissions` | dados ✅ | 4 permissões novas, vinculadas ao ADMIN no boot: `system_settings.view`, `system_settings.edit`, `companies.email_config.view` e `companies.email_config.edit`. |
| `menu_items` | dados ✅ | Item novo "E-mail do sistema" (`/settings/email`), semeado sozinho em bancos já existentes. |

- Tudo criado pelos serviços no primeiro boot com `DB_AUTO_MIGRATE=true` (o `api.auth` e o `api.core` precisam subir uma vez), e o `api.auth` registra a versão 5. Não há migração SQL nova nesta versão.
- Com `DB_AUTO_MIGRATE=false`, os serviços **se recusam a subir** até o banco estar na versão 5: rode o `api.auth` e o `api.core` uma vez com a variável ligada (ver a seção 5).

---

## 2. Novidades

### 2.1 Configuração de e-mail (CFG-052) ✅
- Tela **Configurações → E-mail do sistema**: servidor SMTP padrão, endereço do remetente e o **Endereço do sistema** (usado nos links dos e-mails). Botão de **e-mail de teste**.
- Aba **Configuração de E-mail** na Empresa: servidor próprio, opcional. A regra de envio é: a própria empresa, depois a Matriz do grupo, depois o padrão do sistema; nunca uma empresa de outro grupo, e sem nenhum servidor ativo o envio para com uma mensagem clara.
- A senha do SMTP é criptografada e nunca volta pela API. Modelos de e-mail em HTML e texto, com layout comum.

### 2.2 Recuperação de senha (CFG-038) ✅
- Telas **Esqueci a senha** e **Nova senha**. Link de uso único, válido por 30 minutos; um pedido novo invalida o anterior; no máximo 3 pedidos por hora por usuário e freio por IP. A resposta é sempre a mesma, exista ou não o e-mail.

### 2.3 Convite e cadastro de usuários (CFG-054) ✅
- No cadastro de usuário, a chave **Enviar convite por e-mail** (ligada por padrão): o usuário recebe um link, válido por 7 dias, para criar a própria senha. Há reenvio na lista de usuários.
- Com a chave desligada, quem cadastra informa a senha inicial, e o usuário recebe por e-mail o login, a senha provisória e o botão de acesso. A opção **Exigir troca no primeiro acesso** leva o usuário à tela **Troque sua senha** antes de usar o sistema.
- Se o e-mail falhar, o usuário é criado do mesmo jeito e a tela avisa.

### 2.4 Login por usuário ou e-mail (CFG-053) ✅
- O campo do login passa a se chamar **Usuário ou e-mail** e aceita os dois, sem diferenciar maiúsculas.

### 2.5 Ajustes de cadastro e listas ✅
- Cadastro de usuário: "Nome completo" antes de "Usuário", com **sugestão do usuário** a partir do nome (primeiro e último nome, sem acento), sem sobrescrever o que já foi digitado.
- Listas: as ações com confirmação (Excluir, Reenviar convite) mostram o nome ao passar o mouse.
- Texto novo do e-mail de convite ✅.

---

## 3. Correções

- **Desativar um usuário derruba o acesso na hora (CFG-057, segurança)** ✅: o sistema confere no banco, a cada requisição, se a conta continua existente e ativa. Antes, o acesso durava até o token de 15 minutos expirar.
- **Login seguro (CFG-047, segurança)** ✅: o login responde sempre "usuário ou senha inválidos" (usuário inexistente e senha errada, com o mesmo tempo de resposta); "usuário inativo" só aparece para quem acertou a senha. O `/auth/login` aceita 20 tentativas por IP a cada 15 minutos; depois disso, a tela mostra o aviso para aguardar.
- **Erros no campo certo** ✅: erro de um campo (como "e-mail já está em uso") agora marca o campo certo no cadastro e na edição de usuário; as respostas 400 da API trazem `fieldErrors`.
- **Rótulo dos campos** ✅: formulários de duas colunas (usuário, Configurações → Integrações) passam a ter o rótulo sempre em cima.
- **Tela de login** ✅: removida a caixa "Lembrar-me", que não fazia nada (a sessão já dura até 7 dias pelo cookie de renovação).
- Confirmação de "Reenviar convite" ✅: o botão agora diz "Reenviar" (antes dizia "Excluir", em vermelho).

---

## 4. Infraestrutura

- Pacote novo `common/mail` (SMTP e modelos embutidos), `common/internalmail` (pedido de envio entre serviços) e token de serviço de 30 segundos para as chamadas internas. A rota interna de e-mail do `api.core` **não passa pelo gateway**.
- O `api.auth` passa a chamar o `api.core` direto (variável `CORE_SERVICE_URL`).
- Swagger do gateway regerado (recuperação de senha, convite, configuração de e-mail, `fieldErrors`).
- Nenhuma mudança nos Dockerfiles nem no Caddy.

---

## 5. Configuração necessária no servidor

- **Atualizar o código, reconstruir as imagens e reiniciar** `gateway`, `api-auth`, `api-core`, `api-integrations` e `web`. Sem recriar o banco: os dados são mantidos.
- **Banco:** com `DB_AUTO_MIGRATE=true`, `api-auth` e `api-core` criam as tabelas e colunas novas no boot, e o `api-auth` registra a versão 5. Com `DB_AUTO_MIGRATE=false`, suba o `api-auth` e o `api-core` uma vez com a variável ligada antes dos demais.
- **Variável nova:** `CORE_SERVICE_URL` (endereço do `api-core`, a mesma que o gateway usa). O `docker-compose.prod.yml` já a define; quem sobe de outro jeito precisa defini-la no `api-auth`.
- **Depois de subir, para os e-mails saírem:** em Configurações → E-mail do sistema, preencher o servidor SMTP e o **Endereço do sistema**. Sem isso, a recuperação de senha e os convites não enviam e-mail (o cadastro de usuário continua funcionando e avisa).
- Depois de subir, a barra de status deve mostrar "Versão 2026.5.6.0 - BD v5".

---

## 6. Arquivos principais alterados / criados

**Biblioteca comum (`Backend/common`):** `mail/` (`mail.go`, `template.go`, `templates/`), `internalmail/client.go`, `middlewares/` (`auth.go`, `service_auth.go`, `rate_limit.go`), `utils/` (`service_token.go`, `response.go`, `response_validation.go`), `config/config.go`.
**api.auth:** `main.go`, `internal/domain/` (`password_reset.go`, `user.go`), `internal/service/` (`password_reset_service.go`, `auth_service.go`, `user_service.go`), `internal/handlers/` (`password_reset.go`, `auth.go`, `users.go`), `internal/validator/`, `internal/routes/routes.go`, `seed/permission_catalog.go`, `internal/seeder/menu_item_seeder.go`.
**api.core:** `main.go`, `internal/domain/email_config.go`, `internal/service/mail_config_service.go`, `internal/repository/` (`system_settings_repository.go`, `company_email_config_repository.go`), `internal/handlers/mail_config.go`, `internal/validator/smtp_validator.go`, `internal/routes/routes.go`.
**api.gateway:** `main.go`, `docs/` (Swagger regerado).
**Frontend:** `pages/auth/` (login, esqueci a senha, nova senha, troque sua senha), `pages/settings/email/`, `pages/users/`, `components/users/forms/`, `components/companies/forms/` (aba de e-mail), `components/settings/`, `components/common/`, `schemas/`, `services/` (`auth`, `settings`, `users`, `companies`), `types/`, `utils/` (erros de campo), `routes/`.
**Deploy:** `deploy/docker-compose.prod.yml` (`CORE_SERVICE_URL`).

---

## 7. Pendências conhecidas (próximas versões)

| Item | Situação |
|---|---|
| **Entrega dos e-mails**: alguns provedores de e-mail barram mensagens enviadas sem um domínio próprio e sem HTTPS (CFG-043) | ⏳ |
| **VND-001 / COM-001** Telas de Clientes e de Fornecedores | ⏳ |
| **Bloqueio de conta por tentativas** (CFG-055) | ⏳ a decidir |
| **Retenção da auditoria**: limpeza automática de registros antigos de `system_logs`, e auditar login e logout | ⏳ |
| **Integrações, Fase 2** (cliente da Loja Integrada e processamento dos webhooks) | ⏳ pausado |
| **HTTPS e domínio** no ambiente de testes | ⏳ Hoje o acesso é por IP, em HTTP |
