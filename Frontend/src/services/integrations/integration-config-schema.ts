import { z } from "zod";

export const integrationSettingsSchema = z.object({
  li_api_key: z.string(),
  li_app_key: z.string(),
  li_webhook_secret: z.string(),
  codtipnot: z.string(),
  codlocarm_oficial: z.string(),
  codlocarm_reserva: z.string(),
  codemp: z.string(),
});

export type IntegrationSettings = z.infer<typeof integrationSettingsSchema>;

export type UpdateIntegrationSettingsData = Partial<IntegrationSettings>;
