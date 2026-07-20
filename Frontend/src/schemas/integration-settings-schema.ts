import { z } from "zod";

// Espelha UpdateIntegrationSettingsRequest — todos os campos são opcionais
// (upsert em lote no backend, campos omitidos não são alterados).
export const integrationSettingsSchema = z.object({
  li_api_key: z.string().optional(),
  li_app_key: z.string().optional(),
  li_webhook_secret: z.string().optional(),
  codtipnot: z.string().optional(),
  codlocarm_oficial: z.string().optional(),
  codlocarm_reserva: z.string().optional(),
  codemp: z.string().optional(),
});

export type IntegrationSettingsFormValues = z.infer<typeof integrationSettingsSchema>;
