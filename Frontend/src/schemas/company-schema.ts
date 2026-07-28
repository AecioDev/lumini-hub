import { z } from "zod";

export const companySchema = z.object({
  parent_id: z.number().nullable(),
  legal_name: z.string().min(3, "Informe a razão social"),
  trade_name: z.string().optional(),
  tax_id: z.string().min(1, "Informe o CNPJ"),
  is_active: z.boolean(),
});

export type CompanyFormValues = z.infer<typeof companySchema>;
