import { z } from "zod";

const hexColor = z
  .string()
  .regex(/^#[0-9A-Fa-f]{6}$/, "Informe uma cor no formato hexadecimal #RRGGBB");

const optionalHexColor = z
  .string()
  .regex(/^#[0-9A-Fa-f]{6}$/, "Informe uma cor no formato hexadecimal #RRGGBB")
  .optional()
  .or(z.literal(""));

// Upload de logo fica fora deste schema de propósito — vive como state
// irmão no componente, mesmo padrão do certificado em
// company-fiscal-config-schema.ts.
export const companyVisualConfigSchema = z.object({
  primary_color: hexColor,
  secondary_color: optionalHexColor,
  accent_color: optionalHexColor,
  text_color: optionalHexColor,
});

export type CompanyVisualConfigFormValues = z.infer<typeof companyVisualConfigSchema>;
