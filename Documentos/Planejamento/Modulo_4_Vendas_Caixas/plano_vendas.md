# Módulo 4 - Vendas, Caixas & Faturamento de Balcão

Especializado no fluxo de atendimento de balcão para o nicho de lojas de iluminação e decoração.

## 🛠️ Especificações Gerais
1. **Vendas de Balcão e Projetos (`api.sales` - Porta 4005):**
   - Agrupamento de orçamentos por **Ambientes da Obra** (Sala, Cozinha, Quarto, etc.) e relatórios direcionados.
   - Dossiê da venda com upload e associação de projetos arquitetônicos em PDF ou CAD.
   - Orçamentos de entrega futura e **Faturamento Parcial** (faturar apenas produtos com estoque físico disponível, mantendo o saldo do orçamento aberto).
2. **Controle de Caixas (`api.financial` - Porta 4006):**
   - Abertura, suprimentos, sangrias e fechamento de caixas.
   - Central de recebimento e faturamento presencial.
3. **Fluxo de Expedição:** Venda (Balcão) $\rightarrow$ Caixa (Faturamento) $\rightarrow$ Separação (Almoxarifado) $\rightarrow$ Balcão de Retirada.
