# Módulo 7 - Contabilidade Gerencial

Responsável por registrar os reflexos econômicos e financeiros do ERP sob a ótica contábil.

## 🛠️ Especificações Gerais
1. **Inteligência Contábil:**
   - Plano de Contas flexível por empresa (modelagem detalhada na seção **📐 Modelagem** abaixo).
   - **Lançamentos em Partida Dobrada (Débito e Crédito):** Gerados automaticamente nas baixas de estoque, faturamento de vendas e compras, e liquidações de contas a pagar/receber.
2. **Relatórios Contábeis:** Razão, Balancete de Verificação e Demonstração do Resultado do Exercício (DRE).
3. **Exportador Contábil:** Layouts de integração com sistemas contábeis parceiros e geração de arquivos fiscais.

---

## 📐 Modelagem

> Modelagem movida de `plano_empresa.md` (Módulo 0 - Configuração) em 2026-09-13, decisão do usuário — Plano de Contas pertence de verdade ao Módulo 7 (Contabilidade), não ao Módulo 0 (Empresa), mesmo raciocínio já aplicado ao CFG-2 (Numeração Fiscal, movido pro Módulo 6). Conteúdo trazido integralmente, sem simplificação.

### `ChartOfAccounts` (N:1 com Empresa, self-referencing)
**Definido 2026-07-21**: template global padrão (o mais comumente usado no Brasil), clonado pra cada empresa na criação — facilita a vida do cliente, que já começa operando sem montar plano de contas do zero. Depois de clonado, cada empresa pode customizar o próprio livremente.

| Campo | Tipo | Observação |
|---|---|---|
| `CompanyID` | uint | FK — cada empresa tem sua própria cópia, clonada do template padrão na criação |
| `ParentID` | *uint (nullable) | self-referencing — hierarquia tipo "3.1.01.001 Receita de Vendas" sob "3.1 Receitas Operacionais" sob "3 Contas de Resultado" |
| `Code` | string | ex: `3.1.01.001` |
| `Name` | string | |
| `Type` | enum | Ativo / Passivo / Receita / Despesa / Resultado (a validar terminologia contábil exata) |

---

## ❓ Decisões em Aberto

- [ ] **O conteúdo real do Plano de Contas padrão.** Preciso da lista de contas em si (código + nome + tipo, hierarquia completa) pra ter o que clonar em cada empresa nova. Isso é dado contábil/de negócio, não uma decisão de engenharia — o usuário tem um plano de contas de referência (do contador, de outro sistema, etc.) pra eu usar como seed, ou monto uma sugestão padrão simplificada pra começar e ajustamos depois?
