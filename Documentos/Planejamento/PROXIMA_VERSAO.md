# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/Release_<versão>.md` e zera esta lista.

- **Versão prevista:** a definir (próxima publicação: patch ou nova sequência se não mexer no banco; versão do banco +1 se mexer)
- **Última versão gerada:** `2026.4.5.0` (versão do banco 4), release em `releases/Release_2026.4.5.0.md`. Ainda não publicada no ambiente de testes.
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

- **CFG-052**: duas tabelas novas no `api.core`, criadas pelo `AutoMigrate` (`system_settings`, com uma linha única, e `company_email_configs`). Sem script SQL. Quem roda com `DB_AUTO_MIGRATE=false` precisa subir o `api.core` uma vez com a flag ligada. O menu ganha o item "E-mail do sistema" (o seed insere só o que falta) e o catálogo ganha 4 permissões, vinculadas ao ADMIN no próximo boot do `api.auth`.

## Novidades e melhorias

- **CFG-052** (aguardando teste): configuração de e-mail. Tela **Configurações → E-mail do sistema** (servidor SMTP padrão) e aba **Configuração de E-mail** na Empresa (servidor próprio, opcional). Uma empresa sem servidor próprio herda o da Matriz e, por fim, o do sistema; nunca usa empresa de outro grupo. Senha criptografada, botão de e-mail de teste. É a base da recuperação de senha (CFG-038).

## Correções

(nenhuma)

## Interno (não afeta o servidor)

(nenhuma)
