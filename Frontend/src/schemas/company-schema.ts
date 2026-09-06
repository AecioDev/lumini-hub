import { z } from "zod";

// CNPJ só é obrigatório pra quem é Matriz (parent_id nulo) — uma empresa
// vinculada pode não ter CNPJ próprio quando a parte fiscal fica a cargo da
// Matriz (caso real do usuário, decidido 2026-09-05). Espelha a mesma regra
// de CompanyValidator.ValidateForCreation/ValidateForUpdate no backend.
export const companySchema = z
  .object({
    parent_id: z.number().nullable(),
    legal_name: z.string().min(3, "Informe a razão social"),
    trade_name: z.string().optional(),
    tax_id: z.string().optional(),
    is_active: z.boolean(),
  })
  .superRefine((values, ctx) => {
    if (values.parent_id === null && !values.tax_id?.trim()) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "CNPJ é obrigatório para a empresa Matriz",
        path: ["tax_id"],
      });
    }
  });

export type CompanyFormValues = z.infer<typeof companySchema>;
