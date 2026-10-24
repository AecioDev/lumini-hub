# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/Release_<versão>.md` e zera esta lista.

- **Versão prevista:** a definir (próxima publicação: `2026.2.3.0` se não mexer no banco; `2026.3.3.0` se mexer)
- **Última versão enviada:** `2026.2.2.0` (versão do banco 2), ambiente da OCI em 2026-10-22
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

- **CFG-014**: textos e módulo de permissões do catálogo corrigidos no código. Só vale para bancos novos; bancos já semeados (inclusive o da OCI) mantêm os textos antigos, porque o seed nunca altera permissão existente.

## Novidades e melhorias

- **CFG-026**: o painel de permissões do usuário (editar usuário e o resumo) passa a mostrar sempre a descrição completa de cada permissão, em vez de "Ver"/"Criar"/"Editar"/"Excluir" quando o módulo tinha só uma.

## Correções

- **CFG-014**: `roles.*` deixam de se chamar "usuários" (agora "perfis"); `permissions.*` e `integrations.*` em minúsculas como as demais; `finance.receive_boleto` vira "Receber boletos"; `dashboard.view_default` passa do módulo `dashboard` para `Dashboard`, que aparece em "Outros" na edição de usuário.

## Interno (não afeta o servidor)

_Nada ainda._
