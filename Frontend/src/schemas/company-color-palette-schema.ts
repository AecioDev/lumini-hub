import { z } from "zod";

// Mesmo formato de cor exigido em company-visual-config-schema.ts.
const hexColor = z
  .string()
  .regex(/^#[0-9A-Fa-f]{6}$/, "Informe uma cor no formato hexadecimal #RRGGBB");

export const companyColorPaletteSchema = z.object({
  name: z.string().min(1, "Informe um nome pra paleta"),
  primary_color: hexColor,
  secondary_color: hexColor.optional().or(z.literal("")),
  accent_color: hexColor.optional().or(z.literal("")),
});

export type CompanyColorPaletteFormValues = z.infer<typeof companyColorPaletteSchema>;
