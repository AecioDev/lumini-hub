# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/Release_<versão>.md` e zera esta lista.

- **Versão prevista:** a definir (próxima publicação: patch ou nova sequência se não mexer no banco; versão do banco +1 se mexer)
- **Última versão gerada:** `2026.4.5.0` (versão do banco 4), release em `releases/Release_2026.4.5.0.md`. Ainda não publicada no ambiente de testes.
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

- **CFG-038**: tabela nova `password_reset_tokens` no `api.auth` e coluna nova `app_base_url` em `system_settings` (`api.core`), ambas criadas pelo `AutoMigrate`, sem script SQL. Quem roda com `DB_AUTO_MIGRATE=false` precisa subir o `api.auth` e o `api.core` uma vez com a flag ligada. Configuração do servidor: o `api.auth` passa a chamar o `api.core` direto, pela variável `CORE_SERVICE_URL` (a mesma do gateway; o padrão é o endereço local, e o `docker-compose.prod.yml` já a define). Depois de atualizar, preencher o **Endereço do sistema** em Configurações > E-mail do sistema.
- **CFG-052**: duas tabelas novas no `api.core`, criadas pelo `AutoMigrate` (`system_settings`, com uma linha única, e `company_email_configs`). Sem script SQL. Quem roda com `DB_AUTO_MIGRATE=false` precisa subir o `api.core` uma vez com a flag ligada. O menu ganha o item "E-mail do sistema" (o seed insere só o que falta) e o catálogo ganha 4 permissões, vinculadas ao ADMIN no próximo boot do `api.auth`.

## Novidades e melhorias

- **CFG-053**: o login aceita o nome de usuário ou o e-mail da conta.
- **CFG-038**: recuperação de senha por e-mail. Telas **Esqueci a senha** e **Nova senha**; link de uso único válido por 30 minutos, com limite de pedidos por usuário e por IP. Primeiro template de e-mail do sistema. Campo **Endereço do sistema** na tela E-mail do sistema (usado nos links dos e-mails).
- **CFG-052**: configuração de e-mail. Tela **Configurações → E-mail do sistema** (servidor SMTP padrão) e aba **Configuração de E-mail** na Empresa (servidor próprio, opcional). Uma empresa sem servidor próprio herda o da Matriz e, por fim, o do sistema; nunca usa empresa de outro grupo. Senha criptografada, botão de e-mail de teste. É a base da recuperação de senha (CFG-038).

## Correções

- Tela de login: removida a caixa "Lembrar-me", que não fazia nada (a sessão já dura até 7 dias pelo cookie de renovação, com ou sem ela).

## Interno (não afeta o servidor)

(nenhuma)
