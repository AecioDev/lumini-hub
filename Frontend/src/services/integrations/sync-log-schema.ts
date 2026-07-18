import { z } from "zod";

export const syncLogSchema = z.object({
  id: z.number().int(),
  direction: z.string(),
  entity_type: z.string(),
  reference_id: z.string(),
  status: z.string(),
  message: z.string(),
  created_at: z.string(),
});

export type SyncLog = z.infer<typeof syncLogSchema>;

export interface SyncLogList {
  data: SyncLog[];
  pagination: {
    totalPages: number;
    totalRows: number;
  };
}
