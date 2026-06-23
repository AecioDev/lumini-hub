# Apresentação do Planejamento Geral - Lumini Hub ERP

Este documento apresenta, de forma direta e sem termos técnicos complexos, o planejamento estratégico do desenvolvimento do **Lumini Hub**, o novo sistema de gestão (ERP) focado nas necessidades específicas de lojas de iluminação e decoração.

---

## 🎯 O que é o Lumini Hub?

O Lumini Hub é uma plataforma integrada de gestão empresarial projetada para centralizar e impulsionar a operação comercial da empresa, englobando desde a abordagem inicial ao cliente, passando pelas vendas de balcão e controle físico de depósitos, até o faturamento fiscal e a escrituração financeira e contábil automática.

---

## 💡 Diferenciais Estratégicos do Projeto

1. **Gestão de Múltiplas Filiais (Multi-empresa):** O sistema permitirá gerenciar matriz e filiais (com mesmo CNPJ ou diferentes) a partir de uma única tela. Gestores e diretores podem visualizar dados consolidados do grupo econômico com um clique.
2. **Integração com E-commerce (Loja Integrada):** Sincronização automatizada de estoques e captação de pedidos do site para o sistema físico.
3. **CRM de Alta Conversão:** Uma ferramenta para o vendedor que funciona como assistente de vendas, sugerindo contatos de pós-venda, parabenizações em datas comemorativas e ofertas de renovação baseadas na durabilidade estimada do produto de iluminação.
4. **Orçamentos por Ambientes da Obra:** Uma funcionalidade diferenciada que permite organizar e faturar os itens do orçamento divididos por cômodo (Sala de Estar, Cozinha, Área Gourmet), o que facilita a comunicação com arquitetos, decoradores e instaladores.
5. **Automação Fiscal e Emissão Segura:** Uso de uma biblioteca homologada de mercado (ACBr) para emissão de notas (NFe, NFSe, NFCe, MDFe) e geração de boletos, com um módulo fiscal independente para facilitar a ativação ou suspensão de filiais.

---

## 📅 As Fases de Desenvolvimento

O projeto será entregue em 7 fases lógicas de negócio:

### 📍 Fase 1: Integrações Iniciais, Controle de Filiais e CRM
- **Integração da Loja Integrada:** Conectar o site ao banco de dados atual (SQL Server) para evitar problemas no fluxo de vendas imediato.
- **Estruturação de Filiais:** Cadastro das empresas controladas e painéis de alternância rápida.
- **CRM do Vendedor:** Painel visual (Kanban) para controle de negócios, agenda de tarefas sugeridas, alerta de clientes inativos (churn), ações de pós-venda e atalhos rápidos para envio de mensagens via WhatsApp/E-mail.

### 📍 Fase 2: Catálogo de Produtos e Fotos
- **Cadastro Unificado:** Organização do portfólio de produtos contendo marcas, categorias e durabilidade estimada.
- **Mídia:** Suporte nativo para inclusão de imagens dos produtos no catálogo comercial.

### 📍 Fase 3: Gestão de Estoque e Almoxarifado
- **Controle por Depósitos:** Separação do saldo físico por almoxarifado de cada filial.
- **Endereçamento de Estoque:** Indicação exata de onde o produto está localizado fisicamente (Rua, Prateleira, Gaveta) para que o estoquista realize a separação de forma rápida e sem erros.

### 📍 Fase 4: Compras e Entrada Automatizada de Mercadorias
- **Manifesto SEFAZ:** O sistema monitora de forma automática todas as notas fiscais emitidas por fornecedores contra o CNPJ da empresa, permitindo a manifestação rápida e importação direta das compras sem necessidade de digitação manual de itens.
- **Importação de XML:** Entrada de notas a partir do arquivo XML do fornecedor.

### 📍 Fase 5: Vendas de Balcão, Caixa e Expedição
- **Orçamentos por Ambientes:** Divisão de orçamentos e faturamento por cômodos (Sala, Quarto) com opção de anexar projetos arquitetônicos.
- **Faturamento Parcial:** Emissão de nota fiscal e recebimento apenas de produtos com estoque imediato, mantendo o pedido inicial aberto para os itens sob encomenda.
- **Frente de Caixa (PDV):** Abertura, fechamento, sangrias, recebimentos em cartões, pix ou dinheiro, e faturamento expresso.
- **Esteira de Expedição:** Venda (Balcão) $\rightarrow$ Caixa (Faturamento/Recebimento) $\rightarrow$ Separação no Estoque $\rightarrow$ Balcão de Retirada.

### 📍 Fase 6: Controle Financeiro, Bancos e Contabilidade
- **Tesouraria:** Contas a pagar/receber, saldos de contas bancárias e conciliação fácil de extratos.
- **Boletos:** Emissão de boletos integrados aos principais bancos.
- **Contabilidade Integrada:** Plano de contas contábeis estruturado, lançamentos automáticos por partida dobrada, geração de Balancete de Verificação e Demonstração do Resultado do Exercício (DRE) automática por filial.

### 📍 Fase 7: Painéis de Indicadores (Dashboards) Customizados
- Relatórios e gráficos específicos para a função de cada colaborador:
  - *Vendedores:* Metas pessoais, comissões acumuladas e contatos pendentes.
  - *Caixas:* Controle de fluxo diário de gaveta.
  - *Estoquistas:* Lista de mercadorias a separar e a receber.
  - *Sócios e Diretores:* Visão gerencial financeira unificada, DRE consolidada, margens de lucro de vendas e saúde financeira geral.

---

## 📈 Conclusão e Próximos Passos

Esta divisão permite o desenvolvimento modular e controlado do sistema. O projeto está preparado para nascer robusto, atendendo a todos os requisitos de segurança fiscal e financeira exigidos, enquanto provê uma interface de usuário ágil e moderna. 

O desenvolvimento começará com a **Fase 1** (Instalação e configuração visual da interface no padrão corporativo, e os conectores iniciais de importação de dados).
