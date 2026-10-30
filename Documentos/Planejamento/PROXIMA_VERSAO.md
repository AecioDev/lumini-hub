# Próxima versão

Alterações feitas desde o último envio ao servidor. O Claude acrescenta uma linha aqui ao terminar cada tarefa. Quando você enviar a versão, peça o release: o Claude lê este arquivo, gera `releases/Release_<versão>.md` e zera esta lista.

- **Versão prevista:** a definir (próxima publicação: patch ou nova sequência se não mexer no banco; versão do banco +1 se mexer)
- **Última versão publicada:** `2026.2.3.0` (versão do banco 2). O ambiente de testes foi atualizado para ela em 2026-10-29 (CFG-033, conferido no navegador).
- **Formato:** `ano.versao_do_banco.sequencia_versao.sequencia_patch`. A versão do banco sobe uma unidade por **publicação** que contém mudanças de banco (quantas tabelas, campos ou triggers forem, na mesma publicação conta uma só); a seção "Mudanças que afetam o banco" abaixo lista o que mudou nessa versão.

## Mudanças que afetam o banco ou a configuração do servidor

_Nada ainda._

## Novidades e melhorias

- **CFG-028** (Fases A, B, C e D, na branch `feature/cfg-028-layout`): o logo da marca sobe para o topo do menu lateral, na mesma altura do cabeçalho; a empresa ativa fica no pé do menu (ícone de troca, que abre a gaveta) e o nome dela também aparece na barra de status; o título sai do cabeçalho e o caminho da tela (breadcrumb) ganha uma faixa entre o cabeçalho e o corpo; a busca do cabeçalho sai; todas as páginas passam a usar o mesmo cabeçalho de página (título, descrição, botões de ação e Voltar), fixo no topo ao rolar; o sistema passa a funcionar em tablet e celular: o menu lateral vira só ícones em tela média e uma gaveta em tela estreita, as tabelas rolam por dentro do cartão e a barra de status cabe numa linha; no celular o rótulo dos campos dos formulários fica sempre em cima, os botões de ação das tabelas ficam maiores para o toque, e os botões do cabeçalho ganham rótulos de acessibilidade; as bordas do menu seguem o tema claro e escuro; o menu recolhido é lembrado ao recarregar.

## Correções

_Nada ainda._

## Interno (não afeta o servidor)

_Nada ainda._
