# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/Release_<versão>.md` e zera esta lista.

- **Versão prevista:** a definir (próxima publicação: `2026.2.3.0` se não mexer no banco; `2026.3.3.0` se mexer)
- **Última versão enviada:** `2026.2.2.0` (versão do banco 2), ambiente da OCI em 2026-10-22
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

- **CFG-014**: textos e módulo de permissões do catálogo corrigidos no código. Só vale para bancos novos; bancos já semeados (inclusive o da OCI) mantêm os textos antigos, porque o seed nunca altera permissão existente.

## Novidades e melhorias

- **CFG-024**: barra de status no rodapé do sistema com a versão do sistema, a versão do banco (com aviso se estiver atrás do exigido) e o ambiente. Endpoint novo `GET /auth/version`, protegido por login.
- **CFG-026**: o painel de permissões do usuário (editar usuário e o resumo) passa a mostrar sempre a descrição completa de cada permissão, em vez de "Ver"/"Criar"/"Editar"/"Excluir" quando o módulo tinha só uma.

## Correções

- **CFG-017**: sem `APP_ENV` definido o sistema passa a subir em modo seguro (log do banco sem valores), em vez de `development`. Quem depende do SQL no log em desenvolvimento precisa ter `APP_ENV=development` no `.env`; o servidor da OCI já usa `staging`.
- **CFG-018**: um erro interno (500) deixa de devolver o texto do erro (que podia trazer nome de tabela, constraint e valores do Postgres) na resposta da API: o erro completo vai para o log do servidor e a resposta leva só a mensagem.
- **CFG-015**: o item de menu "Financeiro > Contas a Pagar > Fornecedores" passa a usar `suppliers.view` (antes `purchases.view`), como o de Compras; vale só para bancos novos.
- **CFG-014**: `roles.*` deixam de se chamar "usuários" (agora "perfis"); `permissions.*` e `integrations.*` em minúsculas como as demais; `finance.receive_boleto` vira "Receber boletos"; `dashboard.view_default` passa do módulo `dashboard` para `Dashboard`, que aparece em "Outros" na edição de usuário.

## Interno (não afeta o servidor)

- **CFG-022**: `CLAUDE.md` e as skills `lumini_hub_backend_architecture` e `lumini_hub_entity_creation` traduzidos para português, com três pontos desatualizados das skills corrigidos. O `CLAUDE.md` também deixou de descrever o banco legado e o login local passou a depender do `BOOTSTRAP_ADMIN_PASSWORD`.
