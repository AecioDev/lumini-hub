import { z } from "zod";

export const webhookEventSchema = z.object({
  id: z.number().int(),
  source: z.string(),
  event_type: z.string(),
  status: z.string(),
  received_at: z.string(),
  processed_at: z.string().nullable(),
  error_message: z.string(),
});

export type WebhookEvent = z.infer<typeof webhookEventSchema>;

export interface WebhookEventList {
  data: WebhookEvent[];
  pagination: {
    totalPages: number;
    totalRows: number;
  };
}
