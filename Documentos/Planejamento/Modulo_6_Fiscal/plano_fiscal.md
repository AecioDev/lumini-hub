# Módulo 6 - Módulo Fiscal Desacoplado (ACBr)

Centraliza as regras fiscais e a comunicação com a SEFAZ de forma isolada das tabelas cadastrais base.

## 🛠️ Especificações Gerais
1. **api.fiscal (Porta 4008):**
   - Configurações tributárias complementares (NCM, CEST, CST, CSOSN, CFOP, alíquotas) vinculadas aos cadastros de produtos e parceiros por IDs.
   - Gerenciamento seguro e criptografado de **Certificados Digitais A1** por filial.
2. **Emissão Fiscal via ACBrLib (DLL):**
   - Assinatura, validação e envio de documentos fiscais eletrônicos: **NFe, NFSe, NFCe e MDFe**.
   - Vantagem de ser um módulo plugável para facilitar a ativação ou suspensão de emissões por cliente.
