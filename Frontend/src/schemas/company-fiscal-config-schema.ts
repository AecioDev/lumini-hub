import { z } from "zod";

// Upload de certificado (arquivo + senha) fica fora deste schema de
// propósito — vive como state irmão no componente, não gerenciado pelo
// react-hook-form, mesmo padrão do seletor de Perfil/Permissões em
// CreateUserForm.tsx (ver CLAUDE.md "Mandatory frontend CRUD pattern").
export const companyFiscalConfigSchema = z.object({
  tax_regime: z.enum(["SIMPLES", "PRESUMIDO", "REAL"], {
    message: "Selecione o tipo de tributação",
  }),
  accountant_name: z.string().max(150, "Máximo 150 caracteres").optional(),
  accountant_document: z.string().max(20, "Máximo 20 caracteres").optional(),
  accountant_contact: z.string().max(150, "Máximo 150 caracteres").optional(),
});

export type CompanyFiscalConfigFormValues = z.infer<typeof companyFiscalConfigSchema>;
